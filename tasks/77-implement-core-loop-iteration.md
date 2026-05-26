# Implement core loop iteration

Goal: Iterate remaining tasks running an executor-then-reviewer cycle each.

Context:
- Executor/reviewer prompts are directed agent turns, not a chatroom; reuse adapters.
- Reviewer output must be parseable for pass/fail + message (define a simple convention, e.g. first line PASS/FAIL then message).

Do:
1. In internal/loop/engine.go define Engine holding config, project dir, plan path, adapters (executor, reviewer), ntfy client, master+task loggers, and a UI bridge.
2. On Run: check IsRepo and IsDirty -> if dirty, ntfy DirtyTree + halt.
3. ntfy LoopStart + master event.
4. For each remaining entry: ntfy/log TaskStart; load prompt via PromptFor; run executor (headless) capturing output to task log + TUI; then run reviewer with task prompt + git Diff asking pass/fail + a commit message.
5. On reviewer pass: gitx.Commit(message); if EmptyDiff treat as fail; on success plan.Flip + ntfy TaskEnd + master event.

Done when:
- go build ./... succeeds.
- Happy path with echo agents: each task executes, reviews, commits, flips.
- Dirty tree halts before any task.

Files: internal/loop/engine.go
