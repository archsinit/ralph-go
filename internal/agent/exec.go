package agent

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// runStreaming spawns name with args, optionally writes stdin, and streams
// stdout as Token values. The Token channel closes before a single Result is
// sent on the result channel. Non-zero exit wraps the last 512 bytes of stderr
// in Result.Err. ctx cancellation kills the process via CommandContext.
func runStreaming(ctx context.Context, name string, args []string, stdin string) (<-chan Token, <-chan Result) {
	tokens := make(chan Token)
	results := make(chan Result, 1)

	go func() {
		defer close(tokens)
		defer close(results)

		cmd := exec.CommandContext(ctx, name, args...)

		var stderrBuf bytes.Buffer
		cmd.Stderr = &stderrBuf

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			results <- Result{Err: fmt.Errorf("stdout pipe: %w", err)}
			return
		}

		if stdin != "" {
			cmd.Stdin = strings.NewReader(stdin)
		}

		if err := cmd.Start(); err != nil {
			results <- Result{Err: fmt.Errorf("start: %w", err)}
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
			case tokens <- Token{Text: line + "\n"}:
			}
		}

		waitErr := cmd.Wait()
		if waitErr != nil {
			tail := stderrTail(stderrBuf.String(), 512)
			results <- Result{Err: fmt.Errorf("exit: %w; stderr: %s", waitErr, tail)}
			return
		}
		results <- Result{}
	}()

	return tokens, results
}

func stderrTail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
