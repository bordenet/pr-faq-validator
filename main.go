// Package main provides the CLI entry point for the PR-FAQ validator tool.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bordenet/pr-faq-validator/internal/llm"
	"github.com/bordenet/pr-faq-validator/internal/parser"
	"github.com/bordenet/pr-faq-validator/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	inputFile := flag.String("file", "", "Path to the PR-FAQ markdown file")
	reportFile := flag.String("report", "", "Optional: Output markdown report file (default: interactive TUI)")
	noTUI := flag.Bool("no-tui", false, "Disable interactive TUI and output to stdout")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Please provide a markdown file with -file")
	}

	sections, err := parser.ParsePRFAQ(*inputFile)
	if err != nil {
		log.Fatalf("Failed to parse PR-FAQ: %v", err)
	}

	// If markdown report is requested, generate and save it
	if *reportFile != "" {
		report := parser.GenerateMarkdownReport(sections, sections.PRScore)
		err := writeReportToFile(*reportFile, report)
		if err != nil {
			log.Fatalf("Failed to write report: %v", err)
		}
		fmt.Printf("Report generated: %s\n", *reportFile)
		fmt.Printf("Overall Score: %d/100\n", sections.PRScore.OverallScore)
		return
	}

	// If TUI is disabled, output to stdout (legacy mode)
	if *noTUI {
		runLegacyOutput(*sections)
		return
	}

	// Run interactive TUI
	runInteractiveTUI(*sections)
}

// runInteractiveTUI starts the interactive TUI interface.
func runInteractiveTUI(sections parser.SpecSections) {
	// Initialize TUI model
	model := ui.NewModel(sections)

	// Create Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the TUI
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// runLegacyOutput provides the original stdout-based output.
func runLegacyOutput(sections parser.SpecSections) {
	// Generate comprehensive markdown report
	report := parser.GenerateMarkdownReport(&sections, sections.PRScore)
	fmt.Print(report)

	// Original detailed analysis follows for reference
	fmt.Printf("\n---\n\n== Detailed Analysis ==\n\n")
	fmt.Printf("== PR-FAQ Title ==\n%s\n\n", sections.Title)

	// Display comprehensive PR scoring results
	if sections.PressRelease != "" {
		fmt.Printf("== Press Release Quality Score: %d/100 ==\n\n", sections.PRScore.OverallScore)

		// Quality breakdown
		breakdown := sections.PRScore.QualityBreakdown
		fmt.Println("== Quality Breakdown ==")
		fmt.Printf("Structure & Hook:      %d/30 points\n", breakdown.HeadlineScore+breakdown.HookScore+breakdown.ReleaseDateScore)
		fmt.Printf("  - Headline Quality:   %d/10\n", breakdown.HeadlineScore)
		fmt.Printf("  - Newsworthy Hook:    %d/15\n", breakdown.HookScore)
		fmt.Printf("  - Release Date:       %d/5\n", breakdown.ReleaseDateScore)
		fmt.Printf("Content Quality:       %d/35 points\n", breakdown.FiveWsScore+breakdown.CredibilityScore+breakdown.StructureScore)
		fmt.Printf("  - 5 Ws Coverage:      %d/15\n", breakdown.FiveWsScore)
		fmt.Printf("  - Credibility:        %d/10\n", breakdown.CredibilityScore)
		fmt.Printf("  - Structure:          %d/10\n", breakdown.StructureScore)
		fmt.Printf("Professional Quality:  %d/20 points\n", breakdown.ToneScore+breakdown.FluffScore)
		fmt.Printf("  - Tone & Readability: %d/10\n", breakdown.ToneScore)
		fmt.Printf("  - Fluff Avoidance:    %d/10\n", breakdown.FluffScore)
		fmt.Printf("Customer Evidence:     %d/15 points\n", breakdown.QuoteScore)
		fmt.Printf("  - Quote Quality:      %d/15\n\n", breakdown.QuoteScore)

		// Strengths
		if len(breakdown.Strengths) > 0 {
			fmt.Println("== Strengths ==")
			for _, strength := range breakdown.Strengths {
				fmt.Printf("✓ %s\n", strength)
			}
			fmt.Println()
		}

		// Issues to address
		if len(breakdown.Issues) > 0 {
			fmt.Println("== Areas for Improvement ==")
			for _, issue := range breakdown.Issues {
				fmt.Printf("⚠ %s\n", issue)
			}
			fmt.Println()
		}

		// Detailed quote analysis if present
		if len(sections.PRScore.MetricDetails) > 0 {
			fmt.Printf("== Quote Analysis (%d quotes found) ==\n", sections.PRScore.TotalQuotes)
			for i, detail := range sections.PRScore.MetricDetails {
				fmt.Printf("\nQuote %d (Score: %d/10):\n", i+1, detail.Score)
				fmt.Printf("\"%s\"\n", detail.Quote)
				if len(detail.Metrics) > 0 {
					fmt.Printf("Metrics detected: %v\n", detail.Metrics)
					fmt.Printf("Metric types: %v\n", detail.MetricTypes)
				} else {
					fmt.Println("No quantitative metrics detected")
				}
			}
			fmt.Println()
		}

		fmt.Println("Analyzing Press Release...")
		feedback, err := llm.AnalyzeSection("Press Release", sections.PressRelease)
		if err != nil {
			log.Fatalf("LLM error: %v", err)
		}
		fmt.Printf("== Feedback for Press Release ==\n%s\n\n", feedback.Comments)
	}

	if sections.FAQs != "" {
		fmt.Println("Analyzing FAQs...")
		feedback, err := llm.AnalyzeSection("FAQs", sections.FAQs)
		if err != nil {
			log.Fatalf("LLM error: %v", err)
		}
		fmt.Printf("== Feedback for FAQs ==\n%s\n\n", feedback.Comments)
	}
}

func writeReportToFile(filename, content string) error {
	file, err := os.Create(filename) //nolint:gosec // filename is user-provided CLI argument
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	_, err = file.WriteString(content)
	return err
}
