package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	replModel       string
	replTemperature float64
	history         []Message
)

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

			history = append(history, Message{Role: "user", Content: input})

			fmt.Println("Sending request to Ollama with:")
			fmt.Printf("Model: %s, Temp: %.2f\n", replModel, replTemperature)
			fmt.Printf("History length: %d\n", len(history))

			requestBody, _ := json.Marshal(map[string]interface{}{
				"model":       replModel,
				"temperature": replTemperature,
				"messages":    history,
			})

			resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				fmt.Printf("API error: %s\n", string(bodyBytes))
				continue
			}

			// Buffer the body for potential reuse
			rawBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Failed to read response:", err)
				continue
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

			decoder := json.NewDecoder(bytes.NewReader(rawBody))
			var responseBuilder strings.Builder
			var gotResponse bool

			for {
				var chunk map[string]interface{}
				if err := decoder.Decode(&chunk); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println("Decode error:", err)
					break
				}

				// Check for response field (legacy format)
				if content, ok := chunk["response"].(string); ok {
					gotResponse = true
					responseBuilder.WriteString(content)
					fmt.Print(content)
					continue
				}

				// Check for message.role=assistant content (Ollama chat format)
				if msg, ok := chunk["message"].(map[string]interface{}); ok {
					if content, ok := msg["content"].(string); ok {
						gotResponse = true
						responseBuilder.WriteString(content)
						fmt.Print(content)
					}
				}
			}
			fmt.Println()

			if !gotResponse {
				fmt.Println("Unexpected response format:")
				fmt.Println(string(rawBody))
				continue
			}

			history = append(history, Message{Role: "assistant", Content: responseBuilder.String()})
		}
	},
}

func init() {
	replCmd.Flags().StringVarP(&replModel, "model", "m", "mistral", "Model to use")
	replCmd.Flags().Float64VarP(&replTemperature, "temperature", "t", 0.7, "Sampling temperature")
	rootCmd.AddCommand(replCmd)
}
