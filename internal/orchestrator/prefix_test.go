package orchestrator

import (
	"testing"
)

func TestParsePrefix(t *testing.T) {
	agents := []string{"claude", "codex", "gemini"}

	tests := []struct {
		input       string
		agents      []string
		wantTarget  string
		wantBody    string
		description string
	}{
		{
			input:       "codex: do x",
			agents:      agents,
			wantTarget:  "codex",
			wantBody:    "do x",
			description: "exact match with colon",
		},
		{
			input:       "hello",
			agents:      agents,
			wantTarget:  "",
			wantBody:    "hello",
			description: "no prefix",
		},
		{
			input:       "foo: x",
			agents:      agents,
			wantTarget:  "",
			wantBody:    "foo: x",
			description: "unknown prefix",
		},
		{
			input:       "CLAUDE: something",
			agents:      agents,
			wantTarget:  "claude",
			wantBody:    "something",
			description: "case-insensitive match",
		},
		{
			input:       "claude:",
			agents:      agents,
			wantTarget:  "claude",
			wantBody:    "",
			description: "prefix only",
		},
		{
			input:       "claude: ",
			agents:      agents,
			wantTarget:  "claude",
			wantBody:    "",
			description: "prefix with space",
		},
		{
			input:       "no colon",
			agents:      agents,
			wantTarget:  "",
			wantBody:    "no colon",
			description: "missing colon",
		},
		{
			input:       "",
			agents:      agents,
			wantTarget:  "",
			wantBody:    "",
			description: "empty input",
		},
		{
			input:       ":something",
			agents:      agents,
			wantTarget:  "",
			wantBody:    ":something",
			description: "colon at start",
		},
		{
			input:       "   codex:  hello world  ",
			agents:      agents,
			wantTarget:  "codex",
			wantBody:    "hello world",
			description: "whitespace handling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			target, body := ParsePrefix(tt.input, tt.agents)
			if target != tt.wantTarget {
				t.Errorf("target = %q, want %q", target, tt.wantTarget)
			}
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestParsePrefixEmptyAgentList(t *testing.T) {
	target, body := ParsePrefix("something: text", []string{})
	if target != "" || body != "something: text" {
		t.Errorf("with empty agent list: target=%q (want empty), body=%q (want unchanged)", target, body)
	}
}
