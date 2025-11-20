// Package parser provides functionality for parsing and analyzing PR-FAQ documents.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// SpecSections represents the parsed sections of a PR-FAQ document.
type SpecSections struct {
	Title         string
	PressRelease  string
	FAQs          string
	Metrics       string
	OtherSections map[string]string
	PRScore       *PRScore
}

// PRScore contains the overall quality score and metrics for a press release.
type PRScore struct {
	TotalQuotes       int
	QuotesWithMetrics int
	MetricDetails     []MetricInfo
	OverallScore      int // 0-100
	QualityBreakdown  PRQualityBreakdown
}

// MetricInfo contains details about metrics found in a customer quote.
type MetricInfo struct {
	Quote       string
	Metrics     []string
	MetricTypes []string // percentage, number, ratio, etc.
	Score       int      // 0-10 for this quote
}

// PRQualityBreakdown provides detailed scoring across multiple quality dimensions.
type PRQualityBreakdown struct {
	// Structure & Hook (30 points) - increased to accommodate release date
	HeadlineScore    int // 0-10: Clear, compelling headline
	HookScore        int // 0-15: Newsworthy hook with specificity
	ReleaseDateScore int // 0-5: Release date in top lines

	// Content Quality (35 points)
	FiveWsScore      int // 0-15: Who, what, when, where, why coverage
	CredibilityScore int // 0-10: Supporting details, data, context
	StructureScore   int // 0-10: Inverted pyramid structure

	// Professional Quality (20 points) - reduced to maintain 100 total
	ToneScore  int // 0-10: Professional but readable
	FluffScore int // 0-10: Absence of marketing fluff/hype (reduced from 15)

	// Customer Evidence (15 points) - existing quote scoring
	QuoteScore int // 0-15: Quality customer quotes with metrics

	// Detailed feedback
	Issues    []string
	Strengths []string
}

// GenerateMarkdownReport creates a comprehensive markdown report with scoring table.
func GenerateMarkdownReport(sections *SpecSections, prScore *PRScore) string {
	var report strings.Builder

	// Header
	report.WriteString("# PR-FAQ Analysis Report\n\n")
	if sections.Title != "" {
		report.WriteString("**Document:** " + sections.Title + "\n")
	}
	report.WriteString("**Analysis Date:** " + time.Now().Format("January 2, 2006") + "\n")
	report.WriteString("**Overall Score:** " + fmt.Sprintf("%d/100", prScore.OverallScore) + "\n\n")

	// Executive Summary
	report.WriteString("## Executive Summary\n\n")
	if prScore.OverallScore >= 80 {
		report.WriteString("🟢 **Excellent** - This press release meets high journalistic standards and is ready for media distribution.\n\n")
	} else if prScore.OverallScore >= 60 {
		report.WriteString("🟡 **Good** - This press release has solid foundations but could benefit from targeted improvements.\n\n")
	} else if prScore.OverallScore >= 40 {
		report.WriteString("🟠 **Needs Improvement** - This press release requires significant enhancements before media distribution.\n\n")
	} else {
		report.WriteString("🔴 **Major Issues** - This press release needs substantial revision to meet professional standards.\n\n")
	}

	// Results Table
	breakdown := prScore.QualityBreakdown
	report.WriteString("## Scoring Results\n\n")
	report.WriteString("| Category | Score | Max | Status | Priority |\n")
	report.WriteString("|----------|-------|-----|--------|----------|\n")

	// Structure & Hook (now 30 points)
	structureTotal := breakdown.HeadlineScore + breakdown.HookScore + breakdown.ReleaseDateScore
	structureStatus := getScoreStatus(structureTotal, 30)
	structurePriority := getPriority(structureTotal, 30)
	report.WriteString(fmt.Sprintf("| **Structure & Hook** | %d | 30 | %s | %s |\n",
		structureTotal, structureStatus, structurePriority))
	report.WriteString(fmt.Sprintf("| ├─ Headline Quality | %d | 10 | %s | %s |\n",
		breakdown.HeadlineScore, getScoreStatus(breakdown.HeadlineScore, 10), getPriority(breakdown.HeadlineScore, 10)))
	report.WriteString(fmt.Sprintf("| ├─ Newsworthy Hook | %d | 15 | %s | %s |\n",
		breakdown.HookScore, getScoreStatus(breakdown.HookScore, 15), getPriority(breakdown.HookScore, 15)))
	report.WriteString(fmt.Sprintf("| └─ Release Date | %d | 5 | %s | %s |\n",
		breakdown.ReleaseDateScore, getScoreStatus(breakdown.ReleaseDateScore, 5), getPriority(breakdown.ReleaseDateScore, 5)))

	// Content Quality
	contentTotal := breakdown.FiveWsScore + breakdown.CredibilityScore + breakdown.StructureScore
	contentStatus := getScoreStatus(contentTotal, 35)
	contentPriority := getPriority(contentTotal, 35)
	report.WriteString(fmt.Sprintf("| **Content Quality** | %d | 35 | %s | %s |\n",
		contentTotal, contentStatus, contentPriority))
	report.WriteString(fmt.Sprintf("| ├─ 5 Ws Coverage | %d | 15 | %s | %s |\n",
		breakdown.FiveWsScore, getScoreStatus(breakdown.FiveWsScore, 15), getPriority(breakdown.FiveWsScore, 15)))
	report.WriteString(fmt.Sprintf("| ├─ Credibility | %d | 10 | %s | %s |\n",
		breakdown.CredibilityScore, getScoreStatus(breakdown.CredibilityScore, 10), getPriority(breakdown.CredibilityScore, 10)))
	report.WriteString(fmt.Sprintf("| └─ Structure | %d | 10 | %s | %s |\n",
		breakdown.StructureScore, getScoreStatus(breakdown.StructureScore, 10), getPriority(breakdown.StructureScore, 10)))

	// Professional Quality (now 20 points)
	professionalTotal := breakdown.ToneScore + breakdown.FluffScore
	professionalStatus := getScoreStatus(professionalTotal, 20)
	professionalPriority := getPriority(professionalTotal, 20)
	report.WriteString(fmt.Sprintf("| **Professional Quality** | %d | 20 | %s | %s |\n",
		professionalTotal, professionalStatus, professionalPriority))
	report.WriteString(fmt.Sprintf("| ├─ Tone & Readability | %d | 10 | %s | %s |\n",
		breakdown.ToneScore, getScoreStatus(breakdown.ToneScore, 10), getPriority(breakdown.ToneScore, 10)))
	report.WriteString(fmt.Sprintf("| └─ Fluff Avoidance | %d | 10 | %s | %s |\n",
		breakdown.FluffScore, getScoreStatus(breakdown.FluffScore, 10), getPriority(breakdown.FluffScore, 10)))

	// Customer Evidence
	report.WriteString(fmt.Sprintf("| **Customer Evidence** | %d | 15 | %s | %s |\n",
		breakdown.QuoteScore, getScoreStatus(breakdown.QuoteScore, 15), getPriority(breakdown.QuoteScore, 15)))
	report.WriteString(fmt.Sprintf("| └─ Quote Quality | %d | 15 | %s | %s |\n",
		breakdown.QuoteScore, getScoreStatus(breakdown.QuoteScore, 15), getPriority(breakdown.QuoteScore, 15)))

	// Total
	report.WriteString(fmt.Sprintf("| **TOTAL SCORE** | **%d** | **100** | %s | - |\n\n",
		prScore.OverallScore, getOverallStatus(prScore.OverallScore)))

	// Strengths
	if len(breakdown.Strengths) > 0 {
		report.WriteString("## ✅ Strengths\n\n")
		for _, strength := range breakdown.Strengths {
			report.WriteString("- " + strength + "\n")
		}
		report.WriteString("\n")
	}

	// Priority Improvements
	report.WriteString("## 🎯 Priority Improvements\n\n")
	improvements := getPriorityImprovements(breakdown)
	if len(improvements) == 0 {
		report.WriteString("No critical issues identified. Consider the suggestions below for further optimization.\n\n")
	} else {
		for i, improvement := range improvements {
			report.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, improvement.Title))
			report.WriteString("**Impact:** " + improvement.Impact + "\n\n")
			report.WriteString("**Action Steps:**\n")
			for _, step := range improvement.Steps {
				report.WriteString("- " + step + "\n")
			}
			report.WriteString("\n")
		}
	}

	// All Issues
	if len(breakdown.Issues) > 0 {
		report.WriteString("## ⚠️ Detailed Issues to Address\n\n")
		categoryIssues := categorizeIssues(breakdown.Issues)

		for category, issues := range categoryIssues {
			report.WriteString("### " + category + "\n\n")
			for _, issue := range issues {
				report.WriteString("- " + issue + "\n")
			}
			report.WriteString("\n")
		}
	}

	// Quote Analysis
	if len(prScore.MetricDetails) > 0 {
		report.WriteString("## 📊 Customer Quote Analysis\n\n")
		report.WriteString(fmt.Sprintf("**Total Quotes:** %d | **Quotes with Metrics:** %d\n\n",
			prScore.TotalQuotes, prScore.QuotesWithMetrics))

		for i, detail := range prScore.MetricDetails {
			score := detail.Score
			scoreEmoji := "🔴"
			if score >= 7 {
				scoreEmoji = "🟢"
			} else if score >= 4 {
				scoreEmoji = "🟡"
			}

			report.WriteString(fmt.Sprintf("### Quote %d %s (%d/10 points)\n\n", i+1, scoreEmoji, score))
			report.WriteString("> \"" + detail.Quote + "\"\n\n")

			if len(detail.Metrics) > 0 {
				report.WriteString("**Metrics Detected:**\n")
				for j, metric := range detail.Metrics {
					report.WriteString("- " + metric + " (" + detail.MetricTypes[j] + ")\n")
				}
			} else {
				report.WriteString("**⚠️ No quantitative metrics detected**\n\n")
				report.WriteString("**Suggestions:**\n")
				report.WriteString("- Add specific percentages (e.g., \"reduced costs by 30%\")\n")
				report.WriteString("- Include time savings (e.g., \"saves 2 hours per day\")\n")
				report.WriteString("- Mention scale improvements (e.g., \"processes 10x more data\")\n")
				report.WriteString("- Add customer count or revenue impact\n")
			}
			report.WriteString("\n")
		}
	}

	// Footer
	report.WriteString("---\n\n")
	report.WriteString("*Report generated by pr-faq-validator*\n")
	report.WriteString("*For questions about scoring methodology, see the documentation*\n")

	return report.String()
}

func getScoreStatus(score, maxScore int) string {
	percentage := float64(score) / float64(maxScore)
	if percentage >= 0.8 {
		return "🟢 Excellent"
	} else if percentage >= 0.6 {
		return "🟡 Good"
	} else if percentage >= 0.4 {
		return "🟠 Needs Work"
	} else {
		return "🔴 Critical"
	}
}

func getPriority(score, maxScore int) string {
	percentage := float64(score) / float64(maxScore)
	if percentage >= 0.8 {
		return "Low"
	} else if percentage >= 0.6 {
		return "Medium"
	} else if percentage >= 0.4 {
		return "High"
	} else {
		return "Critical"
	}
}

func getOverallStatus(score int) string {
	if score >= 80 {
		return "🟢 Ready"
	} else if score >= 60 {
		return "🟡 Good"
	} else if score >= 40 {
		return "🟠 Needs Work"
	} else {
		return "🔴 Major Issues"
	}
}

// Improvement represents a suggested improvement with actionable steps.
type Improvement struct {
	Title  string
	Impact string
	Steps  []string
}

func getPriorityImprovements(breakdown PRQualityBreakdown) []Improvement {
	var improvements []Improvement

	// Critical issues (score < 40% of max)
	if breakdown.HeadlineScore < 4 {
		improvements = append(improvements, Improvement{
			Title:  "Create Compelling Headline",
			Impact: "Headlines are the first thing journalists see. Poor headlines lead to immediate rejection.",
			Steps: []string{
				"Write 6-12 word headline with strong action verbs",
				"Include specific metrics or outcomes in the headline",
				"Avoid generic terms like 'innovative' or 'cutting-edge'",
				"Test: Can someone understand the news in 5 seconds?",
			},
		})
	}

	if breakdown.HookScore < 6 {
		improvements = append(improvements, Improvement{
			Title:  "Strengthen Opening Hook",
			Impact: "Journalists need immediate relevance. Weak hooks get press releases ignored.",
			Steps: []string{
				"Start with specific, timely announcement",
				"Include quantifiable outcomes (percentages, metrics)",
				"Clearly identify problem being solved",
				"Avoid emotional language ('excited', 'pleased')",
			},
		})
	}

	if breakdown.QuoteScore < 6 {
		improvements = append(improvements, Improvement{
			Title:  "Add Quantitative Customer Evidence",
			Impact: "Metrics in quotes provide credible proof points that journalists can use in their stories.",
			Steps: []string{
				"Replace generic enthusiasm with specific outcomes",
				"Add percentages: 'reduced processing time by 40%'",
				"Include scale metrics: 'handles 10x more transactions'",
				"Mention ROI or cost savings with numbers",
			},
		})
	}

	if breakdown.FiveWsScore < 9 {
		improvements = append(improvements, Improvement{
			Title:  "Complete the 5 Ws",
			Impact: "Missing WHO, WHAT, WHEN, WHERE, WHY makes press releases unusable for journalists.",
			Steps: []string{
				"Ensure first paragraph answers all 5 Ws",
				"Add specific date and location",
				"Clearly identify your company and what you're announcing",
				"Explain why this matters to the target audience",
			},
		})
	}

	if breakdown.FluffScore < 10 {
		improvements = append(improvements, Improvement{
			Title:  "Eliminate Marketing Fluff",
			Impact: "Hyperbolic language reduces credibility with journalists and readers.",
			Steps: []string{
				"Remove words like 'revolutionary', 'groundbreaking', 'world-class'",
				"Replace vague claims with specific proof points",
				"Back all claims with data or evidence",
				"Focus on concrete benefits rather than emotional language",
			},
		})
	}

	return improvements
}

func categorizeIssues(issues []string) map[string][]string {
	categories := make(map[string][]string)

	for _, issue := range issues {
		category := "General"
		issueLower := strings.ToLower(issue)

		if strings.Contains(issueLower, "headline") || strings.Contains(issueLower, "title") {
			category = "Headline & Title"
		} else if strings.Contains(issueLower, "hook") || strings.Contains(issueLower, "opening") || strings.Contains(issueLower, "first sentence") {
			category = "Opening Hook"
		} else if strings.Contains(issueLower, "who") || strings.Contains(issueLower, "what") || strings.Contains(issueLower, "when") || strings.Contains(issueLower, "where") || strings.Contains(issueLower, "why") {
			category = "5 Ws Coverage"
		} else if strings.Contains(issueLower, "quote") || strings.Contains(issueLower, "metric") {
			category = "Customer Evidence"
		} else if strings.Contains(issueLower, "fluff") || strings.Contains(issueLower, "marketing") || strings.Contains(issueLower, "hyperbolic") {
			category = "Professional Tone"
		} else if strings.Contains(issueLower, "structure") || strings.Contains(issueLower, "paragraph") || strings.Contains(issueLower, "transition") {
			category = "Document Structure"
		} else if strings.Contains(issueLower, "sentence") || strings.Contains(issueLower, "readability") || strings.Contains(issueLower, "passive") {
			category = "Writing Quality"
		}

		categories[category] = append(categories[category], issue)
	}

	return categories
}

// isPressReleaseContent analyzes content to determine if it looks like a press release.
func isPressReleaseContent(content string) bool {
	content = strings.ToLower(content)

	// Check for date patterns commonly found in press releases
	datePatterns := []string{
		`\b(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+\d{1,2},?\s+\d{4}`,
		`\b\d{1,2}/\d{1,2}/\d{4}`,
		`\b\d{4}-\d{1,2}-\d{1,2}`,
	}

	for _, pattern := range datePatterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			// Also check for announcement language
			announceWords := []string{"announce", "today", "excited", "pleased", "proud", "launch", "introduce", "unveil", "reveal"}
			for _, word := range announceWords {
				if strings.Contains(content, word) {
					return true
				}
			}
		}
	}

	// Check for press release structure indicators
	prIndicators := []string{
		"business wire", "pr newswire", "press release",
		"for immediate release", "contact:", "about ",
		"announces", "today announced", "is excited to announce",
		"is pleased to announce", "is proud to announce",
	}

	for _, indicator := range prIndicators {
		if strings.Contains(content, indicator) {
			return true
		}
	}

	return false
}

// isFAQSection checks if a section header indicates FAQ content.
func isFAQSection(header string) bool {
	header = strings.ToLower(strings.TrimSpace(header))

	faqPatterns := []string{
		"faq", "faqs", "frequently asked questions",
		"questions and answers", "q&a", "q & a",
		"common questions", "questions", "internal faq",
	}

	for _, pattern := range faqPatterns {
		if strings.Contains(header, pattern) {
			return true
		}
	}

	return false
}

// isNumberedFAQQuestion checks if a section header is a numbered FAQ question.
func isNumberedFAQQuestion(header string) bool {
	header = strings.TrimSpace(header)

	// Check for patterns like "1. Question here", "2. Another question", etc.
	// Also handle variations like "1) Question" or "Q1. Question"
	patterns := []string{
		`^\d+\.\s+.+`,      // "1. Question here"
		`^\d+\)\s+.+`,      // "1) Question here"
		`^Q\d+[\.\)]\s+.+`, // "Q1. Question" or "Q1) Question"
		`^Question\s+\d+`,  // "Question 1"
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, header); matched {
			return true
		}
	}

	return false
}

// extractQuotes finds customer quotes in press release content.
func extractQuotes(content string) []string {
	var quotes []string

	// Look for quoted text patterns - use Unicode code points for curly quotes
	quotePatterns := []string{
		`"(.+?)"`,           // Standard double quotes
		"\u201C(.+?)\u201D", // Curly quotes (U+201C and U+201D)
		`'(.+?)'`,           // Single quotes
		"\u2018(.+?)\u2019", // Curly single quotes (U+2018 and U+2019)
	}

	for _, pattern := range quotePatterns {
		re := regexp.MustCompile(`(?s)` + pattern) // (?s) enables multiline mode
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				quote := strings.TrimSpace(match[1])
				// Filter out very short quotes (likely not customer testimonials)
				if len(quote) > 20 {
					quotes = append(quotes, quote)
				}
			}
		}
	}
	return quotes
}

// detectMetricsInText finds quantitative metrics in text.
func detectMetricsInText(text string) ([]string, []string) {
	var metrics []string
	var metricTypes []string

	// Percentage patterns
	percentagePatterns := []string{
		`\d+(?:\.\d+)?%`,                       // 50%, 12.5%
		`\d+(?:\.\d+)?\s*percent`,              // 50 percent
		`\d+(?:\.\d+)?\s*percentage\s*points?`, // 12 percentage points
	}

	for _, pattern := range percentagePatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			metrics = append(metrics, match)
			metricTypes = append(metricTypes, "percentage")
		}
	}

	// Ratio and multiplier patterns
	ratioPatterns := []string{
		`\d+x`,                        // 2x, 10x improvement
		`\d+(?:\.\d+)?:\d+(?:\.\d+)?`, // 2:1, 3.5:1 ratios
		`\d+(?:\.\d+)?\s*times`,       // 3 times faster
	}

	for _, pattern := range ratioPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			metrics = append(metrics, match)
			metricTypes = append(metricTypes, "ratio")
		}
	}

	// Absolute number patterns with business context
	numberPatterns := []string{
		`\$\d+(?:,\d{3})*(?:\.\d+)?(?:\s*(?:million|billion|thousand|k|m|b))?`,        // $1.5M, $500K
		`\d+(?:,\d{3})*(?:\.\d+)?\s*(?:milliseconds?|seconds?|minutes?|hours?|days?)`, // 50ms, 2.5 seconds
		`\d+(?:,\d{3})*(?:\.\d+)?\s*(?:customers?|users?|transactions?)`,              // 1000 customers
		`\d+(?:,\d{3})*(?:\.\d+)?\s*(?:basis\s*points?)`,                              // 200 basis points
	}

	for _, pattern := range numberPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			metrics = append(metrics, match)
			metricTypes = append(metricTypes, "absolute")
		}
	}

	// NPS, score-based metrics
	scorePatterns := []string{
		`nps\s*(?:score\s*)?(?:by\s*)?\d+(?:\.\d+)?\s*points?`,                 // NPS by 12 points
		`\d+(?:\.\d+)?\s*(?:point|points)\s*(?:improvement|increase|decrease)`, // 12 points improvement
		`score\s*(?:of\s*)?\d+(?:\.\d+)?`,                                      // score of 9.2
	}

	for _, pattern := range scorePatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			metrics = append(metrics, match)
			metricTypes = append(metricTypes, "score")
		}
	}

	return metrics, metricTypes
}

// scoreQuote evaluates the quality of a customer quote based on metrics.
func scoreQuote(metrics []string, metricTypes []string) int {
	if len(metrics) == 0 {
		return 0 // No metrics = 0 points
	}

	score := 2 // Base score for having any metrics

	// Bonus points for different metric types
	typeBonus := make(map[string]bool)
	for _, metricType := range metricTypes {
		if !typeBonus[metricType] {
			typeBonus[metricType] = true
			switch metricType {
			case "percentage":
				score += 3 // Percentages are highly valuable
			case "ratio":
				score += 2 // Ratios show relative improvement
			case "absolute":
				score += 2 // Absolute numbers show scale
			case "score":
				score += 1 // Score improvements are good
			}
		}
	}

	// Bonus for multiple metrics in one quote
	if len(metrics) >= 2 {
		score += 2
	}
	if len(metrics) >= 3 {
		score += 1
	}

	// Cap at 10
	if score > 10 {
		score = 10
	}

	return score
}

// analyzePRQuotes evaluates customer quotes in press release content.
func analyzePRQuotes(prContent string) *PRScore {
	if prContent == "" {
		return &PRScore{OverallScore: 0}
	}

	quotes := extractQuotes(prContent)

	// Debug: print the actual content being analyzed (disabled)
	// fmt.Printf("DEBUG: Analyzing PR content with %d characters\n", len(prContent))
	// fmt.Printf("DEBUG: Found %d quotes\n", len(quotes))

	score := &PRScore{
		TotalQuotes:   len(quotes),
		MetricDetails: make([]MetricInfo, 0),
	}

	totalQuoteScore := 0
	quotesWithMetrics := 0

	for _, quote := range quotes {
		metrics, metricTypes := detectMetricsInText(quote)
		quoteScore := scoreQuote(metrics, metricTypes)

		if len(metrics) > 0 {
			quotesWithMetrics++
		}

		totalQuoteScore += quoteScore

		score.MetricDetails = append(score.MetricDetails, MetricInfo{
			Quote:       quote,
			Metrics:     metrics,
			MetricTypes: metricTypes,
			Score:       quoteScore,
		})
	}

	score.QuotesWithMetrics = quotesWithMetrics

	// Calculate overall score (0-100)
	if len(quotes) == 0 {
		score.OverallScore = 0
	} else {
		// Base score: 20 points for having quotes
		baseScore := 20

		// Metric bonus: up to 60 points based on quote quality
		metricBonus := 0
		if len(quotes) > 0 {
			avgQuoteScore := totalQuoteScore / len(quotes)
			metricBonus = (avgQuoteScore * 60) / 10 // Scale 0-10 to 0-60
		}

		// Coverage bonus: up to 20 points for having multiple quotes with metrics
		coverageBonus := 0
		if quotesWithMetrics > 0 {
			coverageBonus = 10
			if quotesWithMetrics > 1 {
				coverageBonus = 20
			}
		}

		score.OverallScore = baseScore + metricBonus + coverageBonus
		if score.OverallScore > 100 {
			score.OverallScore = 100
		}
	}

	return score
}

// analyzeHeadlineQuality evaluates headline effectiveness.
func analyzeHeadlineQuality(title string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 0

	if title == "" {
		issues = append(issues, "Missing headline/title")
		return 0, issues, strengths
	}

	// Length analysis (ideal: 6-12 words, 50-80 characters)
	words := len(strings.Fields(title))
	chars := len(title)

	if chars >= 50 && chars <= 80 && words >= 6 && words <= 12 {
		score += 3
		strengths = append(strengths, "Headline length is optimal")
	} else if chars > 100 || words > 15 {
		issues = append(issues, "Headline too long (reduces scannability)")
	} else if chars < 30 || words < 4 {
		issues = append(issues, "Headline too short (lacks specificity)")
	} else {
		score += 1
	}

	// Active voice and strong verbs
	strongVerbs := []string{"launches", "announces", "introduces", "unveils", "delivers", "creates", "develops", "achieves", "reduces", "increases", "improves", "transforms"}
	titleLower := strings.ToLower(title)

	hasStrongVerb := false
	for _, verb := range strongVerbs {
		if strings.Contains(titleLower, verb) {
			hasStrongVerb = true
			break
		}
	}

	if hasStrongVerb {
		score += 2
		strengths = append(strengths, "Uses strong action verbs")
	} else {
		issues = append(issues, "Consider using stronger action verbs")
	}

	// Specificity check (numbers, percentages, specific outcomes)
	hasSpecifics := false
	specificityPatterns := []string{`\d+%`, `\d+x`, `\d+(?:,\d{3})*`, `\$\d+`, `by \d+`, `up to \d+`}

	for _, pattern := range specificityPatterns {
		if matched, _ := regexp.MatchString(pattern, title); matched {
			hasSpecifics = true
			break
		}
	}

	if hasSpecifics {
		score += 3
		strengths = append(strengths, "Includes specific metrics or outcomes")
	} else {
		issues = append(issues, "Consider adding specific metrics to the headline")
	}

	// Avoid generic/weak language
	weakLanguage := []string{"new", "innovative", "cutting-edge", "revolutionary", "world-class", "leading", "comprehensive", "robust"}
	hasWeakLanguage := false

	for _, weak := range weakLanguage {
		if strings.Contains(titleLower, weak) {
			hasWeakLanguage = true
			break
		}
	}

	if hasWeakLanguage {
		issues = append(issues, "Avoid generic marketing language in headlines")
	} else {
		score += 2
		strengths = append(strengths, "Avoids generic marketing language")
	}

	return score, issues, strengths
}

// analyzeNewswortyHook evaluates the opening for immediate relevance and impact.
func analyzeNewswortyHook(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 0

	// Get first paragraph (hook)
	paragraphs := strings.Split(content, "\n\n")
	if len(paragraphs) == 0 {
		issues = append(issues, "No content to analyze")
		return 0, issues, strengths
	}

	hook := strings.TrimSpace(paragraphs[0])
	if hook == "" && len(paragraphs) > 1 {
		hook = strings.TrimSpace(paragraphs[1])
	}

	if hook == "" {
		issues = append(issues, "Missing opening hook")
		return 0, issues, strengths
	}

	hookLower := strings.ToLower(hook)

	// Check for timeliness indicators
	timelinessWords := []string{"today", "this week", "announces", "launched", "released", "unveiled", "now available"}
	hasTimeliness := false

	for _, word := range timelinessWords {
		if strings.Contains(hookLower, word) {
			hasTimeliness = true
			break
		}
	}

	if hasTimeliness {
		score += 3
		strengths = append(strengths, "Opens with timely announcement")
	} else {
		issues = append(issues, "Hook lacks immediate timeliness")
	}

	// Check for specificity (metrics, outcomes, concrete details)
	specificityIndicators := []string{`\d+%`, `\d+x`, `cuts .+ by`, `improves .+ by`, `reduces .+ by`, `increases .+ by`}
	hasSpecificity := false

	for _, pattern := range specificityIndicators {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, hook); matched {
			hasSpecificity = true
			break
		}
	}

	if hasSpecificity {
		score += 4
		strengths = append(strengths, "Hook includes specific, measurable outcomes")
	} else {
		issues = append(issues, "Hook lacks specific metrics or outcomes")
	}

	// Check for industry relevance/pain point addressing
	problemWords := []string{"solves", "addresses", "tackles", "eliminates", "reduces", "improves", "streamlines", "automates"}
	addressesProblem := false

	for _, word := range problemWords {
		if strings.Contains(hookLower, word) {
			addressesProblem = true
			break
		}
	}

	if addressesProblem {
		score += 3
		strengths = append(strengths, "Addresses clear problem or improvement")
	} else {
		issues = append(issues, "Hook doesn't clearly address a problem or need")
	}

	// Check for company/product clarity (who is doing what)
	sentences := strings.Split(hook, ".")
	if len(sentences) > 0 {
		firstSentence := sentences[0]
		// Should mention company and action
		if strings.Contains(firstSentence, ",") && (strings.Contains(strings.ToLower(firstSentence), "announce") || strings.Contains(strings.ToLower(firstSentence), "launch")) {
			score += 2
			strengths = append(strengths, "Clear company identification and action")
		} else {
			issues = append(issues, "First sentence should clearly identify who is doing what")
		}
	}

	// Avoid fluff language in hook
	fluffWords := []string{"excited", "pleased", "proud", "thrilled", "delighted", "revolutionary", "groundbreaking", "cutting-edge"}
	hasFluff := false

	for _, fluff := range fluffWords {
		if strings.Contains(hookLower, fluff) {
			hasFluff = true
			break
		}
	}

	if hasFluff {
		issues = append(issues, "Hook contains marketing fluff - focus on concrete value")
		score -= 1
	} else {
		score += 3
		strengths = append(strengths, "Hook avoids marketing fluff")
	}

	return score, issues, strengths
}

// analyzeFiveWs checks coverage of who, what, when, where, why.
func analyzeFiveWs(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 0

	// Get first 2-3 paragraphs for analysis
	paragraphs := strings.Split(content, "\n\n")
	leadContent := ""
	for i := 0; i < min(3, len(paragraphs)); i++ {
		leadContent += paragraphs[i] + " "
	}
	leadContentLower := strings.ToLower(leadContent)

	// WHO: Company/organization clearly identified
	companyPatterns := []string{`\b[A-Z][a-z]+\s+(?:Inc|Corp|Company|LLC|Ltd)`, `[A-Z][a-zA-Z]+\s+announced`, `[A-Z][a-zA-Z]+\s+today`}
	hasWho := false

	for _, pattern := range companyPatterns {
		if matched, _ := regexp.MatchString(pattern, leadContent); matched {
			hasWho = true
			break
		}
	}

	if hasWho {
		score += 3
		strengths = append(strengths, "Clearly identifies WHO (company/organization)")
	} else {
		issues = append(issues, "WHO: Company/organization not clearly identified in lead")
	}

	// WHAT: Product/service/action clearly described
	actionWords := []string{"announces", "launches", "introduces", "unveils", "releases", "develops", "creates"}
	hasWhat := false

	for _, action := range actionWords {
		if strings.Contains(leadContentLower, action) {
			hasWhat = true
			break
		}
	}

	if hasWhat {
		score += 3
		strengths = append(strengths, "Clearly describes WHAT (action/product/service)")
	} else {
		issues = append(issues, "WHAT: Action or offering not clearly described")
	}

	// WHEN: Timing/date mentioned
	timePatterns := []string{`\b(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+\d`, `today`, `this week`, `this month`, `\d{4}`, `yesterday`, `recently`}
	hasWhen := false

	for _, pattern := range timePatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, leadContent); matched {
			hasWhen = true
			break
		}
	}

	if hasWhen {
		score += 3
		strengths = append(strengths, "Includes WHEN (timing/date)")
	} else {
		issues = append(issues, "WHEN: Timing or date not specified")
	}

	// WHERE: Location/market mentioned
	wherePatterns := []string{`[A-Z][a-z]+,\s+[A-Z]{2}`, `[A-Z][a-z]+\s+\([A-Z][a-z]+\s+Wire\)`, `headquarters`, `market`, `globally`, `worldwide`, `nation`}
	hasWhere := false

	for _, pattern := range wherePatterns {
		if matched, _ := regexp.MatchString(pattern, leadContent); matched {
			hasWhere = true
			break
		}
	}

	if hasWhere {
		score += 2
		strengths = append(strengths, "Mentions WHERE (location/market)")
	} else {
		issues = append(issues, "WHERE: Location or market context could be clearer")
	}

	// WHY: Reason/problem/benefit explained
	whyIndicators := []string{"because", "to help", "to address", "to solve", "to improve", "to reduce", "to increase", "enables", "allows", "provides"}
	hasWhy := false

	for _, indicator := range whyIndicators {
		if strings.Contains(leadContentLower, indicator) {
			hasWhy = true
			break
		}
	}

	if hasWhy {
		score += 4
		strengths = append(strengths, "Explains WHY (reason/benefit/problem solved)")
	} else {
		issues = append(issues, "WHY: Reason or benefit not clearly explained")
	}

	return score, issues, strengths
}

// analyzeToneAndReadability evaluates professional tone and accessibility.
func analyzeToneAndReadability(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 5 // Start with neutral score

	contentLower := strings.ToLower(content)

	// Check sentence length (ideal: 15-20 words average)
	sentences := regexp.MustCompile(`[.!?]+`).Split(content, -1)
	totalWords := 0
	longSentences := 0

	for _, sentence := range sentences {
		words := len(strings.Fields(strings.TrimSpace(sentence)))
		totalWords += words
		if words > 25 {
			longSentences++
		}
	}

	if len(sentences) > 1 {
		avgWordsPerSentence := totalWords / len(sentences)
		if avgWordsPerSentence >= 15 && avgWordsPerSentence <= 20 {
			score += 2
			strengths = append(strengths, "Good sentence length for readability")
		} else if avgWordsPerSentence > 25 {
			issues = append(issues, "Sentences too long - break into shorter, clearer statements")
		}
	}

	if longSentences > len(sentences)/3 {
		issues = append(issues, "Too many overly long sentences - impacts readability")
		score -= 1
	}

	// Check for passive voice overuse
	passiveIndicators := []string{"is being", "was being", "are being", "were being", "has been", "have been", "had been", "will be"}
	passiveCount := 0

	for _, passive := range passiveIndicators {
		passiveCount += strings.Count(contentLower, passive)
	}

	if passiveCount > len(sentences)/4 {
		issues = append(issues, "Overuse of passive voice - use active voice for clarity")
		score -= 1
	} else {
		score += 1
		strengths = append(strengths, "Good use of active voice")
	}

	// Check for jargon density
	techJargon := []string{"synergies", "paradigm", "leverage", "ecosystem", "scalable", "turnkey", "best-in-class", "enterprise-grade"}
	jargonCount := 0

	for _, jargon := range techJargon {
		if strings.Contains(contentLower, jargon) {
			jargonCount++
		}
	}

	if jargonCount > 3 {
		issues = append(issues, "Too much technical jargon - write for broader audience")
		score -= 1
	} else if jargonCount == 0 {
		score += 1
		strengths = append(strengths, "Avoids unnecessary jargon")
	}

	// Check for quotation variety and quality
	quotes := extractQuotes(content)
	executiveFluff := []string{"excited", "pleased", "proud", "thrilled", "honored", "delighted"}
	fluffyQuotes := 0

	for _, quote := range quotes {
		quoteLower := strings.ToLower(quote)
		for _, fluff := range executiveFluff {
			if strings.Contains(quoteLower, fluff) {
				fluffyQuotes++
				break
			}
		}
	}

	if len(quotes) > 0 {
		if fluffyQuotes < len(quotes)/2 {
			score += 1
			strengths = append(strengths, "Quotes provide substantive insight")
		} else {
			issues = append(issues, "Too many generic 'excited' quotes - add substantive insights")
		}
	}

	return score, issues, strengths
}

// analyzeMarketingFluff detects and penalizes excessive promotional language.
func analyzeMarketingFluff(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 10 // Start with full points, deduct for fluff

	contentLower := strings.ToLower(content)

	// Hyperbolic adjectives
	hypeWords := []string{
		"revolutionary", "groundbreaking", "cutting-edge", "world-class",
		"industry-leading", "best-in-class", "state-of-the-art", "next-generation",
		"breakthrough", "game-changing", "disruptive", "unprecedented",
		"ultimate", "premier", "superior", "exceptional", "outstanding",
	}

	hypeCount := 0
	for _, hype := range hypeWords {
		if strings.Contains(contentLower, hype) {
			hypeCount++
		}
	}

	if hypeCount > 3 {
		score -= 3
		issues = append(issues, "Excessive hyperbolic language reduces credibility")
	} else if hypeCount > 1 {
		score -= 1
		issues = append(issues, "Consider reducing promotional adjectives")
	} else if hypeCount == 0 {
		strengths = append(strengths, "Avoids hyperbolic marketing language")
	}

	// Emotional fluff in quotes
	emotionalFluff := []string{"excited", "thrilled", "delighted", "pleased", "proud", "honored"}
	quotes := extractQuotes(content)
	fluffyQuotes := 0

	for _, quote := range quotes {
		quoteLower := strings.ToLower(quote)
		for _, fluff := range emotionalFluff {
			if strings.Contains(quoteLower, fluff) {
				fluffyQuotes++
				break
			}
		}
	}

	if len(quotes) > 0 {
		fluffRatio := float64(fluffyQuotes) / float64(len(quotes))
		if fluffRatio > 0.7 {
			score -= 3
			issues = append(issues, "Most quotes are generic emotional responses")
		} else if fluffRatio > 0.3 {
			score -= 1
			issues = append(issues, "Some quotes lack substantive content")
		} else {
			strengths = append(strengths, "Quotes provide meaningful insights")
		}
	}

	// Vague benefits without proof
	vagueTerms := []string{"comprehensive solution", "robust platform", "seamless integration", "enhanced productivity", "improved efficiency", "optimal performance"}
	vagueCount := 0

	for _, vague := range vagueTerms {
		if strings.Contains(contentLower, vague) {
			vagueCount++
		}
	}

	if vagueCount > 2 {
		score -= 2
		issues = append(issues, "Vague benefit claims need specific proof points")
	} else if vagueCount == 0 {
		strengths = append(strengths, "Avoids vague, unsubstantiated claims")
	}

	// Check for proof backing claims
	proofIndicators := []string{`\d+%`, `\d+x`, `study shows`, `research indicates`, `data reveals`, `according to`, `measured`, `demonstrated`}
	hasProof := false

	for _, pattern := range proofIndicators {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, content); matched {
			hasProof = true
			break
		}
	}

	if hasProof {
		strengths = append(strengths, "Backs claims with data or evidence")
	} else {
		score -= 1
		issues = append(issues, "Claims would be stronger with supporting data")
	}

	return score, issues, strengths
}

// analyzeStructure evaluates inverted pyramid and logical flow.
func analyzeStructure(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 0

	paragraphs := strings.Split(content, "\n\n")
	if len(paragraphs) < 3 {
		issues = append(issues, "Press release too short for proper structure analysis")
		return 2, issues, strengths
	}

	// First paragraph should contain key info (lead)
	firstPara := strings.TrimSpace(paragraphs[0])
	if firstPara == "" && len(paragraphs) > 1 {
		firstPara = strings.TrimSpace(paragraphs[1])
	}

	// Lead should be substantial but not too long
	leadWords := len(strings.Fields(firstPara))
	if leadWords >= 25 && leadWords <= 50 {
		score += 3
		strengths = append(strengths, "Lead paragraph has appropriate length")
	} else if leadWords > 60 {
		issues = append(issues, "Lead paragraph too long - should be concise")
	} else if leadWords < 20 {
		issues = append(issues, "Lead paragraph too brief - lacks key details")
	}

	// Check for supporting details in middle paragraphs
	middleContent := ""
	startIdx := 1
	endIdx := len(paragraphs) - 2
	if endIdx <= startIdx {
		endIdx = len(paragraphs) - 1
	}

	for i := startIdx; i < endIdx; i++ {
		middleContent += paragraphs[i] + " "
	}

	if len(middleContent) > 0 {
		// Should contain supporting details, context, or additional quotes
		supportingElements := []string{"according to", "the company", "additionally", "furthermore", "the solution", "customers"}
		hasSupporting := false

		for _, element := range supportingElements {
			if strings.Contains(strings.ToLower(middleContent), element) {
				hasSupporting = true
				break
			}
		}

		if hasSupporting {
			score += 3
			strengths = append(strengths, "Includes supporting details and context")
		} else {
			issues = append(issues, "Middle content lacks supporting details")
		}
	}

	// Last paragraph should contain boilerplate (about company)
	if len(paragraphs) >= 3 {
		lastPara := strings.ToLower(paragraphs[len(paragraphs)-1])
		boilerplateIndicators := []string{"about ", "founded", "headquartered", "company", "organization", "learn more"}
		hasBoilerplate := false

		for _, indicator := range boilerplateIndicators {
			if strings.Contains(lastPara, indicator) {
				hasBoilerplate = true
				break
			}
		}

		if hasBoilerplate {
			score += 2
			strengths = append(strengths, "Includes proper company boilerplate")
		} else {
			issues = append(issues, "Missing company boilerplate information")
		}
	}

	// Check for logical flow and transitions
	transitionWords := []string{"additionally", "furthermore", "moreover", "however", "meanwhile", "as a result"}
	hasTransitions := false

	for _, transition := range transitionWords {
		if strings.Contains(strings.ToLower(content), transition) {
			hasTransitions = true
			break
		}
	}

	if hasTransitions {
		score += 2
		strengths = append(strengths, "Uses transitions for logical flow")
	} else if len(paragraphs) > 4 {
		issues = append(issues, "Consider adding transitions between sections")
	}

	return score, issues, strengths
}

// analyzeReleaseDate checks for proper date formatting in the opening lines.
func analyzeReleaseDate(content string) (int, []string, []string) {
	var issues []string
	var strengths []string
	score := 0

	// Get the first few lines (first 200 characters) to look for release date
	firstLines := content
	if len(content) > 200 {
		firstLines = content[:200]
	}

	// Common date patterns for press releases
	datePatterns := []string{
		// Month Day, Year format: "Aug 20, 2024", "August 20, 2024", "Jan 1, 2025"
		`(?i)\b(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2},?\s+\d{4}\b`,
		// Month Day Year format: "August 20 2024"
		`(?i)\b(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2}\s+\d{4}\b`,
		// Day Month Year format: "20 August 2024"
		`(?i)\b\d{1,2}\s+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{4}\b`,
		// MM/DD/YYYY format: "08/20/2024"
		`\b\d{1,2}/\d{1,2}/\d{4}\b`,
		// MM-DD-YYYY format: "08-20-2024"
		`\b\d{1,2}-\d{1,2}-\d{4}\b`,
		// YYYY-MM-DD format: "2024-08-20"
		`\b\d{4}-\d{1,2}-\d{1,2}\b`,
		// Full date with day: "Monday, August 20, 2024"
		`(?i)\b(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday),?\s+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2},?\s+\d{4}\b`,
	}

	hasDate := false
	for _, pattern := range datePatterns {
		if matched, _ := regexp.MatchString(pattern, firstLines); matched {
			hasDate = true
			break
		}
	}

	if hasDate {
		score = 5
		strengths = append(strengths, "Includes release date in opening lines")

		// Check if it follows the standard press release format (Date. Location. Company...)
		locationPattern := `(?i)\b[A-Z][a-z]+,?\s+[A-Z]{2,}\b` // City, State/Country pattern
		if matched, _ := regexp.MatchString(locationPattern, firstLines); matched {
			strengths = append(strengths, "Follows standard press release dateline format")
		}
	} else {
		score = 0
		issues = append(issues, "Missing release date in opening lines")
		issues = append(issues, "Add date and location (e.g., 'Aug 20, 2024. Seattle, WA.')")
	}

	return score, issues, strengths
}

// comprehensivePRAnalysis combines all quality metrics.
func comprehensivePRAnalysis(prContent string, title string, quoteScore int) *PRScore {
	if prContent == "" {
		return &PRScore{OverallScore: 0}
	}

	// Analyze each component
	headlineScore, headlineIssues, headlineStrengths := analyzeHeadlineQuality(title)
	hookScore, hookIssues, hookStrengths := analyzeNewswortyHook(prContent)
	releaseDateScore, releaseDateIssues, releaseDateStrengths := analyzeReleaseDate(prContent)
	fiveWsScore, fiveWsIssues, fiveWsStrengths := analyzeFiveWs(prContent)
	structureScore, structIssues, structStrengths := analyzeStructure(prContent)
	toneScore, toneIssues, toneStrengths := analyzeToneAndReadability(prContent)
	fluffScore, fluffIssues, fluffStrengths := analyzeMarketingFluff(prContent)

	// Combine all issues and strengths
	allIssues := append(headlineIssues, hookIssues...)
	allIssues = append(allIssues, releaseDateIssues...)
	allIssues = append(allIssues, fiveWsIssues...)
	allIssues = append(allIssues, structIssues...)
	allIssues = append(allIssues, toneIssues...)
	allIssues = append(allIssues, fluffIssues...)

	allStrengths := append(headlineStrengths, hookStrengths...)
	allStrengths = append(allStrengths, releaseDateStrengths...)
	allStrengths = append(allStrengths, fiveWsStrengths...)
	allStrengths = append(allStrengths, structStrengths...)
	allStrengths = append(allStrengths, toneStrengths...)
	allStrengths = append(allStrengths, fluffStrengths...)

	// Calculate overall score (100 points total)
	// New scoring: Structure & Hook (30), Content Quality (35), Professional Quality (20), Customer Evidence (15)
	totalScore := headlineScore + hookScore + releaseDateScore + fiveWsScore + structureScore + toneScore + fluffScore + quoteScore

	breakdown := PRQualityBreakdown{
		HeadlineScore:    headlineScore,
		HookScore:        hookScore,
		ReleaseDateScore: releaseDateScore,
		FiveWsScore:      fiveWsScore,
		CredibilityScore: toneScore, // Use tone score for credibility
		StructureScore:   structureScore,
		ToneScore:        toneScore,
		FluffScore:       fluffScore,
		QuoteScore:       quoteScore,
		Issues:           allIssues,
		Strengths:        allStrengths,
	}

	// Get quote analysis from existing function
	quoteAnalysis := analyzePRQuotes(prContent)

	// Add quote count feedback
	var quoteCountIssues []string
	if quoteAnalysis.TotalQuotes > 4 {
		quoteCountIssues = append(quoteCountIssues, "Consider reducing quotes - press releases work best with 3-4 focused customer testimonials")
	}

	// Combine quote count feedback with other issues
	allIssues = append(allIssues, quoteCountIssues...)

	// Update the breakdown with the complete issue list
	breakdown.Issues = allIssues

	return &PRScore{
		TotalQuotes:       quoteAnalysis.TotalQuotes,
		QuotesWithMetrics: quoteAnalysis.QuotesWithMetrics,
		MetricDetails:     quoteAnalysis.MetricDetails,
		OverallScore:      totalScore,
		QualityBreakdown:  breakdown,
	}
}

// ParsePRFAQ reads a markdown file and extracts key sections.
func ParsePRFAQ(path string) (*SpecSections, error) {
	file, err := os.Open(path) //nolint:gosec // path is user-provided CLI argument
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't override return error
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	sections := &SpecSections{
		OtherSections: make(map[string]string),
	}

	type sectionInfo struct {
		name    string
		content string
	}

	var currentSection string
	var sectionBuffer strings.Builder
	var titleSet bool
	var allSections []sectionInfo

	// Define common section headers once
	commonHeaders := []string{
		"Press Release", "FAQ", "FAQs", "Frequently Asked Questions",
		"Q&A", "Questions and Answers", "Success Metrics", "Key Metrics",
		"Metrics", "Internal FAQ", "Questions", "Answers",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Extract the title (first H1)
		if !titleSet && strings.HasPrefix(line, "# ") {
			titleText := strings.TrimPrefix(line, "# ")
			sections.Title = titleText
			titleSet = true

			// Check if this H1 is also a section header (like "# Press Release")
			for _, header := range commonHeaders {
				if strings.EqualFold(titleText, header) {
					// This H1 is both a title and a section header
					currentSection = titleText
					continue
				}
			}
			continue
		}

		// Detect section heading (both markdown ## and plain text)
		isMarkdownHeader := strings.HasPrefix(line, "## ")
		isPlainTextHeader := false

		// Check for plain text headers (common section names that stand alone)
		trimmedLine := strings.TrimSpace(line)

		for _, header := range commonHeaders {
			if strings.EqualFold(trimmedLine, header) {
				isPlainTextHeader = true
				break
			}
		}

		if isMarkdownHeader || isPlainTextHeader {
			// Save the previous section's content
			if currentSection != "" {
				content := strings.TrimSpace(sectionBuffer.String())
				allSections = append(allSections, sectionInfo{
					name:    currentSection,
					content: content,
				})
				sectionBuffer.Reset()
			}

			if isMarkdownHeader {
				currentSection = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			} else {
				currentSection = trimmedLine
			}
			continue
		}

		// Accumulate section content
		if currentSection != "" {
			sectionBuffer.WriteString(line + "\n")
		}
	}

	// Capture last section
	if currentSection != "" {
		content := strings.TrimSpace(sectionBuffer.String())
		allSections = append(allSections, sectionInfo{
			name:    currentSection,
			content: content,
		})
	}

	// Process sections with fuzzy logic and handle FAQ numbering
	var faqContent strings.Builder
	var inFAQSection bool

	for _, section := range allSections {
		// Check for FAQ sections first (more specific)
		if isFAQSection(section.name) {
			sections.FAQs = section.content
			faqContent.WriteString(section.content + "\n\n")
			inFAQSection = true
			continue
		}

		// Check if this is a numbered FAQ question (part of FAQ section)
		if inFAQSection && isNumberedFAQQuestion(section.name) {
			faqContent.WriteString("## " + section.name + "\n\n")
			faqContent.WriteString(section.content + "\n\n")
			continue
		} else if inFAQSection {
			// We've left the FAQ section, finalize it
			sections.FAQs = strings.TrimSpace(faqContent.String())
			inFAQSection = false
		}

		// Check for explicit press release header
		if strings.ToLower(section.name) == "press release" {
			sections.PressRelease = section.content
			continue
		}

		// Check for metrics sections
		lowerName := strings.ToLower(section.name)
		if strings.Contains(lowerName, "success metrics") || strings.Contains(lowerName, "key metrics") {
			sections.Metrics = section.content
			continue
		}

		// Use fuzzy logic to detect press release content
		if sections.PressRelease == "" && isPressReleaseContent(section.content) {
			sections.PressRelease = section.content
			continue
		}

		// Default to other sections
		sections.OtherSections[section.name] = section.content
	}

	// Handle case where FAQ section continues to end of document
	if inFAQSection && faqContent.Len() > 0 {
		sections.FAQs = strings.TrimSpace(faqContent.String())
	}

	// Analyze PR with comprehensive quality metrics
	if sections.PressRelease != "" {
		quoteAnalysis := analyzePRQuotes(sections.PressRelease)
		quoteScore := (quoteAnalysis.OverallScore * 15) / 100 // Scale to 15 points max
		sections.PRScore = comprehensivePRAnalysis(sections.PressRelease, sections.Title, quoteScore)
	} else {
		sections.PRScore = &PRScore{OverallScore: 0}
	}

	return sections, nil
}
