package llm

import (
	"context"
	"fmt"
)

type OllamaClient struct {
	model string
}

func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{model: model}
}

func (c *OllamaClient) Ask(ctx context.Context, prompt string) (string, error) {
	// TODO: Replace with actual Ollama HTTP API call
	return fmt.Sprintf("[Ollama (%s) Response to: %s]", c.model, prompt), nil
}

