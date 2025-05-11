package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lngramos/chopper/internal/ollama"
	"github.com/lngramos/chopper/internal/tools"
	"github.com/spf13/cobra"
)

var chatModel string
var chatTemperature float64

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a single prompt to Ollama and get a response",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		client := ollama.NewClient("http://localhost:11434")

		systemMessage := ollama.Message{
			Role: "system",
			Content: `You are Chopper, a command-line assistant.

You can call the following tools by replying with a JSON object:
{
  "tool_call": {
    "name": "run",
    "arguments": {
      "command": "ls -la"
    }
  }
}

Available tools:
- run(command: string): Execute a shell command and return output.
- read_file(path: string): Read contents of a file at the given path.
Only return a valid JSON object when calling a tool.`,
		}

		messages := []ollama.Message{
			systemMessage,
			{Role: "user", Content: prompt},
		}

		response, err := client.Chat(chatModel, chatTemperature, messages)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonStart := strings.Index(response, "{")
		if jsonStart >= 0 {
			jsonPart := response[jsonStart:]
			var toolCheck ollama.ToolCheck
			if err := json.Unmarshal([]byte(jsonPart), &toolCheck); err == nil {
				toolCheck.Debug()
				if toolCheck.ToolCall != nil {
					result, err := tools.CallTool(toolCheck.ToolCall.Name, toolCheck.ToolCall.Arguments, !unsafeMode)
					if err != nil {
						fmt.Println("Tool error:", err)
						os.Exit(1)
					}
					fmt.Println(result)
					return
				} else if len(toolCheck.ToolCalls) > 0 {
					for _, call := range toolCheck.ToolCalls {
						result, err := tools.CallTool(call.Name, call.Arguments, !unsafeMode)
						if err != nil {
							fmt.Printf("Tool error [%s]: %v\n", call.Name, err)
							continue
						}
						fmt.Println(result)
					}
					return
				}
			}
		}

		fmt.Println(response)
	},
}
