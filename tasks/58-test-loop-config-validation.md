# TEST: loop config validation

Goal: Unit-test the new loop validation rules.

Do:
1. Extend internal/config/config_test.go with loop cases: zero timeout, negative retries, missing executor/reviewer, empty ntfy topic, bad server URL -> errors; valid -> nil.

Done when:
- go test ./internal/config/... passes.
- Each loop rule asserted.

Files: internal/config/config_test.go
