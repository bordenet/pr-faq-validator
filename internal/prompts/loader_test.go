package prompts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader(t *testing.T) {
	t.Run("creates loader with default prompts directory", func(t *testing.T) {
		loader := NewLoader("")
		if loader == nil {
			t.Fatal("expected non-nil loader")
		}
		if loader.promptsDir == "" {
			t.Error("expected non-empty promptsDir")
		}
	})

	t.Run("creates loader with custom directory", func(t *testing.T) {
		loader := NewLoader("/custom/path")
		if loader == nil {
			t.Fatal("expected non-nil loader")
		}
		if loader.promptsDir != "/custom/path" {
			t.Errorf("expected promptsDir to be /custom/path, got %s", loader.promptsDir)
		}
	})
}

func TestLoadPrompt(t *testing.T) {
	// Use actual prompts directory
	loader := NewLoader("../../prompts")

	t.Run("loads section review prompt", func(t *testing.T) {
		tmpl, err := loader.Load("analysis/section_review.yaml")
		if err != nil {
			t.Fatalf("failed to load prompt: %v", err)
		}

		if tmpl.Name != "section-review" {
			t.Errorf("expected name 'section-review', got '%s'", tmpl.Name)
		}

		if tmpl.Version != "1.0.0" {
			t.Errorf("expected version '1.0.0', got '%s'", tmpl.Version)
		}

		if tmpl.SystemPrompt == "" {
			t.Error("expected non-empty system prompt")
		}

		if tmpl.UserPromptTemplate == "" {
			t.Error("expected non-empty user prompt template")
		}
	})

	t.Run("loads pr-faq generation prompt", func(t *testing.T) {
		tmpl, err := loader.Load("generation/pr_faq_generation.yaml")
		if err != nil {
			t.Fatalf("failed to load prompt: %v", err)
		}

		if tmpl.Name != "pr-faq-generation" {
			t.Errorf("expected name 'pr-faq-generation', got '%s'", tmpl.Name)
		}

		if len(tmpl.QualityCriteria) == 0 {
			t.Error("expected quality criteria to be present")
		}
	})

	t.Run("caches loaded prompts", func(t *testing.T) {
		tmpl1, err := loader.Load("analysis/section_review.yaml")
		if err != nil {
			t.Fatalf("failed to load prompt: %v", err)
		}

		tmpl2, err := loader.Load("analysis/section_review.yaml")
		if err != nil {
			t.Fatalf("failed to load prompt second time: %v", err)
		}

		// Should be same instance from cache
		if tmpl1 != tmpl2 {
			t.Error("expected cached prompt to be same instance")
		}
	})

	t.Run("returns error for non-existent prompt", func(t *testing.T) {
		_, err := loader.Load("nonexistent/prompt.yaml")
		if err == nil {
			t.Error("expected error for non-existent prompt")
		}
	})
}

func TestRenderSystemPrompt(t *testing.T) {
	loader := NewLoader("../../prompts")
	tmpl, err := loader.Load("analysis/section_review.yaml")
	if err != nil {
		t.Fatalf("failed to load prompt: %v", err)
	}

	t.Run("renders with variables", func(t *testing.T) {
		vars := map[string]interface{}{
			"section_name": "Press Release",
			"content":      "Test content",
		}

		rendered, err := tmpl.RenderSystemPrompt(vars)
		if err != nil {
			t.Fatalf("failed to render system prompt: %v", err)
		}

		if rendered == "" {
			t.Error("expected non-empty rendered prompt")
		}
	})
}

func TestRenderUserPrompt(t *testing.T) {
	loader := NewLoader("../../prompts")
	tmpl, err := loader.Load("analysis/section_review.yaml")
	if err != nil {
		t.Fatalf("failed to load prompt: %v", err)
	}

	t.Run("renders with variables", func(t *testing.T) {
		vars := map[string]interface{}{
			"section_name": "Press Release",
			"content":      "Test content for press release",
		}

		rendered, err := tmpl.RenderUserPrompt(vars)
		if err != nil {
			t.Fatalf("failed to render user prompt: %v", err)
		}

		if rendered == "" {
			t.Error("expected non-empty rendered prompt")
		}

		// Check that variables were substituted
		if !contains(rendered, "Press Release") {
			t.Error("expected rendered prompt to contain section name")
		}

		if !contains(rendered, "Test content for press release") {
			t.Error("expected rendered prompt to contain content")
		}
	})
}

func TestGetParameter(t *testing.T) {
	tmpl := &PromptTemplate{
		Parameters: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  2000,
		},
	}

	t.Run("returns existing parameter", func(t *testing.T) {
		val := tmpl.GetParameter("temperature", 1.0)
		if val != 0.7 {
			t.Errorf("expected 0.7, got %v", val)
		}
	})

	t.Run("returns default for missing parameter", func(t *testing.T) {
		val := tmpl.GetParameter("missing", "default")
		if val != "default" {
			t.Errorf("expected 'default', got %v", val)
		}
	})
}

func TestClearCache(t *testing.T) {
	loader := NewLoader("../../prompts")

	// Load a prompt to populate cache
	_, err := loader.Load("analysis/section_review.yaml")
	if err != nil {
		t.Fatalf("failed to load prompt: %v", err)
	}

	if len(loader.cache) == 0 {
		t.Error("expected cache to be populated")
	}

	loader.ClearCache()

	if len(loader.cache) != 0 {
		t.Error("expected cache to be empty after clear")
	}
}

func TestLoadFromFilesystem(t *testing.T) {
	// Create temporary directory with test prompt
	tmpDir := t.TempDir()
	promptDir := filepath.Join(tmpDir, "test_prompts")
	// #nosec G301 - test directory permissions are acceptable for tests
	if err := os.MkdirAll(promptDir, 0750); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Write test prompt
	testPrompt := `name: "test-prompt"
version: "1.0.0"
description: "Test prompt"
context: "Test context"
system_prompt: "Test system prompt"
user_prompt_template: "Test user prompt"
parameters:
  temperature: 0.5
`
	promptFile := filepath.Join(promptDir, "test.yaml")
	// #nosec G306 - test file permissions are acceptable for tests
	if err := os.WriteFile(promptFile, []byte(testPrompt), 0600); err != nil {
		t.Fatalf("failed to write test prompt: %v", err)
	}

	loader := NewLoader(promptDir)

	tmpl, err := loader.Load("test.yaml")
	if err != nil {
		t.Fatalf("failed to load prompt from filesystem: %v", err)
	}

	if tmpl.Name != "test-prompt" {
		t.Errorf("expected name 'test-prompt', got '%s'", tmpl.Name)
	}
}

// Test RenderSystemPrompt with invalid template
func TestRenderSystemPromptError(t *testing.T) {
	tmpl := &PromptTemplate{
		SystemPrompt: "Hello {{.invalid_syntax",
	}

	_, err := tmpl.RenderSystemPrompt(map[string]interface{}{})
	if err == nil {
		t.Error("expected error for invalid template syntax")
	}
}

// Test RenderUserPrompt with invalid template
func TestRenderUserPromptError(t *testing.T) {
	tmpl := &PromptTemplate{
		UserPromptTemplate: "Hello {{.invalid_syntax",
	}

	_, err := tmpl.RenderUserPrompt(map[string]interface{}{})
	if err == nil {
		t.Error("expected error for invalid template syntax")
	}
}

// Test NewLoader with non-existent directory
func TestNewLoaderNonExistentDir(t *testing.T) {
	loader := NewLoader("/nonexistent/path/that/does/not/exist")
	if loader == nil {
		t.Fatal("expected non-nil loader even with non-existent path")
	}
	// The loader should still be created, but Load will fail
	_, err := loader.Load("test.yaml")
	if err == nil {
		t.Error("expected error when loading from non-existent directory")
	}
}

// Test Load with invalid YAML
func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Write invalid YAML
	invalidYAML := `name: "test
	invalid yaml content
	missing closing quote`
	promptFile := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(promptFile, []byte(invalidYAML), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	loader := NewLoader(tmpDir)
	_, err := loader.Load("invalid.yaml")
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
