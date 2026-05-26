# Implement config validation

Goal: Validate semantic correctness of a loaded Config and return clear aggregated errors.

Context:
- Known CLIs set should live as a var in internal/agent or internal/config; pick one and reference it.

Do:
1. Add func (c *Config) Validate() error.
2. Checks: agent Names unique and non-empty; CLI is one of a known set (claude, codex, gemini, opencode); Plan.TurnOrder entries each resolve to an enabled agent or the literal "user"; Plan.PlanAgent resolves to an enabled agent; Loop.Executor and Loop.Reviewer resolve to enabled agents; Loop.MaxRetries >= 0; Loop.TaskTimeout > 0; Ntfy.Topic non-empty.
3. Aggregate all failures into one error (join messages), not just the first.
4. Call Validate() after Load in cmd.

Done when:
- go build ./... succeeds.
- A config violating any rule yields a clear message identifying the field.
- A valid config passes.

Files: internal/config/config.go, cmd/root.go
