# Implement diff capture for reviewer

Goal: Expose the staged/working diff so the reviewer can judge and describe changes.

Context:
- Feeds the reviewer prompt and the commit-message generation.

Do:
1. Add func Diff(dir string) (string, error) returning 'git diff' (unstaged) plus 'git diff --staged' as needed, or 'git add -A' then 'git diff --staged'.
2. Cap very large diffs (truncate with a note) to bound prompt size.

Done when:
- go build ./... succeeds.
- Diff returns changes after a task run.
- Large diffs truncated with a marker.

Files: internal/gitx/git.go
