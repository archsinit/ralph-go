package agent

import (
	"context"
)

type geminiAdapter struct{}

// NewGeminiAdapter returns an Adapter backed by the gemini CLI.
func NewGeminiAdapter() Adapter {
	return &geminiAdapter{}
}

func (a *geminiAdapter) Name() string { return "gemini" }

func (a *geminiAdapter) Capabilities() Capabilities {
	return Capabilities{SupportsResume: false}
}

// buildGeminiArgs constructs CLI args for a non-interactive gemini invocation.
// Gemini CLI does not expose a stable resume flag; transcript is included via replay.
// TODO: verify exact non-interactive flags for current gemini CLI version.
func buildGeminiArgs(req Request) []string {
	args := []string{
		"--no-interactive",
	}
	if req.SystemPrompt != "" {
		// TODO: verify system-prompt flag name for gemini CLI
		args = append(args, "--system", req.SystemPrompt)
	}
	msg := req.NewMessage
	// gemini has no resume; include transcript context via replay rendering
	if len(req.Transcript) > 0 {
		msg = renderReplay(req)
	}
	args = append(args, msg)
	return args
}

func (a *geminiAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		rawTokens, rawResults := streamCommand(ctx, "gemini", buildGeminiArgs(req), "")

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
