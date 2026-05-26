# Implement per-task commit

Goal: Stage all changes and commit with a provided message, capturing output.

Context:
- Commit message authored by the reviewer from output + diff.

Do:
1. Add func Commit(dir, message string) (output string, err error): run 'git add -A' then 'git commit -m <message>' (use -F - via stdin to allow multi-line messages).
2. Detect 'nothing to commit' and return a sentinel error EmptyDiff so the caller can treat an empty change set as a review failure, not a silent pass.
3. Return combined stdout/stderr.

Done when:
- go build ./... succeeds.
- Clean change set commits.
- Empty diff yields EmptyDiff sentinel.
- Multi-line messages supported.

Files: internal/gitx/git.go
