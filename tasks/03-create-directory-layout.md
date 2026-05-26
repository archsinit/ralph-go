# Create directory layout

Goal: Establish the package directory structure with placeholder files so imports resolve.

Context:
- Package list maps to design components: config, agent adapters, session, orchestrator (plan turn engine), plan (shared format), tui, loop engine, ntfy client, git integration, logging.

Do:
1. Create dirs: cmd/, internal/config/, internal/agent/, internal/session/, internal/orchestrator/, internal/plan/, internal/tui/, internal/loop/, internal/ntfy/, internal/gitx/, internal/logx/.
2. In each internal/<pkg>/ add a doc.go declaring the package with a one-line comment.
3. In cmd/ add doc.go with package cmd.

Do NOT:
- Do not implement logic, only package declarations.

Done when:
- go build ./... succeeds.
- Every internal package directory has at least one .go file.
- go vet ./... passes.

Files: cmd/doc.go, internal/config/doc.go, internal/agent/doc.go, internal/session/doc.go, internal/orchestrator/doc.go, internal/plan/doc.go, internal/tui/doc.go, internal/loop/doc.go, internal/ntfy/doc.go, internal/gitx/doc.go, internal/logx/doc.go
