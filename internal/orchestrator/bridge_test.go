package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/archsinit/ralph-go/internal/tui"
)

// TestBridgeGetSubmitCallback verifies the submit callback works.
func TestBridgeGetSubmitCallback(t *testing.T) {
	b := NewTUIBridge(nil)
	cb := b.GetSubmitCallback()

	if cb == nil {
		t.Fatal("GetSubmitCallback returned nil")
	}

	// Call the callback with text
	cb("test input")

	// Verify the input is available on the channel
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	input, err := b.RequestUserInput(ctx)
	if err != nil {
		t.Fatalf("RequestUserInput: %v", err)
	}
	if input != "test input" {
		t.Errorf("expected 'test input', got %q", input)
	}
}

// TestBridgeGetCancelCallback verifies the cancel callback works.
func TestBridgeGetCancelCallback(t *testing.T) {
	b := NewTUIBridge(nil)
	cb := b.GetCancelCallback()

	if cb == nil {
		t.Fatal("GetCancelCallback returned nil")
	}

	// Call the callback to signal cancel
	cb()

	// Verify the cancel signal is available
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cancelSig := b.RequestCancelSignal(ctx)
	select {
	case <-cancelSig:
		// Cancel signal received as expected
	case <-ctx.Done():
		t.Fatal("cancel signal not received")
	}
}

// TestBridgeStreamStartWithoutHandle verifies stream calls are safe when handle is nil.
func TestBridgeStreamStartWithoutHandle(t *testing.T) {
	b := NewTUIBridge(nil)

	// These should not panic even though handle is nil
	err := b.StreamStart("agent")
	if err != nil {
		t.Errorf("StreamStart returned error: %v", err)
	}

	err = b.StreamToken("test")
	if err != nil {
		t.Errorf("StreamToken returned error: %v", err)
	}

	err = b.StreamEnd()
	if err != nil {
		t.Errorf("StreamEnd returned error: %v", err)
	}
}

// TestBridgeContextCancellation verifies context cancellation is honored.
func TestBridgeContextCancellation(t *testing.T) {
	b := NewTUIBridge(nil)

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// RequestUserInput should return context.Canceled
	_, err := b.RequestUserInput(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestBridgeSetHandle verifies handle can be set and is protected by mutex.
func TestBridgeSetHandle(t *testing.T) {
	b := NewTUIBridge(nil)

	// Create a temporary TUI handle
	handle, err := tui.Run()
	if err != nil {
		t.Fatalf("tui.Run: %v", err)
	}
	defer handle.Quit()
	defer handle.Wait()

	// Set the handle
	b.SetHandle(handle)

	// Verify the handle was set by calling a method that uses it
	err = b.StreamStart("test")
	if err != nil {
		t.Errorf("StreamStart with handle: %v", err)
	}
}

// TestBridgeBufferedCancelSignal verifies cancel signals are not dropped.
func TestBridgeBufferedCancelSignal(t *testing.T) {
	b := NewTUIBridge(nil)
	cb := b.GetCancelCallback()

	// Send multiple cancel signals
	for i := 0; i < 5; i++ {
		cb()
	}

	// All signals should be received (buffered channel has capacity 10)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	for i := 0; i < 5; i++ {
		sig := b.RequestCancelSignal(ctx)
		select {
		case <-sig:
			// Signal received
		case <-ctx.Done():
			t.Fatalf("expected signal %d, got context done", i)
		}
	}
}

// TestBridgeAddMessageWithoutHandle verifies AddMessage is safe without handle.
func TestBridgeAddMessageWithoutHandle(t *testing.T) {
	b := NewTUIBridge(nil)

	// Should not panic
	err := b.AddMessage("user", "test")
	if err != nil {
		t.Errorf("AddMessage returned error: %v", err)
	}
}
