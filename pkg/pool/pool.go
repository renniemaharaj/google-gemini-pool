package pool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer"
	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer/gemi"
)

type Instance struct {
	Channel chan *transformer.API
	once    sync.Once // ensures initialization happens only once

}

// Initialize API key pool
func (p *Instance) HydrateChannels(keys []transformer.API) {
	p.once.Do(func() {
		p.Channel = make(chan *transformer.API, len(keys))
		for _, key := range keys {
			p.Channel <- &key
		}
		log.Printf("API Key Pool Initialized with %d keys", len(keys))
	})
}

// LoadGeminiAPIPool loads API keys from an environment variable.
func (p *Instance) LoadEnv_GEMINI_API_KEYS_POOL(envVar string) ([]transformer.API, error) {
	jsonStr := os.Getenv(envVar)
	if jsonStr == "" {
		return nil, fmt.Errorf("environment variable %s is empty", envVar)
	}

	var keys []transformer.API
	if err := json.Unmarshal([]byte(jsonStr), &keys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no API keys found in environment variable %s", envVar)
	}

	log.Printf("Loaded %d API keys from environment variable %s", len(keys), envVar)

	// Initialize the API key pool
	p.HydrateChannels(keys)

	return keys, nil
}

// InitializePool initializes the API key pool.
func (p *Instance) InitializePool() {
	keys, err := p.LoadEnv_GEMINI_API_KEYS_POOL("GEMINI_API_KEYS_POOL")
	if err != nil {
		log.Println(err)
		return
	}
	p.HydrateChannels(keys)
}

// QueuedEVS queues, exponential backoff, validating model responses.
func (p *Instance) QueuedEVS(ctx context.Context, input gemi.Input, validate func(resp string) error, queueTries int, backoff int) (string, error) {
	timeStart := time.Now()
	log.Println("Starting queued-based EVS with exponential backoff and validation...")
	for i := 0; i < queueTries; i++ {
		log.Printf("Attempt %d of %d", i+1, queueTries)
		session, cleanup, err := p.Queue(ctx)

		if err != nil {
			log.Printf("Failed to get session: %v", err)
			continue
		}

		resp, err := session.ExponentiallyValidateSend(ctx, &input, validate, backoff)
		cleanup()

		if err != nil {
			log.Printf("Failed to validate response: %v", err)
			continue
		}

		log.Printf("Success after %d attempts, took %v", i+1, time.Since(timeStart))
		return resp, nil
	}

	log.Printf("failed to validate pool response over %dqueues * %v = %v (since)", queueTries, backoff, time.Since(timeStart))
	return "", fmt.Errorf("failed to validate pool response over %dqueues * %v = %v (since)", queueTries, backoff, time.Since(timeStart))
}

// Queue returns a session from the pool of available API keys.
func (p *Instance) Queue(ctx context.Context) (*gemi.Session, func(), error) {
	log.Println("Waiting for available key...")

	// Non-blocking key retrieval
	select {
	case api := <-p.Channel:
		log.Printf("Using key: %s", api.Key)

		// New configuration
		cfx := transformer.Configuration{
			Key:        *api,
			Parameters: api.Parameters(),
		}

		// Set system instruction
		cfx.Parameters.SystemInstruction.Parts = []genai.Part{
			genai.Text(genai.Text(transformer.GetProgramming())),
		}

		log.Println("Creating model...")
		model, cleanup, err := gemi.Model(ctx, cfx)
		if err != nil {
			p.Channel <- api                              // return key if model creation fails
			log.Printf("Freeing key (Fail): %s", api.Key) // Log key release
			return nil, nil, fmt.Errorf("error creating model: %w", err)
		}

		// Custom cleanup to return the key back when done
		cleanupFunc := func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in cleanup: %v", r)
				}
				log.Printf("Freeing key (Finish): %s", api.Key) // Log key release
				p.Channel <- api                                // Always return key
			}()

			cleanup() // Call original cleanup
		}

		session := gemi.Session{Model: model}

		return &session, cleanupFunc, nil

	default:
		return nil, nil, fmt.Errorf("no API keys available")
	}
}
