package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lngramos/chopper/internal/ollama"
	"github.com/lngramos/chopper/internal/tools"
	"github.com/spf13/cobra"
)

var (
	replModel       string
	replTemperature float64
	history         []ollama.Message
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start an interactive chat session with Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
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

			history = append(history, ollama.Message{Role: "user", Content: input})

			fmt.Println("Sending request to Ollama with:")
			fmt.Printf("Model: %s, Temp: %.2f\n", replModel, replTemperature)
			fmt.Printf("History length: %d\n", len(history))

			messages := append([]ollama.Message{systemMessage}, history...)
			response, err := client.Chat(replModel, replTemperature, messages)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// try to parse tool_call
			var toolCheck struct {
				ToolCall *struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments"`
				} `json:"tool_call"`
			}
			if err := json.Unmarshal([]byte(response), &toolCheck); err == nil && toolCheck.ToolCall != nil {
				result, err := tools.CallTool(toolCheck.ToolCall.Name, toolCheck.ToolCall.Arguments)
				if err != nil {
					fmt.Println("Tool error:", err)
				} else {
					fmt.Println(result)
					history = append(history, ollama.Message{Role: "assistant", Content: result})
				}
				continue
			}

			fmt.Println(response)
			history = append(history, ollama.Message{Role: "assistant", Content: response})
		}
	},
}
