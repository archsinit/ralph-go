# REVIEW: logs readable under streaming

Goal: Confirm logs are coherent, not truncated or interleaved, during real streaming.

Context:
- Verification task.

Do:
1. Drive a fake task with streamed tokens through the tee into a task log.
2. Confirm no truncation, no interleaving within a single stream, sections clearly delimited.
3. Confirm master.log captures the event sequence.

Done when:
- Logs readable and complete.
- No garbled interleaving.
- Tests pass.
