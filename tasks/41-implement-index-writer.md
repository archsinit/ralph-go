# Implement index writer

Goal: Write plan.md as the index of checkbox lines referencing task files.

Context:
- outDir is the active project root (config.Paths.OutDir, default ".").

Do:
1. In internal/plan/write.go add func Write(outDir string, p *Plan) error.
2. Assign Number+Slug per task; write tasks/NN-slug.md via RenderTask.
3. Write plan.md: one line per task '- [ ] <title> → tasks/NN-slug.md'; if Agent set, prefix title with 'executor: '/'reviewer: '.
4. Create tasks/ dir under outDir; fail if files exist unless an overwrite flag passed (default: error to avoid clobber).

Done when:
- go build ./... succeeds.
- plan.md and tasks/*.md written under outDir.
- Each index line references an existing file.

Files: internal/plan/write.go
