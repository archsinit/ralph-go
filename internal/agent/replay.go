package agent

import "strings"

const replayDelimiter = "---"

// renderReplay renders a full conversation into a single prompt string for
// CLIs that lack native session resume. The format is:
//
//	SYSTEM: <system prompt>
//	---
//	<Author>: <text>
//	...
//	---
//	<NewMessage>
func renderReplay(req Request) string {
	var b strings.Builder

	if req.SystemPrompt != "" {
		b.WriteString("SYSTEM: ")
		b.WriteString(req.SystemPrompt)
		b.WriteString("\n")
		b.WriteString(replayDelimiter)
		b.WriteString("\n")
	}

	for _, m := range req.Transcript {
		b.WriteString(m.Author)
		b.WriteString(": ")
		b.WriteString(m.Text)
		b.WriteString("\n")
	}

	if len(req.Transcript) > 0 {
		b.WriteString(replayDelimiter)
		b.WriteString("\n")
	}

	b.WriteString(req.NewMessage)

	return b.String()
}
