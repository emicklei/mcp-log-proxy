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

// TODO make struct
func isErrorMessage(m map[string]any) bool {
	if _, ok := m["error"]; ok {
		return true
	}
	if result, ok := m["result"].(map[string]any); ok {
		if list, ok := result["content"].([]any); ok {
			for _, item := range list {
				if itemDoc, ok := item.(map[string]any); ok {
					if err, ok := itemDoc["isError"]; ok {
						if is, ok := err.(bool); ok {
							return is
						}
					}
				}
			}
		}
	}
	return false
}

// https://modelcontextprotocol.io/specification/2025-03-26
const (
	RESOURCE_NOT_FOUND = float64(-32002)
	METHOD_NOT_FOUND   = float64(-32601)
)

func isWarnMessage(m map[string]any) bool {
	if doc, ok := m["error"]; ok {
		if docMap, ok := doc.(map[string]any); ok {
			if code, ok := docMap["code"]; ok {
				switch code {
				case RESOURCE_NOT_FOUND, METHOD_NOT_FOUND:
					return true
				default:
					return false
				}
			}
		}
	}
	return false
}
