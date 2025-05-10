package tools

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

type RunTool struct{}

func (t *RunTool) Name() string {
	return "run"
}

func (t *RunTool) Call(args map[string]interface{}) (string, error) {
	cmdStr, ok := args["command"].(string)
	if !ok || cmdStr == "" {
		return "", errors.New("missing or invalid 'command' argument")
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %w\n%s", err, out.String())
	}
	return out.String(), nil
}
