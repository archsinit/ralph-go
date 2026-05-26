# REVIEW: phase 4.5 readiness for orchestrator

Goal: Verify the Phase 3.5/4 follow-ups are complete before Phase 5 starts wiring the plan orchestrator.

Context:
- This review covers the gaps found after inspecting the completed session hardening and TUI implementation.

Do:
1. Run:
   - `go test ./internal/session/...`
   - `go test ./internal/tui/...`
   - `go test ./...`
   - `go vet ./...`
   - `go build ./...`
2. Inspect `internal/session/session.go` and confirm:
   - resumed `nextSeq` follows the documented last-valid-vs-max sequence contract,
   - `Append` commits memory only after all required persistence operations succeed,
   - `SetAgentSession` cannot leave memory ahead of disk on rewrite failure,
   - transcript loading tolerates only trailing corrupt JSONL content,
   - large valid messages are handled according to the documented contract.
3. Inspect `internal/tui` from the perspective of another package and confirm external orchestration code can, using exported APIs only:
   - run the TUI and retain a send/control handle,
   - receive user submissions,
   - render resumed/finalized messages,
   - feed stream start/token/end events.
4. Confirm submit ownership is clear: either the TUI appends submitted user messages immediately or the orchestrator echoes persisted messages back, but not both.
5. Run the temporary TUI harness and manually verify typing, submit, scrolling, autoscroll suppression/resume, resize, fake live streaming, status text, and spinner animation.
6. Confirm TUI tests include meaningful assertions for submit callbacks/channels, stream accumulation, viewport scrolling/autoscroll, spinner ticking, stable author styling, wrapping, and public API usability.
7. Document any remaining limitations here or create another follow-up before Phase 5.

Done when:
- Tests, vet, and build pass.
- The reviewer is satisfied Phase 5 can wire the orchestrator without package-private TUI access or session consistency surprises.
- No known Phase 3.5/4 gap remains untracked.
