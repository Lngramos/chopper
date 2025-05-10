package cmd

import (
	"fmt"
	"os"

	"github.com/lngramos/artoo/internal/ollama"
	"github.com/spf13/cobra"
)

var (
	chatModel       string
	chatTemperature float64
)

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a single prompt to Ollama and get a response",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		client := ollama.NewClient("http://localhost:11434")

		history := []ollama.Message{
			{Role: "user", Content: prompt},
		}

		response, err := client.Chat(chatModel, chatTemperature, history)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println(response)
	},
}

func init() {
	chatCmd.Flags().StringVarP(&chatModel, "model", "m", "mistral", "Model to use")
	chatCmd.Flags().Float64VarP(&chatTemperature, "temperature", "t", 0.7, "Sampling temperature")
	rootCmd.AddCommand(chatCmd)
}
