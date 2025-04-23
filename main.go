package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/emicklei/nanny"
)

var (
	targetCommand = flag.String("command", "", "full command with arguments")
	errLog        = flag.String("log", "mcp-log-proxy.log", "file to append errors to")
	httPort       = flag.Int("port", 5656, "port to listen on")
	pageTitle     = flag.String("title", "mcp-log-proxy", "title of the web page")
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
	rec := nanny.NewRecorder()
	reclog := slog.New(nanny.NewLogHandler(rec, logHandler, slog.LevelInfo))
	slog.SetDefault(reclog)
	http.Handle("/", nanny.NewBrowser(rec, nanny.BrowserOptions{PageSize: 100, PageTitle: *pageTitle}))

	// serve nanny
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", *httPort), nil)
	}()

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

	// client -> proxy -> target
	runClientToTarget(stdin)

	// target -> proxy -> client
	runTargetToClient(stdout)

	// run target command
	if err := cmd.Start(); err != nil {
		slog.Error("failed to start target command", "error", err, "command", parts[0], "args", parts[1:])
		os.Exit(1)
	}
	// wait for target command to finish
	if err := cmd.Wait(); err != nil {
		slog.Error("failed to wait for target command", "error", err, "command", parts[0], "args", parts[1:])
		os.Exit(1)
	}
}

func runTargetToClient(stdout io.ReadCloser) {
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log(" ... client <= proxy <= target", line)
			fmt.Fprintln(os.Stdout, line)
		}
	}()
}

func log(flow, line string) {
	level := slog.LevelError
	id := "?"
	msg, ok := parseJSONMessage(line)
	if ok {
		if isErrorMessage(msg) {
			// maybe warning
			if isWarnMessage(msg) {
				level = slog.LevelWarn
			}
		} else {
			level = slog.LevelInfo
		}
		id = getMessageID(msg)
	}
	slog.Log(context.Background(), level, fmt.Sprintf("%s:%s", id, flow), "line", line)
}

func runClientToTarget(stdin io.WriteCloser) {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				slog.Error("failed to read from stdin", "error", err)
				os.Exit(1)
			}
			log(" client => proxy => target", line)
			fmt.Fprintln(stdin, line)
		}
	}()
}
