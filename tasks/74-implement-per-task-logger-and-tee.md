# Implement per-task logger and tee

Goal: Write a per-task log capturing full agent and git output, also shown in TUI.

Context:
- Slug/number from internal/plan helpers for consistent filenames.

Do:
1. Add TaskLog with Open(logDir, idx, slug) creating tasks/NN-slug.log under logDir.
2. func (t *TaskLog) Write(p []byte) so it satisfies io.Writer; flush per write.
3. Provide a Tee(io.Writer...) helper so engine can fan agent token streams to both the TUI and the task log.
4. Capture executor output, reviewer output, and git output into the task log with section headers.

Done when:
- go build ./... succeeds.
- Per-task log contains executor+reviewer+git sections.
- Tee writes to all targets without reordering within a stream.

Files: internal/logx/logx.go
