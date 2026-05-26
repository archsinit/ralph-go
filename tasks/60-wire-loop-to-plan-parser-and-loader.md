# Wire loop to plan parser and loader

Goal: Use internal/plan to load the index and per-task prompts for the loop.

Context:
- No new parsing logic; reuse Phase 5.5.

Do:
1. In internal/loop/queue.go add func LoadQueue(planPath string) (entries []plan.Entry, remaining []plan.Entry, err error) using plan.Parse and plan.Remaining.
2. Add func PromptFor(planDir string, e plan.Entry) (string, error) wrapping plan.LoadTask.
3. Resolve each task's agent: Entry.Agent if set else config.Loop.Executor for execution; reviewer always config.Loop.Reviewer.

Done when:
- go build ./... succeeds.
- LoadQueue returns ordered entries and unchecked remainder.
- Missing referenced task file -> clear error.

Files: internal/loop/queue.go
