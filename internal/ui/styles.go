package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	primaryColor = lipgloss.Color("#7C3AED") // Purple
	successColor = lipgloss.Color("#10B981") // Green
	warningColor = lipgloss.Color("#F59E0B") // Orange
	errorColor   = lipgloss.Color("#EF4444") // Red
	mutedColor   = lipgloss.Color("#6B7280") // Gray
	textColor    = lipgloss.Color("#F9FAFB") // Light gray

	// TitleStyle is the style for the main title.
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center).
			Width(25)

	// SubtitleStyle is the style for subtitles.
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			MarginBottom(1)

	// ScoreStyle is the style for high scores.
	ScoreStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	// ScoreLowStyle is the style for low scores.
	ScoreLowStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(0, 1)

	// ScoreMediumStyle is the style for medium scores.
	ScoreMediumStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true).
				Padding(0, 1)

	// TableHeaderStyle is the style for table headers.
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(mutedColor).
				Padding(0, 1)

	// TableRowStyle is the style for table rows.
	TableRowStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	// TableRowAltStyle is the style for alternate table rows.
	TableRowAltStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(lipgloss.Color("#374151")).
				Padding(0, 1)

	// ProgressBarStyle is the style for progress bars.
	ProgressBarStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Padding(0, 1)

	ProgressFillStyle = lipgloss.NewStyle().
				Background(successColor).
				Foreground(lipgloss.Color("#000000"))

	// ProgressEmptyStyle is the style for empty progress bar sections.
	ProgressEmptyStyle = lipgloss.NewStyle().
				Background(mutedColor)

	// CardStyle is the style for card containers.
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2).
			MarginBottom(1)

	SuccessCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(successColor).
				Padding(1, 2).
				MarginBottom(1)

	// WarningCardStyle is the style for warning cards.
	WarningCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(warningColor).
				Padding(1, 2).
				MarginBottom(1)

	// ListItemStyle is the style for list items.
	ListItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	SuccessListItemStyle = lipgloss.NewStyle().
				Foreground(successColor).
				PaddingLeft(2)

	// WarningListItemStyle is the style for warning list items.
	WarningListItemStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				PaddingLeft(2)

	// StatusStyle is the style for status messages.
	StatusStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// HelpStyle is the style for help text.
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			MarginTop(1)

	// ActiveTabStyle is the style for the active tab.
	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2).
			Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "│",
			BottomRight: "│",
		}).
		BorderForeground(primaryColor)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 2).
				Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "│",
			BottomRight: "│",
		}).
		BorderForeground(mutedColor)
)

// GetScoreStyle returns the appropriate style based on score
func GetScoreStyle(score int) lipgloss.Style {
	if score >= 70 {
		return ScoreStyle
	} else if score >= 40 {
		return ScoreMediumStyle
	}
	return ScoreLowStyle
}

// CreateProgressBar creates a styled progress bar.
func CreateProgressBar(current, maxScore int, width int) string {
	if maxScore == 0 {
		return ""
	}

	percentage := float64(current) / float64(maxScore)
	fillWidth := int(float64(width) * percentage)
	emptyWidth := width - fillWidth

	fill := ProgressFillStyle.Width(fillWidth).Render("")
	empty := ProgressEmptyStyle.Width(emptyWidth).Render("")

	return ProgressBarStyle.Render(fill + empty)
}

// FormatScore formats a score with appropriate styling.
func FormatScore(score, maxScore int) string {
	style := GetScoreStyle(score)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d", score)),
		"/",
		fmt.Sprintf("%d", maxScore),
	))
}
