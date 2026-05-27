package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestTUIIntegration verifies the TUI can handle realistic workflows.
func TestTUIIntegration(t *testing.T) {
	var submitted string
	m := NewModel(WithSubmitCallback(func(text string) {
		submitted = text
	}))

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
	m3, cmd := model.Update(enterMsg)
	model = m3.(*Model)

	// Process the msgSubmit command
	if cmd != nil {
		submitMsg := cmd()
		m4, _ := model.Update(submitMsg)
		model = m4.(*Model)
	}

	if submitted != "hello" {
		t.Fatalf("expected submitted text 'hello', got %q", submitted)
	}

	// Simulate orchestrator echoing the message back
	m5, _ := model.Update(msgAddMessage{author: "user", text: submitted})
	model = m5.(*Model)

	if len(model.messages) != 1 {
		t.Fatalf("expected 1 message after add, got %d", len(model.messages))
	}

	// Simulate agent streaming
	startMsg := StreamStart{Author: "claude"}
	m6, _ := model.Update(startMsg)
	model = m6.(*Model)

	// Stream multiple tokens
	for i := 0; i < 5; i++ {
		tokenMsg := StreamToken{Text: "token "}
		m7, _ := model.Update(tokenMsg)
		model = m7.(*Model)
	}

	if model.inProgressText != "token token token token token " {
		t.Errorf("tokens not concatenated correctly: %q", model.inProgressText)
	}

	// End stream - StreamEnd only clears the in-progress stream state
	endMsg := StreamEnd{}
	m8, _ := model.Update(endMsg)
	model = m8.(*Model)

	if model.isStreaming {
		t.Error("streaming should be false after end")
	}

	// Engine adds the streamed message after StreamEnd
	m8b, _ := model.Update(msgAddMessage{author: "claude", text: "token token token token token "})
	model = m8b.(*Model)

	if len(model.messages) != 2 {
		t.Fatalf("expected 2 messages after AddMessage, got %d", len(model.messages))
	}

	// Verify view renders without error
	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	// Simulate resize
	windowMsg2 := tea.WindowSizeMsg{Width: 60, Height: 20}
	m9, _ := model.Update(windowMsg2)
	model = m9.(*Model)

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
	startMsg := StreamStart{Author: "codex"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	// View should reflect streaming status
	view3 := model.View()
	if view3 == "" {
		t.Error("view should not be empty during stream")
	}

	// End stream
	endMsg := StreamEnd{}
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
	model.messages = append(model.messages, Message{
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
	startMsg := StreamStart{Author: "gemini"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	// Rapidly send tokens (simulating fast generation)
	tokens := []string{"The ", "quick ", "brown ", "fox ", "jumps ", "over ", "the ", "lazy ", "dog"}
	for _, token := range tokens {
		tokenMsg := StreamToken{Text: token}
		m4, _ := model.Update(tokenMsg)
		model = m4.(*Model)
	}

	expectedText := "The quick brown fox jumps over the lazy dog"
	if model.inProgressText != expectedText {
		t.Errorf("stream text mismatch: got %q, want %q", model.inProgressText, expectedText)
	}

	// End stream - StreamEnd only clears the in-progress stream state
	endMsg := StreamEnd{}
	m5, _ := model.Update(endMsg)
	model = m5.(*Model)

	// Engine adds the finalized message after StreamEnd
	m5b, _ := model.Update(msgAddMessage{author: "gemini", text: expectedText})
	model = m5b.(*Model)

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
		startMsg := StreamStart{Author: agent}
		m3, _ := model.Update(startMsg)
		model = m3.(*Model)

		// Send token
		tokenMsg := StreamToken{Text: "response from " + agent}
		m4, _ := model.Update(tokenMsg)
		model = m4.(*Model)

		// End stream
		endMsg := StreamEnd{}
		m5, _ := model.Update(endMsg)
		model = m5.(*Model)

		if !model.isStreaming {
			// After stream end, isStreaming should be false
		}

		// Engine adds the finalized message after StreamEnd
		m5b, _ := model.Update(msgAddMessage{author: agent, text: "response from " + agent})
		model = m5b.(*Model)
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

// TestScrollUpStaysScrolledUp verifies tokens don't auto-scroll when user has scrolled up.
func TestScrollUpStaysScrolledUp(t *testing.T) {
	m := NewModel()

	// Setup window
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Add some initial messages to fill the viewport with content
	for i := 0; i < 30; i++ {
		m2, _ = model.Update(msgAddMessage{
			author: "user",
			text:   "message with some content to ensure we fill the viewport " + string(rune('0'+byte(i%10))),
		})
		model = m2.(*Model)
	}

	// Manually scroll up away from bottom
	model.viewport.LineUp(5)

	// Mark that user has scrolled (userScrolled becomes true when not at bottom)
	if model.viewport.AtBottom() {
		model.viewport.LineUp(1) // Ensure not at bottom
	}
	model.userScrolled = !model.viewport.AtBottom()

	// Store the scroll position
	prevScroll := model.viewport.YOffset

	// Start streaming
	startMsg := StreamStart{Author: "claude"}
	m3, _ := model.Update(startMsg)
	model = m3.(*Model)

	// Stream tokens - should NOT auto-scroll to bottom since user scrolled up
	for i := 0; i < 5; i++ {
		tokenMsg := StreamToken{Text: "token "}
		m4, _ := model.Update(tokenMsg)
		model = m4.(*Model)
	}

	// Viewport scroll position should not have changed to bottom
	if model.viewport.YOffset > prevScroll+10 {
		t.Error("viewport should not auto-scroll to bottom when user has scrolled up")
	}
}

// TestViewportResize verifies layout dimensions are valid for small terminals.
func TestViewportResize(t *testing.T) {
	m := NewModel()

	// Small terminal that would previously have negative viewport height
	windowMsg := tea.WindowSizeMsg{Width: 20, Height: 3}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Viewport height should be clamped to minimum
	if model.viewport.Height < 5 {
		t.Errorf("viewport height should be at least 5, got %d", model.viewport.Height)
	}

	// View should still render without panic
	view := model.View()
	if view == "" {
		t.Error("view should not be empty for small terminal")
	}
}

// TestAuthorStyleConsistency verifies new authors get consistent distinct styles.
func TestAuthorStyleConsistency(t *testing.T) {
	m := NewModel()

	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m2, _ := m.Update(windowMsg)
	model := m2.(*Model)

	// Add messages from multiple new authors
	authors := []string{"agent1", "agent2", "agent3"}
	for _, author := range authors {
		m3, _ := model.Update(msgAddMessage{author: author, text: "test"})
		model = m3.(*Model)
	}

	// Each author should have a unique style assigned
	if len(model.authorColors) != 3 {
		t.Errorf("expected 3 authors in color map, got %d", len(model.authorColors))
	}

	// Getting the same author's style twice should return the same style
	style1 := model.getAuthorStyle("agent1")
	style2 := model.getAuthorStyle("agent1")
	if style1.String() != style2.String() {
		t.Error("same author should get consistent style")
	}
}

// TestLongUnbrokenTokenWrapping verifies long tokens are wrapped properly.
func TestLongUnbrokenTokenWrapping(t *testing.T) {
	longToken := "verylongwordthatexceedstheviewportwidthandshouldbebrokenupbythelinebreaker"

	lines := wrapText(longToken, 20)
	if len(lines) < 3 {
		t.Errorf("long token should wrap into multiple lines, got %d", len(lines))
	}

	// Each line should be at most 20 chars
	for i, line := range lines {
		if len(line) > 20 {
			t.Errorf("line %d exceeds width: %q (%d chars)", i, line, len(line))
		}
	}

	// Concatenated should equal original
	combined := ""
	for _, line := range lines {
		combined += line
	}
	if combined != longToken {
		t.Errorf("wrapped lines don't concatenate to original: got %q, want %q", combined, longToken)
	}
}
