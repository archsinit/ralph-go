package agent

import (
	"context"
)

type codexAdapter struct{}

// NewCodexAdapter returns an Adapter backed by the codex CLI.
func NewCodexAdapter() Adapter {
	return &codexAdapter{}
}

func (a *codexAdapter) Name() string { return "codex" }

func (a *codexAdapter) Capabilities() Capabilities {
	// codex CLI does not expose a stable resume flag; fall back to transcript replay.
	return Capabilities{SupportsResume: false}
}

// buildCodexArgs constructs CLI args for a non-interactive codex invocation.
// Kept pure for unit testing.
func buildCodexArgs(req Request) []string {
	// TODO: verify exact non-interactive flag for current codex CLI version
	args := []string{
		"--quiet",
	}
	// codex has no resume; include transcript context via replay rendering
	msg := req.NewMessage
	if len(req.Transcript) > 0 {
		msg = renderReplay(req)
	} else if req.SystemPrompt != "" {
		// TODO: verify system-prompt flag name for codex CLI
		args = append(args, "--system-prompt", req.SystemPrompt)
	}
	args = append(args, msg)
	return args
}

func (a *codexAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		rawTokens, rawResults := runStreaming(ctx, "codex", buildCodexArgs(req), "")

		for t := range rawTokens {
			select {
			case <-ctx.Done():
				drain(rawTokens)
				return
			case tokens <- t:
			}
		}

		r := <-rawResults
		// codex CLI does not expose a session ID in its output
		results <- r
	}()

	return tokens, results
}
