# REVIEW: phase 2.5 adapter readiness

Goal: Verify adapter behavior is ready for Phase 3 session persistence and later orchestration work.

Context:
- This is a verification task. Fix only issues found during the review.
- Focus on behavior that future orchestrators will rely on: channel ordering, terminal results, cancellation, replay fallback, and real CLI argument validity.

Do:
1. Run `gofmt -l .` and confirm it produces no output.
2. Run `go test ./...`, `go vet ./...`, and `go build ./...`.
3. Confirm every Adapter implementation sends exactly one terminal Result after the Token channel has closed.
4. Confirm cancellation produces a non-nil Result.Err for runStreaming, claude, codex, gemini, opencode, and echo paths.
5. Confirm transcript replay is used whenever native resume is unavailable or ResumeSessionID is empty and prior transcript context exists.
6. Confirm codex args match `codex exec --help` and do not include unsupported flags.
7. Confirm claude stream-json output is parsed without dropping streamed text, duplicating final content, or losing session id when present.
8. If authenticated real CLI smoke tests are unavailable, document exactly which checks were skipped and why.

Done when:
- All automated checks pass.
- Echo adapter can be used as a reliable network-free test adapter.
- Known CLI validation messages match accepted CLI names.
- No Phase 2 adapter blockers remain before Phase 3.
