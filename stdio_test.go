package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestRunTargetToClient(t *testing.T) {
	discardSlogging()
	targetOutput := io.NopCloser(strings.NewReader(`{"id":1}
`))
	output, _ := captureOutput(func() error {
		runTargetToClient(context.Background(), targetOutput)
		return nil
	})
	if got, want := output, `{"id":1}
`; got != want {
		t.Fatalf("got [%s] want [%s]", got, want)
	}
}

func discardSlogging() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}
