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
	// codex exec supports resume via `codex exec resume <session_id>`
	return Capabilities{SupportsResume: true}
}

// buildCodexArgs constructs CLI args for a non-interactive codex invocation.
// Uses `codex exec` subcommand and `codex exec resume` for session resumption.
// Transcript replay is used when transcript is present without a ResumeSessionID.
func buildCodexArgs(req Request) []string {
	var args []string

	if req.ResumeSessionID != "" {
		// Resume an existing session
		args = []string{"exec", "resume", req.ResumeSessionID}
	} else {
		// Start a new session
		args = []string{"exec"}
		// codex exec does not expose --system-prompt; replay carries context.
		if req.SystemPrompt != "" || len(req.Transcript) > 0 {
			args = append(args, renderReplay(req))
		} else {
			args = append(args, req.NewMessage)
		}
		return args
	}

	// For resume, the prompt is provided as additional argument
	args = append(args, req.NewMessage)
	return args
}

func (a *codexAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		rawTokens, rawResults := streamCommand(ctx, "codex", buildCodexArgs(req), "")

		var cancelled bool
		for t := range rawTokens {
			if !sendAdapterToken(ctx, rawTokens, tokens, t) {
				cancelled = true
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
		// codex CLI does not expose a session ID in its output stream
		results <- r
		close(results)
	}()

	return tokens, results
}
