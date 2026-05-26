# Implement per-task timeout

Goal: Enforce a hard time limit per task, killing the agent and alerting on exceed.

Context:
- Per-task timeout only; no global iteration cap.

Do:
1. Wrap each task (executor+reviewer) in a context with deadline config.Loop.TaskTimeout.
2. On deadline: cancel the adapter context (kills process), ntfy Timeout + master event, halt the loop.
3. Ensure partial output already captured is flushed to the task log.

Done when:
- go build ./... succeeds.
- A slow echo agent exceeding the timeout is killed.
- Timeout triggers ntfy + halt, log retains partial output.

Files: internal/loop/engine.go
