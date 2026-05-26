# REVIEW: generated artifacts are clean

Goal: Hand-inspect a generated plan.md and a task file for prompt quality.

Context:
- Verification task; this format is the contract between plan and loop.

Do:
1. Generate a small plan via the writer (or a test fixture) into a temp dir.
2. Open plan.md and one tasks/NN-slug.md; confirm flush-left formatting, correct refs, readable prompt, agent prefix where expected.
3. Confirm a checkbox flip leaves a minimal, sensible diff.

Done when:
- Artifacts read cleanly and match the contract.
- Flip diff minimal.
- Tests pass.
