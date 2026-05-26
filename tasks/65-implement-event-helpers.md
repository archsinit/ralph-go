# Implement event helpers

Goal: Provide typed helpers for each loop event with sensible titles/priorities/tags.

Context:
- Events per design: start, task start/end, retry, failure, timeout, dirty tree, loop end.

Do:
1. Add methods: LoopStart, TaskStart(idx,total,title), TaskEnd(title,commit), Retry(title,attempt,max), Failure(title,reason), Timeout(title,dur), DirtyTree(), LoopEnd(done,total).
2. Each composes a clear title/body and tags (e.g. warning/rotating_light for failures).
3. All call Publish; log on error, never return fatal to the loop.

Done when:
- go build ./... succeeds.
- Each helper produces a distinct, readable notification.
- Errors are swallowed-with-log.

Files: internal/ntfy/ntfy.go
