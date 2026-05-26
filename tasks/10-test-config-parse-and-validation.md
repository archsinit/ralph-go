# TEST: config parse and validation

Goal: Unit-test config loading and validation across valid and invalid inputs.

Do:
1. Create internal/config/config_test.go.
2. Test: load a valid temp TOML -> no error, fields populated.
3. Test: missing file -> error mentions path.
4. Test: malformed TOML -> error.
5. Table test for Validate: duplicate names, unknown cli, turn_order ref missing, plan_agent missing, negative max_retries, zero timeout, empty ntfy topic -> each errors; a fully valid config -> nil.
6. Use t.TempDir() for files.

Done when:
- go test ./internal/config/... passes.
- Each invalid case asserts an error; valid case asserts nil.

Files: internal/config/config_test.go
