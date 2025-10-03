package main

import (
	"context"
	"log"
	"codeline/llm"
	"codeline/tui"
)

func main() {
	ctx := context.Background()

	client, err := llm.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	tui.StartChat(ctx, client)
}
