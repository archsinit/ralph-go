package agent

import "context"

// Capabilities describes optional features an adapter supports.
type Capabilities struct {
	SupportsResume bool
}

// Token is a single streamed chunk of text from the agent.
type Token struct {
	Text string
}

// Result is delivered once when a turn completes. Err is non-nil on failure.
type Result struct {
	SessionID string
	Err       error
}

// Message is one entry in a conversation transcript.
type Message struct {
	Author string
	Role   string
	Text   string
}

// Request is the input to a single agent turn.
type Request struct {
	SystemPrompt    string
	Transcript      []Message
	NewMessage      string
	ResumeSessionID string
}

// Adapter is the contract every agent backend must satisfy.
//
// Streaming contract: Invoke returns two channels. The Token channel receives
// zero or more chunks and is closed when the agent is done producing output.
// After the Token channel closes, exactly one Result is sent on the Result
// channel and that channel is then closed. Callers must drain both channels.
type Adapter interface {
	Name() string
	Capabilities() Capabilities
	Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result)
}
