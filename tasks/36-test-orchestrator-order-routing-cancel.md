# TEST: orchestrator order, routing, cancel

Goal: Unit-test the engine with a fake UI bridge and echo adapters.

Do:
1. Create internal/orchestrator/engine_test.go with a fake bridge that scripts user inputs and records streamed authors.
2. Test: order claude,codex,user repeats correctly for several cycles.
3. Test: user input "codex: hi" makes codex the next responder.
4. Test: cancel signal mid-turn stops streaming and engine continues to a user turn without deadlock.
5. Test ParsePrefix table cases (may live here or in prefix_test.go).

Done when:
- go test ./internal/orchestrator/... passes.
- Order, routing, and cancel all asserted.
- No goroutine leak/deadlock (use timeouts in test).

Files: internal/orchestrator/engine_test.go, internal/orchestrator/prefix_test.go
