package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/emicklei/mcp-log-proxy/lockedfile"
)

type proxyInstance struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Title   string `json:"title"`
	Command string `json:"command"`
}

func addToOrRemoveFromRegistry(inst proxyInstance, isRemove bool) error {
	err := lockedfile.Transform(getRegistryLocation(), func(stored []byte) ([]byte, error) {
		list := []proxyInstance{}
		if len(stored) > 0 { // file has content
			err := json.Unmarshal(stored, &list)
			if err != nil {
				return nil, err
			}
		}
		if isRemove {
			withRemoval := []proxyInstance{}
			for _, each := range list {
				if !(each.Host == inst.Host && each.Port == inst.Port) {
					withRemoval = append(withRemoval, each)
				}
			}
			list = withRemoval
		} else {
			list = append(list, inst)
		}
		stored, err := json.Marshal(list)
		if err != nil {
			return nil, err
		}
		return stored, nil
	})
	return err
}

func readRegistryEntries() ([]proxyInstance, error) {
	content, err := lockedfile.Read(getRegistryLocation())
	if err != nil {
		return nil, err
	}
	list := []proxyInstance{}
	if len(content) > 0 { // file has content
		err := json.Unmarshal(content, &list)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func getRegistryLocation() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, ".mcp-log-proxy-instances.json")
}
