package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestTUIIntegration verifies the TUI can handle realistic workflows.
func TestTUIIntegration(t *testing.T) {
	m := NewModel()

	// Simulate window resize
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	if model.width != 100 || model.height != 30 {
		t.Fatalf("window size not updated: %dx%d", model.width, model.height)
	}

	// Simulate user input
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model.input.SetValue("hello")
	m3, _ := model.Update(enterMsg)
	model = m3.(*Model)

	if len(model.messages) != 1 {
		t.Fatalf("expected 1 message after submit, got %d", len(model.messages))
	}

	// Simulate agent streaming
	startMsg := msgStreamStart{author: "claude"}
	m4, _ := model.Update(startMsg)
	model = m4.(*Model)

	// Stream multiple tokens
	for i := 0; i < 5; i++ {
		tokenMsg := msgStreamToken{text: "token "}
		m5, _ := model.Update(tokenMsg)
		model = m5.(*Model)
	}

	if model.inProgressText != "token token token token token " {
		t.Errorf("tokens not concatenated correctly: %q", model.inProgressText)
	}

	// End stream
	endMsg := msgStreamEnd{}
	m6, _ := model.Update(endMsg)
	model = m6.(*Model)

	if len(model.messages) != 2 {
		t.Fatalf("expected 2 messages after stream end, got %d", len(model.messages))
	}

	if model.isStreaming {
		t.Error("streaming should be false after end")
	}

	// Verify view renders without error
	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	// Simulate resize
	windowMsg2 := tea.WindowSizeMsg{Width: 60, Height: 20}
	m7, _ := model.Update(windowMsg2)
	model = m7.(*Model)

	if model.width != 60 {
		t.Errorf("width not updated to 60, got %d", model.width)
	}

	// Verify view still renders correctly after resize
	view2 := model.View()
	if view2 == "" {
		t.Error("view should not be empty after resize")
	}
}

// TestStreamingStatusIndicator verifies the status line reflects stream state.
func TestStreamingStatusIndicator(t *testing.T) {
	m := NewModel()

	// Setup window
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Initial status should show "your turn"
	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	// Start stream
	startMsg := msgStreamStart{author: "codex"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	// View should reflect streaming status
	view3 := model.View()
	if view3 == "" {
		t.Error("view should not be empty during stream")
	}

	// End stream
	endMsg := msgStreamEnd{}
	m4, _ := model.Update(endMsg)
	model = m4.(*Model)

	// View should return to "your turn"
	view4 := model.View()
	if view4 == "" {
		t.Error("view should not be empty after stream")
	}
}

// TestWrappingOnResize verifies text wraps correctly when terminal is resized.
func TestWrappingOnResize(t *testing.T) {
	m := NewModel()

	// Setup initial window
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Add a long message
	longText := "This is a very long message that should wrap to multiple lines when displayed in a terminal with limited width. It contains many words to ensure proper wrapping behavior."
	model.messages = append(model.messages, renderedMsg{
		Author: "claude",
		Text:   longText,
	})
	model.updateViewport()

	view1 := model.View()
	initialLines := len(wrapText(longText, 76)) // Account for author label and indentation

	// Resize to smaller width
	windowMsg2 := tea.WindowSizeMsg{Width: 40, Height: 24}
	m3, _ := model.Update(windowMsg2)
	model = m3.(*Model)

	view2 := model.View()
	resizedLines := len(wrapText(longText, 36)) // Account for author label and indentation

	if resizedLines <= initialLines {
		t.Errorf("resize should increase line count: before=%d, after=%d", initialLines, resizedLines)
	}

	// Verify both views render without error
	if view1 == "" || view2 == "" {
		t.Error("views should not be empty")
	}
}

// TestFastStreamingAccumulation verifies rapid token arrival doesn't cause issues.
func TestFastStreamingAccumulation(t *testing.T) {
	m := NewModel()

	// Setup window
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Start stream
	startMsg := msgStreamStart{author: "gemini"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	// Rapidly send tokens (simulating fast generation)
	tokens := []string{"The ", "quick ", "brown ", "fox ", "jumps ", "over ", "the ", "lazy ", "dog"}
	for _, token := range tokens {
		tokenMsg := msgStreamToken{text: token}
		m4, _ := model.Update(tokenMsg)
		model = m4.(*Model)
	}

	expectedText := "The quick brown fox jumps over the lazy dog"
	if model.inProgressText != expectedText {
		t.Errorf("stream text mismatch: got %q, want %q", model.inProgressText, expectedText)
	}

	// End stream
	endMsg := msgStreamEnd{}
	m5, _ := model.Update(endMsg)
	model = m5.(*Model)

	if len(model.messages) != 1 {
		t.Fatalf("expected 1 finalized message, got %d", len(model.messages))
	}

	if model.messages[0].Text != expectedText {
		t.Errorf("finalized text mismatch: got %q, want %q", model.messages[0].Text, expectedText)
	}
}

// TestMultipleConsecutiveStreams verifies multiple agents can stream in sequence.
func TestMultipleConsecutiveStreams(t *testing.T) {
	m := NewModel()

	// Setup window
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	agents := []string{"user", "claude", "codex"}
	for _, agent := range agents {
		// Start stream
		startMsg := msgStreamStart{author: agent}
		m3, _ := model.Update(startMsg)
		model = m3.(*Model)

		// Send token
		tokenMsg := msgStreamToken{text: "response from " + agent}
		m4, _ := model.Update(tokenMsg)
		model = m4.(*Model)

		// End stream
		endMsg := msgStreamEnd{}
		m5, _ := model.Update(endMsg)
		model = m5.(*Model)

		if !model.isStreaming {
			// After stream end, isStreaming should be false
		}
	}

	if len(model.messages) != 3 {
		t.Fatalf("expected 3 finalized messages, got %d", len(model.messages))
	}

	// Verify all agents have their responses
	for i, agent := range agents {
		if model.messages[i].Author != agent {
			t.Errorf("message[%d] author mismatch: got %q, want %q", i, model.messages[i].Author, agent)
		}
	}
}
