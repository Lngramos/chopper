package tools

import (
	"errors"
)

type Tool interface {
	Name() string
	Call(args map[string]interface{}) (string, error)
}

var registry = make(map[string]Tool)

func Register(tool Tool) {
	registry[tool.Name()] = tool
}

func Get(name string) (Tool, bool) {
	t, ok := registry[name]
	return t, ok
}

func CallTool(name string, args map[string]interface{}) (string, error) {
	tool, ok := Get(name)
	if !ok {
		return "", errors.New("tool not found: " + name)
	}
	return tool.Call(args)
}

func ListTools() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
