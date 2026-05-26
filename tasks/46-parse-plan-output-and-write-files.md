# Parse plan output and write files

Goal: Validate the agent JSON and write plan.md + tasks via internal/plan.

Context:
- Atomic-ish: only call Write after Decode succeeds.

Do:
1. Pass the captured response to plan.Decode.
2. On success call plan.Write into config.Paths.OutDir (default '.').
3. On decode/validation failure show a clear in-TUI error and do not write any files (no partial writes).
4. Report the written file paths in the TUI on success.

Done when:
- go build ./... succeeds.
- Valid output writes plan.md + tasks/*.md to outDir.
- Invalid output shows an error and writes nothing.

Files: internal/orchestrator/engine.go
