# Implement graceful stop

Goal: Allow the user to stop the loop; finish or abort the current task cleanly and persist.

Context:
- Consistency invariant: flip and commit happen together or not at all.

Do:
1. Handle a stop signal (ctrl-c / a TUI key): after the current task settles (or aborts safely without a partial commit), stop the loop, write a master event, exit cleanly.
2. Never leave a half-applied commit or a flipped-but-uncommitted task.

Done when:
- go build ./... succeeds.
- Stop ends the loop without corrupting plan.md or git.
- State consistent: each task is either fully done (flipped+committed) or untouched.

Files: internal/loop/engine.go
