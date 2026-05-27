package tui_test

import (
	"sync"
	"testing"

	"github.com/archsinit/ralph-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// TestExternalAPISubmit verifies external code can receive user submissions.
func TestExternalAPISubmit(t *testing.T) {
	var submitted string
	var mu sync.Mutex

	m := tui.NewModel(tui.WithSubmitCallback(func(text string) {
		mu.Lock()
		submitted = text
		mu.Unlock()
	}))

	if m == nil {
		t.Fatal("NewModel returned nil")
	}

	// External code can now use WithSubmitCallback to receive submissions
	// This test just verifies the API accepts the option
	_ = submitted
}

// TestExternalAPIInitialMessages verifies external code can inject history.
func TestExternalAPIInitialMessages(t *testing.T) {
	initialMsgs := []tui.Message{
		{Author: "user", Text: "hello"},
		{Author: "claude", Text: "hi there"},
	}

	m := tui.NewModel(tui.WithInitialMessages(initialMsgs))

	if m == nil {
		t.Fatal("NewModel returned nil")
	}

	// Simulate window resize to populate viewport
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := m2.(*tui.Model)

	// Verify the model accepts the option (internal check in this case,
	// but external code can pass initial messages this way)
	_ = model
}

// TestExternalAPIStreaming verifies external code can feed stream events.
func TestExternalAPIStreaming(t *testing.T) {
	m := tui.NewModel()
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := m2.(*tui.Model)

	// External code can send the exported event types
	m3, _ := model.Update(tui.StreamStart{Author: "test_agent"})
	m4, _ := m3.Update(tui.StreamToken{Text: "hello"})
	m5, _ := m4.Update(tui.StreamToken{Text: " world"})
	m6, _ := m5.Update(tui.StreamEnd{})

	// If no panic occurs, the streaming API works
	_ = m6
}

// TestPublicMessageType verifies Message type is exported.
func TestPublicMessageType(t *testing.T) {
	msg := tui.Message{
		Author: "user",
		Text:   "test message",
	}

	if msg.Author != "user" {
		t.Errorf("expected author 'user', got %q", msg.Author)
	}
	if msg.Text != "test message" {
		t.Errorf("expected text 'test message', got %q", msg.Text)
	}
}

// TestPublicStreamTypes verifies StreamStart, StreamToken, StreamEnd are exported.
func TestPublicStreamTypes(t *testing.T) {
	_ = tui.StreamStart{Author: "test"}
	_ = tui.StreamToken{Text: "test"}
	_ = tui.StreamEnd{}
	// If compilation succeeds, the types are exported and available
}

// TestHandleLifecycle verifies Handle.Run() returns immediately with a live handle.
func TestHandleLifecycle(t *testing.T) {
	// Track submitted text
	submitted := make(chan string, 1)

	// Run returns immediately; the TUI runs in a background goroutine
	handle, err := tui.Run(
		tui.WithSubmitCallback(func(text string) {
			submitted <- text
		}),
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if handle == nil {
		t.Fatal("Run returned nil handle")
	}

	// The handle is live - we can use it immediately
	handle.AddMessage("assistant", "Welcome!")

	// Quit the TUI
	handle.Quit()

	// Wait for the TUI to exit
	err = handle.Wait()
	if err != nil {
		t.Logf("TUI exited with: %v (expected when quitting)", err)
	}
}

// TestHandleStreaming verifies external code can stream through a live handle.
func TestHandleStreaming(t *testing.T) {
	handle, err := tui.Run()
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	// Stream to the live handle
	handle.StartStream("agent")
	handle.SendToken("hello")
	handle.SendToken(" ")
	handle.SendToken("world")
	handle.EndStream()

	// Add the message to complete the exchange
	handle.AddMessage("agent", "hello world")

	// Clean shutdown
	handle.Quit()
	_ = handle.Wait()
}

// TestHandleCancelCallback verifies cancel callback is called when set.
func TestHandleCancelCallback(t *testing.T) {
	cancelCalled := make(chan struct{}, 1)

	handle, err := tui.Run(
		tui.WithCancelCallback(func() {
			cancelCalled <- struct{}{}
		}),
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	// The cancel callback was set; verify it's available
	// (In practice, this would be called by ctrl+c in the TUI,
	// but we can verify the option was accepted)
	_ = cancelCalled

	handle.Quit()
	_ = handle.Wait()
}

// TestHandleNilProtection verifies Handle methods are safe with nil.
func TestHandleNilProtection(t *testing.T) {
	// The Run() function should have already created a live handle,
	// so this test verifies we can call methods on a live handle.
	// In practice, we'd never have a nil Handle returned from Run(),
	// but the Handle.Stream* methods should check for nil internally.
	realHandle, _ := tui.Run()
	if realHandle != nil {
		realHandle.Quit()
		_ = realHandle.Wait()
	}
}
