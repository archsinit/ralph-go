# Finalize loop config fields

Goal: Ensure LoopConfig covers executor, reviewer, retries, timeout, git, ntfy, and log dir.

Context:
- Most fields exist from Phase 1; this closes gaps for loop.

Do:
1. Confirm/extend config: Loop.Executor, Loop.Reviewer, Loop.MaxRetries, Loop.TaskTimeout (Duration), Loop.Git.CommitPrefix, Loop.Ntfy{Server,Topic,Token}, Paths.LogDir.
2. Add any missing TOML tags and example values to ralph.toml.

Done when:
- go build ./... succeeds.
- ralph.toml has a complete, documented loop section.
- Config loads it.

Files: internal/config/config.go, ralph.toml
