# ralph-go — build plan

Index of build tasks. Each line references a prompt file under tasks/. Run these with opencode (hand-fed through Phase 13, then via ralph-go loop for the remainder).

## Phase 0 — scaffold

- [x] Init Go module and repo skeleton → tasks/01-init-go-module-and-repo-skeleton.md
- [x] Add dependencies → tasks/02-add-dependencies.md
- [x] Create directory layout → tasks/03-create-directory-layout.md
- [x] Wire cobra root and subcommands → tasks/04-wire-cobra-root-and-subcommands.md
- [x] reviewer: REVIEW: scaffold compiles and is wired → tasks/05-review-scaffold-compiles-and-is-wired.md

## Phase 1 — config

- [x] Define config structs → tasks/06-define-config-structs.md
- [x] Implement TOML config loader → tasks/07-implement-toml-config-loader.md
- [x] Implement config validation → tasks/08-implement-config-validation.md
- [x] Write example ralph.toml → tasks/09-write-example-ralph-toml.md
- [x] executor: TEST: config parse and validation → tasks/10-test-config-parse-and-validation.md
- [x] reviewer: REVIEW: config error UX → tasks/11-review-config-error-ux.md

## Phase 1.5 — config review follow-up

- [x] Fix config required-field validation and gofmt drift → tasks/11a-fix-config-required-fields-and-gofmt.md
- [x] reviewer: REVIEW: phase 1.5 config readiness → tasks/11b-review-phase-1-5-config-readiness.md

## Phase 2 — agent adapter

- [x] Define Adapter interface and types → tasks/12-define-adapter-interface-and-types.md
- [x] Implement exec runner helper → tasks/13-implement-exec-runner-helper.md
- [x] Implement claude adapter → tasks/14-implement-claude-adapter.md
- [x] Implement codex adapter → tasks/15-implement-codex-adapter.md
- [x] Implement gemini and opencode adapters → tasks/16-implement-gemini-and-opencode-adapters.md
- [x] Implement transcript replay fallback → tasks/17-implement-transcript-replay-fallback.md
- [x] Implement registry and echo adapter → tasks/18-implement-registry-and-echo-adapter.md
- [x] executor: TEST: adapter flag building and echo streaming → tasks/19-test-adapter-flag-building-and-echo-streaming.md
- [x] reviewer: REVIEW: real claude one-shot streams → tasks/20-review-real-claude-one-shot-streams.md

## Phase 2.5 — adapter correctness follow-up

- [x] Fix adapter stream contract and replay fallback → tasks/20a-fix-adapter-stream-contract-and-replay.md
- [x] Align real CLI adapter arguments with installed CLIs → tasks/20b-align-real-cli-adapter-arguments.md
- [x] reviewer: REVIEW: phase 2.5 adapter readiness → tasks/20c-review-phase-2-5-adapter-readiness.md

## Phase 3 — session and persistence

- [x] Define transcript and session model → tasks/21-define-transcript-and-session-model.md
- [x] Implement append and flush writer → tasks/22-implement-append-and-flush-writer.md
- [x] Implement load and resume → tasks/23-implement-load-and-resume.md
- [x] executor: TEST: session round-trip and resume → tasks/24-test-session-round-trip-and-resume.md
- [x] reviewer: REVIEW: crash-safety of session → tasks/25-review-crash-safety-of-session.md

## Phase 3.5 — session crash-safety follow-up

- [x] Harden session persistence edge cases → tasks/25a-harden-session-persistence-edge-cases.md
- [x] reviewer: REVIEW: phase 3.5 session durability readiness → tasks/25b-review-phase-3-5-session-durability-readiness.md

## Phase 4 — plan TUI

- [x] Build TUI model skeleton → tasks/26-build-tui-model-skeleton.md
- [x] Render messages with per-author styling → tasks/27-render-messages-with-per-author-styling.md
- [x] Stream tokens into in-progress message → tasks/28-stream-tokens-into-in-progress-message.md
- [x] Add turn/status indicator → tasks/29-add-turn-status-indicator.md
- [x] executor: TEST: TUI model updates → tasks/30-test-tui-model-updates.md
- [x] reviewer: REVIEW: TUI interaction by hand → tasks/31-review-tui-interaction-by-hand.md

## Phase 4.5 — TUI/session readiness follow-up

- [x] Fix TUI orchestrator-facing API and submit contract → tasks/31a-fix-tui-external-event-and-submit-hooks.md
- [x] Fix TUI viewport, spinner, rendering, and harness behavior → tasks/31b-fix-tui-scrolling-autoscroll-and-spinner-behavior.md
- [x] Harden remaining session sequence and durability edge cases → tasks/31c-harden-remaining-session-durability-edge-cases.md
- [x] reviewer: REVIEW: phase 4.5 readiness for orchestrator → tasks/31d-review-phase-4-5-readiness-for-orchestrator.md

## Phase 5 — plan orchestrator

- [x] Implement prefix parser → tasks/32-implement-prefix-parser.md
- [x] Implement round-robin turn engine → tasks/33-implement-round-robin-turn-engine.md
- [x] Implement cancel current turn → tasks/34-implement-cancel-current-turn.md
- [x] Wire plan subcommand end-to-end → tasks/35-wire-plan-subcommand-end-to-end.md
- [x] executor: TEST: orchestrator order, routing, cancel → tasks/36-test-orchestrator-order-routing-cancel.md
- [x] reviewer: REVIEW: plan loop with two echo agents → tasks/37-review-plan-loop-with-two-echo-agents.md

## Phase 5.5 — plan orchestrator readiness follow-up

- [x] Fix TUI runtime lifecycle, cancel signaling, and bridge safety → tasks/37a-fix-tui-runtime-lifecycle-cancel-and-bridge-safety.md
- [x] Fix engine turn semantics, cancellation, and adapter request contract → tasks/37b-fix-engine-turn-semantics-cancellation-and-adapter-request-contract.md
- [x] Fix plan command wiring and turn-order model ownership → tasks/37c-fix-plan-command-wiring-and-turn-order-model-ownership.md
- [x] executor: TEST: plan orchestrator lifecycle and cancellation regressions → tasks/37d-test-plan-orchestrator-lifecycle-and-cancellation-regressions.md
- [x] reviewer: REVIEW: phase 5.5 plan orchestrator readiness → tasks/37e-review-phase-5-5-plan-orchestrator-readiness.md

## Phase 5.6 — plan format package (shared)

- [x] Define TaskSpec schema and strict decoder → tasks/38-define-taskspec-schema-and-strict-decoder.md
- [x] Implement slug and numbering helpers → tasks/39-implement-slug-and-numbering-helpers.md
- [x] Implement task file renderer → tasks/40-implement-task-file-renderer.md
- [x] Implement index writer → tasks/41-implement-index-writer.md
- [x] Implement index parser, loader, and flip → tasks/42-implement-index-parser-loader-and-flip.md
- [x] executor: TEST: plan format round-trip → tasks/43-test-plan-format-round-trip.md
- [x] reviewer: REVIEW: generated artifacts are clean → tasks/44-review-generated-artifacts-are-clean.md

## Phase 6 — checklist generation

- [x] Implement /plan command in chatroom → tasks/45-implement-plan-command-in-chatroom.md
- [x] Parse plan output and write files → tasks/46-parse-plan-output-and-write-files.md
- [x] Implement /quit command → tasks/47-implement-quit-command.md
- [x] executor: TEST: /plan writes valid files with echo agent → tasks/48-test-plan-writes-valid-files-with-echo-agent.md
- [x] reviewer: REVIEW: real /plan generation quality → tasks/49-review-real-plan-generation-quality.md

## Phase 6.5 — plan chatroom and checklist generation audit follow-up

- [ ] Fix remaining TUI/bridge lifecycle, cancellation, and shutdown semantics → tasks/49a-fix-tui-bridge-lifecycle-cancellation-and-shutdown.md
- [ ] Fix engine turn cancellation and adapter request contracts → tasks/49b-fix-engine-turn-cancellation-and-adapter-request-contracts.md
- [ ] Align /plan generation with the plan package schema and safe artifact writes → tasks/49c-align-plan-generation-schema-and-safe-writes.md
- [ ] executor: TEST: phase 6.5 lifecycle, cancellation, and /plan regressions → tasks/49d-test-phase-6-5-lifecycle-cancellation-and-plan-regressions.md
- [ ] reviewer: REVIEW: phase 6.5 plan chatroom readiness → tasks/49e-review-phase-6-5-plan-chatroom-readiness.md

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
