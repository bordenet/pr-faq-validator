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
