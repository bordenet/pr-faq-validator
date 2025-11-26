package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bordenet/pr-faq-validator/internal/parser"
)

func TestMain_NoArgs(t *testing.T) {
	if os.Getenv("TEST_MAIN_NO_ARGS") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_NoArgs") //nolint:gosec // test code
	cmd.Env = append(os.Environ(), "TEST_MAIN_NO_ARGS=1")
	err := cmd.Run()

	// Should exit with error when no file specified
	if err == nil {
		t.Error("Expected error when no file specified, got nil")
	}
}

func TestMain_InvalidFile(t *testing.T) {
	if os.Getenv("TEST_MAIN_INVALID_FILE") == "1" {
		os.Args = []string{"cmd", "-file", "/nonexistent/file.md", "-no-tui"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_InvalidFile") //nolint:gosec // test code
	cmd.Env = append(os.Environ(), "TEST_MAIN_INVALID_FILE=1")
	err := cmd.Run()

	// Should exit with error for invalid file
	if err == nil {
		t.Error("Expected error for invalid file, got nil")
	}
}

func TestMain_ValidFile_NoTUI(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")

	content := `# Test PR-FAQ

## Press Release

**SEATTLE, WA - November 20, 2025** - Company announces new product.

This is a test press release with proper structure.

"We improved performance by 50%," said the CEO.

## FAQ

**Q: What is this?**
A: A test document.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Build the binary
	binPath := filepath.Join(tmpDir, "pr-faq-validator")
	buildCmd := exec.Command("go", "build", "-o", binPath) //nolint:gosec // test code
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Run with -no-tui flag
	cmd := exec.Command(binPath, "-file", tmpFile, "-no-tui") //nolint:gosec // test code
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Verify output contains expected sections
	expectedSections := []string{
		"PR-FAQ Analysis",
		"Overall Score",
		"Test PR-FAQ",
	}

	for _, section := range expectedSections {
		if !strings.Contains(outputStr, section) {
			t.Errorf("Output missing expected section: %q\nOutput: %s", section, outputStr)
		}
	}
}

func TestMain_ValidFile_Report(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")

	content := `# Test PR-FAQ

## Press Release

Test content.

## FAQ

Q: Test?
A: Yes.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Build the binary
	binPath := filepath.Join(tmpDir, "pr-faq-validator")
	buildCmd := exec.Command("go", "build", "-o", binPath) //nolint:gosec // test code
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Run with -report flag
	reportPath := filepath.Join(tmpDir, "report.md")
	cmd := exec.Command(binPath, "-file", tmpFile, "-report", reportPath) //nolint:gosec // test code
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Read the generated report file
	reportContent, err := os.ReadFile(reportPath) //nolint:gosec // test code with controlled paths
	if err != nil {
		t.Fatalf("Failed to read report file: %v", err)
	}

	reportStr := string(reportContent)

	// Verify markdown report format
	expectedSections := []string{
		"# PR-FAQ Analysis Report",
		"## Executive Summary",
		"## Scoring Results",
	}

	for _, section := range expectedSections {
		if !strings.Contains(reportStr, section) {
			t.Errorf("Report missing expected section: %q", section)
		}
	}
}

func TestWriteReportToFile(t *testing.T) {
	t.Run("writes content to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filename := filepath.Join(tmpDir, "test_report.md")
		content := "# Test Report\n\nThis is a test."

		err := writeReportToFile(filename, content)
		if err != nil {
			t.Fatalf("writeReportToFile failed: %v", err)
		}

		// Verify file was created and contains content
		data, err := os.ReadFile(filename) //nolint:gosec // test code with controlled paths
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if string(data) != content {
			t.Errorf("file content = %q, want %q", string(data), content)
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		err := writeReportToFile("/nonexistent/path/report.md", "content")
		if err == nil {
			t.Error("expected error for invalid path")
		}
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filename := filepath.Join(tmpDir, "existing.md")

		// Create initial file
		err := os.WriteFile(filename, []byte("old content"), 0600)
		if err != nil {
			t.Fatalf("failed to create initial file: %v", err)
		}

		// Overwrite with new content
		newContent := "new content"
		err = writeReportToFile(filename, newContent)
		if err != nil {
			t.Fatalf("writeReportToFile failed: %v", err)
		}

		// Verify new content
		data, err := os.ReadFile(filename) //nolint:gosec // test code with controlled paths
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if string(data) != newContent {
			t.Errorf("file content = %q, want %q", string(data), newContent)
		}
	})
}

func TestLoggerInitialization(t *testing.T) {
	// Verify logger is initialized
	if logger == nil {
		t.Error("logger should be initialized in init()")
	}
}

func TestRunLegacyOutput(t *testing.T) {
	// Create test sections
	sections := parser.SpecSections{
		Title:        "Test PR-FAQ",
		PressRelease: "Test press release content with a quote: \"We saved 50%,\" said the CEO.",
		FAQs:         "Q: Test?\nA: Yes.",
		PRScore: &parser.PRScore{
			OverallScore:      75,
			TotalQuotes:       1,
			QuotesWithMetrics: 1,
			MetricDetails: []parser.MetricInfo{
				{Quote: "We saved 50%", Metrics: []string{"50%"}, MetricTypes: []string{"percentage"}, Score: 8},
			},
			QualityBreakdown: parser.PRQualityBreakdown{
				HeadlineScore:    8,
				HookScore:        12,
				ReleaseDateScore: 5,
				FiveWsScore:      12,
				CredibilityScore: 7,
				StructureScore:   7,
				ToneScore:        8,
				FluffScore:       8,
				QuoteScore:       10,
				Strengths:        []string{"Good headline"},
				Issues:           []string{"Add more metrics"},
			},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function (this will also try to call LLM which will fail without API key)
	runLegacyOutput(sections)

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read captured output
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Verify output contains expected sections
	if !strings.Contains(output, "Test PR-FAQ") {
		t.Error("Output missing title")
	}
	if !strings.Contains(output, "Quality Breakdown") {
		t.Error("Output missing quality breakdown")
	}
}

func TestRunLegacyOutputEmptySections(t *testing.T) {
	// Create minimal sections
	sections := parser.SpecSections{
		Title: "Empty Test",
		PRScore: &parser.PRScore{
			OverallScore: 0,
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	runLegacyOutput(sections)

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read captured output
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Should still produce output
	if output == "" {
		t.Error("Expected non-empty output")
	}
}
