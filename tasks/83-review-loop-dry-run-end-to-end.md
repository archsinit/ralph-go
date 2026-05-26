# REVIEW: loop dry-run end-to-end

Goal: Run a full echo-agent loop and inspect commits, flips, logs, and notifications.

Context:
- Verification task; integrates Phases 8–13.

Do:
1. Prepare a tiny plan.md+tasks in a temp git repo; configure echo executor/reviewer and a real test ntfy topic.
2. Run ralph-go loop; watch tasks flip, commits appear, logs fill, ntfy events arrive.
3. Verify state consistency after a mid-run stop.

Done when:
- End-to-end echo loop behaves per design.
- Commits/flips/logs/ntfy all correct.
- Stop leaves consistent state.
