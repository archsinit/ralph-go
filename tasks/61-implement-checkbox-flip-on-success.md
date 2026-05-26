# Implement checkbox flip on success

Goal: Flip a task to done via the plan writer after a passing review.

Context:
- Flip implemented in Phase 5.5; here it is invoked at the right moment.

Do:
1. In the loop engine, after a passing review and successful commit, call plan.Flip(planPath, entry).
2. Reload or update in-memory entries so the TUI reflects the change.

Done when:
- go build ./... succeeds.
- Flip writes back atomically, minimal diff.
- TUI/state shows the task checked.

Files: internal/loop/engine.go
