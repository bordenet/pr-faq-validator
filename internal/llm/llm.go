// Package llm provides integration with OpenAI's GPT models for qualitative feedback.
package llm

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// GPT4O is the model identifier for OpenAI's GPT-4o model.
const GPT4O = "gpt-4o"

// Feedback contains qualitative analysis feedback from the LLM.
type Feedback struct {
	Section  string
	Comments string
	Score    float64
}

// AnalyzeSection sends a section to the LLM for qualitative feedback.
func AnalyzeSection(sectionName, content string) (*Feedback, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt := fmt.Sprintf(`
You are an expert product reviewer. Review the following section of a PR-FAQ:

## Section: %s

%s

Provide specific, actionable feedback on how to improve this section. Then give it a score from 0–10 based on clarity, completeness, and effectiveness.
`, sectionName, content)

	var resp openai.ChatCompletionResponse
	var err error

	const maxAttempts = 5
	baseDelay := time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: GPT4O,
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: "You are a product manager reviewing PR-FAQs for clarity and quality."},
					{Role: openai.ChatMessageRoleUser, Content: prompt},
				},
			},
		)

		// success
		if err == nil {
			break
		}

		// check if error is retryable
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.HTTPStatusCode {
			case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				// retryable, continue
			default:
				// not retryable
				return nil, fmt.Errorf("LLM error (non-retryable): %w", err)
			}
		} else {
			// unknown or non-API error
			return nil, fmt.Errorf("LLM error: %w", err)
		}

		// backoff
		jitter := time.Duration(rand.Intn(300)) * time.Millisecond //nolint:gosec // weak random is fine for jitter
		delay := baseDelay * (1 << (attempt - 1))                  // exponential
		time.Sleep(delay + jitter)
	}

	// if we failed all attempts
	if err != nil {
		return nil, fmt.Errorf("LLM error: exceeded retries: %w", err)
	}

	text := resp.Choices[0].Message.Content

	return &Feedback{
		Section:  sectionName,
		Comments: text,
		Score:    0, // optional TODO: parse score
	}, nil
}
