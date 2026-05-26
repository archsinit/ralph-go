# Implement task file renderer

Goal: Render a TaskSpec into the flush-left markdown prompt file.

Context:
- This text is sent verbatim to the executor as the prompt.

Do:
1. In internal/plan/render.go add func RenderTask(t TaskSpec) string producing sections: '# <title>', 'Goal: ...', 'Context:' bullets (omit if empty), 'Do:' numbered list, 'Do NOT:' bullets (omit if empty), 'Done when:' bullets, 'Files: a, b' (omit if empty).
2. No leading indentation anywhere.
3. Stable, deterministic output (used in byte-equality test).

Done when:
- go build ./... succeeds.
- Empty optional sections omitted.
- Output flush-left and deterministic.

Files: internal/plan/render.go
