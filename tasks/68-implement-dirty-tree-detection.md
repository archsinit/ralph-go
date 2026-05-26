# Implement dirty-tree detection

Goal: Detect a dirty working tree before the loop starts.

Context:
- Loop halts on dirty tree per design.

Do:
1. In internal/gitx/git.go add func IsDirty(dir string) (bool, error) using 'git status --porcelain' (non-empty output = dirty).
2. Add func IsRepo(dir string) (bool, error).
3. Run commands with the project dir as cwd; capture output.

Done when:
- go build ./... succeeds.
- Dirty repo detected true; clean false.
- Non-repo dir reported clearly.

Files: internal/gitx/git.go
