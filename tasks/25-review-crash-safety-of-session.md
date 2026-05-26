# REVIEW: crash-safety of session

Goal: Confirm a mid-session interruption leaves a readable, consistent transcript.

Context:
- Verification task.

Do:
1. Run go test ./internal/session/....
2. Simulate crash: in a scratch program, Append messages then exit without clean shutdown; reopen and confirm all appended messages present.
3. Confirm agents.json is never left half-written (atomic rename).

Done when:
- Transcript intact after abrupt exit.
- No partial agents.json observed.
- Tests pass.
