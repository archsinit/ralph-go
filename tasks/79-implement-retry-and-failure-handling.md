# Implement retry and failure handling

Goal: Retry the executor on reviewer-fail up to max_retries, then alert and halt.

Context:
- Feedback loop: reviewer feedback informs the next executor attempt.

Do:
1. On reviewer FAIL (or EmptyDiff): ntfy Retry, log, re-run executor with the reviewer's feedback appended to the prompt; increment attempt.
2. After config.Loop.MaxRetries exhausted: ntfy Failure + master event, halt the loop (do not flip, do not commit).
3. Make the halt reason available to the TUI/status.

Done when:
- go build ./... succeeds.
- Fail path retries up to the limit then halts with an alert.
- No flip/commit on ultimate failure.

Files: internal/loop/engine.go
