package tools

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Call(args map[string]interface{}, out io.Writer) error {
	path, ok := args["path"].(string)
	if !ok {
		return errors.New("missing 'path' argument")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(out, string(content))
	return err
}
