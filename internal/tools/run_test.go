package tools

import (
	"strings"
	"testing"
)

func TestRunTool_Echo(t *testing.T) {
	rt := RunTool{}
	output, err := rt.Call(map[string]interface{}{"command": "echo hello"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("expected 'hello' in output, got %q", output)
	}
}

func TestRunTool_InvalidCommand(t *testing.T) {
	rt := RunTool{}
	_, err := rt.Call(map[string]interface{}{"command": "nonexistentcommand123"})
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestRunTool_MissingCommand(t *testing.T) {
	rt := RunTool{}
	_, err := rt.Call(map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error when command is missing")
	}
}
