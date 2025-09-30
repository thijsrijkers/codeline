package main

import (
	"context"
	"fmt"
	"log"

	"llm-example/llm"
)

func main() {
	ctx := context.Background()

	client, err := llm.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Ask(ctx, "What is the capital of France?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

