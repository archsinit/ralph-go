# REVIEW: fresh-machine plan smoke test

Goal: Follow the README only on a clean checkout to validate plan mode.

Context:
- Verification task.

Do:
1. From a clean clone, follow README to build and run plan with echo agents.
2. Confirm no undocumented step is required.
3. Re-run go test ./... and confirm Phases 1–6 still pass.

Done when:
- README sufficient to build and run.
- All tests pass.
- No regressions in earlier phases.
