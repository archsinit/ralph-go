package agent

import (
	"context"
	"encoding/json"
)

type claudeAdapter struct{}

// NewClaudeAdapter returns an Adapter backed by the claude CLI.
func NewClaudeAdapter() Adapter {
	return &claudeAdapter{}
}

func (a *claudeAdapter) Name() string { return "claude" }

func (a *claudeAdapter) Capabilities() Capabilities {
	return Capabilities{SupportsResume: true}
}

// buildClaudeArgs constructs CLI args for a non-interactive, streaming invocation.
// Kept as a pure function so it can be tested without spawning a process.
//
// --verbose is required for --output-format=stream-json to emit events.
func buildClaudeArgs(req Request) []string {
	args := []string{
		"--print",
		"--output-format", "stream-json",
		"--verbose",
	}
	if req.SystemPrompt != "" {
		args = append(args, "--system-prompt", req.SystemPrompt)
	}
	if req.ResumeSessionID != "" {
		args = append(args, "--resume", req.ResumeSessionID)
	}
	args = append(args, req.NewMessage)
	return args
}

// claudeEvent covers the stream-json event shapes emitted by the claude CLI.
// Every event carries an optional session_id.
type claudeEvent struct {
	Type      string         `json:"type"`
	Subtype   string         `json:"subtype,omitempty"`
	SessionID string         `json:"session_id,omitempty"`
	Message   *claudeMessage `json:"message,omitempty"`
}

type claudeMessage struct {
	Content []claudeContent `json:"content"`
}

type claudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (a *claudeAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		rawTokens, rawResults := runStreaming(ctx, "claude", buildClaudeArgs(req), "")

		var sessionID string
		for t := range rawTokens {
			var ev claudeEvent
			if err := json.Unmarshal([]byte(t.Text), &ev); err != nil {
				// Non-JSON line — emit as-is
				select {
				case <-ctx.Done():
					drain(rawTokens)
					return
				case tokens <- t:
				}
				continue
			}

			// Capture session ID from any event that carries it
			if ev.SessionID != "" {
				sessionID = ev.SessionID
			}

			// Emit text from assistant message content blocks
			if ev.Type == "assistant" && ev.Message != nil {
				for _, block := range ev.Message.Content {
					if block.Type == "text" && block.Text != "" {
						select {
						case <-ctx.Done():
							drain(rawTokens)
							return
						case tokens <- Token{Text: block.Text}:
						}
					}
				}
			}
		}

		r := <-rawResults
		r.SessionID = sessionID
		results <- r
	}()

	return tokens, results
}

func drain(ch <-chan Token) {
	for range ch {
	}
}
