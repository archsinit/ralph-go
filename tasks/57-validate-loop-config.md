# Validate loop config

Goal: Add loop-specific validation rules.

Context:
- Reuse the aggregation approach from Phase 1.

Do:
1. Extend Config.Validate: TaskTimeout > 0; MaxRetries >= 0; Executor and Reviewer resolve to enabled agents and are not the same unless explicitly allowed; Ntfy.Topic non-empty; Ntfy.Server is a valid URL.

Done when:
- go build ./... succeeds.
- Each invalid loop setting yields a clear error.
- Valid loop config passes.

Files: internal/config/config.go
