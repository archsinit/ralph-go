# Align /plan generation with the plan package schema and safe artifact writes

Goal: Make the `/plan` command generate, validate, and write the same artifact format that the shared `internal/plan` package and later loop phases expect.

Context:
- The audit found a simplified `internal/plan` implementation using `{id,title,description}` tasks, while the Phase 5.6 task files specify a richer `TaskSpec` schema with title, optional agent, goal, context, do, do_not, done_when, and files.
- The current `/plan` prompt embeds the simplified schema and does not use the final plan package rules.
- There is no canned echo/fake path that can return valid plan JSON for automated `/plan` tests.
- Decode failures avoid writes, but write failures can still leave partially updated artifacts.

Do:
1. Complete or align with the Phase 5.6 `internal/plan` API before relying on `/plan`: strict decoder, slug/number helpers, task renderer, index writer, parser/loader, and checkbox flip compatibility.
2. Update the `/plan` prompt to embed the exact final `TaskSpec` schema and output rules: JSON only, no fences unless the decoder intentionally supports stripping them, required fields, useful atomic task prompts, and optional agent assignment.
3. Include the full prior transcript as context without duplicating `NewMessage`; include plan-agent system prompt and resume ID according to the request contract from task 49b.
4. Decode and validate the full plan response before writing any artifacts.
5. Make artifact writes safe enough for failures: write to a staging/temp location or otherwise ensure a failed write cannot leave a half-updated `plan.md`/`tasks` set in the target outDir.
6. Surface decode/write errors clearly in the TUI and return to user input without corrupting the session.
7. Add a test hook or fake adapter that can return canned valid/invalid JSON without requiring real CLIs or network access.

Do NOT:
- Do not keep two incompatible plan schemas alive.
- Do not require a real terminal or real agent CLI for automated `/plan` generation tests.
- Do not overwrite unrelated files in outDir.

Done when:
- `/plan` writes `plan.md` and `tasks/*.md` matching the Phase 5.6 format, and the plan parser can read them back.
- Invalid JSON/schema produces a visible TUI error and writes no artifacts.
- Write failures leave no partial generated plan set behind.
- The generated files are suitable inputs for the later loop parser tasks.
- `go test ./internal/plan/...`, `go test ./internal/orchestrator/...`, `go test ./...`, `go vet ./...`, `go build ./...`, and `gofmt -l .` pass.

Files: internal/plan/spec.go, internal/plan/slug.go, internal/plan/render.go, internal/plan/write.go, internal/plan/read.go, internal/plan/plan_test.go, internal/orchestrator/engine.go, internal/orchestrator/generate_test.go, internal/agent/echo.go
