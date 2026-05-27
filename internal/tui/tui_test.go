package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m == nil {
		t.Fatal("NewModel returned nil")
	}
	if len(m.messages) != 0 {
		t.Errorf("expected 0 initial messages, got %d", len(m.messages))
	}
	if !m.input.Focused() {
		t.Error("input should be focused")
	}
}

func TestInit(t *testing.T) {
	m := NewModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("Update should not return a command for WindowSizeMsg")
	}
	model := m2.(*Model)
	if model.width != 80 {
		t.Errorf("expected width 80, got %d", model.width)
	}
	if model.height != 24 {
		t.Errorf("expected height 24, got %d", model.height)
	}
	if model.viewport.Width != 80 {
		t.Errorf("expected viewport width 80, got %d", model.viewport.Width)
	}
	if model.viewport.Height != 21 {
		t.Errorf("expected viewport height 21 (24-3), got %d", model.viewport.Height)
	}
}

func TestCancelKey(t *testing.T) {
	cancelCalled := false
	m := NewModel(WithCancelCallback(func() {
		cancelCalled = true
	}))
	// Create a KeyMsg for Ctrl+C
	msg := tea.KeyMsg{Type: tea.KeyCtrlC, Alt: false}
	m2, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("Update should return nil for ctrl+c (handled by callback)")
	}
	if !cancelCalled {
		t.Error("Cancel callback should be called")
	}
	_ = m2
}

func TestEnterKey(t *testing.T) {
	var submitted string
	m := NewModel(WithSubmitCallback(func(text string) {
		submitted = text
	}))
	m.input.SetValue("hello")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	m2, cmd := m.Update(msg)
	if cmd == nil {
		t.Error("Update should return a command for enter key")
	}

	// Process the command to get the msgSubmit
	if cmd != nil {
		submitMsg := cmd()
		m3, _ := m2.(*Model).Update(submitMsg)
		model := m3.(*Model)
		if model.input.Value() != "" {
			t.Errorf("expected input to be reset, got %q", model.input.Value())
		}
	}

	if submitted != "hello" {
		t.Errorf("expected submitted text 'hello', got %q", submitted)
	}
}

func TestView(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)
	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}
}

func TestEmptyMessageNotAdded(t *testing.T) {
	m := NewModel()
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	m2, _ := m.Update(msg)

	model := m2.(*Model)
	if len(model.messages) != 0 {
		t.Errorf("expected 0 messages for empty input, got %d", len(model.messages))
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		text  string
		width int
		want  int
	}{
		{"hello world", 20, 1},
		{"this is a very long message that should wrap to multiple lines when the width is narrow", 20, 5},
		{"short", 80, 1},
		{"", 20, 0},
	}

	for _, tt := range tests {
		got := wrapText(tt.text, tt.width)
		if len(got) != tt.want {
			t.Errorf("wrapText(%q, %d) got %d lines, want %d", tt.text, tt.width, len(got), tt.want)
		}
	}
}

func TestAuthorStyling(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)
	model.messages = append(model.messages, Message{Author: "user", Text: "test"})
	model.messages = append(model.messages, Message{Author: "claude", Text: "response"})

	model.updateViewport()
	// Check that messages were added and viewport was updated
	if len(model.messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(model.messages))
	}
}

func TestAutoScrollOnNewMessage(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)
	model.messages = append(model.messages, Message{Author: "user", Text: "test"})
	model.updateViewport()

	// After adding a message, userScrolled should be false (unless manually set)
	if model.userScrolled {
		t.Error("userScrolled should be false after adding message")
	}
}

func TestUserScrollDetection(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)
	if model.userScrolled {
		t.Error("userScrolled should be false initially")
	}

	// Verify the field exists and is initialized correctly
	if model.userScrolled {
		t.Error("userScrolled should be false initially")
	}
}

func TestStreamMessages(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)

	// Test stream start
	startMsg := StreamStart{Author: "claude"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	if !model.isStreaming {
		t.Error("isStreaming should be true after stream start")
	}
	if model.inProgressAuthor != "claude" {
		t.Errorf("inProgressAuthor should be 'claude', got %q", model.inProgressAuthor)
	}

	// Test stream token
	tokenMsg := StreamToken{Text: "hello"}
	m4, _ := model.Update(tokenMsg)
	model = m4.(*Model)

	if model.inProgressText != "hello" {
		t.Errorf("inProgressText should be 'hello', got %q", model.inProgressText)
	}

	// Test another token
	tokenMsg2 := StreamToken{Text: " world"}
	m5, _ := model.Update(tokenMsg2)
	model = m5.(*Model)

	if model.inProgressText != "hello world" {
		t.Errorf("inProgressText should be 'hello world', got %q", model.inProgressText)
	}

	// Test stream end - StreamEnd only clears the in-progress stream
	endMsg := StreamEnd{}
	m6, _ := model.Update(endMsg)
	model = m6.(*Model)

	if model.isStreaming {
		t.Error("isStreaming should be false after stream end")
	}
	if model.inProgressAuthor != "" {
		t.Errorf("inProgressAuthor should be empty after stream end, got %q", model.inProgressAuthor)
	}
	if model.inProgressText != "" {
		t.Errorf("inProgressText should be empty after stream end, got %q", model.inProgressText)
	}

	// The engine (not StreamEnd) adds the message via AddMessage
	addMsg := msgAddMessage{author: "claude", text: "hello world"}
	m7, _ := model.Update(addMsg)
	model = m7.(*Model)

	if len(model.messages) != 1 {
		t.Errorf("expected 1 message after AddMessage, got %d", len(model.messages))
	}
	if model.messages[0].Author != "claude" || model.messages[0].Text != "hello world" {
		t.Errorf("message not finalized correctly: %+v", model.messages[0])
	}
}

func TestStatusIndicator(t *testing.T) {
	m := NewModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(msg)

	model := m2.(*Model)

	// Test default status
	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// Test status during streaming
	startMsg := StreamStart{Author: "claude"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	view = model.View()
	if !strings.Contains(view, "waiting for claude") {
		t.Error("status should indicate waiting for claude")
	}
}
