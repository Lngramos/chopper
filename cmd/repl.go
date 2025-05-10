package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var replModel string
var replTemperature float64

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start an interactive chat session with Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Artoo REPL - Type 'exit' to quit")
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

			requestBody, _ := json.Marshal(map[string]interface{}{
				"model":       replModel,
				"prompt":      input,
				"temperature": replTemperature,
			})

			resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Println("Error:", err)
				continue
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

				if content, ok := chunk["response"].(string); ok {
					fmt.Print(content)
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	replCmd.Flags().StringVarP(&replModel, "model", "m", "qwen3:14b", "Model to use")
	replCmd.Flags().Float64VarP(&replTemperature, "temperature", "t", 0.7, "Sampling temperature")

	rootCmd.AddCommand(replCmd)
}
