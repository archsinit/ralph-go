# REVIEW: commits land correctly

Goal: Confirm one commit per task with sensible messages in a real repo.

Context:
- Verification task.

Do:
1. In a scratch repo, drive Commit a few times with multi-line messages.
2. Confirm git log shows one commit each with full messages.
3. Confirm dirty-at-start detection halts as intended (manual).

Done when:
- One commit per call.
- Messages intact and multi-line.
- Dirty detection correct.
- Tests pass.
