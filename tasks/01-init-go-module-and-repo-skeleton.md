# Init Go module and repo skeleton

Goal: Create the Go module, git repo, and base files so the project compiles empty.

Context:
- This is the ralph-go tool's own source repo.
- Go 1.22+ assumed.

Do:
1. Run: go mod init github.com/archsinit/ralph-go.
2. Create .gitignore (ignore built binary, *.log, /tmp).
3. Add a README.md stub with project name and one-line description.
4. Create empty main.go in repo root with package main and an empty main().

Do NOT:
- Do not add dependencies yet.
- Do not create subcommands yet.

Done when:
- go build ./... succeeds.
- go.mod has module github.com/archsinit/ralph-go.
- git status is clean after an initial commit by the user (do not commit yourself).

Files: go.mod, main.go, .gitignore, README.md
