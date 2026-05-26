# REVIEW: TUI interaction by hand

Goal: Manually confirm typing, scrolling, live streaming, and resize behave.

Context:
- Verification task.

Do:
1. Run a temporary harness feeding fake streams via the echo adapter.
2. Type, submit, scroll back, resize the terminal; confirm wrapping and auto-scroll behave.
3. Confirm spinner/status correct during/after streams.

Done when:
- Smooth interaction observed.
- No layout corruption on resize.
- Tests pass.
