package main

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func registerAndStartListener(whoami *proxyInstance) {
	if err := addToOrRemoveFromRegistry(whoami, false); err != nil {
		slog.Error("failed to add to registry", "error", err)
	} else {
		slog.Debug("added to registry", "instance", whoami)
	}
	// use the given port to listen, fall back to a free port if it is already in use
	if err := http.ListenAndServe(net.JoinHostPort("localhost", strconv.Itoa(whoami.Port)), nil); err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			// try with a different port
			newPort, err := GetFreePort()
			if err == nil {
				whoami.Port = newPort
				// add new port to registry
				addToOrRemoveFromRegistry(whoami, false)
				http.ListenAndServe(net.JoinHostPort("localhost", strconv.Itoa(whoami.Port)), nil)
			} else {
				slog.Error("failed to get free port, cannot start log service", "error", err)
			}
		} else {
			slog.Error("failed to start HTTP log service", "error", err)
		}
	}
	addToOrRemoveFromRegistry(whoami, true)
}
