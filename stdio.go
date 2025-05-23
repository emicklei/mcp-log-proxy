package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"maps"
)

func runTargetToClient(ctx context.Context, stdout io.ReadCloser) {
	lineReader := bufio.NewReader(stdout)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := lineReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					slog.Error("failed to read from stdout", "error", err)
				}
				return
			}
			// id hack for clients that cannot handle null
			line = strings.ReplaceAll(line, `"id":null`, `"id":""`)

			isResponse := log(" ... client <= proxy <= target", "line", line, false)
			if isResponse {
				io.WriteString(os.Stdout, line) // \n is part of the line
			} else {
				slog.Debug("not a JSON response message", "line", line, "length", len(line))
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
			isRequest := log(" client => proxy => target", "line", line, true)
			if isRequest {
				io.WriteString(stdin, line) // \n is part of the line
			} else {
				slog.Debug("not a JSON request message", "line", line, "length", len(line))
			}
		}
	}
}

// returns true if the message is a JSON message
func log(flow, lineKey, line string, toServer bool) bool {
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
			if m, ok := msg["method"]; ok {
				flow = m.(string)
			}
			if m, ok := msg["result"]; ok {
				mm := m.(map[string]any)
				flow = strings.Join(slices.Collect(maps.Keys(mm)), ", ")
			}
			level = slog.LevelInfo
		}
		id = getMessageID(msg)
	}
	traffic := "request"
	if !toServer {
		traffic = "response"
	}
	slog.Default().With("traffic", traffic).Log(context.Background(), level, fmt.Sprintf("%s:%s", id, flow), lineKey, line)
	return ok
}
