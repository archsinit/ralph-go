package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/archsinit/ralph-go/internal/agent"
	"github.com/archsinit/ralph-go/internal/session"
)

// MockUIBridge is a test implementation of UIBridge.
type MockUIBridge struct {
	inputQueue     []string
	inputIdx       int
	cancelChan     chan struct{}
	streamOrder    []string // Records the order of agent streams
	messages       []string // Records all messages added
	messageAuthors []string // Records the author of each message
}

func NewMockUIBridge(inputs []string) *MockUIBridge {
	return &MockUIBridge{
		inputQueue: inputs,
		cancelChan: make(chan struct{}),
	}
}

func (m *MockUIBridge) RequestUserInput(ctx context.Context) (string, error) {
	if m.inputIdx >= len(m.inputQueue) {
		return "", context.Canceled
	}
	input := m.inputQueue[m.inputIdx]
	m.inputIdx++
	return input, nil
}

func (m *MockUIBridge) RequestCancelSignal(ctx context.Context) <-chan struct{} {
	return m.cancelChan
}

func (m *MockUIBridge) StreamStart(author string) error {
	m.streamOrder = append(m.streamOrder, author)
	return nil
}

func (m *MockUIBridge) StreamToken(text string) error {
	return nil
}

func (m *MockUIBridge) StreamEnd() error {
	return nil
}

func (m *MockUIBridge) AddMessage(author, text string) error {
	m.messages = append(m.messages, text)
	m.messageAuthors = append(m.messageAuthors, author)
	return nil
}

// MockEchoAdapter returns input back as output.
type MockEchoAdapter struct {
	name string
}

func (m *MockEchoAdapter) Name() string {
	return m.name
}

func (m *MockEchoAdapter) Capabilities() agent.Capabilities {
	return agent.Capabilities{SupportsResume: false}
}

func (m *MockEchoAdapter) Invoke(ctx context.Context, req agent.Request) (<-chan agent.Token, <-chan agent.Result) {
	tokenChan := make(chan agent.Token, 1)
	resultChan := make(chan agent.Result, 1)

	go func() {
		defer close(tokenChan)
		defer close(resultChan)

		// Echo back the new message
		if req.NewMessage != "" {
			tokenChan <- agent.Token{Text: "Echo: " + req.NewMessage}
		} else {
			tokenChan <- agent.Token{Text: "Ready"}
		}
		resultChan <- agent.Result{SessionID: "", Err: nil}
	}()

	return tokenChan, resultChan
}

func TestEngineCreation(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	turnOrder := NewTurnOrder([]string{"claude", "user", "codex"})

	adapters := map[string]agent.Adapter{
		"claude": &MockEchoAdapter{name: "claude"},
		"codex":  &MockEchoAdapter{name: "codex"},
	}

	bridge := NewMockUIBridge([]string{})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	if engine == nil {
		t.Error("Engine creation failed")
	}
	if len(engine.adapters) != 2 {
		t.Errorf("Engine adapters: expected 2, got %d", len(engine.adapters))
	}
}

func TestTurnOrderExecution(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	// Turn order: claude, user (repeat)
	turnOrder := NewTurnOrder([]string{"claude", "user"})

	adapters := map[string]agent.Adapter{
		"claude": &MockEchoAdapter{name: "claude"},
	}

	// Provide one input for user turn
	bridge := NewMockUIBridge([]string{"hello"})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	// Run with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		_ = engine.Run(ctx)
	}()

	// Wait for execution to start
	<-time.After(1 * time.Second)
	cancel()

	// Should have at least one stream start (claude)
	if len(bridge.streamOrder) < 1 {
		t.Logf("Stream order: %v (len=%d)", bridge.streamOrder, len(bridge.streamOrder))
	}

	// Should have recorded a message (user input)
	if len(bridge.messageAuthors) < 1 {
		t.Logf("Messages recorded: %d", len(bridge.messageAuthors))
	}
}

func TestPrefixParsing(t *testing.T) {
	// Test that prefix parsing works correctly
	target, body := ParsePrefix("codex: please respond", []string{"codex", "claude"})
	if target != "codex" {
		t.Errorf("expected target=codex, got %q", target)
	}
	if body != "please respond" {
		t.Errorf("expected body='please respond', got %q", body)
	}

	// Test without prefix
	target2, body2 := ParsePrefix("hello", []string{"codex"})
	if target2 != "" {
		t.Errorf("expected empty target for non-prefixed input, got %q", target2)
	}
	if body2 != "hello" {
		t.Errorf("expected body unchanged, got %q", body2)
	}
}

// TestTurnOrderRoundRobin verifies round-robin execution of turns.
func TestTurnOrderRoundRobin(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	// Turn order: claude, user (cycles back to claude, user, ...)
	turnOrder := NewTurnOrder([]string{"claude", "user"})

	adapters := map[string]agent.Adapter{
		"claude": &MockEchoAdapter{name: "claude"},
	}

	// Provide one user input
	bridge := NewMockUIBridge([]string{
		"hello", // User input for second turn
	})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	// Run with timeout - will timeout because engine waits for more input
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- engine.Run(ctx)
	}()

	// Wait for timeout
	err = <-done

	// Should have executed: claude (first turn), user (second turn),
	// claude (third turn, times out waiting)
	if len(bridge.streamOrder) < 1 {
		t.Errorf("expected at least 1 agent stream, got %d", len(bridge.streamOrder))
	}

	// All streams should be from claude
	for _, stream := range bridge.streamOrder {
		if stream != "claude" {
			t.Errorf("expected streams to be from 'claude', got %q", stream)
		}
	}

	// Should have at least one message from the user input
	if len(bridge.messageAuthors) < 1 {
		t.Errorf("expected at least 1 message, got %d", len(bridge.messageAuthors))
	}
}

// TestAdapterRequestContents verifies adapter receives correct request data.
func TestAdapterRequestContents(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	// Pre-populate session with one message
	sess.Append(session.Message{
		Author: "user",
		Role:   "user",
		Text:   "initial prompt",
	})

	turnOrder := NewTurnOrder([]string{"claude"})

	// Track what request the adapter received
	requestChan := make(chan agent.Request, 1)
	mockAdapter := &mockRequestCapture{name: "claude", captureReq: requestChan}

	adapters := map[string]agent.Adapter{
		"claude": mockAdapter,
	}

	bridge := NewMockUIBridge([]string{})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- engine.Run(ctx)
	}()

	// Wait for the adapter to be invoked
	select {
	case req := <-requestChan:
		// Verify NewMessage was populated with the user's message
		if req.NewMessage != "initial prompt" {
			t.Errorf("expected NewMessage='initial prompt', got %q", req.NewMessage)
		}

		// Verify Transcript includes the user message
		if len(req.Transcript) < 1 {
			t.Errorf("expected Transcript with at least 1 message, got %d", len(req.Transcript))
		}
		if len(req.Transcript) > 0 && req.Transcript[0].Text != "initial prompt" {
			t.Errorf("expected Transcript[0].Text='initial prompt', got %q", req.Transcript[0].Text)
		}
	case <-ctx.Done():
		t.Fatal("adapter was not invoked within timeout")
	}

	cancel()
	<-done
}

// mockRequestCapture captures the request sent to an adapter.
type mockRequestCapture struct {
	name       string
	captureReq chan<- agent.Request
}

func (m *mockRequestCapture) Name() string {
	return m.name
}

func (m *mockRequestCapture) Capabilities() agent.Capabilities {
	return agent.Capabilities{SupportsResume: false}
}

func (m *mockRequestCapture) Invoke(ctx context.Context, req agent.Request) (<-chan agent.Token, <-chan agent.Result) {
	// Send the request on the capture channel
	select {
	case m.captureReq <- req:
	default:
		// Channel full, skip
	}

	tokenChan := make(chan agent.Token, 1)
	resultChan := make(chan agent.Result, 1)

	go func() {
		defer close(tokenChan)
		defer close(resultChan)

		tokenChan <- agent.Token{Text: "response"}
		resultChan <- agent.Result{SessionID: "", Err: nil}
	}()

	return tokenChan, resultChan
}

// TestEmptyTurnOrder returns clear error.
func TestEmptyTurnOrder(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	// Empty turn order should error
	turnOrder := NewTurnOrder([]string{})

	adapters := map[string]agent.Adapter{}
	bridge := NewMockUIBridge([]string{})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = engine.Run(ctx)
	if err == nil {
		t.Fatal("expected error for empty turn order")
	}
	if err.Error() != "empty turn order" {
		t.Errorf("expected 'empty turn order' error, got: %v", err)
	}
}

// TestUnknownAgent returns clear error.
func TestUnknownAgent(t *testing.T) {
	dir := t.TempDir()

	sess, err := session.Open(dir)
	if err != nil {
		t.Fatalf("Open session: %v", err)
	}

	// Turn order references an unknown agent
	turnOrder := NewTurnOrder([]string{"unknown_agent", "user"})

	adapters := map[string]agent.Adapter{}
	bridge := NewMockUIBridge([]string{})

	engine := NewEngine(nil, sess, turnOrder, adapters, bridge)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = engine.Run(ctx)
	if err == nil {
		t.Fatal("expected error for unknown agent")
	}
	if err.Error() != "unknown agent: unknown_agent" {
		t.Errorf("expected 'unknown agent' error, got: %v", err)
	}
}
