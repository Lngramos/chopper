package cmd

import (
	"fmt"
	"os"

	"github.com/lngramos/chopper/internal/ollama"
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

		fmt.Println(response)
	},
}

func init() {
	replCmd.Flags().StringVarP(&replModel, "model", "m", "qwen3:14b", "Model to use")
	replCmd.Flags().Float64VarP(&replTemperature, "temperature", "t", 0.7, "Sampling temperature")
	chatCmd.Flags().StringVarP(&chatModel, "model", "m", "qwen3:14b", "Model to use")
	chatCmd.Flags().Float64VarP(&chatTemperature, "temperature", "t", 0.7, "Sampling temperature")

	rootCmd.AddCommand(replCmd)
	rootCmd.AddCommand(chatCmd)
}
