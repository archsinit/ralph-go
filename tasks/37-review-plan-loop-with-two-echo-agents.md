# REVIEW: plan loop with two echo agents

Goal: Confirm a full scripted plan session runs cleanly end-to-end.

Context:
- Verification task; integrates Phases 1–5.

Do:
1. Run ./ralph-go plan with a config using two echo agents plus user.
2. Exchange several messages, use a prefix override, cancel one turn.
3. Confirm transcript on disk matches what was shown and resume reopens it.

Done when:
- Full loop runs without hang or crash.
- Disk transcript consistent.
- Resume works.
