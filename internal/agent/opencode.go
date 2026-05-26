package agent

import (
	"context"
)

type opencodeAdapter struct{}

// NewOpencodeAdapter returns an Adapter backed by the opencode CLI.
func NewOpencodeAdapter() Adapter {
	return &opencodeAdapter{}
}

func (a *opencodeAdapter) Name() string { return "opencode" }

func (a *opencodeAdapter) Capabilities() Capabilities {
	return Capabilities{SupportsResume: false}
}

// buildOpencodeArgs constructs CLI args for a non-interactive opencode invocation.
// Opencode CLI does not expose a stable resume flag; transcript is included via replay.
// TODO: verify exact non-interactive flags for current opencode CLI version.
func buildOpencodeArgs(req Request) []string {
	args := []string{
		"run",
	}
	if req.SystemPrompt != "" {
		// TODO: verify system-prompt flag name for opencode CLI
		args = append(args, "--system-prompt", req.SystemPrompt)
	}
	msg := req.NewMessage
	// opencode has no resume; include transcript context via replay rendering
	if len(req.Transcript) > 0 {
		msg = renderReplay(req)
	}
	args = append(args, msg)
	return args
}

func (a *opencodeAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		rawTokens, rawResults := streamCommand(ctx, "opencode", buildOpencodeArgs(req), "")

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
		results <- r
		close(results)
	}()

	return tokens, results
}
