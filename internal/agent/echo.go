package agent

import (
	"context"
	"strings"
)

type echoAdapter struct{}

// NewEchoAdapter returns a network-free adapter that echoes NewMessage back as
// a stream of word tokens. Used in tests and dry-run scenarios.
func NewEchoAdapter() Adapter {
	return &echoAdapter{}
}

func (a *echoAdapter) Name() string { return "echo" }

func (a *echoAdapter) Capabilities() Capabilities {
	return Capabilities{SupportsResume: false}
}

func (a *echoAdapter) Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		words := strings.Fields(req.NewMessage)
		for _, w := range words {
			select {
			case <-ctx.Done():
				results <- Result{Err: ctx.Err()}
				return
			case tokens <- Token{Text: w + " "}:
			}
		}

		results <- Result{SessionID: "echo-session-0"}
	}()

	return tokens, results
}
