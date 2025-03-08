package transformer

import (
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// LintCodeFences removes ```html from the start and ``` from the end of the input string.
func LintCodeFences(input *string, language string) *string {
	codeFenceStart := fmt.Sprintf("```%v", language)
	const codeFenceEnd = "```"

	// trim the starting "```html"
	*input = strings.TrimPrefix(*input, codeFenceStart)

	// trim any leading/trailing whitespace or newlines to better detect the ending code fence
	*input = strings.TrimSpace(*input)

	// trim the ending "```"
	*input = strings.TrimSuffix(*input, codeFenceEnd)

	// trim excess whitespace again
	trimmedInput := strings.TrimSpace(*input)

	return &trimmedInput
}

func PartsToString(parts []genai.Part) string {
	var sb strings.Builder
	for _, part := range parts {
		switch v := part.(type) {
		case genai.Text:
			sb.WriteString(string(v))
		}
	}
	return sb.String()
}
