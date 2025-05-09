package main

import (
	"net/http"
	"strconv"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><h1>Dashboard</h1></body></html>"))
	w.Write([]byte("<h2>Instances</h2>"))
	instances, err := readRegistryEntries()
	if err != nil {
		w.Write([]byte("<p>Error reading registry: " + err.Error() + "</p>"))
		return
	}
	w.Write([]byte("<ul>"))
	for _, instance := range instances {
		w.Write([]byte("<li>" + instance.Host + ":" + strconv.Itoa(instance.Port) + "</li>"))
	}
	w.Write([]byte("</ul>"))
}
