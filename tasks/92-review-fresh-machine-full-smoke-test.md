# REVIEW: fresh-machine full smoke test

Goal: From a clean checkout, follow the README to run both plan and loop.

Context:
- Verification task; final gate.

Do:
1. Clean clone; build per README; run plan to generate a tiny plan; run loop to completion with echo agents.
2. Confirm no undocumented steps.
3. Run go test ./... and go vet ./...; confirm all phases pass.

Done when:
- README sufficient for both modes.
- All tests/vet pass.
- No regressions.
