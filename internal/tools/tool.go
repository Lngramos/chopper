package tools

import (
	"fmt"
	"strings"
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

func CallTool(name string, args map[string]interface{}, safeMode bool) (string, error) {
	// If safe mode is on, confirm before running sensitive tools
	if safeMode && name == "run" {
		fmt.Printf("⚠️ Confirm execution of tool '%s' with arguments: %v [y/N]: ", name, args)
		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))
		if input != "y" && input != "yes" {
			return "Tool execution cancelled.", nil
		}
	}

	tool, ok := Get(name)
	if !ok {
		return "", fmt.Errorf("tool '%s' not found", name)
	}
	return tool.Call(args)
}
