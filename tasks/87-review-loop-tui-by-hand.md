# REVIEW: loop TUI by hand

Goal: Manually watch a long task and confirm live state and graceful stop.

Context:
- Verification task.

Do:
1. Run ralph-go loop with a slow echo task; watch the timer, retries, and output stream.
2. Trigger graceful stop; confirm the status shows the reason and state stays consistent.

Done when:
- Live state updates correctly.
- Stop is graceful and clearly indicated.
- Tests pass.
