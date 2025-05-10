package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Chat(model string, temperature float64, messages []Message) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       model,
		"temperature": temperature,
		"messages":    messages,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(bodyBytes))
	}

	decoder := json.NewDecoder(resp.Body)
	var responseBuilder strings.Builder
	var gotResponse bool

	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if content, ok := chunk["response"].(string); ok {
			gotResponse = true
			responseBuilder.WriteString(content)
			continue
		}

		if msg, ok := chunk["message"].(map[string]interface{}); ok {
			if content, ok := msg["content"].(string); ok {
				gotResponse = true
				responseBuilder.WriteString(content)
			}
		}
	}

	if !gotResponse {
		return "", errors.New("unexpected response format")
	}

	return responseBuilder.String(), nil
}
