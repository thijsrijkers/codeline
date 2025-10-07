package llm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"bufio"
	"encoding/json"
	"net/http"
)

type OllamaClient struct {
	model string
	url   string
}

func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{
		model: model,
		url:   "http://localhost:11434/api/generate",
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

type ollamaStreamResponse struct {
	Model    string `json:"model,omitempty"`
	Response string `json:"response,omitempty"`
	Done     bool   `json:"done,omiterpty"`
}

func (c *OllamaClient) Ask(ctx context.Context, prompt string) (string, error) {
    reqBody, err := json.Marshal(ollamaRequest{
        Model:  c.model,
        Prompt: prompt,
        Stream: false,
    })
    if err != nil {
        return "", err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(reqBody))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var parsed ollamaResponse
    if err := json.Unmarshal(body, &parsed); err != nil {
        return "", fmt.Errorf("failed to parse Ollama response: %w", err)
    }

    return parsed.Response, nil
}

func (c *OllamaClient) AskStream(ctx context.Context, prompt string) (<-chan string, error) {
	reqBody, err := json.Marshal(ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var parsed ollamaStreamResponse
			if err := json.Unmarshal(line, &parsed); err != nil {
				continue 
			}

			if parsed.Response != "" {
				select {
				case <-ctx.Done():
					return
				case ch <- parsed.Response:
				}
			}

			if parsed.Done {
				return
			}
		}
	}()

	return ch, nil
}

