# End-to-end real loop run

Goal: Run loop on a small real checklist with real executor and reviewer (cheap model).

Context:
- Integration; needs real CLIs + network.

Do:
1. Generate or hand-write a tiny plan.md (2-3 atomic tasks) in a real git repo.
2. Configure a cheap model (e.g. haiku) for executor and reviewer, real ntfy topic.
3. Run ralph-go loop; confirm tasks complete, commits made, flips applied, logs and notifications correct.

Done when:
- Real loop completes the checklist.
- Commits/flips/logs/ntfy all correct.
- No manual intervention needed for the happy path.
