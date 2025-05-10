package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with a local LLM via Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		prompt := "Say hello like R2D2."
		if len(args) > 0 {
			prompt = args[0]
		}

		requestBody, _ := json.Marshal(map[string]string{
			"model":  "qwen3:14b",
			"prompt": prompt,
		})

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk map[string]interface{}
			if err := decoder.Decode(&chunk); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Decode error:", err)
				break
			}

			// Print only the "response" field if it exists
			if content, ok := chunk["response"].(string); ok {
				fmt.Print(content) // no newline; model usually ends with \n
			}
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
