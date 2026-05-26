# REVIEW: scaffold compiles and is wired

Goal: Verify the scaffold builds, subcommands register, and there is no dead or unused code.

Context:
- This is a verification task. Do not add features; only fix issues found.

Do:
1. Run go build ./..., go vet ./..., gofmt -l . (expect no output).
2. Run ./ralph-go --help and confirm both subcommands shown.
3. Confirm every internal package compiles and has a doc.go.

Done when:
- All build/vet/fmt checks clean.
- Both subcommands invocable.
- No unused imports or files.
