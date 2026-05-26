# TEST: loop queue parse, load, flip

Goal: Unit-test queue loading and flip behavior with fixtures.

Do:
1. Create internal/loop/queue_test.go using t.TempDir() and a fixture plan.md + tasks/.
2. Test: LoadQueue returns expected order and remaining excludes checked.
3. Test: PromptFor returns the task file contents.
4. Test: flipping the first remaining entry updates plan.md and Remaining shrinks.
5. Test: missing task file -> error.

Done when:
- go test ./internal/loop/... passes (queue parts).
- Order/remaining/flip asserted.
- Missing-file error asserted.

Files: internal/loop/queue_test.go
