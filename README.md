# Google Gemini Pool Manager

A Go-based pool manager for handling multiple Google Gemini API keys efficiently. This package provides thread-safe API key management and session handling for Google's Gemini AI model.

## Features

- Thread-safe API key pool management
- Session queueing and cleanup
- Chat history management
- Multiple usage patterns (Channel, Queue, QueuedEVS)
- Retry mechanism with validation

## Installation

```bash
go get github.com/renniemaharaj/google-gemini-pool
```

## Configuration

Set up your API keys in your application:

```go
// Store json object, [] of transformer.API: key and base for pool
var GEMINI_API_KEYS_POOL = []transformer.API{
    {
        Key:  "YOUR_API_KEY",
        Base: "gemini-20-pro-exp-0205", // or your preferred model
    },
}

// Loads from GEMINI_API_KEYS_POOL environment variable & pushes all to pool
pool.InitializePool() 

// Push one (1) sigle transformer.API key to channel
myGeminiKey := tranformer.API{
    {
        Key:  "YOUR_API_KEY",
        Base: "gemini-20-pro-exp-0205", // or your preferred model
    },
}

// Push directly to exposed channel
pool.Channel <- myGeminiKey
```

## Usage Examples

### 1. Basic Queue Usage

```go
session, cleanup, err := pool.Queue(context.Background())
if err != nil {
    log.Fatalf("Failed to queue session: %v", err)
}
defer cleanup()

// Send a simple string
response, err := session.SendString(context.Background(), "Hello World")
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Println(response)
```

### 2. Chat Application

```go
var HISTORY = []*genai.Content{}

func chatApp() {
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
		session, cleanup, err := pool.Queue(context.Background())
		if err != nil {
			log.Printf("Failed to queue session: %v", err)
			continue
		}
		defer cleanup()

		// Send message and get response
		response, err := session.SendInput(context.Background(), gemi.Input{
			Current: genai.Text(userInput),
			History: HISTORY,
			Context: []map[string]string{},
		})

		if err != nil {
			log.Printf("Error getting response: %v", err)
			continue
		}

		fmt.Printf("AI: %s\n", response)

		HISTORY = append(HISTORY, &genai.Content{Parts: []genai.Part{genai.Text(userInput)}, Role: "user"})
		HISTORY = append(HISTORY, &genai.Content{Parts: []genai.Part{genai.Text(response)}, Role: "model"})
	}
}

func main() {
	// Initialize API pool
	pool.InitializePool()

    	// Start chatting
	chatApp()

	// Wait for the example to finish
	fmt.Scanln()
}
```

### 3. QueuedEVS (Queue with Validation and Retries)

```go
validate := func(resp string) error {
    // Add your validation logic
    return nil
}

response, err := pool.QueuedEVS(
    context.Background(),
    input,
    validate,
    3, // max retries
    2, // concurrent requests
)
```

## API Reference

### Types

- `transformer.API`: Configuration for API keys
- `gemi.Input`: Input structure for sending messages
- `pool.Session`: Manages individual API sessions

### Main Functions

- `pool.InitializePool()`: Initialize the API key pool
- `pool.Queue(ctx)`: Get a session from the pool
- `pool.QueuedEVS(ctx, input, validate, maxRetries, concurrent)`: Queue with validation
- `session.SendString(ctx, text)`: Send a simple string message
- `session.SendInput(ctx, input)`: Send a structured input with history

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
