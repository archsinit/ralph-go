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
		var result Result

		cmd := exec.CommandContext(ctx, name, args...)

		var stderrBuf bytes.Buffer
		cmd.Stderr = &stderrBuf

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			result = Result{Err: fmt.Errorf("stdout pipe: %w", err)}
			close(tokens)
			results <- result
			close(results)
			return
		}

		if stdin != "" {
			cmd.Stdin = strings.NewReader(stdin)
		}

		if err := cmd.Start(); err != nil {
			result = Result{Err: fmt.Errorf("start: %w", err)}
			close(tokens)
			results <- result
			close(results)
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)

		cancelled := false
	scanLoop:
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				cancelled = true
				break scanLoop
			case tokens <- Token{Text: line + "\n"}:
			}
		}

		scanErr := scanner.Err()
		if scanErr != nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
		}

		close(tokens)

		waitErr := cmd.Wait()
		if ctxErr := ctx.Err(); ctxErr != nil {
			result = Result{Err: ctxErr}
		} else if cancelled {
			result = Result{Err: context.Canceled}
		} else if scanErr != nil && waitErr != nil {
			tail := stderrTail(stderrBuf.String(), 512)
			result = Result{Err: fmt.Errorf("scanner: %w; exit: %v; stderr: %s", scanErr, waitErr, tail)}
		} else if scanErr != nil {
			result = Result{Err: fmt.Errorf("scanner: %w", scanErr)}
		} else if waitErr != nil {
			tail := stderrTail(stderrBuf.String(), 512)
			result = Result{Err: fmt.Errorf("exit: %w; stderr: %s", waitErr, tail)}
		}

		results <- result
		close(results)
	}()

	return tokens, results
}

func stderrTail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

var streamCommand = runStreaming
