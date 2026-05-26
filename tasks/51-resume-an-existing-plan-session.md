# Resume an existing plan session

Goal: Verify resuming a saved session restores context and continues correctly.

Do:
1. Run ./ralph-go plan --session <dir> on a prior session.
2. Confirm transcript restored, agents resume via stored session ids (or replay), new turns append correctly.

Done when:
- Resume restores full transcript.
- Agents continue coherently.
- New messages persist.
