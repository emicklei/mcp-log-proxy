package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func runTargetToClient(ctx context.Context, stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if scanner.Scan() {
				if err := scanner.Err(); err != nil {
					if err != io.EOF {
						slog.Error("failed to read from stdout", "error", err)
					}
					return
				}
				line := scanner.Text()

				// id hack https://github.com/mark3labs/mcp-go/issues/201
				line = strings.ReplaceAll(line, `"id":null`, `"id":""`)

				isResponse := log(" ... client <= proxy <= target", "jsonresponse", line)
				if isResponse {
					io.WriteString(os.Stdout, line)
					io.WriteString(os.Stdout, "\n")
				} else {
					slog.Debug("not a JSON response message", "line", line, "length", len(line))
				}
			}
		}
	}
}

func runClientToTarget(ctx context.Context, stdin io.WriteCloser) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					slog.Error("failed to read from stdin", "error", err)
				}
				return
			}
			isRequest := log(" client => proxy => target", "jsonrequest", line)
			if isRequest {
				io.WriteString(stdin, line) // \n is part of the line
			} else {
				slog.Debug("not a JSON request message", "line", line, "length", len(line))
			}
		}
	}
}

// returns true if the message is a JSON message
func log(flow, key, line string) bool {
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
	slog.Log(context.Background(), level, fmt.Sprintf("%s:%s", id, flow), key, line)
	return ok
}
