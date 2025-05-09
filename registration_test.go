package main

import (
	"os"
	"testing"
)

func TestUpdateInstances(t *testing.T) {
	registryLocation = "/tmp/test-instances.log"

	pi := proxyInstance{
		Host:    "localhost",
		Port:    5656,
		Title:   "mcp-log-proxy",
		Command: "echo hello",
	}
	err := updateRegistry(pi, false)
	if err != nil {
		t.Fatalf("failed to register instance: %v", err)
	}
	one, err := readRegistry()
	if err != nil {
		t.Fatalf("failed to read instance: %v", err)
	}
	if len(one) != 1 {
		t.Fatalf("expected 1 instance, got %d", len(one))
	}
	err = updateRegistry(pi, true)
	if err != nil {
		t.Fatalf("failed to register instance: %v", err)
	}
	none, err := readRegistry()
	if err != nil {
		t.Fatalf("failed to read registry: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("expected 0 registry, got %d", len(none))
	}
	all, _ := os.ReadFile("instances.json")
	if len(all) != 2 {
		t.Fatalf("expected [] file, got %d bytes", len(all))
	}
}
