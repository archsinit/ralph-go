# REVIEW: config error UX

Goal: Confirm bad configs fail gracefully with actionable messages, never a panic.

Context:
- Verification task; fix only issues found.

Do:
1. Run go test ./internal/config/....
2. Manually run ./ralph-go plan --config /nonexistent.toml and a deliberately broken toml; confirm clear stderr message and nonzero exit, no stack trace.
3. Confirm Validate aggregates multiple errors.

Done when:
- No panics on bad input.
- Messages name the offending field or path.
- Tests pass.
