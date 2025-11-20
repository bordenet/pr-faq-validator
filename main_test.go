package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
