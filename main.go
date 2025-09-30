package main

import (
	"context"
	"fmt"
	"log"
	"codeline/llm"
)

func main() {
	ctx := context.Background()

	client, err := llm.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Ask(context.Background(), "What is the Go programming language?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

