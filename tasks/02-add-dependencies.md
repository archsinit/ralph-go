# Add dependencies

Goal: Add all third-party deps the project needs and pin them in go.mod.

Context:
- Deps: cobra (CLI), bubbletea/bubbles/lipgloss (TUI), BurntSushi/toml (config).

Do:
1. go get github.com/spf13/cobra@latest.
2. go get github.com/charmbracelet/bubbletea@latest.
3. go get github.com/charmbracelet/bubbles@latest.
4. go get github.com/charmbracelet/lipgloss@latest.
5. go get github.com/BurntSushi/toml@latest.
6. Run go mod tidy.

Done when:
- go build ./... succeeds.
- go.mod lists cobra, bubbletea, bubbles, lipgloss, toml.
- go.sum present and consistent.

Files: go.mod, go.sum
