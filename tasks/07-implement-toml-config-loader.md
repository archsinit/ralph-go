# Implement TOML config loader

Goal: Load and decode the TOML config file into Config, wired to the --config flag.

Context:
- Default path ./ralph.toml from the global flag.

Do:
1. Add func Load(path string) (*Config, error) in internal/config/config.go.
2. Read file; on missing/unreadable return error wrapping the path.
3. Decode TOML into Config; on decode error wrap with path.
4. If using a custom Duration type, ensure UnmarshalText parses values like "15m".
5. In cmd/root.go resolve --config path; have plan and loop call config.Load before running; on error print to stderr and exit 1 (use a helper).

Do NOT:
- Do not add semantic validation here (separate task).

Done when:
- go build ./... succeeds.
- Valid config loads without error.
- Missing or malformed file prints a clear error naming the path and exits nonzero.

Files: internal/config/config.go, cmd/root.go
