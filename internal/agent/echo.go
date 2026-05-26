package agent

import (
	"context"
)

type echoAdapter struct{}

// NewEchoAdapter returns a network-free adapter that echoes the rendered prompt
// back as deterministic text chunks. Used in tests and dry-run scenarios.
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
		msg := req.NewMessage
		if req.SystemPrompt != "" || len(req.Transcript) > 0 {
			msg = renderReplay(req)
		}

		if err := ctx.Err(); err != nil {
			close(tokens)
			results <- Result{Err: err}
			close(results)
			return
		}

		for _, chunk := range splitEchoChunks(msg) {
			select {
			case <-ctx.Done():
				close(tokens)
				results <- Result{Err: ctx.Err()}
				close(results)
				return
			case tokens <- Token{Text: chunk}:
			}
		}

		close(tokens)
		if err := ctx.Err(); err != nil {
			results <- Result{Err: err}
		} else {
			results <- Result{SessionID: "echo-session-0"}
		}
		close(results)
	}()

	return tokens, results
}

func splitEchoChunks(s string) []string {
	var chunks []string
	for start := 0; start < len(s); {
		end := start
		if isEchoSpace(s[end]) {
			for end < len(s) && isEchoSpace(s[end]) {
				end++
			}
		} else {
			for end < len(s) && !isEchoSpace(s[end]) {
				end++
			}
			for end < len(s) && isEchoSpace(s[end]) {
				end++
			}
		}
		chunks = append(chunks, s[start:end])
		start = end
	}
	return chunks
}

func isEchoSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}
