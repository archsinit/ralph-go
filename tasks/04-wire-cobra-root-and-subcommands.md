# Wire cobra root and subcommands

Goal: Create the root command plus plan and loop subcommands; loop prints a not-implemented notice.

Context:
- Binary name ralph-go.
- --config is global, shared by both subcommands.

Do:
1. In cmd/root.go define rootCmd (use: "ralph-go") and Execute().
2. Add persistent --config flag (default ./ralph.toml), stored in a package var.
3. In cmd/plan.go define planCmd; Run prints "plan: not implemented" for now.
4. In cmd/loop.go define loopCmd; Run prints "loop: not implemented".
5. Register both on rootCmd in init().
6. main.go calls cmd.Execute().

Done when:
- go build ./... succeeds.
- ./ralph-go --help lists plan and loop.
- ./ralph-go plan and ./ralph-go loop each print their placeholder line and exit 0.

Files: cmd/root.go, cmd/plan.go, cmd/loop.go, main.go
