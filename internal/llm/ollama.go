package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaClient struct {
	Host string
}

func NewOllamaClient(host string) *OllamaClient {
	return &OllamaClient{Host: host}
}

func (c *OllamaClient) Chat(model string, temperature float64, messages []Message) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":       model,
		"temperature": temperature,
		"messages":    messages,
	})

	resp, err := http.Post(c.Host+"/api/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error: %s", string(raw))
	}

	// Read stream fully
	var result bytes.Buffer
	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		if msg, ok := chunk["message"].(map[string]interface{}); ok {
			if content, ok := msg["content"].(string); ok {
				result.WriteString(content)
			}
		}
	}
	return result.String(), nil
}
