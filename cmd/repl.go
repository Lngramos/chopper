package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lngramos/chopper/internal/llm"
	"github.com/lngramos/chopper/internal/tools"
	"github.com/spf13/cobra"
)

func NewReplCommand(client llm.Client) *cobra.Command {
	var model string
	var temperature float64
	var history []llm.Message

	cmd := &cobra.Command{
		Use:   "repl",
		Short: "Start an interactive chat session with LLM",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

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

			fmt.Println("Chopper REPL - Type 'exit' to quit")
			for {
				fmt.Print(">> ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if input == "" {
					continue
				} else if input == "exit" {
					fmt.Println("Goodbye.")
					break
				}

				history = append(history, llm.Message{Role: "user", Content: input})

				fmt.Println("Sending request to LLM with:")
				fmt.Printf("Model: %s, Temp: %.2f\n", model, temperature)
				fmt.Printf("History length: %d\n", len(history))

				messages := append([]llm.Message{systemMessage}, history...)
				response, err := client.Chat(model, temperature, messages)
				if err != nil {
					fmt.Println("Error:", err)
					continue
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
							}
							history = append(history, llm.Message{Role: "assistant", Content: ""})
							continue
						} else if len(toolCheck.ToolCalls) > 0 {
							for _, call := range toolCheck.ToolCalls {
								err := tools.CallTool(call.Name, call.Arguments, !unsafeMode, os.Stdout)
								if err != nil {
									fmt.Printf("Tool error [%s]: %v\n", call.Name, err)
								}
							}
							history = append(history, llm.Message{Role: "assistant", Content: ""})
							continue
						}
					}
				}

				fmt.Println(response)
				history = append(history, llm.Message{Role: "assistant", Content: response})
			}
		},
	}

	cmd.Flags().StringVarP(&model, "model", "m", "qwen3:14b", "Model to use")
	cmd.Flags().Float64VarP(&temperature, "temperature", "t", 0.7, "Sampling temperature")
	return cmd
}
