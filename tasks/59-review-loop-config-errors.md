# REVIEW: loop config errors

Goal: Confirm loop misconfig produces actionable messages.

Context:
- Verification task.

Do:
1. Run the tests; manually run ralph-go loop with a bad timeout and empty ntfy topic.
2. Confirm clear stderr and nonzero exit.

Done when:
- Clear messages, no panic.
- Tests pass.
