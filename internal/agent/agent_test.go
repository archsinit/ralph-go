package agent

import (
	"context"
	"strings"
	"testing"
	"time"
)

// --- buildClaudeArgs ---

func TestBuildClaudeArgs_Basic(t *testing.T) {
	req := Request{NewMessage: "hello"}
	args := buildClaudeArgs(req)
	assertContains(t, args, "--print")
	assertContains(t, args, "hello")
	assertNotContains(t, args, "--resume")
}

func TestBuildClaudeArgs_WithSystemPrompt(t *testing.T) {
	req := Request{SystemPrompt: "be brief", NewMessage: "hi"}
	args := buildClaudeArgs(req)
	assertContains(t, args, "--system-prompt")
	assertContains(t, args, "be brief")
}

func TestBuildClaudeArgs_WithResume(t *testing.T) {
	req := Request{NewMessage: "continue", ResumeSessionID: "ses_abc"}
	args := buildClaudeArgs(req)
	assertContains(t, args, "--resume")
	assertContains(t, args, "ses_abc")
}

func TestBuildClaudeArgs_NoResume(t *testing.T) {
	req := Request{NewMessage: "start"}
	args := buildClaudeArgs(req)
	assertNotContains(t, args, "--resume")
}

// --- buildCodexArgs ---

func TestBuildCodexArgs_Basic(t *testing.T) {
	req := Request{NewMessage: "explain this"}
	args := buildCodexArgs(req)
	assertContains(t, args, "--quiet")
	assertContains(t, args, "explain this")
}

func TestBuildCodexArgs_WithTranscript(t *testing.T) {
	req := Request{
		SystemPrompt: "be helpful",
		Transcript:   []Message{{Author: "user", Role: "user", Text: "prior turn"}},
		NewMessage:   "new message",
	}
	args := buildCodexArgs(req)
	// When transcript is present, the message arg should be the replay blob
	last := args[len(args)-1]
	if !strings.Contains(last, "prior turn") {
		t.Errorf("expected transcript in message arg, got: %s", last)
	}
	if !strings.Contains(last, "new message") {
		t.Errorf("expected new message in replay blob, got: %s", last)
	}
}

// --- renderReplay ---

func TestRenderReplay_Empty(t *testing.T) {
	req := Request{NewMessage: "hello"}
	out := renderReplay(req)
	if out != "hello" {
		t.Errorf("expected bare message, got: %q", out)
	}
}

func TestRenderReplay_WithSystemAndTranscript(t *testing.T) {
	req := Request{
		SystemPrompt: "system instructions",
		Transcript: []Message{
			{Author: "user", Text: "first message"},
			{Author: "claude", Text: "first reply"},
		},
		NewMessage: "second message",
	}
	out := renderReplay(req)

	checks := []string{
		"SYSTEM: system instructions",
		"user: first message",
		"claude: first reply",
		"second message",
		replayDelimiter,
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("renderReplay output missing %q:\n%s", want, out)
		}
	}

	// system prompt must appear before transcript
	sysIdx := strings.Index(out, "SYSTEM:")
	firstMsgIdx := strings.Index(out, "user: first message")
	if sysIdx > firstMsgIdx {
		t.Error("system prompt should appear before transcript")
	}

	// transcript must appear before new message
	lastMsgIdx := strings.LastIndex(out, "claude: first reply")
	newMsgIdx := strings.Index(out, "second message")
	if lastMsgIdx > newMsgIdx {
		t.Error("transcript should appear before new message")
	}
}

// --- echo adapter ---

func TestEchoAdapter_StreamsAndResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	a := NewEchoAdapter()
	if a.Name() != "echo" {
		t.Fatalf("expected Name() == echo, got %s", a.Name())
	}

	req := Request{NewMessage: "one two three"}
	tokenCh, resultCh := a.Invoke(ctx, req)

	var collected strings.Builder
	for t := range tokenCh {
		collected.WriteString(t.Text)
	}

	r := <-resultCh
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.SessionID == "" {
		t.Error("expected non-empty SessionID from echo adapter")
	}

	got := strings.TrimSpace(collected.String())
	if got != "one two three" {
		t.Errorf("reconstructed message mismatch: got %q", got)
	}
}

func TestEchoAdapter_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	a := NewEchoAdapter()
	req := Request{NewMessage: "should not stream"}
	tokenCh, resultCh := a.Invoke(ctx, req)

	for range tokenCh {
	}
	r := <-resultCh
	if r.Err == nil {
		t.Error("expected error on cancelled context")
	}
}

// --- helpers ---

func assertContains(t *testing.T, args []string, want string) {
	t.Helper()
	for _, a := range args {
		if a == want {
			return
		}
	}
	t.Errorf("args %v missing %q", args, want)
}

func assertNotContains(t *testing.T, args []string, unwanted string) {
	t.Helper()
	for _, a := range args {
		if a == unwanted {
			t.Errorf("args %v should not contain %q", args, unwanted)
			return
		}
	}
}
