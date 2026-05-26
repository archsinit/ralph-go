# Implement master logger

Goal: Append timestamped loop events to a master log, flushed per step.

Context:
- Plain text per design.

Do:
1. In internal/logx/logx.go add Master { f *os.File } with Open(logDir) creating master.log (append).
2. func (m *Master) Event(format string, args...) writing 'TS  message\n' and Sync each write.
3. func (m *Master) Close().

Done when:
- go build ./... succeeds.
- Events appear immediately on disk.
- Timestamps present.

Files: internal/logx/logx.go
