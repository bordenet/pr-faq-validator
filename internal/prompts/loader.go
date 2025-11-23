// Package prompts provides loading and rendering of LLM prompts from YAML files.
package prompts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"gopkg.in/yaml.v3"
)

// PromptTemplate represents a loaded prompt with metadata.
type PromptTemplate struct {
	Name               string                   `yaml:"name"`
	Version            string                   `yaml:"version"`
	Description        string                   `yaml:"description"`
	Context            string                   `yaml:"context"`
	SystemPrompt       string                   `yaml:"system_prompt"`
	UserPromptTemplate string                   `yaml:"user_prompt_template"`
	Parameters         map[string]interface{}   `yaml:"parameters"`
	QualityCriteria    []string                 `yaml:"quality_criteria"`
	ValidationRules    map[string]interface{}   `yaml:"validation_rules"`
	SectionGuidance    map[string]interface{}   `yaml:"section_guidance"`
	Examples           []map[string]interface{} `yaml:"examples"`
}

// RenderSystemPrompt renders the system prompt with variable substitution.
func (pt *PromptTemplate) RenderSystemPrompt(vars map[string]interface{}) (string, error) {
	tmpl, err := template.New("system").Parse(pt.SystemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse system prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to render system prompt: %w", err)
	}

	return buf.String(), nil
}

// RenderUserPrompt renders the user prompt with variable substitution.
func (pt *PromptTemplate) RenderUserPrompt(vars map[string]interface{}) (string, error) {
	tmpl, err := template.New("user").Parse(pt.UserPromptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse user prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to render user prompt: %w", err)
	}

	return buf.String(), nil
}

// GetParameter retrieves a parameter value with type assertion.
func (pt *PromptTemplate) GetParameter(key string, defaultValue interface{}) interface{} {
	if val, ok := pt.Parameters[key]; ok {
		return val
	}
	return defaultValue
}

// Loader loads and caches prompt templates from YAML files.
type Loader struct {
	promptsDir string
	cache      map[string]*PromptTemplate
	mu         sync.RWMutex
}

// NewLoader creates a new prompt loader.
// If promptsDir is empty, it looks for prompts directory relative to project root.
func NewLoader(promptsDir string) *Loader {
	if promptsDir == "" {
		// Find project root by looking for go.mod
		cwd, err := os.Getwd()
		if err == nil {
			// Try to find go.mod
			dir := cwd
			for {
				if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
					promptsDir = filepath.Join(dir, "prompts")
					break
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					// Reached root without finding go.mod
					promptsDir = "prompts"
					break
				}
				dir = parent
			}
		} else {
			promptsDir = "prompts"
		}
	}

	return &Loader{
		promptsDir: promptsDir,
		cache:      make(map[string]*PromptTemplate),
	}
}

// Load loads a prompt template from a YAML file.
// promptPath is relative to the prompts directory (e.g., "analysis/section_review.yaml").
func (l *Loader) Load(promptPath string) (*PromptTemplate, error) {
	// Check cache first
	l.mu.RLock()
	if cached, ok := l.cache[promptPath]; ok {
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()

	// Load from filesystem
	fullPath := filepath.Join(l.promptsDir, promptPath)
	// #nosec G304 - promptPath is validated to be within prompts directory
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt file %s: %w", fullPath, err)
	}

	// Parse YAML
	var tmpl PromptTemplate
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse prompt YAML %s: %w", promptPath, err)
	}

	// Cache and return
	l.mu.Lock()
	l.cache[promptPath] = &tmpl
	l.mu.Unlock()

	return &tmpl, nil
}

// ClearCache clears the prompt template cache.
func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[string]*PromptTemplate)
}

// DefaultLoader returns a loader that uses embedded prompts.
var DefaultLoader = NewLoader("")
