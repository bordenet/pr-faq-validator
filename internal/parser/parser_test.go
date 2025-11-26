package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePRFAQ(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantTitle string
		wantPR    bool
		wantFAQ   bool
		wantErr   bool
	}{
		{
			name: "valid pr-faq with standard headers",
			content: `# Test PR-FAQ

## Press Release
This is a press release.

## FAQ
Q: Question?
A: Answer.
`,
			wantTitle: "Test PR-FAQ",
			wantPR:    true,
			wantFAQ:   true,
			wantErr:   false,
		},
		{
			name: "pr-faq with alternate headers",
			content: `# Product Launch

## Announcement
Product announcement here.

## Q&A
Q: What is this?
A: A product.
`,
			wantTitle: "Product Launch",
			wantPR:    true,
			wantFAQ:   true,
			wantErr:   false,
		},
		{
			name: "pr-only document",
			content: `# Press Release Only

## Press Release
Just a press release.
`,
			wantTitle: "Press Release Only",
			wantPR:    true,
			wantFAQ:   false,
			wantErr:   false,
		},
		{
			name:      "empty document",
			content:   "",
			wantTitle: "",
			wantPR:    false,
			wantFAQ:   false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0600); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			// Parse
			sections, err := ParsePRFAQ(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePRFAQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify results
			if sections.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", sections.Title, tt.wantTitle)
			}

			hasPR := sections.PressRelease != ""
			if hasPR != tt.wantPR {
				t.Errorf("Has press release = %v, want %v", hasPR, tt.wantPR)
			}

			hasFAQ := sections.FAQs != ""
			if hasFAQ != tt.wantFAQ {
				t.Errorf("Has FAQ = %v, want %v", hasFAQ, tt.wantFAQ)
			}

			// Verify score is calculated
			if sections.PRScore == nil {
				t.Error("PRScore is nil")
			}
		})
	}
}

func TestParsePRFAQ_InvalidPath(t *testing.T) {
	_, err := ParsePRFAQ("/nonexistent/file.md")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestParsePRFAQ_EmptyPath(t *testing.T) {
	_, err := ParsePRFAQ("")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestDetectMetricsInText(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		wantMetricCount int
		wantTypes       []string
	}{
		{
			name:            "percentage metric",
			text:            "We improved performance by 75%",
			wantMetricCount: 1,
			wantTypes:       []string{"percentage"},
		},
		{
			name:            "multiple percentages",
			text:            "Increased by 50% and reduced costs by 30%",
			wantMetricCount: 2,
			wantTypes:       []string{"percentage", "percentage"},
		},
		{
			name:            "ratio metric",
			text:            "Performance improved 10x",
			wantMetricCount: 1,
			wantTypes:       []string{"ratio"},
		},
		{
			name:            "absolute number",
			text:            "Saved $1.5M in costs",
			wantMetricCount: 1,
			wantTypes:       []string{"absolute"},
		},
		{
			name:            "mixed metrics",
			text:            "Reduced latency by 50% and improved throughput 3x, saving $500K",
			wantMetricCount: 3,
			wantTypes:       []string{"percentage", "ratio", "absolute"},
		},
		{
			name:            "no metrics",
			text:            "This is a great product that customers love",
			wantMetricCount: 0,
			wantTypes:       []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics, types := detectMetricsInText(tt.text)

			if len(metrics) != tt.wantMetricCount {
				t.Errorf("detectMetricsInText() got %d metrics, want %d", len(metrics), tt.wantMetricCount)
			}

			if len(types) != len(tt.wantTypes) {
				t.Errorf("detectMetricsInText() got %d types, want %d", len(types), len(tt.wantTypes))
			}

			// Verify types match (order-independent for mixed metrics)
			typeMap := make(map[string]int)
			for _, typ := range types {
				typeMap[typ]++
			}
			wantTypeMap := make(map[string]int)
			for _, typ := range tt.wantTypes {
				wantTypeMap[typ]++
			}

			for typ, count := range wantTypeMap {
				if typeMap[typ] != count {
					t.Errorf("detectMetricsInText() type %q count = %d, want %d", typ, typeMap[typ], count)
				}
			}
		})
	}
}

func TestScoreQuote(t *testing.T) {
	tests := []struct {
		name        string
		metrics     []string
		metricTypes []string
		wantMin     int
		wantMax     int
	}{
		{
			name:        "no metrics",
			metrics:     []string{},
			metricTypes: []string{},
			wantMin:     0,
			wantMax:     0,
		},
		{
			name:        "single percentage",
			metrics:     []string{"75%"},
			metricTypes: []string{"percentage"},
			wantMin:     5,
			wantMax:     5,
		},
		{
			name:        "multiple different types",
			metrics:     []string{"50%", "10x", "$1M"},
			metricTypes: []string{"percentage", "ratio", "absolute"},
			wantMin:     9,
			wantMax:     10,
		},
		{
			name:        "many metrics capped at 10",
			metrics:     []string{"50%", "75%", "10x", "5x", "$1M", "$2M"},
			metricTypes: []string{"percentage", "percentage", "ratio", "ratio", "absolute", "absolute"},
			wantMin:     10,
			wantMax:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scoreQuote(tt.metrics, tt.metricTypes)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("scoreQuote() = %d, want between %d and %d", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestExtractQuotes(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantCount  int
		wantQuotes []string
	}{
		{
			name:       "single quote",
			text:       `The CEO said, "This is an amazing product launch announcement."`,
			wantCount:  1,
			wantQuotes: []string{"This is an amazing product launch announcement."},
		},
		{
			name:      "multiple quotes",
			text:      `"This first quote is long enough to be extracted," said John. "This second quote is also long enough," said Jane.`,
			wantCount: 2,
		},
		{
			name:      "no quotes",
			text:      "This text has no quotes at all.",
			wantCount: 0,
		},
		{
			name:       "quote with metrics",
			text:       `"We improved performance by 75%," said the CTO.`,
			wantCount:  1,
			wantQuotes: []string{"We improved performance by 75%,"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes := extractQuotes(tt.text)

			if len(quotes) != tt.wantCount {
				t.Errorf("extractQuotes() got %d quotes, want %d", len(quotes), tt.wantCount)
			}

			for i, want := range tt.wantQuotes {
				if i >= len(quotes) {
					break
				}
				if quotes[i] != want {
					t.Errorf("extractQuotes() quote[%d] = %q, want %q", i, quotes[i], want)
				}
			}
		})
	}
}

func TestAnalyzeHeadlineQuality(t *testing.T) {
	tests := []struct {
		name     string
		headline string
		wantMin  int
		wantMax  int
	}{
		{
			name:     "optimal length with action verb",
			headline: "Company Launches Product That Reduces Costs by 50% for Enterprise Customers",
			wantMin:  7,
			wantMax:  10,
		},
		{
			name:     "too short",
			headline: "New Product",
			wantMin:  0,
			wantMax:  5,
		},
		{
			name:     "too long",
			headline: "Company Announces the Launch of a Revolutionary New Product That Will Transform the Industry and Change Everything Forever",
			wantMin:  0,
			wantMax:  5,
		},
		{
			name:     "empty headline",
			headline: "",
			wantMin:  0,
			wantMax:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeHeadlineQuality(tt.headline)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeHeadlineQuality(%q) = %d, want between %d and %d", tt.headline, score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAnalyzeMarketingFluff(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantMin int
		wantMax int
	}{
		{
			name:    "contains revolutionary",
			text:    "Our revolutionary new product",
			wantMin: 9,
			wantMax: 9,
		},
		{
			name:    "contains game-changing",
			text:    "This game-changing solution",
			wantMin: 9,
			wantMax: 9,
		},
		{
			name:    "clean professional text",
			text:    "The product reduces latency by 50%",
			wantMin: 10,
			wantMax: 10,
		},
		{
			name:    "multiple fluff words",
			text:    "Revolutionary game-changing breakthrough disruptive innovative solution",
			wantMin: 0,
			wantMax: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeMarketingFluff(tt.text)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeMarketingFluff(%q) = %d, want between %d and %d", tt.text, score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAnalyzePRQuotes(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		wantQuoteCount   int
		wantMetricsCount int
		wantScoreMin     int
		wantScoreMax     int
	}{
		{
			name:             "no quotes",
			content:          "This is a press release with no quotes.",
			wantQuoteCount:   0,
			wantMetricsCount: 0,
			wantScoreMin:     0,
			wantScoreMax:     0,
		},
		{
			name:             "quote with metrics",
			content:          `"We reduced operational costs by 50% in the first quarter," said the CEO.`,
			wantQuoteCount:   1,
			wantMetricsCount: 1,
			wantScoreMin:     5,
			wantScoreMax:     15,
		},
		{
			name:             "quote without metrics",
			content:          `"This is a great product that our team really enjoys using," said the customer.`,
			wantQuoteCount:   1,
			wantMetricsCount: 0,
			wantScoreMin:     0,
			wantScoreMax:     5,
		},
		{
			name: "multiple quotes with mixed metrics",
			content: `"We improved performance by 75% compared to the previous version," said Alice.
"The system is 10x faster than our legacy infrastructure," said Bob.
"Great product that solved our main challenges," said Charlie.`,
			wantQuoteCount:   3,
			wantMetricsCount: 2,
			wantScoreMin:     5,
			wantScoreMax:     15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := analyzePRQuotes(tt.content)

			if score.TotalQuotes != tt.wantQuoteCount {
				t.Errorf("TotalQuotes = %d, want %d", score.TotalQuotes, tt.wantQuoteCount)
			}

			if score.QuotesWithMetrics != tt.wantMetricsCount {
				t.Errorf("QuotesWithMetrics = %d, want %d", score.QuotesWithMetrics, tt.wantMetricsCount)
			}

			quoteScore := score.QualityBreakdown.QuoteScore
			if quoteScore < tt.wantScoreMin || quoteScore > tt.wantScoreMax {
				t.Errorf("QuoteScore = %d, want between %d and %d", quoteScore, tt.wantScoreMin, tt.wantScoreMax)
			}
		})
	}
}

func TestGenerateMarkdownReport(t *testing.T) {
	sections := &SpecSections{
		Title:        "Test PR-FAQ",
		PressRelease: "Test content",
	}

	score := &PRScore{
		OverallScore: 75,
		QualityBreakdown: PRQualityBreakdown{
			HeadlineScore:    8,
			HookScore:        12,
			ReleaseDateScore: 5,
			FiveWsScore:      12,
			CredibilityScore: 8,
			StructureScore:   7,
			ToneScore:        8,
			FluffScore:       9,
			QuoteScore:       6,
		},
	}

	report := GenerateMarkdownReport(sections, score)

	// Verify report contains key sections
	requiredSections := []string{
		"# PR-FAQ Analysis Report",
		"## Executive Summary",
		"## Scoring Results",
		"Test PR-FAQ",
		"75/100",
	}

	for _, section := range requiredSections {
		if !contains(report, section) {
			t.Errorf("Report missing required section: %q", section)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAnalyzeNewswortyHook(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMin int
		wantMax int
	}{
		{
			name:    "complete hook with date and location",
			content: "SEATTLE, WA - November 20, 2025 - Company announces new product.",
			wantMin: 7,
			wantMax: 10,
		},
		{
			name:    "partial hook with date only",
			content: "November 20, 2025 - Company announces new product.",
			wantMin: 5,
			wantMax: 12,
		},
		{
			name:    "no hook elements",
			content: "Company has a new product.",
			wantMin: 0,
			wantMax: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeNewswortyHook(tt.content)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeNewswortyHook() = %d, want between %d and %d", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAnalyzeFiveWs(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMin int
		wantMax int
	}{
		{
			name: "complete 5Ws coverage",
			content: `Company announces new product today.
The product helps customers reduce costs.
It works by optimizing processes.
Available in Seattle starting next month.
Customers can sign up at website.com.`,
			wantMin: 5,
			wantMax: 8,
		},
		{
			name:    "minimal coverage",
			content: "New product available.",
			wantMin: 0,
			wantMax: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeFiveWs(tt.content)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeFiveWs() = %d, want between %d and %d", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAnalyzeToneAndReadability(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMin int
		wantMax int
	}{
		{
			name:    "professional tone",
			content: "The company announced a new product. It reduces costs by 50%. Customers benefit from improved efficiency.",
			wantMin: 6,
			wantMax: 10,
		},
		{
			name:    "overly complex",
			content: "The aforementioned organization has promulgated a revolutionary paradigm-shifting solution.",
			wantMin: 5,
			wantMax: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeToneAndReadability(tt.content)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeToneAndReadability() = %d, want between %d and %d", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAnalyzeReleaseDate(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMin int
		wantMax int
	}{
		{
			name:    "contains full date",
			content: "SEATTLE, WA - November 20, 2025 - Company announces product.",
			wantMin: 4,
			wantMax: 5,
		},
		{
			name:    "contains month and year",
			content: "In November 2025, the company will launch.",
			wantMin: 0,
			wantMax: 2,
		},
		{
			name:    "no date",
			content: "Company announces new product.",
			wantMin: 0,
			wantMax: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeReleaseDate(tt.content)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("analyzeReleaseDate() = %d, want between %d and %d", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestComprehensivePRAnalysis(t *testing.T) {
	tests := []struct {
		name         string
		prContent    string
		wantScoreMin int
		wantScoreMax int
	}{
		{
			name: "high quality PR",
			prContent: `# Company Launches New Product

SEATTLE, WA - November 20, 2025 - Company announces new product that reduces costs.

The product helps customers save money and improve efficiency. It works by optimizing processes.

"We reduced costs by 50%," said the CEO.

Available starting next month at website.com.`,
			wantScoreMin: 40,
			wantScoreMax: 100,
		},
		{
			name:         "minimal PR",
			prContent:    "New product.",
			wantScoreMin: 0,
			wantScoreMax: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := comprehensivePRAnalysis(tt.prContent, "Test Title", 5)

			if score.OverallScore < tt.wantScoreMin || score.OverallScore > tt.wantScoreMax {
				t.Errorf("comprehensivePRAnalysis() OverallScore = %d, want between %d and %d",
					score.OverallScore, tt.wantScoreMin, tt.wantScoreMax)
			}

			// Verify breakdown is populated
			if score.QualityBreakdown.HeadlineScore < 0 {
				t.Error("HeadlineScore should be >= 0")
			}
		})
	}
}

// Test categorizeIssues function
func TestCategorizeIssues(t *testing.T) {
	tests := []struct {
		name     string
		issues   []string
		expected map[string][]string
	}{
		{
			name:   "headline issue",
			issues: []string{"Headline is too long"},
			expected: map[string][]string{
				"Headline & Title": {"Headline is too long"},
			},
		},
		{
			name:   "hook issue",
			issues: []string{"Opening hook is weak"},
			expected: map[string][]string{
				"Opening Hook": {"Opening hook is weak"},
			},
		},
		{
			name:   "5Ws issue",
			issues: []string{"Missing who in the story"},
			expected: map[string][]string{
				"5 Ws Coverage": {"Missing who in the story"},
			},
		},
		{
			name:   "quote issue",
			issues: []string{"Quote lacks metrics"},
			expected: map[string][]string{
				"Customer Evidence": {"Quote lacks metrics"},
			},
		},
		{
			name:   "fluff issue",
			issues: []string{"Too much marketing fluff"},
			expected: map[string][]string{
				"Professional Tone": {"Too much marketing fluff"},
			},
		},
		{
			name:   "structure issue",
			issues: []string{"Poor paragraph structure"},
			expected: map[string][]string{
				"Document Structure": {"Poor paragraph structure"},
			},
		},
		{
			name:   "readability issue",
			issues: []string{"Sentence too long"},
			expected: map[string][]string{
				"Writing Quality": {"Sentence too long"},
			},
		},
		{
			name:   "general issue",
			issues: []string{"Some other problem"},
			expected: map[string][]string{
				"General": {"Some other problem"},
			},
		},
		{
			name:   "multiple issues",
			issues: []string{"Headline too long", "Missing who", "Some other issue"},
			expected: map[string][]string{
				"Headline & Title": {"Headline too long"},
				"5 Ws Coverage":    {"Missing who"},
				"General":          {"Some other issue"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := categorizeIssues(tt.issues)
			for category, expectedIssues := range tt.expected {
				if len(result[category]) != len(expectedIssues) {
					t.Errorf("category %q: got %d issues, want %d", category, len(result[category]), len(expectedIssues))
				}
			}
		})
	}
}

// Test isPressReleaseContent function
func TestIsPressReleaseContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "with date and announce",
			content: "January 15, 2025 - Company announces new product today",
			want:    true,
		},
		{
			name:    "with for immediate release",
			content: "FOR IMMEDIATE RELEASE: New product available",
			want:    true,
		},
		{
			name:    "with business wire",
			content: "BUSINESS WIRE - Company launches product",
			want:    true,
		},
		{
			name:    "with today announced",
			content: "Company today announced a new initiative",
			want:    true,
		},
		{
			name:    "with is excited to announce",
			content: "We are is excited to announce our new product",
			want:    true,
		},
		{
			name:    "plain text without PR indicators",
			content: "This is just some regular text without any special indicators",
			want:    false,
		},
		{
			name:    "date without announce words",
			content: "January 15, 2025 - Regular meeting notes",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPressReleaseContent(tt.content)
			if got != tt.want {
				t.Errorf("isPressReleaseContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test isNumberedFAQQuestion function
func TestIsNumberedFAQQuestion(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{
			name:   "numbered with period",
			header: "1. What is this product?",
			want:   true,
		},
		{
			name:   "numbered with parenthesis",
			header: "2) How does it work?",
			want:   true,
		},
		{
			name:   "Q prefix with period",
			header: "Q1. What are the benefits?",
			want:   true,
		},
		{
			name:   "Q prefix with parenthesis",
			header: "Q2) Who is this for?",
			want:   true,
		},
		{
			name:   "Question prefix",
			header: "Question 3: When is it available?",
			want:   true,
		},
		{
			name:   "plain question",
			header: "What is this product?",
			want:   false,
		},
		{
			name:   "just a number",
			header: "1",
			want:   false,
		},
		{
			name:   "empty string",
			header: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNumberedFAQQuestion(tt.header)
			if got != tt.want {
				t.Errorf("isNumberedFAQQuestion(%q) = %v, want %v", tt.header, got, tt.want)
			}
		})
	}
}

// Benchmark tests for performance-critical functions
func BenchmarkDetectMetricsInText(b *testing.B) {
	text := "We improved performance by 75% and reduced costs by 50%, achieving 10x throughput with $1.5M in savings."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectMetricsInText(text)
	}
}

func BenchmarkAnalyzePRQuotes(b *testing.B) {
	content := `"We reduced costs by 50%," said the CEO.
"Performance improved 10x," said the CTO.
"Customers love it," said a user.`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzePRQuotes(content)
	}
}

func BenchmarkComprehensivePRAnalysis(b *testing.B) {
	content := `# Company Launches New Product

SEATTLE, WA - November 20, 2025 - Company announces new product that reduces costs.

The product helps customers save money and improve efficiency. It works by optimizing processes.

"We reduced costs by 50%," said the CEO.

Available starting next month at website.com.`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comprehensivePRAnalysis(content, "Company Launches New Product", 8)
	}
}

// Test getScoreStatus function
func TestGetScoreStatus(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		maxScore int
		want     string
	}{
		{"excellent", 9, 10, "🟢 Excellent"},
		{"good", 7, 10, "🟡 Good"},
		{"needs work", 5, 10, "🟠 Needs Work"},
		{"critical", 2, 10, "🔴 Critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getScoreStatus(tt.score, tt.maxScore)
			if got != tt.want {
				t.Errorf("getScoreStatus(%d, %d) = %q, want %q", tt.score, tt.maxScore, got, tt.want)
			}
		})
	}
}

// Test getPriority function
func TestGetPriority(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		maxScore int
		want     string
	}{
		{"low priority", 9, 10, "Low"},
		{"medium priority", 7, 10, "Medium"},
		{"high priority", 5, 10, "High"},
		{"critical priority", 2, 10, "Critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPriority(tt.score, tt.maxScore)
			if got != tt.want {
				t.Errorf("getPriority(%d, %d) = %q, want %q", tt.score, tt.maxScore, got, tt.want)
			}
		})
	}
}

// Test getOverallStatus function
func TestGetOverallStatus(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{"ready", 85, "🟢 Ready"},
		{"good", 65, "🟡 Good"},
		{"needs work", 45, "🟠 Needs Work"},
		{"major issues", 25, "🔴 Major Issues"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOverallStatus(tt.score)
			if got != tt.want {
				t.Errorf("getOverallStatus(%d) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

// Test GenerateMarkdownReport with more edge cases
func TestGenerateMarkdownReportEdgeCases(t *testing.T) {
	// Test with low score (triggers different status)
	sections := &SpecSections{
		Title:        "Test PR-FAQ",
		PressRelease: "Test press release content",
		FAQs:         "Q: Test?\nA: Yes.",
	}

	lowScore := &PRScore{
		OverallScore: 30,
		QualityBreakdown: PRQualityBreakdown{
			HeadlineScore: 2,
			HookScore:     3,
		},
	}

	report := GenerateMarkdownReport(sections, lowScore)
	if report == "" {
		t.Error("GenerateMarkdownReport with low score returned empty string")
	}

	// Test with empty sections but valid PRScore
	emptySections := &SpecSections{}
	emptyScore := &PRScore{OverallScore: 50}
	emptyReport := GenerateMarkdownReport(emptySections, emptyScore)
	if emptyReport == "" {
		t.Error("GenerateMarkdownReport with empty sections returned empty string")
	}
}

// Test analyzeStructure function with additional cases
func TestAnalyzeStructureAdditional(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantMinScore int
	}{
		{
			name:         "empty content",
			content:      "",
			wantMinScore: 0,
		},
		{
			name: "good structure with quotes",
			content: `SEATTLE, WA - November 20, 2025 - Company announces new product.

The product helps customers save time and money. It works by optimizing processes.

"We reduced costs by 50%," said the CEO.

Available starting next month at website.com.`,
			wantMinScore: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _, _ := analyzeStructure(tt.content)
			if score < tt.wantMinScore {
				t.Errorf("analyzeStructure() score = %d, want >= %d", score, tt.wantMinScore)
			}
		})
	}
}
