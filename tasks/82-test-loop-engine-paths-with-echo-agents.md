# TEST: loop engine paths with echo agents

Goal: Unit/integration test the engine's pass, fail/retry, timeout, and resume paths.

Context:
- Inject a fake ntfy and fake bridge for assertions.

Do:
1. Create internal/loop/engine_test.go with a fake UI bridge, a temp git repo, a fixture plan.md+tasks, and echo adapters returning scripted outputs.
2. Test pass path: task executes, reviewer returns PASS+message, checkbox flips, a commit exists.
3. Test fail path: reviewer FAIL repeated -> retries up to max then halt; no flip/commit.
4. Test timeout: a deliberately slow echo executor exceeds a tiny timeout -> killed + halt + timeout event recorded (use a fake ntfy capturing events).
5. Test resume: pre-check one task in the fixture -> engine skips it.
6. Use generous test-level timeouts to avoid hangs.

Done when:
- go test ./internal/loop/... passes.
- Pass, fail/retry, timeout, resume all asserted.
- No deadlocks.

Files: internal/loop/engine_test.go
