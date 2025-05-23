package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/emicklei/nanny"
)

var (
	targetCommand = flag.String("command", "", "full command with arguments")
	errLog        = flag.String("log", "mcp-log-proxy.log", "file to append errors to")
	httPort       = flag.Int("port", 5656, "port to listen on")
	pageTitle     = flag.String("title", "mcp-log-proxy", "title of the web page")
	isDebug       = flag.Bool("debug", false, "enable debug logging")
)

func main() {
	flag.Parse()
	// check if target command is provided
	if *targetCommand == "" {
		flag.Usage()
		return
	}
	// open error log file
	logFile, err := os.OpenFile(*errLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open error log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// setup nanny
	logHandler := slog.NewTextHandler(logFile, nil)
	rec := nanny.NewRecorder(nanny.RecorderOptions{
		GroupMarkers: []string{"traffic"},
	})
	lvl := slog.LevelInfo
	if *isDebug {
		lvl = slog.LevelDebug
	}

	whoami := &proxyInstance{
		Host:    "localhost",
		Port:    *httPort,
		Title:   *pageTitle,
		Command: *targetCommand,
	}
	selector := &instanceSelector{
		current: whoami,
	}

	reclog := slog.New(nanny.NewLogHandler(rec, logHandler, lvl))
	slog.SetDefault(reclog)
	options := nanny.BrowserOptions{
		PageSize:            100,
		PageTitle:           *pageTitle,
		BeforeHTMLTableFunc: selector.beforeTableHTML,
		EndHTMLHeadFunc:     selector.endHeadHTML,
	}
	http.Handle("/", nanny.NewBrowser(rec, options))

	// run target command
	parts := strings.Split(*targetCommand, " ")
	cmd := exec.Command(parts[0], parts[1:]...)

	// set up pipes for stdin and stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		slog.Error("failed to get stdin pipe", "error", err)
		os.Exit(1)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to get stdout pipe", "error", err)
		os.Exit(1)
	}

	// to stop stdio
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		slog.Info("received termination signal, shutting down")
		cancel()
		abort(whoami, 0)
	}()

	// serve nanny
	go registerAndStartListener(whoami)

	// client -> proxy -> target
	go func() {
		runClientToTarget(ctx, stdin)
		cancel()
		abort(whoami, 0)
	}()

	// target -> proxy -> client
	go func() {
		runTargetToClient(ctx, stdout)
		cancel()
		abort(whoami, 0)
	}()

	// run target command
	if err := cmd.Start(); err != nil {
		slog.Error("failed to start target command", "error", err, "command", parts[0], "args", parts[1:])
		abort(whoami, 1)
	}
	// wait for target command to finish
	if err := cmd.Wait(); err != nil {
		slog.Error("failed to wait for target command", "error", err, "command", parts[0], "args", parts[1:])
		abort(whoami, 1)
	}
}

func abort(pi *proxyInstance, code int) {
	if err := addToOrRemoveFromRegistry(pi, true); err != nil {
		slog.Error("failed to remote from registry on abort", "error", err)
	}
	os.Exit(code)
}
