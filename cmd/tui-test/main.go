package main

import (
	"log"
	"time"

	"github.com/archsinit/ralph-go/internal/tui"
)

func main() {
	// Create initial history to demonstrate resume
	initialMsgs := []tui.Message{
		{Author: "user", Text: "Hello, can you help me?"},
		{Author: "claude", Text: "Of course! I'd be happy to help. What do you need?"},
	}

	// Track streaming state for cancel handling
	var currentStreamCh chan struct{}

	// Declare handle first so it can be used in closures
	var handle *tui.Handle

	// Run TUI returns immediately with a handle; we can interact with it while it runs
	var err error
	handle, err = tui.Run(
		tui.WithInitialMessages(initialMsgs),
		tui.WithSubmitCallback(func(text string) {
			// Echo the user message back
			handle.AddMessage("user", text)

			// Simulate agent response after a short delay
			cancelCh := make(chan struct{})
			currentStreamCh = cancelCh
			go simulateStreamResponse(handle, cancelCh)
		}),
		tui.WithCancelCallback(func() {
			// User pressed ctrl+c; signal cancel to the stream if active
			if currentStreamCh != nil {
				select {
				case currentStreamCh <- struct{}{}:
				default:
				}
			}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// The handle is live and the TUI is running; we can send messages to it now
	// but we'll let the submit callback handle interaction.

	// Wait for the TUI to exit (user quits or closes terminal)
	if err := handle.Wait(); err != nil {
		log.Printf("TUI error: %v", err)
	}
}

func simulateStreamResponse(handle *tui.Handle, cancelCh chan struct{}) {
	// Small delay before starting stream
	time.Sleep(100 * time.Millisecond)

	handle.StartStream("echo")
	responses := []string{
		"You said: ",
		"\"",
		"hello",
		"world",
		"\"",
		"\n",
		"That's great!",
	}

	for _, token := range responses {
		select {
		case <-cancelCh:
			// Stream was canceled
			handle.EndStream()
			// Still add the partial response to show what was streamed
			handle.AddMessage("echo", "You said: \"helloworld\"\n(cancelled)")
			return
		default:
		}

		handle.SendToken(token)
		time.Sleep(50 * time.Millisecond)
	}

	handle.EndStream()
	// Persist the full response
	handle.AddMessage("echo", "You said: \"helloworld\"\nThat's great!")
}
