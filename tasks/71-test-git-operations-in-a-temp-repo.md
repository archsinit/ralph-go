# TEST: git operations in a temp repo

Goal: Unit-test dirty detection, commit, empty-diff, and diff capture.

Do:
1. Create internal/gitx/git_test.go: init a temp repo (git init, set user.email/name), create a file.
2. Test IsDirty true before commit, false after.
3. Test Commit creates a commit; a second Commit with no changes returns EmptyDiff.
4. Test Diff returns expected file content after a change.
5. Skip tests if git binary absent (t.Skip).

Done when:
- go test ./internal/gitx/... passes where git present.
- All four behaviors asserted.

Files: internal/gitx/git_test.go
