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

// Parameters returns the parameters for a api key by matching api.Base or defaulting to the default parameters
func (api *API) Parameters() Parameters {
	switch api.Base {
	case "gemini-20-pro-exp-0205":
		return Parameters{
			Temperature:      1,
			TopK:             64,
			TopP:             0.95,
			MaxOutputTokens:  8192,
			ResponseMIMEType: "text/plain",
			SystemInstruction: &genai.Content{
				Parts: []genai.Part{
					genai.Text(" "),
				},
			},
		}
	default:
		return Parameters{
			Temperature:      1,
			TopK:             64,
			TopP:             0.95,
			MaxOutputTokens:  8192,
			ResponseMIMEType: "text/plain",
			SystemInstruction: &genai.Content{
				Parts: []genai.Part{
					genai.Text(" "),
				},
			},
		}
	}

}
