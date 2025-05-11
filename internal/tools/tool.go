package tools

import (
	"fmt"
	"io"
)

// Tool defines a common interface for CLI-accessible tools.
type Tool interface {
	Name() string
	Call(args map[string]interface{}, out io.Writer) error
}

var registry = make(map[string]Tool)

func Register(t Tool) {
	registry[t.Name()] = t
}

// CallTool dispatches the tool by name, optionally prompting in safe mode.
// Output is written to the provided io.Writer.
func CallTool(name string, args map[string]interface{}, safe bool, out io.Writer) error {
	tool, ok := registry[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}

	if safe {
		if !confirmToolExecution(name, args) {
			fmt.Fprintln(out, "Tool execution cancelled.")
			return nil
		}
	}

	return tool.Call(args, out)
}

func confirmToolExecution(name string, args map[string]interface{}) bool {
	fmt.Printf("⚠️ Confirm execution of tool '%s' with arguments: %+v [y/N]: ", name, args)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	switch response {
	case "y", "Y", "yes", "YES":
		return true
	default:
		return false
	}
}
