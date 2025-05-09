package main

import "testing"

func TestUpdateInstance(t *testing.T) {
	pi := proxyInstance{
		host:    "localhost",
		port:    5656,
		title:   "mcp-log-proxy",
		command: "echo hello",
	}
	err := register(pi)
	if err != nil {
		t.Fatalf("failed to register instance: %v", err)
	}
}
