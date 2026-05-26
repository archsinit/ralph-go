# Define config structs

Goal: Define Go structs mapping the TOML config for both plan and loop modes.

Context:
- TOML lib: github.com/BurntSushi/toml.
- TaskTimeout is a human string like "15m"; decide representation (recommend a custom Duration type that implements UnmarshalText).

Do:
1. In internal/config/config.go define Config with: Agents []Agent, Plan PlanConfig, Loop LoopConfig, Paths Paths.
2. Agent: Name, CLI, SystemPrompt, Enabled bool.
3. PlanConfig: TurnOrder []string (names plus the literal "user"), PlanAgent string (which agent answers /plan).
4. LoopConfig: Executor, Reviewer string, MaxRetries int, TaskTimeout duration (string in TOML, parsed), Git GitConfig, Ntfy NtfyConfig.
5. GitConfig: nothing required yet besides a CommitPrefix string (optional).
6. NtfyConfig: Server, Topic, Token.
7. Paths: SessionDir, LogDir, OutDir.
8. Add TOML tags on every field.

Do NOT:
- Do not write loading or validation logic here.
- Do not parse files.

Done when:
- go build ./... succeeds.
- All structs documented with short comments.
- TaskTimeout stored as time.Duration after parse (see config loader task; here just declare a string field plus a parsed field, or use a custom type).

Files: internal/config/config.go
