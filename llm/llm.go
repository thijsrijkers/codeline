package llm

import (
	"context"
	"os"
)

type LLM interface {
	Ask(ctx context.Context, prompt string) (string, error)
}

func NewFromEnv() (LLM, error) {
	provider := os.Getenv("LLM_PROVIDER")
	model := os.Getenv("LLM_MODEL")

	switch provider {
	case "ollama":
		if model == "" {
			model = "llama2"
		}
		return NewOllamaClient(model), nil

	default:
		return NewOllamaClient("llama2"), nil
	}
}
