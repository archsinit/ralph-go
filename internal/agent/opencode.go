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
func buildOpencodeArgs(req Request) []string {
	// TODO: verify exact non-interactive flags for current opencode CLI version
	args := []string{
		"run",
	}
	if req.SystemPrompt != "" {
		// TODO: verify system-prompt flag name for opencode CLI
		args = append(args, "--system-prompt", req.SystemPrompt)
	}
	args = append(args, req.NewMessage)
	return args
}

func (a *opencodeAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		rawTokens, rawResults := runStreaming(ctx, "opencode", buildOpencodeArgs(req), "")

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
