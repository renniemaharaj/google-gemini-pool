package gemi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer"
)

// Session manages the generative model interactions
type Session struct {
	Model *genai.GenerativeModel
}

type Input struct {
	Current genai.Part          `json:"current"`
	History []*genai.Content    `json:"history"`
	Context []map[string]string `json:"context"`
	Errors  []error             `json:"errors"`
}

func (i *Input) SendError(err error) {
	i.Errors = append(i.Errors, err)
}

func (i *Input) String() string {
	return transformer.PartsToString([]genai.Part{i.Current})
}

// ExponentiallyValidateSend sends input to the AI model with retries, validation, and caching
func (s *Session) ExponentiallyValidateSend(ctx context.Context, input Input, validate func(resp string) error, maxTries int) (string, error) {
	startTime := time.Now()

	log.Println("Starting model-based interaction with exponential backoff and validation...")

	for i := 0; i < maxTries; i++ {
		log.Printf("Attempt %d/%d\n", i+1, maxTries)

		// Send input to AI
		resp, err := s.SendInput(ctx, input)
		if err != nil {
			input.SendError(err)

			log.Printf("API request failed: %v", err)
			time.Sleep(time.Second << i) // Exponential backoff
			continue
		}

		linted := transformer.LintCodeFences(&resp, "json")

		// Validate response
		err = validate(*linted)
		if err != nil {
			input.SendError(err)

			log.Printf("Validation failed: %v", err)
			time.Sleep(time.Second << i) // Exponential backoff
			continue
		}

		log.Printf("Response validated in %s", time.Since(startTime))
		return resp, nil
	}

	// Log final failure after max attempts
	log.Printf("Failed to validate response after %d interactions, (%s elapsed)", maxTries, time.Since(startTime))
	return "", fmt.Errorf("failed to validate response")
}

// SendInput sends a message to the AI model and returns the response
func (s *Session) SendInput(ctx context.Context, input Input) (string, error) {
	session := s.Model.StartChat()

	session.History = input.History

	structInput := struct {
		Current genai.Part          `json:"current"`
		Context []map[string]string `json:"context"`
		// Errors  []error             `json:"errors"`
	}{
		Current: input.Current,
		Context: input.Context,
		// Errors:  input.Errors,
	}
	// if structInput.Errors == nil {
	// 	structInput.Errors = []error{}
	// }

	structInputBytes, err := json.Marshal(structInput)
	if err != nil {
		return "", fmt.Errorf("error marshalling input: %v", err)
	}

	log.Printf("Sending message to model...%v\n", string(structInputBytes))
	resp, err := session.SendMessage(ctx, genai.Text(string(structInputBytes)))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	response := transformer.PartsToString(resp.Candidates[0].Content.Parts)
	return response, nil
}

// SendString sends a string message to the AI model and returns the response
func (s *Session) SendString(ctx context.Context, message string) (string, error) {
	session := s.Model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	response := transformer.PartsToString(resp.Candidates[0].Content.Parts)
	return response, nil
}
