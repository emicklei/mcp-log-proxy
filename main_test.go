package main

import "testing"

func Test_parseSequenceID(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid JSON with string ID",
			args: args{line: `{"method":"tools/list","jsonrpc":"2.0","id":"12345"}`},
			want: "12345",
		},
		{
			name: "Valid JSON with integer ID",
			args: args{line: `{"method":"tools/list","jsonrpc":"2.0","id":67890}`},
			want: "67890",
		},
		{
			name: "Valid JSON without ID",
			args: args{line: `{"method":"tools/list","jsonrpc":"2.0"}`},
			want: "?",
		},
		{
			name: "Invalid JSON",
			args: args{line: `{"method":"tools/list","jsonrpc":"2.0","id":}`},
			want: "?",
		},
		{
			name: "Empty string",
			args: args{line: ``},
			want: "?",
		},
		{
			name: "Non-JSON string",
			args: args{line: `This is not JSON`},
			want: "?",
		},
		{
			name: "JSON with non-string/non-integer ID",
			args: args{line: `{"method":"tools/list","jsonrpc":"2.0","id":true}`},
			want: "?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSequenceID(tt.args.line); got != tt.want {
				t.Errorf("parseSequenceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
