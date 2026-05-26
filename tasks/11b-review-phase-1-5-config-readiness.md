# REVIEW: phase 1.5 config readiness

Goal: Verify the phase 1 config package is ready for phase 2 adapter work.

Context:
- This is a verification task. Fix only issues found during the review.
- The review should specifically cover the phase 1.5 validation and formatting follow-up.

Do:
1. Run `gofmt -l .` and confirm it produces no output.
2. Run `go test ./...`, `go vet ./...`, and `go build ./...`.
3. Build the CLI and confirm `ralph-go --help` lists `plan` and `loop`.
4. Confirm `ralph-go plan --config ./ralph.toml` reaches the placeholder output with the example config.
5. Confirm missing, malformed, and semantically invalid configs exit nonzero with clear stderr and no panic.
6. Confirm validation aggregates multiple errors instead of stopping at the first one.

Done when:
- All checks pass.
- Missing required config fields are rejected.
- Error messages name either the config path or the invalid field.
- No phase 1.5 blockers remain before Phase 2.
