package main

import (
	"encoding/json"

	"github.com/emicklei/mcp-log-proxy/lockedfile"
)

type proxyInstance struct {
	host    string
	port    int
	title   string
	command string
}

func register(inst proxyInstance) error {
	err := lockedfile.Transform("instances.json", func(stored []byte) ([]byte, error) {
		list := []proxyInstance{}
		if len(stored) > 0 { // file was created
			err := json.Unmarshal(stored, &list)
			if err != nil {
				return nil, err
			}
		}
		list = append(list, inst)
		stored, err := json.Marshal(list)
		if err != nil {
			return nil, err
		}
		return stored, nil
	})
	return err
}
