# REVIEW: regression sweep across all phases

Goal: Final regression pass ensuring earlier phases still work after all changes.

Context:
- Verification task; closes the build.

Do:
1. Run the full test suite and vet.
2. Spot-check plan mode and loop mode each still run end-to-end with echo agents.
3. Confirm config, session, plan-format, adapters, ntfy, git, logging packages all green.

Done when:
- Entire suite green.
- Both modes run.
- No known regressions remain.
