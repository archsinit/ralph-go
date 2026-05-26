# TEST: /plan writes valid files with echo agent

Goal: Unit/integration test that /plan produces valid artifacts using a canned JSON response.

Context:
- May require a small hook on the echo adapter to return canned content.

Do:
1. In a test, configure an echo agent whose response is a fixed valid TaskSpec JSON (extend echo to return a preset payload, or inject via the bridge).
2. Drive the engine /plan path; assert plan.md and tasks/*.md created in a temp outDir and that Parse reads them back.
3. Test invalid JSON path: assert no files written and an error surfaced.

Done when:
- go test passes for the generation path.
- Valid -> files exist and reparse.
- Invalid -> no files.

Files: internal/orchestrator/generate_test.go, internal/agent/echo.go
