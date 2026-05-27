package orchestrator

import (
	"context"
	"sync"

	"github.com/archsinit/ralph-go/internal/tui"
)

// TUIBridge implements UIBridge by bridging to a TUI Handle.
type TUIBridge struct {
	handle         *tui.Handle
	userInputCh    chan string
	cancelSignalCh chan struct{}
	mu             sync.RWMutex
}

// NewTUIBridge creates a new bridge to the TUI.
func NewTUIBridge(handle *tui.Handle) *TUIBridge {
	return &TUIBridge{
		handle:         handle,
		userInputCh:    make(chan string, 1),
		cancelSignalCh: make(chan struct{}, 10), // Buffered to avoid dropped signals
	}
}

// GetSubmitCallback returns a function to use with TUI's WithSubmitCallback option.
func (b *TUIBridge) GetSubmitCallback() func(string) {
	return func(text string) {
		select {
		case b.userInputCh <- text:
		default:
			// Channel full, skip
		}
	}
}

// GetCancelCallback returns a function to use with TUI's WithCancelCallback option.
func (b *TUIBridge) GetCancelCallback() func() {
	return func() {
		b.SendCancelSignal()
	}
}

// SetHandle sets the TUI handle after it's been created. Safe to call from any goroutine.
func (b *TUIBridge) SetHandle(handle *tui.Handle) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handle = handle
	// Inject the cancel callback into the TUI model
	if handle != nil {
		// We need to set the cancel callback on the model
		// But we don't have direct access to the model from here
		// The cancel callback was already set when we created the TUI with WithCancelCallback
	}
}

// RequestUserInput blocks until the user submits text.
func (b *TUIBridge) RequestUserInput(ctx context.Context) (string, error) {
	select {
	case input := <-b.userInputCh:
		return input, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// RequestCancelSignal returns a channel that signals when the user presses cancel.
func (b *TUIBridge) RequestCancelSignal(ctx context.Context) <-chan struct{} {
	return b.cancelSignalCh
}

// StreamStart notifies the TUI that a stream is starting.
func (b *TUIBridge) StreamStart(author string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.handle == nil {
		return nil // TUI not yet started, silently ignore
	}
	b.handle.StartStream(author)
	return nil
}

// StreamToken sends a token to the TUI.
func (b *TUIBridge) StreamToken(text string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.handle == nil {
		return nil // TUI not yet started, silently ignore
	}
	b.handle.SendToken(text)
	return nil
}

// StreamEnd notifies the TUI that a stream has ended.
func (b *TUIBridge) StreamEnd() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.handle == nil {
		return nil // TUI not yet started, silently ignore
	}
	b.handle.EndStream()
	return nil
}

// AddMessage adds a message to the TUI transcript.
func (b *TUIBridge) AddMessage(author, text string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.handle == nil {
		return nil // TUI not yet started, silently ignore
	}
	b.handle.AddMessage(author, text)
	return nil
}

// SendCancelSignal sends a cancel signal to the engine.
func (b *TUIBridge) SendCancelSignal() {
	select {
	case b.cancelSignalCh <- struct{}{}:
	default:
	}
}
