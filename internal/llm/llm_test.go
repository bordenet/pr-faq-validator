package llm

import (
	"os"
	"testing"
)

func TestAnalyzeSection_NoAPIKey(t *testing.T) {
	// Save original API key
	originalKey := os.Getenv("OPENAI_API_KEY")
	defer func() {
		if originalKey != "" {
			_ = os.Setenv("OPENAI_API_KEY", originalKey)
		}
	}()

	// Unset API key
	_ = os.Unsetenv("OPENAI_API_KEY")

	_, err := AnalyzeSection("Test Section", "Test content")
	if err == nil {
		t.Error("Expected error when OPENAI_API_KEY is not set, got nil")
	}

	expectedMsg := "OPENAI_API_KEY not set"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestAnalyzeSection_EmptyContent(t *testing.T) {
	// Skip if no API key (this would require actual API call)
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	// This test would make an actual API call, so we skip it in normal test runs
	t.Skip("Skipping integration test that requires API call")
}

func TestFeedbackStruct(t *testing.T) {
	feedback := &Feedback{
		Section:  "Test Section",
		Comments: "Test comments",
		Score:    8.5,
	}

	if feedback.Section != "Test Section" {
		t.Errorf("Section = %q, want %q", feedback.Section, "Test Section")
	}

	if feedback.Comments != "Test comments" {
		t.Errorf("Comments = %q, want %q", feedback.Comments, "Test comments")
	}

	if feedback.Score != 8.5 {
		t.Errorf("Score = %f, want %f", feedback.Score, 8.5)
	}
}

func TestGPT4OConstant(t *testing.T) {
	expected := "gpt-4o"
	if GPT4O != expected {
		t.Errorf("GPT4O = %q, want %q", GPT4O, expected)
	}
}
