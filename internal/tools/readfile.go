package tools

import (
	"errors"
	"os"
)

type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Call(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return "", errors.New("missing or invalid 'path' argument")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func init() {
	Register(&RunTool{})
	Register(&ReadFileTool{})
}
