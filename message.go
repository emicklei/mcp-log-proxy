package main

import (
	"encoding/json"
	"fmt"
)

// {"method":"tools/list","jsonrpc":"2.0","id":1}
func parseJSONMessage(line string) (map[string]any, bool) {
	if (line == "") || (line[0] != '{') {
		return nil, false
	}
	m := map[string]any{}
	err := json.Unmarshal([]byte(line), &m)
	if err != nil {
		return nil, false
	}
	return m, true
}

func getMessageID(m map[string]any) string {
	if id, ok := m["id"]; ok {
		switch id := id.(type) {
		case string:
			return id
		case int:
			return fmt.Sprintf("%d", id)
		case float64:
			return fmt.Sprintf("%v", id)
		}
	}
	return "?"
}

func isErrorMessage(m map[string]any) bool {
	if _, ok := m["error"]; ok {
		return true
	}
	if _, ok := m["result"]; ok {
		return false
	}
	if _, ok := m["content"]; ok {
		cm := m["content"].(map[string]any)
		if v, ok := cm["isError"]; ok {
			return v == true
		}
	}
	return false
}

func isWarnMessage(m map[string]any) bool {
	if doc, ok := m["error"]; ok {
		if docMap, ok := doc.(map[string]any); ok {
			if code, ok := docMap["code"]; ok {
				if code == -32603 {
					return true
				}
			}
		}
	}
	return false
}
