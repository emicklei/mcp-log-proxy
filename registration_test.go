package main

import (
	"os"
	"testing"
	"time"

	"github.com/emicklei/mcp-log-proxy/lockedfile"
)

func TestUpdateInstances(t *testing.T) {
	*registryLocation = "/tmp/test-instances.log"

	pi := proxyInstance{
		Host:    "localhost",
		Port:    5656,
		Title:   "mcp-log-proxy",
		Command: "echo hello",
	}
	err := addToOrRemoveFromRegistry(pi, false)
	if err != nil {
		t.Fatalf("failed to register instance: %v", err)
	}
	one, err := readRegistryEntries()
	if err != nil {
		t.Fatalf("failed to read instance: %v", err)
	}
	if len(one) != 1 {
		t.Fatalf("expected 1 instance, got %d", len(one))
	}
	err = addToOrRemoveFromRegistry(pi, true)
	if err != nil {
		t.Fatalf("failed to register instance: %v", err)
	}
	none, err := readRegistryEntries()
	if err != nil {
		t.Fatalf("failed to read registry: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("expected 0 registry, got %d", len(none))
	}
	all, _ := os.ReadFile(getRegistryLocation())
	if len(all) != 2 {
		t.Fatalf("expected '[]' in file, got %d bytes : '%s'", len(all), string(all))
	}
}

func TestLock(t *testing.T) {
	f1, err := lockedfile.OpenFile("instances.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("failed to open locked file: %v", err)
	}
	defer f1.Close()
	go func() {
		time.Sleep(1 * time.Second)
		f1.Close()
	}()
	t.Log("waiting for lock to be released")
	f2, err := lockedfile.OpenFile("instances.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("failed to open locked file: %v", err)
	}
	defer f2.Close()
	t.Log("got second lock")
}
