package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type LLM interface {
	Ask(ctx context.Context, prompt string) (string, error)
}

func NewFromEnv() (LLM, error) {
	provider := strings.ToLower(os.Getenv("LLM_PROVIDER"))
	model := os.Getenv("LLM_MODEL")
	apiKey := os.Getenv("LLM_API_KEY")

	switch provider {
	case "ollama":
		if model == "" {
			model = "llama2"
		}
		return NewOllamaClient(model), nil

	default:
		return nil, fmt.Errorf("unknown LLM_PROVIDER: %s", provider)
	}
}

