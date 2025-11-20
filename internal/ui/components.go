// Package ui provides terminal user interface components for the PR-FAQ validator.
package ui

import (
	"fmt"
	"strings"

	"github.com/bordenet/pr-faq-validator/internal/parser"
	"github.com/charmbracelet/lipgloss"
)

// RenderHeader creates a styled header section.
func RenderHeader(title string, score int) string {
	var parts []string

	// Main title in a simple border
	titleBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 1).
		Align(lipgloss.Center).
		Render("🔍 PR-FAQ Validator")
	parts = append(parts, titleBox)

	// Document title
	if title != "" {
		parts = append(parts, SubtitleStyle.Render("📄 "+title))
	}

	// Overall score in simple format
	scoreText := GetScoreStyle(score).Render(fmt.Sprintf("%d/100", score))
	scoreDisplay := fmt.Sprintf("Overall Score: %s", scoreText)
	parts = append(parts, scoreDisplay)

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

// RenderScoreBreakdown creates a styled score breakdown table.
func RenderScoreBreakdown(breakdown parser.PRQualityBreakdown) string {
	var rows []string

	// Header
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		TableHeaderStyle.Width(25).Render("Category"),
		TableHeaderStyle.Width(12).Render("Score"),
		TableHeaderStyle.Width(12).Render("Max"),
		TableHeaderStyle.Width(30).Render("Progress"),
	)
	rows = append(rows, header)

	// Structure & Hook section (now 30 points)
	structureTotal := breakdown.HeadlineScore + breakdown.HookScore + breakdown.ReleaseDateScore
	rows = append(rows, renderScoreRow("Structure & Hook", structureTotal, 30, false))
	rows = append(rows, renderScoreRow("  Headline Quality", breakdown.HeadlineScore, 10, true))
	rows = append(rows, renderScoreRow("  Newsworthy Hook", breakdown.HookScore, 15, true))
	rows = append(rows, renderScoreRow("  Release Date", breakdown.ReleaseDateScore, 5, true))

	// Content Quality section
	contentTotal := breakdown.FiveWsScore + breakdown.CredibilityScore + breakdown.StructureScore
	rows = append(rows, renderScoreRow("Content Quality", contentTotal, 35, false))
	rows = append(rows, renderScoreRow("  5 Ws Coverage", breakdown.FiveWsScore, 15, true))
	rows = append(rows, renderScoreRow("  Credibility", breakdown.CredibilityScore, 10, true))
	rows = append(rows, renderScoreRow("  Structure", breakdown.StructureScore, 10, true))

	// Professional Quality section (now 20 points)
	professionalTotal := breakdown.ToneScore + breakdown.FluffScore
	rows = append(rows, renderScoreRow("Professional Quality", professionalTotal, 20, false))
	rows = append(rows, renderScoreRow("  Tone & Readability", breakdown.ToneScore, 10, true))
	rows = append(rows, renderScoreRow("  Fluff Avoidance", breakdown.FluffScore, 10, true))

	// Customer Evidence section
	rows = append(rows, renderScoreRow("Customer Evidence", breakdown.QuoteScore, 15, false))

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return CardStyle.Width(85).Render(content)
}

// renderScoreRow creates a single row in the score breakdown table.
func renderScoreRow(category string, score, maxScore int, isSubcategory bool) string {
	style := TableRowStyle
	if isSubcategory {
		style = TableRowAltStyle
	}

	scoreText := GetScoreStyle(score).Render(fmt.Sprintf("%d", score))
	progressBar := CreateProgressBar(score, maxScore, 20)

	return lipgloss.JoinHorizontal(lipgloss.Left,
		style.Width(25).Render(category),
		style.Width(12).Render(scoreText),
		style.Width(12).Render(fmt.Sprintf("%d", maxScore)),
		style.Width(30).Render(progressBar),
	)
}

// RenderStrengths creates a styled strengths section.
func RenderStrengths(strengths []string) string {
	if len(strengths) == 0 {
		return ""
	}

	var items []string
	items = append(items, SubtitleStyle.Render("✅ Strengths"))

	for _, strength := range strengths {
		item := SuccessListItemStyle.Render("• " + strength)
		items = append(items, item)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return SuccessCardStyle.Width(65).Render(content)
}

// RenderImprovements creates a styled improvements section.
func RenderImprovements(issues []string) string {
	if len(issues) == 0 {
		return ""
	}

	var items []string
	items = append(items, SubtitleStyle.Render("WARNING: Areas for Improvement "))

	for _, issue := range issues {
		item := WarningListItemStyle.Render("• " + issue)
		items = append(items, item)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return WarningCardStyle.Width(65).Align(lipgloss.Left).Render(content)
}

// RenderQuoteAnalysis creates a styled quote analysis section.
func RenderQuoteAnalysis(score parser.PRScore) string {
	if len(score.MetricDetails) == 0 {
		return ""
	}

	var items []string
	items = append(items, SubtitleStyle.Render(fmt.Sprintf("💬 Quote Analysis (%d quotes found)", score.TotalQuotes)))

	for i, detail := range score.MetricDetails {
		var quoteItems []string

		// Quote header with score
		header := lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Quote %d", i+1)),
			" ",
			GetScoreStyle(detail.Score).Render(fmt.Sprintf("%d/10", detail.Score)),
		)
		quoteItems = append(quoteItems, header)

		// Quote text (truncated if too long)
		quote := detail.Quote
		if len(quote) > 100 {
			quote = quote[:100] + "..."
		}
		quoteItems = append(quoteItems, lipgloss.NewStyle().Italic(true).Render("\""+quote+"\""))

		// Metrics
		if len(detail.Metrics) > 0 {
			metricsText := "Metrics: " + strings.Join(detail.Metrics, ", ")
			quoteItems = append(quoteItems, SuccessListItemStyle.Render(metricsText))

			typesText := "Types: " + strings.Join(detail.MetricTypes, ", ")
			quoteItems = append(quoteItems, ListItemStyle.Render(typesText))
		} else {
			quoteItems = append(quoteItems, WarningListItemStyle.Render("No quantitative metrics detected"))
		}

		items = append(items, lipgloss.NewStyle().Margin(1, 0).Render(
			lipgloss.JoinVertical(lipgloss.Left, quoteItems...),
		))
	}

	return CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// RenderLLMFeedback creates a styled LLM feedback section.
func RenderLLMFeedback(title, feedback string) string {
	if feedback == "" {
		return ""
	}

	var items []string
	items = append(items, SubtitleStyle.Render("🤖 AI Analysis: "+title))
	items = append(items, ListItemStyle.Render(feedback))

	return CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// RenderTabs creates a styled tab interface.
func RenderTabs(tabs []string, activeTab int) string {
	var renderedTabs []string

	for i, tab := range tabs {
		if i == activeTab {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, InactiveTabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// RenderHelp creates a styled help section.
func RenderHelp() string {
	helpText := `
Navigation:
  ←/→ or h/l    Switch tabs
  ↑/↓ or j/k    Scroll content
  q or esc      Quit
  ?             Toggle help
`
	return HelpStyle.Render(helpText)
}

// RenderStatus creates a styled status line.
func RenderStatus(message string) string {
	return StatusStyle.Render(message)
}
