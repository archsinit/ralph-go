# REVIEW: real /plan generation quality

Goal: Run a real session, generate a plan, and eyeball it.

Context:
- Verification task; needs a real CLI.

Do:
1. With a real plan agent, hold a short design chat, run /plan.
2. Inspect generated plan.md + tasks for schema conformance and useful, atomic prompts.
3. Confirm /quit saved the session and the plan files are in the project root.

Done when:
- Real generation produces valid, readable artifacts.
- Files land in outDir.
- Session saved.
