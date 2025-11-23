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

	"github.com/bordenet/pr-faq-validator/internal/prompts"
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

	// Load prompt template from YAML
	loader := prompts.DefaultLoader
	promptTemplate, err := loader.Load("analysis/section_review.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt template: %w", err)
	}

	// Render prompts with variables
	vars := map[string]interface{}{
		"section_name": sectionName,
		"content":      content,
	}

	systemPrompt, err := promptTemplate.RenderSystemPrompt(vars)
	if err != nil {
		return nil, fmt.Errorf("failed to render system prompt: %w", err)
	}

	userPrompt, err := promptTemplate.RenderUserPrompt(vars)
	if err != nil {
		return nil, fmt.Errorf("failed to render user prompt: %w", err)
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	var resp openai.ChatCompletionResponse
	var apiErr error

	const maxAttempts = 5
	baseDelay := time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, apiErr = client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: GPT4O,
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
					{Role: openai.ChatMessageRoleUser, Content: userPrompt},
				},
			},
		)

		// success
		if apiErr == nil {
			break
		}

		// check if error is retryable
		var openaiErr *openai.APIError
		if errors.As(apiErr, &openaiErr) {
			switch openaiErr.HTTPStatusCode {
			case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				// retryable, continue
			default:
				// not retryable
				return nil, fmt.Errorf("LLM error (non-retryable): %w", apiErr)
			}
		} else {
			// unknown or non-API error
			return nil, fmt.Errorf("LLM error: %w", apiErr)
		}

		// backoff
		jitter := time.Duration(rand.Intn(300)) * time.Millisecond //nolint:gosec // weak random is fine for jitter
		delay := baseDelay * (1 << (attempt - 1))                  // exponential
		time.Sleep(delay + jitter)
	}

	// if we failed all attempts
	if apiErr != nil {
		return nil, fmt.Errorf("LLM error: exceeded retries: %w", apiErr)
	}

	text := resp.Choices[0].Message.Content

	return &Feedback{
		Section:  sectionName,
		Comments: text,
		Score:    0, // optional TODO: parse score
	}, nil
}
