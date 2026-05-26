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
func buildGeminiArgs(req Request) []string {
	// TODO: verify exact non-interactive flags for current gemini CLI version
	args := []string{
		"--no-interactive",
	}
	if req.SystemPrompt != "" {
		// TODO: verify system-prompt flag name for gemini CLI
		args = append(args, "--system", req.SystemPrompt)
	}
	args = append(args, req.NewMessage)
	return args
}

func (a *geminiAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		rawTokens, rawResults := runStreaming(ctx, "gemini", buildGeminiArgs(req), "")

		for t := range rawTokens {
			select {
			case <-ctx.Done():
				drain(rawTokens)
				return
			case tokens <- t:
			}
		}

		results <- <-rawResults
	}()

	return tokens, results
}
