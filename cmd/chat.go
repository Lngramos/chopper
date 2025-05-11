package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lngramos/chopper/internal/llm"
	"github.com/lngramos/chopper/internal/tools"
	"github.com/spf13/cobra"
)

func NewChatCommand(client llm.Client) *cobra.Command {
	var model string
	var temperature float64

	cmd := &cobra.Command{
		Use:   "chat [prompt]",
		Short: "Send a single prompt to LLM and get a response",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			prompt := args[0]

			systemMessage := llm.Message{
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

			messages := []llm.Message{
				systemMessage,
				{Role: "user", Content: prompt},
			}

			response, err := client.Chat(model, temperature, messages)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			jsonStart := strings.Index(response, "{")
			if jsonStart >= 0 {
				jsonPart := response[jsonStart:]
				var toolCheck llm.ToolCheck
				if err := json.Unmarshal([]byte(jsonPart), &toolCheck); err == nil {
					toolCheck.Debug()

					if toolCheck.ToolCall != nil {
						err := tools.CallTool(toolCheck.ToolCall.Name, toolCheck.ToolCall.Arguments, !unsafeMode, os.Stdout)
						if err != nil {
							fmt.Println("Tool error:", err)
							os.Exit(1)
						}
						return
					} else if len(toolCheck.ToolCalls) > 0 {
						for _, call := range toolCheck.ToolCalls {
							err := tools.CallTool(call.Name, call.Arguments, !unsafeMode, os.Stdout)
							if err != nil {
								fmt.Printf("Tool error [%s]: %v\n", call.Name, err)
							}
						}
						return
					}
				}
			}

			fmt.Println(response)
		},
	}

	cmd.Flags().StringVarP(&model, "model", "m", "qwen3:14b", "Model to use")
	cmd.Flags().Float64VarP(&temperature, "temperature", "t", 0.7, "Sampling temperature")
	return cmd
}
