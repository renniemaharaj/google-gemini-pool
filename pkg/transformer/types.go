package transformer

import (
	"github.com/google/generative-ai-go/genai"
)

// API is a struct that holds the API key and base for the model
type API struct {
	Key  string `json:"key"`
	Base string `json:"base"`
}

// Parameters is a struct that holds the parameters for the model
type Parameters struct {
	Temperature       float32
	TopK              int32
	TopP              float32
	MaxOutputTokens   int32
	ResponseMIMEType  string
	SystemInstruction *genai.Content
}

// SetSystemInstructions sets the system instructions for the model
func (p *Parameters) SetSystemInstructions(i **genai.Content) {
	p.SystemInstruction = *i
}

// Configuration is a struct that holds the API key and parameters for the model
type Configuration struct {
	Key        API
	Parameters Parameters
}

// SetParameters sets the parameters for the model
func (c *Configuration) SetParameters(p *Parameters) {
	c.Parameters = *p
}

// SetKey sets the API key for the model
func (c *Configuration) SetKey(k *API) {
	c.Key = *k
}

var noContent = genai.Content{Parts: []genai.Part{genai.Text(" ")}}

// Define common parameter sets
var defaultTopK40 = Parameters{
	Temperature:       1,
	TopK:              40,
	TopP:              0.95,
	MaxOutputTokens:   8192,
	ResponseMIMEType:  "text/plain",
	SystemInstruction: &noContent,
}

var defaultTopK64 = Parameters{
	Temperature:       1,
	TopK:              64,
	TopP:              0.95,
	MaxOutputTokens:   8192,
	ResponseMIMEType:  "text/plain",
	SystemInstruction: &noContent,
}

// Map of models to their corresponding parameter sets
var paramMap = map[string]Parameters{
	"gemini-2.0-flash":                    defaultTopK40,
	"gemini-2.0-flash-exp":                defaultTopK40,
	"gemini-2.0-flash-lite":               defaultTopK40,
	"gemini-2.0-pro-exp-02-05":            defaultTopK64,
	"gemini-2.0-flash-thinking-exp-01-21": {Temperature: 0.7, TopK: 64, TopP: 0.95, MaxOutputTokens: 8192, ResponseMIMEType: "text/plain", SystemInstruction: &noContent},
	"learnlm-1.5-pro-experimental":        defaultTopK64,
	"gemini-1.5-pro":                      defaultTopK40,
	"gemini-1.5-flash":                    defaultTopK40,
	"gemini-1.5-flash-8b":                 defaultTopK40,
}

// Parameters returns the parameters for an API key by matching api.Base or defaulting to defaultTopK40
func (api *API) Parameters() Parameters {
	if params, exists := paramMap[api.Base]; exists {
		return params
	}
	return defaultTopK40 // Fallback default
}
