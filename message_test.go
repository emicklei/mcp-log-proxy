package main

import "testing"

func TestWarnMessage(t *testing.T) {
	msg := `{
		"jsonrpc": "2.0",
		"id": 13,
		"error": {
			"code": -32601,
			"message": "resources not supported"
		}}`
	m, ok := parseJSONMessage(msg)
	if !ok {
		t.Errorf("failed to parse message: %s", msg)
	}
	is := isWarnMessage(m)
	if !is {
		t.Errorf("expected isWarnMessage to be true, got false")
	}
}

func TestErrorMessage(t *testing.T) {
	msg := `{
		"result": {
			"content": [
				{
					"type": "text",
					"text": "bla",
					"isError": true
				}
			]
		},
		"jsonrpc": "2.0",
		"id": 46
	}`
	m, ok := parseJSONMessage(msg)
	if !ok {
		t.Errorf("failed to parse message: %s", msg)
	}
	if !isErrorMessage(m) {
		t.Errorf("expected error message to be true, got false")
	}
}
