package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/emicklei/nanny"
)

var targetCommand = flag.String("command", "", "full command with arguments")

func main() {
	flag.Parse()
	if *targetCommand == "" {
		flag.Usage()
		return
	}
	errOut, _ := os.OpenFile(
		"/Users/ernestmicklei/Documents/mcp-log-proxy.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	defer errOut.Close()

	logHandler := slog.NewTextHandler(errOut, nil)
	rec := nanny.NewRecorder()
	reclog := slog.New(nanny.NewLogHandler(rec, logHandler, slog.LevelInfo))
	slog.SetDefault(reclog)
	http.Handle("/", nanny.NewBrowser(rec, nanny.BrowserOptions{PageSize: 100}))

	go func() {
		http.ListenAndServe(":5656", nil)
	}()

	parts := strings.Split(*targetCommand, " ")

	// msg := map[string]any{
	// 	"command": parts[0],
	// 	"args":    parts[1:],
	// }
	// json.NewEncoder(os.Stderr).Encode(msg)

	cmd := exec.Command(parts[0], parts[1:]...)
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
	// stderr?

	// client -> proxy -> target
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				slog.Error("failed to read from stdin", "error", err)
				os.Exit(1)
			}
			slog.Info("client -> proxy -> target", "line", line)
			fmt.Fprintln(stdin, line)
		}
	}()
	// target -> proxy -> client
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			slog.Info("target -> proxy -> client", "line", line)
			fmt.Fprintln(os.Stdout, line)
		}
	}()
	if err := cmd.Start(); err != nil {
		slog.Error("failed to start target command", "error", err, "command", parts[0], "args", parts[1:])
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		slog.Error("failed to wait for target command", "error", err, "command", parts[0], "args", parts[1:])
		os.Exit(1)
	}
}
