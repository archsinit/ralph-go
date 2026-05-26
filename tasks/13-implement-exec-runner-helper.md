# Implement exec runner helper

Goal: Create a reusable helper that spawns a CLI process and streams stdout as tokens.

Context:
- Used by every CLI adapter.
- Keep parsing of CLI-specific session IDs out of here; adapters post-process.

Do:
1. In internal/agent/exec.go add func runStreaming(ctx, name string, args []string, stdin string) (<-chan Token, <-chan Result).
2. Start exec.CommandContext; write stdin if non-empty then close.
3. Scan stdout incrementally (bufio scanner with a word/rune split or line split) emitting Token per chunk.
4. Collect stderr; on nonzero exit return Result.Err wrapping stderr tail.
5. Respect ctx cancellation (process killed via CommandContext).
6. Close Token channel before sending Result.

Done when:
- go build ./... succeeds.
- ctx cancellation kills the process.
- Nonzero exit surfaces stderr in Result.Err.

Files: internal/agent/exec.go
