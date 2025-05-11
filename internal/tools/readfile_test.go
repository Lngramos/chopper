package tools

import (
	"os"
	"testing"
)

func TestReadFileTool_ValidPath(t *testing.T) {
	// Setup: create a temp file
	tmpFile, err := os.CreateTemp("", "chopper_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	expected := "hello world"
	tmpFile.WriteString(expected)
	tmpFile.Close()

	rt := ReadFileTool{}
	content, err := rt.Call(map[string]interface{}{"path": tmpFile.Name()})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestReadFileTool_MissingPath(t *testing.T) {
	rt := ReadFileTool{}
	_, err := rt.Call(map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error when path is missing")
	}
}

func TestReadFileTool_NonExistentPath(t *testing.T) {
	rt := ReadFileTool{}
	_, err := rt.Call(map[string]interface{}{"path": "/tmp/this_should_not_exist.chopper"})
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}
