package tools

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunTool_Echo(t *testing.T) {
	rt := RunTool{}
	var buf bytes.Buffer

	err := rt.Call(map[string]interface{}{"command": "echo hello"}, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "hello") {
		t.Errorf("expected 'hello' in output, got %q", output)
	}
}

func TestRunTool_InvalidCommand(t *testing.T) {
	rt := RunTool{}
	var buf bytes.Buffer

	err := rt.Call(map[string]interface{}{"command": "nonexistentcommand123"}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestRunTool_MissingCommand(t *testing.T) {
	rt := RunTool{}
	var buf bytes.Buffer

	err := rt.Call(map[string]interface{}{}, &buf)
	if err == nil {
		t.Fatal("expected error when command is missing")
	}
}
