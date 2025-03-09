package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
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
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		if userInput == "exit" {
			break
		}

		// Queue a session
		// session, cleanup, err := p.Queue(context.Background())
		response, err := p.QueuedEVS(context.Background(), gemi.Input{
			Current: genai.Text(userInput),
			History: HISTORY,
			Context: []map[string]string{},
		}, func(resp string) error { return nil }, 1, 2)

		if err != nil {
			log.Printf("Failed to queue session: %v", err)
			continue
		}
		// defer cleanup()

		// Send message and get response
		// response, err := session.SendInput(context.Background(), gemi.Input{
		// 	Current: genai.Text(userInput),
		// 	History: HISTORY,
		// 	Context: []map[string]string{},
		// })

		// if err != nil {
		// 	log.Printf("Error getting response: %v", err)
		// 	continue
		// }

		fmt.Printf("AI: %s\n", response)

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
