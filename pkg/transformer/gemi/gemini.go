package gemi

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/renniemaharaj/google-gemini-pool/pkg/transformer"
)

// Model creates a new generative model from a configuration
func Model(ctx context.Context, cfx transformer.Configuration) (*genai.GenerativeModel, func(), error) {
	log.Println("Creating google gemini client...")

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfx.Key.Key))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating client: %v", err)
	}

	model := client.GenerativeModel(cfx.Key.Base)
	model.SetTemperature(cfx.Parameters.Temperature)
	model.SetTopK(cfx.Parameters.TopK)
	model.SetTopP(cfx.Parameters.TopP)
	model.SetMaxOutputTokens(cfx.Parameters.MaxOutputTokens)
	model.ResponseMIMEType = cfx.Parameters.ResponseMIMEType
	model.SystemInstruction = cfx.Parameters.SystemInstruction

	return model, func() { client.Close() }, nil
}
