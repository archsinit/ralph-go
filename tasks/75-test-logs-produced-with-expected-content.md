# TEST: logs produced with expected content

Goal: Unit-test master and per-task logging including tee.

Do:
1. Create internal/logx/logx_test.go using t.TempDir().
2. Test Master.Event writes timestamped lines that persist.
3. Test TaskLog captures written sections; Tee duplicates a byte stream to two buffers identically.

Done when:
- go test ./internal/logx/... passes.
- Content and tee duplication asserted.

Files: internal/logx/logx_test.go
