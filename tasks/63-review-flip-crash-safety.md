# REVIEW: flip crash-safety

Goal: Confirm a flip interrupted mid-write cannot corrupt plan.md.

Context:
- Verification task.

Do:
1. Inspect that Flip uses temp+rename.
2. Simulate interruption around the rename in a scratch test; confirm plan.md is either old or new, never partial.

Done when:
- Flip is atomic.
- No corruption observed.
- Tests pass.
