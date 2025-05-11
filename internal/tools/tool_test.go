package tools

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCallTool_UnknownTool(t *testing.T) {
	var buf bytes.Buffer
	err := CallTool("not_a_tool", map[string]interface{}{}, false, &buf)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestCallTool_DispatchReadFile(t *testing.T) {
	Register(&ReadFileTool{}) // pointer value

	tmpFile, _ := os.CreateTemp("", "chopper_test")
	tmpFile.WriteString("test")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	var buf bytes.Buffer
	err := CallTool("read_file", map[string]interface{}{"path": tmpFile.Name()}, false, &buf)
	if err != nil {
		t.Fatalf("expected no error for read_file, got: %v", err)
	}
	if !strings.Contains(buf.String(), "test") {
		t.Errorf("expected file content, got: %s", buf.String())
	}
}

func TestCallTool_SafeMode_Confirm(t *testing.T) {
	Register(&RunTool{})

	restore := mockStdin("y\n")
	defer restore()

	var buf bytes.Buffer
	err := CallTool("run", map[string]interface{}{"command": "echo test"}, true, &buf)
	if err != nil {
		t.Fatalf("expected run to succeed in safe mode with 'yes', got: %v", err)
	}
	if !strings.Contains(buf.String(), "test") {
		t.Errorf("expected 'test' in output, got: %s", buf.String())
	}
}

func TestCallTool_SafeMode_Reject(t *testing.T) {
	Register(&RunTool{})

	restore := mockStdin("n\n")
	defer restore()

	var buf bytes.Buffer
	err := CallTool("run", map[string]interface{}{"command": "echo test"}, true, &buf)
	if err != nil {
		t.Fatalf("expected no error, just cancelled, got: %v", err)
	}
	if !strings.Contains(buf.String(), "cancelled") {
		t.Errorf("expected cancellation message, got: %s", buf.String())
	}
}

// helper to mock os.Stdin
func mockStdin(input string) func() {
	r, w, _ := os.Pipe()
	_, _ = w.WriteString(input)
	w.Close()

	orig := os.Stdin
	os.Stdin = r
	return func() {
		os.Stdin = orig
	}
}
