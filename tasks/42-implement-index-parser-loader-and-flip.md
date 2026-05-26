# Implement index parser, loader, and flip

Goal: Read plan.md, load task prompts, and flip checkboxes durably.

Context:
- Loop mode depends on these. Title may carry an 'executor:'/'reviewer:' prefix to populate Entry.Agent.

Do:
1. In internal/plan/read.go add Entry { Checked bool; Title, Agent, FileRef string; Line int }.
2. func Parse(planPath string) ([]Entry, error): read lines, match '- [ ] '/'- [x] ' with a ' → ' file ref; ignore non-checkbox lines (headers/blanks); error if a checkbox line lacks a ref.
3. func LoadTask(planDir string, e Entry) (string, error): read the referenced file, return contents verbatim; error if missing.
4. func Flip(planPath string, e Entry) error: set that line to '- [x] ', write back via temp+rename preserving all other lines.
5. func Remaining(entries) []Entry: unchecked only, in order.

Done when:
- go build ./... succeeds.
- Parse handles mixed checked/unchecked and ignores headers.
- Flip changes one line, minimal diff, atomic.
- Missing ref or missing file -> clear error.

Files: internal/plan/read.go
