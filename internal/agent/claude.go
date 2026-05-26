package agent

import (
	"context"
	"encoding/json"
	"strings"
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
// --include-partial-messages enables true partial-message streaming.
// If ResumeSessionID is empty but Transcript is present, use transcript replay.
func buildClaudeArgs(req Request) []string {
	args := []string{
		"--print",
		"--output-format", "stream-json",
		"--verbose",
		"--include-partial-messages",
	}
	if req.SystemPrompt != "" {
		args = append(args, "--system-prompt", req.SystemPrompt)
	}
	if req.ResumeSessionID != "" {
		args = append(args, "--resume", req.ResumeSessionID)
	}
	msg := req.NewMessage
	if len(req.Transcript) > 0 && req.ResumeSessionID == "" {
		msg = renderReplay(req)
	}
	args = append(args, msg)
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
		rawTokens, rawResults := streamCommand(ctx, "claude", buildClaudeArgs(req), "")

		var sessionID string
		var emittedText string
		var cancelled bool
		for t := range rawTokens {
			var ev claudeEvent
			if err := json.Unmarshal([]byte(t.Text), &ev); err != nil {
				// Non-JSON line — emit as-is
				if !sendAdapterToken(ctx, rawTokens, tokens, t) {
					cancelled = true
					break
				}
				continue
			}

			// Capture session ID from any event that carries it
			if ev.SessionID != "" {
				sessionID = ev.SessionID
			}

			// Emit text from assistant message content blocks
			if ev.Type == "assistant" && ev.Message != nil {
				text := claudeText(ev.Message)
				if text == "" {
					continue
				}
				delta := claudeDelta(emittedText, text)
				if delta == "" {
					continue
				}
				if !sendAdapterToken(ctx, rawTokens, tokens, Token{Text: delta}) {
					cancelled = true
					break
				}
				emittedText += delta
			}

			if cancelled {
				break
			}
		}

		close(tokens)

		var r Result
		if ctxErr := ctx.Err(); ctxErr != nil {
			r = Result{Err: ctxErr}
		} else if !cancelled {
			r = <-rawResults
		} else {
			r = Result{Err: context.Canceled}
		}
		r.SessionID = sessionID
		results <- r
		close(results)
	}()

	return tokens, results
}

func drain(ch <-chan Token) {
	for range ch {
	}
}

func sendAdapterToken(ctx context.Context, rawTokens <-chan Token, tokens chan<- Token, token Token) bool {
	select {
	case <-ctx.Done():
		drain(rawTokens)
		return false
	default:
	}

	select {
	case <-ctx.Done():
		drain(rawTokens)
		return false
	case tokens <- token:
		return true
	}
}

func claudeText(msg *claudeMessage) string {
	var b strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			b.WriteString(block.Text)
		}
	}
	return b.String()
}

func claudeDelta(emitted, current string) string {
	switch {
	case strings.HasPrefix(current, emitted):
		return current[len(emitted):]
	case strings.HasPrefix(emitted, current):
		return ""
	default:
		return current
	}
}
