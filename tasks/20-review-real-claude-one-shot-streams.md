# REVIEW: real claude one-shot streams

Goal: Manually confirm a real claude invocation streams and errors surface cleanly.

Context:
- Verification task. Network/CLI required. If unavailable, document what was checked and defer.

Do:
1. With claude CLI installed and authed, write a tiny throwaway main or use a debug flag to invoke claudeAdapter once with a trivial prompt.
2. Confirm tokens stream and Result has no error and a session id.
3. Force an error (bad flag/unauth) and confirm Result.Err carries a useful message.
4. Remove any throwaway debug code afterward.

Done when:
- Real streaming observed.
- Errors produce actionable messages, not silent failure.
- go build/test still pass; no debug cruft left.
