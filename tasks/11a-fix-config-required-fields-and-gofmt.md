# Fix config required-field validation and gofmt drift

Goal: Close the phase 1 review gaps before starting adapter work.

Context:
- Phase 1 validation currently accepts configs with missing required agent references, which would leave plan or loop mode unusable later.
- `gofmt -l .` reports `cmd/plan.go` and `cmd/loop.go`.

Do:
1. Run gofmt on `cmd/plan.go` and `cmd/loop.go`.
2. Update `Config.Validate()` so these fields are required and must resolve to enabled agents:
   - `plan.plan_agent`
   - `loop.executor`
   - `loop.reviewer`
3. Require `plan.turn_order` to contain at least one entry.
4. Add focused tests proving missing `plan.plan_agent`, empty `plan.turn_order`, missing `loop.executor`, and missing `loop.reviewer` each return validation errors.
5. Keep validation errors aggregated with the existing error style.

Do NOT:
- Do not implement agent adapters.
- Do not change the config schema unless required for these checks.

Done when:
- `gofmt -l .` produces no output.
- `go test ./...` passes.
- `go vet ./...` passes.
- `go build ./...` passes.
- A config missing required plan or loop agent fields fails validation with messages naming the offending fields.

Files: internal/config/config.go, internal/config/config_test.go, cmd/plan.go, cmd/loop.go
