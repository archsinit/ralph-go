# TEST: ntfy request building and resilience

Goal: Unit-test request construction against a mock server and failure tolerance.

Do:
1. Create internal/ntfy/ntfy_test.go with httptest.Server capturing requests.
2. Test: Publish sets path to /topic, Title/Priority/Tags headers, and Authorization when Token set.
3. Test: server returns 500 or is unreachable -> Publish returns error but helper-level call does not panic and is non-fatal.

Done when:
- go test ./internal/ntfy/... passes.
- Headers asserted.
- Failure tolerance asserted.

Files: internal/ntfy/ntfy_test.go
