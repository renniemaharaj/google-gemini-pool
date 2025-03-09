package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/renniemaharaj/google-gemini-pool/internals/chat"
	"github.com/renniemaharaj/google-gemini-pool/pkg/pool"
	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer"
	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer/gemi"
)

var GEMINI_API_KEYS_POOL_SCHEMA = []transformer.API{
	{
		Key:  "API_KEY_HERE",
		Base: "gemini-20-pro-exp-0205 or other base",
	},
}

var HISTORY = []*genai.Content{}

func chatApp(p *pool.Instance) {
	log.Println("Starting chat application (type 'exit' to quit)")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		if userInput == "exit" {
			break
		}

		// Queue a session
		response, err := p.QueuedEVS(context.Background(), gemi.Input{
			Current: genai.Text(userInput),
			History: HISTORY,
			Context: []map[string]string{},
		}, func(resp string) error { return nil }, 1, 2)

		if err != nil {
			log.Printf("Failed to queue session: %v", err)
			continue
		}

		// Clear the screen before displaying chat history
		chat.ClearScreen()

		for _, part := range HISTORY {
			if part.Role == "user" {
				fmt.Print("<--/--> ➜ You: \n\n")
				fmt.Printf("-----/--/-->> ➜ %v \n\n", part.Parts[0])
			} else {
				fmt.Print("<--/--> ➜ Model: \n\n")
				fmt.Printf("-----/--/-->> ➜ %v \n\n", part.Parts[0])
			}
		}

		// Display the latest exchange with extra spacing
		fmt.Print("<--/--> ➜ User: \n\n")
		fmt.Printf("-----/--/-->> ➜ %v \n\n", userInput)

		fmt.Print("<--/--> ➜ Model: \n\n")
		fmt.Printf("-----/--/-->> ➜ %v \n\n", response)

		// Append user and model responses to history
		HISTORY = append(HISTORY, &genai.Content{Parts: []genai.Part{genai.Text(userInput)}, Role: "user"})
		HISTORY = append(HISTORY, &genai.Content{Parts: []genai.Part{genai.Text(response)}, Role: "model"})
	}
}

func main() {
	// Initialize API pool
	pool := pool.Instance{}
	pool.InitializePool()

	// Start chatting
	chatApp(&pool)
	// useQueuedEVS(&pool)
	// useChannel(&pool)
	// useQueue(&pool)

	// Wait for the example to finish
	fmt.Scanln()
}

// func useChannel(p *pool.Pool) {
// 	log.Println("Using Channel Example")
// 	p.Channel <- GEMINI_API_KEYS_POOL_SCHEMA[0]
// }

// func useQueue(p *pool.Pool) {
// 	log.Println("Using Queue Example")

// 	// Queue a session
// 	session, cleanup, err := p.Queue(context.Background())
// 	if err != nil {
// 		log.Fatalf("Failed to queue session: %v", err)
// 	}
// 	defer cleanup()

// 	// Create input
// 	input := gemi.Input{
// 		Current: genai.Text("Hello World"),
// 		History: []*genai.Content{},
// 		Context: []map[string]string{},
// 	}

// 	// Send input to AI model
// 	resp, err := session.SendInput(context.Background(), input)
// 	if err != nil {
// 		log.Fatalf("Error sending input: %v", err)
// 	}
// 	log.Println(resp)

// 	// Send a string directly
// 	resp2, err2 := session.SendString(context.Background(), "Hello World")
// 	if err2 != nil {
// 		log.Fatalf("Error sending string: %v", err2)
// 	}
// 	log.Println(resp2)
// }

// func useQueuedEVS(p *pool.Pool) {
// 	log.Println("Using QueuedEVS Example")

// 	// Create input
// 	input := gemi.Input{
// 		Current: genai.Text("Hello World"),
// 		History: []*genai.Content{},
// 		Context: []map[string]string{},
// 	}

// 	// Validation function
// 	validate := func(resp string) error {
// 		return nil // Replace with actual validation logic
// 	}

// 	// Queue a resp with retries and validation
// 	resp, err := p.QueuedEVS(context.Background(), input, validate, 3, 2)
// 	if err != nil {
// 		log.Fatalf("Error in QueuedEVS: %v", err)
// 	}

// 	log.Println(resp)
// }
