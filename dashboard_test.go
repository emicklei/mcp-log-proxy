package main

import (
	"strings"
	"testing"
)

func TestBeforeTableHTMLUsesCurrentHost(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	current := &proxyInstance{
		Host:    "0.0.0.0",
		Port:    5656,
		Title:   "mcp-log-proxy",
		Command: "echo hello",
	}
	other := &proxyInstance{
		Host:    "localhost",
		Port:    5657,
		Title:   "other",
		Command: "echo other",
	}
	if err := addToOrRemoveFromRegistry(current, false); err != nil {
		t.Fatalf("failed to register current instance: %v", err)
	}
	if err := addToOrRemoveFromRegistry(other, false); err != nil {
		t.Fatalf("failed to register other instance: %v", err)
	}

	html := (&instanceSelector{current: current}).beforeTableHTML()
	if !strings.Contains(html, `<option value="http://0.0.0.0:5656" selected>`) {
		t.Fatalf("expected current host to be selected, got %s", html)
	}
	if strings.Contains(html, `<option value="http://localhost:5657" selected>`) {
		t.Fatalf("expected other host to remain unselected, got %s", html)
	}
}
