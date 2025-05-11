package tools

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type RunTool struct{}

func (t *RunTool) Name() string {
	return "run"
}

func (t *RunTool) Call(args map[string]interface{}, out io.Writer) error {
	cmdStrRaw, ok := args["command"]
	if !ok {
		return fmt.Errorf("missing 'command' argument")
	}
	cmdStr, ok := cmdStrRaw.(string)
	if !ok {
		return fmt.Errorf("'command' must be a string")
	}

	// 🔧 Fix double-escaped % symbols
	cmdStr = strings.ReplaceAll(cmdStr, "%%", "%")

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	_, _ = out.Write(output)

	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
