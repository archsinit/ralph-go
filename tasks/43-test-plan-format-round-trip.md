# TEST: plan format round-trip

Goal: Unit-test decode, render, write, parse, load, and flip.

Do:
1. Create internal/plan/plan_test.go using t.TempDir().
2. Test Decode: valid JSON ok; unknown field, missing Do/DoneWhen/Title, empty tasks each error.
3. Test RenderTask byte-equality against an expected string (include omitted optional sections).
4. Test Write then Parse: index entries match titles, refs point to existing files, agent prefixes parsed back into Entry.Agent.
5. Test LoadTask returns rendered content verbatim.
6. Test Flip: flip one entry, reparse shows it checked, other lines byte-identical except the flipped one.

Done when:
- go test ./internal/plan/... passes.
- Round-trip and byte-equality asserted.
- Error cases asserted.

Files: internal/plan/plan_test.go
