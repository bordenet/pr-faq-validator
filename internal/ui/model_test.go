package ui

import (
	"testing"

	"github.com/bordenet/pr-faq-validator/internal/parser"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	sections := parser.SpecSections{
		Title:        "Test PR-FAQ",
		PressRelease: "Test content",
		PRScore: &parser.PRScore{
			OverallScore: 75,
		},
	}

	model := NewModel(sections)

	if model.activeTab != TabOverview {
		t.Errorf("activeTab = %v, want %v", model.activeTab, TabOverview)
	}

	if len(model.tabs) != 4 {
		t.Errorf("tabs length = %d, want 4", len(model.tabs))
	}

	if model.sections.Title != "Test PR-FAQ" {
		t.Errorf("sections.Title = %q, want %q", model.sections.Title, "Test PR-FAQ")
	}
}

func TestModel_Init(t *testing.T) {
	sections := parser.SpecSections{
		Title: "Test",
		PRScore: &parser.PRScore{
			OverallScore: 50,
		},
	}

	model := NewModel(sections)
	cmd := model.Init()

	// Init should return a command to start AI analysis
	if cmd == nil {
		t.Error("Init() should return a command, got nil")
	}
}

func TestModel_Update_Quit(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test quit with 'q'
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}

	// Verify model is returned
	if _, ok := updatedModel.(Model); !ok {
		t.Error("Update should return Model type")
	}
}

func TestModel_Update_TabNavigation(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test tab navigation with right arrow
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	m := updatedModel.(Model)

	if m.activeTab != TabBreakdown {
		t.Errorf("After right arrow, activeTab = %v, want %v", m.activeTab, TabBreakdown)
	}

	// Test tab navigation with left arrow
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updatedModel.(Model)

	if m.activeTab != TabOverview {
		t.Errorf("After left arrow, activeTab = %v, want %v", m.activeTab, TabOverview)
	}
}

func TestModel_Update_HelpToggle(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test help toggle with '?'
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m := updatedModel.(Model)

	if !m.showHelp {
		t.Error("Expected showHelp to be true after pressing '?'")
	}

	// Toggle again
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updatedModel.(Model)

	if m.showHelp {
		t.Error("Expected showHelp to be false after pressing '?' again")
	}
}

func TestModel_Update_WindowSize(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test window size message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m := updatedModel.(Model)

	if m.windowWidth != 120 {
		t.Errorf("windowWidth = %d, want 120", m.windowWidth)
	}

	if m.windowHeight != 40 {
		t.Errorf("windowHeight = %d, want 40", m.windowHeight)
	}
}

func TestSetFeedbackMsg(t *testing.T) {
	msg := SetFeedbackMsg{
		Section:  "Test",
		Feedback: "Test feedback",
	}

	if msg.Section != "Test" {
		t.Errorf("Section = %q, want %q", msg.Section, "Test")
	}

	if msg.Feedback != "Test feedback" {
		t.Errorf("Feedback = %q, want %q", msg.Feedback, "Test feedback")
	}
}

// Test RenderHeader function
func TestRenderHeader(t *testing.T) {
	tests := []struct {
		name  string
		title string
		score int
	}{
		{"high score", "Test Document", 85},
		{"medium score", "Another Doc", 55},
		{"low score", "Poor Doc", 25},
		{"empty title", "", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderHeader(tt.title, tt.score)
			if result == "" {
				t.Error("RenderHeader returned empty string")
			}
		})
	}
}

// Test RenderScoreBreakdown function
func TestRenderScoreBreakdown(t *testing.T) {
	breakdown := parser.PRQualityBreakdown{
		HeadlineScore:    8,
		HookScore:        12,
		ReleaseDateScore: 5,
		FiveWsScore:      15,
		CredibilityScore: 8,
		StructureScore:   8,
		ToneScore:        8,
		FluffScore:       8,
		QuoteScore:       10,
	}

	result := RenderScoreBreakdown(breakdown)
	if result == "" {
		t.Error("RenderScoreBreakdown returned empty string")
	}
}

// Test RenderStrengths function
func TestRenderStrengths(t *testing.T) {
	tests := []struct {
		name      string
		strengths []string
	}{
		{"with strengths", []string{"Good headline", "Strong metrics"}},
		{"empty strengths", []string{}},
		{"single strength", []string{"One strength"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderStrengths(tt.strengths)
			// Should not panic
			_ = result
		})
	}
}

// Test RenderImprovements function
func TestRenderImprovements(t *testing.T) {
	tests := []struct {
		name         string
		improvements []string
	}{
		{"with improvements", []string{"Add metrics", "Improve hook"}},
		{"empty improvements", []string{}},
		{"single improvement", []string{"One improvement"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderImprovements(tt.improvements)
			// Should not panic
			_ = result
		})
	}
}

// Test RenderQuoteAnalysis function
func TestRenderQuoteAnalysis(t *testing.T) {
	tests := []struct {
		name  string
		score parser.PRScore
	}{
		{
			name: "with quotes",
			score: parser.PRScore{
				TotalQuotes:       2,
				QuotesWithMetrics: 1,
				MetricDetails: []parser.MetricInfo{
					{
						Quote:       "We saved 50% on costs",
						Metrics:     []string{"50%"},
						MetricTypes: []string{"percentage"},
						Score:       8,
					},
				},
			},
		},
		{
			name:  "empty quotes",
			score: parser.PRScore{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderQuoteAnalysis(tt.score)
			// Should not panic
			_ = result
		})
	}
}

// Test RenderLLMFeedback function
func TestRenderLLMFeedback(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		feedback string
	}{
		{
			name:     "with feedback",
			title:    "Press Release",
			feedback: "Good structure",
		},
		{
			name:     "empty feedback",
			title:    "FAQ",
			feedback: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderLLMFeedback(tt.title, tt.feedback)
			// Should not panic
			_ = result
		})
	}
}

// Test RenderTabs function
func TestRenderTabs(t *testing.T) {
	tabs := []string{"Overview", "Breakdown", "Quotes", "AI Feedback"}

	for activeTab := 0; activeTab < len(tabs); activeTab++ {
		result := RenderTabs(tabs, activeTab)
		if result == "" {
			t.Errorf("RenderTabs returned empty string for activeTab=%d", activeTab)
		}
	}
}

// Test RenderHelp function
func TestRenderHelp(t *testing.T) {
	result := RenderHelp()
	if result == "" {
		t.Error("RenderHelp returned empty string")
	}
}

// Test RenderStatus function
func TestRenderStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"with status", "Loading..."},
		{"empty status", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderStatus(tt.status)
			// Should not panic
			_ = result
		})
	}
}

// Test GetScoreStyle function
func TestGetScoreStyle(t *testing.T) {
	tests := []struct {
		name  string
		score int
	}{
		{"high score", 85},
		{"medium score", 55},
		{"low score", 25},
		{"boundary high", 70},
		{"boundary medium", 40},
		{"zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := GetScoreStyle(tt.score)
			// Should return a valid style
			_ = style
		})
	}
}

// Test CreateProgressBar function
func TestCreateProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		maxScore int
		width    int
	}{
		{"full bar", 10, 10, 20},
		{"half bar", 5, 10, 20},
		{"empty bar", 0, 10, 20},
		{"zero max", 5, 0, 20},
		{"narrow width", 5, 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateProgressBar(tt.current, tt.maxScore, tt.width)
			// Should not panic
			_ = result
		})
	}
}

// Test FormatScore function
func TestFormatScore(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		maxScore int
	}{
		{"high score", 85, 100},
		{"medium score", 55, 100},
		{"low score", 25, 100},
		{"partial max", 8, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatScore(tt.score, tt.maxScore)
			if result == "" {
				t.Error("FormatScore returned empty string")
			}
		})
	}
}

// Test Model View function
func TestModel_View(t *testing.T) {
	sections := parser.SpecSections{
		Title:        "Test PR-FAQ",
		PressRelease: "Test content",
		FAQs:         "Q: Test?\nA: Yes.",
		PRScore: &parser.PRScore{
			OverallScore:      75,
			TotalQuotes:       1,
			QuotesWithMetrics: 1,
			MetricDetails: []parser.MetricInfo{
				{Quote: "Test quote", Metrics: []string{"50%"}, Score: 5},
			},
			QualityBreakdown: parser.PRQualityBreakdown{
				HeadlineScore: 8,
				HookScore:     12,
				Strengths:     []string{"Good headline"},
				Issues:        []string{"Add metrics"},
			},
		},
	}

	model := NewModel(sections)
	model.windowWidth = 80
	model.windowHeight = 24

	// Test View for each tab
	for tab := TabOverview; tab <= TabFeedback; tab++ {
		model.activeTab = tab
		result := model.View()
		if result == "" {
			t.Errorf("View returned empty string for tab %d", tab)
		}
	}

	// Test with help shown
	model.showHelp = true
	result := model.View()
	if result == "" {
		t.Error("View returned empty string with help shown")
	}
}

// Test SetFeedback function
func TestSetFeedback(t *testing.T) {
	cmd := SetFeedback("Press Release", "Good structure")
	if cmd == nil {
		t.Error("SetFeedback returned nil command")
	}

	// Execute the command and check the message
	msg := cmd()
	feedbackMsg, ok := msg.(SetFeedbackMsg)
	if !ok {
		t.Error("SetFeedback command did not return SetFeedbackMsg")
	}
	if feedbackMsg.Section != "Press Release" {
		t.Errorf("Section = %q, want %q", feedbackMsg.Section, "Press Release")
	}
	if feedbackMsg.Feedback != "Good structure" {
		t.Errorf("Feedback = %q, want %q", feedbackMsg.Feedback, "Good structure")
	}
}

// Test SetStatus function
func TestSetStatus(t *testing.T) {
	cmd := SetStatus("Loading...")
	if cmd == nil {
		t.Error("SetStatus returned nil command")
	}

	// Execute the command and check the message
	msg := cmd()
	statusMsg, ok := msg.(SetStatusMsg)
	if !ok {
		t.Error("SetStatus command did not return SetStatusMsg")
	}
	if string(statusMsg) != "Loading..." {
		t.Errorf("Status = %q, want %q", string(statusMsg), "Loading...")
	}
}

// Test SetLoading function
func TestSetLoading(t *testing.T) {
	cmd := SetLoading(true)
	if cmd == nil {
		t.Error("SetLoading returned nil command")
	}

	// Execute the command and check the message
	msg := cmd()
	loadingMsg, ok := msg.(SetLoadingMsg)
	if !ok {
		t.Error("SetLoading command did not return SetLoadingMsg")
	}
	if !bool(loadingMsg) {
		t.Error("Expected loading to be true")
	}
}

// Test Model Update with SetFeedbackMsg
func TestModel_Update_SetFeedbackMsg(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test SetFeedbackMsg for Press Release
	updatedModel, _ := model.Update(SetFeedbackMsg{Section: "Press Release", Feedback: "Good PR"})
	m := updatedModel.(Model)

	// The model stores feedback internally - we just verify no panic
	_ = m
}

// Test Model Update with SetStatusMsg
func TestModel_Update_SetStatusMsg(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test SetStatusMsg
	updatedModel, _ := model.Update(SetStatusMsg("Processing..."))
	m := updatedModel.(Model)

	if m.status != "Processing..." {
		t.Errorf("status = %q, want %q", m.status, "Processing...")
	}
}

// Test Model Update with SetLoadingMsg
func TestModel_Update_SetLoadingMsg(t *testing.T) {
	sections := parser.SpecSections{
		PRScore: &parser.PRScore{},
	}

	model := NewModel(sections)

	// Test SetLoadingMsg
	updatedModel, _ := model.Update(SetLoadingMsg(true))
	m := updatedModel.(Model)

	if !m.loading {
		t.Error("Expected loading to be true")
	}
}
