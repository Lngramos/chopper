package llm

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client interface {
	Chat(model string, temperature float64, messages []Message) (string, error)
}
