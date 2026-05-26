# ralph-go — build plan

Index of build tasks. Each line references a prompt file under tasks/. Run these with opencode (hand-fed through Phase 13, then via ralph-go loop for the remainder).

## Phase 0 — scaffold

- [x] Init Go module and repo skeleton → tasks/01-init-go-module-and-repo-skeleton.md
- [x] Add dependencies → tasks/02-add-dependencies.md
- [x] Create directory layout → tasks/03-create-directory-layout.md
- [x] Wire cobra root and subcommands → tasks/04-wire-cobra-root-and-subcommands.md
- [x] reviewer: REVIEW: scaffold compiles and is wired → tasks/05-review-scaffold-compiles-and-is-wired.md

## Phase 1 — config

- [ ] Define config structs → tasks/06-define-config-structs.md
- [ ] Implement TOML config loader → tasks/07-implement-toml-config-loader.md
- [ ] Implement config validation → tasks/08-implement-config-validation.md
- [ ] Write example ralph.toml → tasks/09-write-example-ralph-toml.md
- [ ] executor: TEST: config parse and validation → tasks/10-test-config-parse-and-validation.md
- [ ] reviewer: REVIEW: config error UX → tasks/11-review-config-error-ux.md

## Phase 2 — agent adapter

- [ ] Define Adapter interface and types → tasks/12-define-adapter-interface-and-types.md
- [ ] Implement exec runner helper → tasks/13-implement-exec-runner-helper.md
- [ ] Implement claude adapter → tasks/14-implement-claude-adapter.md
- [ ] Implement codex adapter → tasks/15-implement-codex-adapter.md
- [ ] Implement gemini and opencode adapters → tasks/16-implement-gemini-and-opencode-adapters.md
- [ ] Implement transcript replay fallback → tasks/17-implement-transcript-replay-fallback.md
- [ ] Implement registry and echo adapter → tasks/18-implement-registry-and-echo-adapter.md
- [ ] executor: TEST: adapter flag building and echo streaming → tasks/19-test-adapter-flag-building-and-echo-streaming.md
- [ ] reviewer: REVIEW: real claude one-shot streams → tasks/20-review-real-claude-one-shot-streams.md

## Phase 3 — session and persistence

- [ ] Define transcript and session model → tasks/21-define-transcript-and-session-model.md
- [ ] Implement append and flush writer → tasks/22-implement-append-and-flush-writer.md
- [ ] Implement load and resume → tasks/23-implement-load-and-resume.md
- [ ] executor: TEST: session round-trip and resume → tasks/24-test-session-round-trip-and-resume.md
- [ ] reviewer: REVIEW: crash-safety of session → tasks/25-review-crash-safety-of-session.md

## Phase 4 — plan TUI

- [ ] Build TUI model skeleton → tasks/26-build-tui-model-skeleton.md
- [ ] Render messages with per-author styling → tasks/27-render-messages-with-per-author-styling.md
- [ ] Stream tokens into in-progress message → tasks/28-stream-tokens-into-in-progress-message.md
- [ ] Add turn/status indicator → tasks/29-add-turn-status-indicator.md
- [ ] executor: TEST: TUI model updates → tasks/30-test-tui-model-updates.md
- [ ] reviewer: REVIEW: TUI interaction by hand → tasks/31-review-tui-interaction-by-hand.md

## Phase 5 — plan orchestrator

- [ ] Implement prefix parser → tasks/32-implement-prefix-parser.md
- [ ] Implement round-robin turn engine → tasks/33-implement-round-robin-turn-engine.md
- [ ] Implement cancel current turn → tasks/34-implement-cancel-current-turn.md
- [ ] Wire plan subcommand end-to-end → tasks/35-wire-plan-subcommand-end-to-end.md
- [ ] executor: TEST: orchestrator order, routing, cancel → tasks/36-test-orchestrator-order-routing-cancel.md
- [ ] reviewer: REVIEW: plan loop with two echo agents → tasks/37-review-plan-loop-with-two-echo-agents.md

## Phase 5.5 — plan format package (shared)

- [ ] Define TaskSpec schema and strict decoder → tasks/38-define-taskspec-schema-and-strict-decoder.md
- [ ] Implement slug and numbering helpers → tasks/39-implement-slug-and-numbering-helpers.md
- [ ] Implement task file renderer → tasks/40-implement-task-file-renderer.md
- [ ] Implement index writer → tasks/41-implement-index-writer.md
- [ ] Implement index parser, loader, and flip → tasks/42-implement-index-parser-loader-and-flip.md
- [ ] executor: TEST: plan format round-trip → tasks/43-test-plan-format-round-trip.md
- [ ] reviewer: REVIEW: generated artifacts are clean → tasks/44-review-generated-artifacts-are-clean.md

## Phase 6 — checklist generation

- [ ] Implement /plan command in chatroom → tasks/45-implement-plan-command-in-chatroom.md
- [ ] Parse plan output and write files → tasks/46-parse-plan-output-and-write-files.md
- [ ] Implement /quit command → tasks/47-implement-quit-command.md
- [ ] executor: TEST: /plan writes valid files with echo agent → tasks/48-test-plan-writes-valid-files-with-echo-agent.md
- [ ] reviewer: REVIEW: real /plan generation quality → tasks/49-review-real-plan-generation-quality.md

## Phase 7 — plan integration and polish

- [ ] reviewer: End-to-end real plan session → tasks/50-end-to-end-real-plan-session.md
- [ ] reviewer: Resume an existing plan session → tasks/51-resume-an-existing-plan-session.md
- [ ] Plan error UX hardening → tasks/52-plan-error-ux-hardening.md
- [ ] Cross-platform build → tasks/53-cross-platform-build.md
- [ ] README for plan mode → tasks/54-readme-for-plan-mode.md
- [ ] reviewer: REVIEW: fresh-machine plan smoke test → tasks/55-review-fresh-machine-plan-smoke-test.md

## Phase 8 — loop config and shared

- [ ] Finalize loop config fields → tasks/56-finalize-loop-config-fields.md
- [ ] Validate loop config → tasks/57-validate-loop-config.md
- [ ] executor: TEST: loop config validation → tasks/58-test-loop-config-validation.md
- [ ] reviewer: REVIEW: loop config errors → tasks/59-review-loop-config-errors.md

## Phase 9 — checklist parsing for loop

- [ ] Wire loop to plan parser and loader → tasks/60-wire-loop-to-plan-parser-and-loader.md
- [ ] Implement checkbox flip on success → tasks/61-implement-checkbox-flip-on-success.md
- [ ] executor: TEST: loop queue parse, load, flip → tasks/62-test-loop-queue-parse-load-flip.md
- [ ] reviewer: REVIEW: flip crash-safety → tasks/63-review-flip-crash-safety.md

## Phase 10 — ntfy client

- [ ] Implement ntfy client → tasks/64-implement-ntfy-client.md
- [ ] Implement event helpers → tasks/65-implement-event-helpers.md
- [ ] executor: TEST: ntfy request building and resilience → tasks/66-test-ntfy-request-building-and-resilience.md
- [ ] reviewer: REVIEW: real ntfy ping → tasks/67-review-real-ntfy-ping.md

## Phase 11 — git integration

- [ ] Implement dirty-tree detection → tasks/68-implement-dirty-tree-detection.md
- [ ] Implement per-task commit → tasks/69-implement-per-task-commit.md
- [ ] Implement diff capture for reviewer → tasks/70-implement-diff-capture-for-reviewer.md
- [ ] executor: TEST: git operations in a temp repo → tasks/71-test-git-operations-in-a-temp-repo.md
- [ ] reviewer: REVIEW: commits land correctly → tasks/72-review-commits-land-correctly.md

## Phase 12 — logging

- [ ] Implement master logger → tasks/73-implement-master-logger.md
- [ ] Implement per-task logger and tee → tasks/74-implement-per-task-logger-and-tee.md
- [ ] executor: TEST: logs produced with expected content → tasks/75-test-logs-produced-with-expected-content.md
- [ ] reviewer: REVIEW: logs readable under streaming → tasks/76-review-logs-readable-under-streaming.md

## Phase 13 — loop engine

- [ ] Implement core loop iteration → tasks/77-implement-core-loop-iteration.md
- [ ] Implement per-task timeout → tasks/78-implement-per-task-timeout.md
- [ ] Implement retry and failure handling → tasks/79-implement-retry-and-failure-handling.md
- [ ] Implement graceful stop → tasks/80-implement-graceful-stop.md
- [ ] Wire loop subcommand end-to-end → tasks/81-wire-loop-subcommand-end-to-end.md
- [ ] executor: TEST: loop engine paths with echo agents → tasks/82-test-loop-engine-paths-with-echo-agents.md
- [ ] reviewer: REVIEW: loop dry-run end-to-end → tasks/83-review-loop-dry-run-end-to-end.md

## Phase 14 — loop TUI

- [ ] Build loop TUI two-pane model → tasks/84-build-loop-tui-two-pane-model.md
- [ ] Show progress, retries, timing, last event → tasks/85-show-progress-retries-timing-last-event.md
- [ ] executor: TEST: loop TUI model updates → tasks/86-test-loop-tui-model-updates.md
- [ ] reviewer: REVIEW: loop TUI by hand → tasks/87-review-loop-tui-by-hand.md

## Phase 15 — loop integration and polish

- [ ] reviewer: End-to-end real loop run → tasks/88-end-to-end-real-loop-run.md
- [ ] Loop failure UX hardening → tasks/89-loop-failure-ux-hardening.md
- [ ] reviewer: Confirm resume after stop → tasks/90-confirm-resume-after-stop.md
- [ ] Complete README for loop → tasks/91-complete-readme-for-loop.md
- [ ] reviewer: REVIEW: fresh-machine full smoke test → tasks/92-review-fresh-machine-full-smoke-test.md
- [ ] reviewer: REVIEW: regression sweep across all phases → tasks/93-review-regression-sweep-across-all-phases.md
