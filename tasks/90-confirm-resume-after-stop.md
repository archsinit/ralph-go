# Confirm resume after stop

Goal: Verify that re-running loop after a stop continues from the first unchecked task.

Do:
1. Stop a loop mid-checklist, then re-run ralph-go loop on the same plan.md.
2. Confirm checked tasks are skipped and execution resumes at the next unchecked task with no double commits.

Done when:
- Resume skips done tasks.
- No duplicate commits.
- State consistent.
