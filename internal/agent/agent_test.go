package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// --- buildClaudeArgs ---

func TestBuildClaudeArgs_Basic(t *testing.T) {
	req := Request{NewMessage: "hello"}
	args := buildClaudeArgs(req)
	assertContains(t, args, "--print")
	assertContains(t, args, "--include-partial-messages")
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
	assertEqualArgs(t, args, []string{"exec", "explain this"})
	assertNotContains(t, args, "--quiet")
	assertNotContains(t, args, "--system-prompt")
}

func TestBuildCodexArgs_WithTranscript(t *testing.T) {
	req := Request{
		SystemPrompt: "be helpful",
		Transcript:   []Message{{Author: "user", Role: "user", Text: "prior turn"}},
		NewMessage:   "new message",
	}
	args := buildCodexArgs(req)
	// When transcript is present, the message arg should be the replay blob
	assertContains(t, args, "exec")
	last := args[len(args)-1]
	if !strings.Contains(last, "prior turn") {
		t.Errorf("expected transcript in message arg, got: %s", last)
	}
	if !strings.Contains(last, "new message") {
		t.Errorf("expected new message in replay blob, got: %s", last)
	}
}

func TestBuildCodexArgs_WithSystemPrompt(t *testing.T) {
	req := Request{
		SystemPrompt: "be concise",
		NewMessage:   "new message",
	}
	args := buildCodexArgs(req)
	assertContains(t, args, "exec")
	assertNotContains(t, args, "--system-prompt")
	last := args[len(args)-1]
	if !strings.Contains(last, "SYSTEM: be concise") {
		t.Errorf("expected system prompt in replay blob, got: %s", last)
	}
	if !strings.Contains(last, "new message") {
		t.Errorf("expected new message in replay blob, got: %s", last)
	}
}

func TestBuildCodexArgs_WithResume(t *testing.T) {
	req := Request{
		NewMessage:      "continue",
		ResumeSessionID: "ses_abc",
	}
	args := buildCodexArgs(req)
	assertEqualArgs(t, args, []string{"exec", "resume", "ses_abc", "continue"})
	assertNotContains(t, args, "--quiet")
	assertNotContains(t, args, "--system-prompt")
}

func TestBuildGeminiArgs_WithTranscriptUsesReplay(t *testing.T) {
	req := Request{
		Transcript:      []Message{{Author: "user", Text: "prior turn"}},
		NewMessage:      "new message",
		ResumeSessionID: "ignored-session",
	}
	args := buildGeminiArgs(req)
	last := args[len(args)-1]
	if !strings.Contains(last, "prior turn") || !strings.Contains(last, "new message") {
		t.Errorf("expected replay prompt in gemini args, got: %s", last)
	}
}

func TestBuildOpencodeArgs_WithTranscriptUsesReplay(t *testing.T) {
	req := Request{
		Transcript:      []Message{{Author: "user", Text: "prior turn"}},
		NewMessage:      "new message",
		ResumeSessionID: "ignored-session",
	}
	args := buildOpencodeArgs(req)
	last := args[len(args)-1]
	if !strings.Contains(last, "prior turn") || !strings.Contains(last, "new message") {
		t.Errorf("expected replay prompt in opencode args, got: %s", last)
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

	req := Request{NewMessage: " one  two\nthree\t"}
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

	got := collected.String()
	if got != req.NewMessage {
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

func TestEchoAdapter_ContextCancelWithEmptyMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	a := NewEchoAdapter()
	tokenCh, resultCh := a.Invoke(ctx, Request{})

	for range tokenCh {
	}
	r := readOneResult(t, resultCh)
	if !errors.Is(r.Err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", r.Err)
	}
}

func TestEchoAdapter_WithTranscriptUsesReplay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := Request{
		SystemPrompt: "system instructions",
		Transcript:   []Message{{Author: "user", Text: "prior turn"}},
		NewMessage:   "new message",
	}
	tokenCh, resultCh := NewEchoAdapter().Invoke(ctx, req)

	var collected strings.Builder
	for tok := range tokenCh {
		collected.WriteString(tok.Text)
	}
	r := readOneResult(t, resultCh)
	if r.Err != nil {
		t.Fatalf("unexpected result error: %v", r.Err)
	}
	if got, want := collected.String(), renderReplay(req); got != want {
		t.Fatalf("echo replay mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

// --- channel contract tests ---

func TestEchoAdapter_TokenChannelClosesBeforeResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	a := NewEchoAdapter()
	req := Request{NewMessage: "one two three"}
	tokenCh, resultCh := a.Invoke(ctx, req)

	// Drain all tokens
	for range tokenCh {
	}

	// After token channel closes, we should be able to receive exactly one Result
	r, ok := <-resultCh
	if !ok {
		t.Fatal("result channel closed without sending Result")
	}
	if r.Err != nil {
		t.Errorf("unexpected error: %v", r.Err)
	}

	// Result channel should be closed after the one Result
	_, ok = <-resultCh
	if ok {
		t.Error("result channel should be closed after Result")
	}
}

func TestEchoAdapter_CancellationSendsTerminalResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	a := NewEchoAdapter()
	req := Request{NewMessage: "one two three"}
	tokenCh, resultCh := a.Invoke(ctx, req)

	// Consume one token
	<-tokenCh

	// Cancel the context
	cancel()

	// Drain remaining tokens
	for range tokenCh {
	}

	// Should get exactly one Result with an error
	r, ok := <-resultCh
	if !ok {
		t.Fatal("result channel closed without sending Result on cancellation")
	}
	if r.Err == nil {
		t.Error("expected error Result on cancellation")
	}

	// Result channel should be closed after the one Result
	_, ok = <-resultCh
	if ok {
		t.Error("result channel should be closed after Result")
	}
}

func TestBuildClaudeArgs_WithTranscriptAndNoResume(t *testing.T) {
	req := Request{
		SystemPrompt: "be helpful",
		Transcript:   []Message{{Author: "user", Text: "prior turn"}},
		NewMessage:   "new message",
	}
	args := buildClaudeArgs(req)
	// When transcript is present but ResumeSessionID is empty, the message arg should be the replay blob
	last := args[len(args)-1]
	if !strings.Contains(last, "prior turn") {
		t.Errorf("expected transcript in message arg, got: %s", last)
	}
	if !strings.Contains(last, "new message") {
		t.Errorf("expected new message in replay blob, got: %s", last)
	}
}

// --- runStreaming ---

func TestRunStreaming_SuccessSendsExactlyOneResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tokenCh, resultCh := runStreaming(ctx, os.Args[0], helperProcessArgs("success"), "")

	var collected strings.Builder
	for tok := range tokenCh {
		collected.WriteString(tok.Text)
	}
	if got, want := collected.String(), "alpha\nbeta\n"; got != want {
		t.Fatalf("streamed tokens mismatch: got %q, want %q", got, want)
	}

	r := readOneResult(t, resultCh)
	if r.Err != nil {
		t.Fatalf("unexpected result error: %v", r.Err)
	}
}

func TestRunStreaming_CancellationSendsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tokenCh, resultCh := runStreaming(ctx, os.Args[0], helperProcessArgs("sleep"), "")

	tok, ok := <-tokenCh
	if !ok {
		t.Fatal("token channel closed before first token")
	}
	if tok.Text != "ready\n" {
		t.Fatalf("unexpected first token: %q", tok.Text)
	}

	cancel()
	for range tokenCh {
	}

	r := readOneResult(t, resultCh)
	if !errors.Is(r.Err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", r.Err)
	}
}

func TestRunStreaming_ProcessErrorIncludesStderrTail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tokenCh, resultCh := runStreaming(ctx, os.Args[0], helperProcessArgs("exit-error"), "")
	for range tokenCh {
	}

	r := readOneResult(t, resultCh)
	if r.Err == nil {
		t.Fatal("expected process error")
	}
	if !strings.Contains(r.Err.Error(), "exit:") {
		t.Fatalf("expected exit error, got %v", r.Err)
	}
	if !strings.Contains(r.Err.Error(), "stderr-tail-marker") {
		t.Fatalf("expected stderr tail in error, got %v", r.Err)
	}
}

func TestRunStreaming_ScannerErrorIsResultError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tokenCh, resultCh := runStreaming(ctx, os.Args[0], helperProcessArgs("long-line"), "")
	for range tokenCh {
	}

	r := readOneResult(t, resultCh)
	if r.Err == nil {
		t.Fatal("expected scanner error")
	}
	if !strings.Contains(r.Err.Error(), "scanner:") {
		t.Fatalf("expected scanner error, got %v", r.Err)
	}
}

// --- CLI adapter wrappers ---

func TestClaudeAdapter_ParsesPartialMessagesWithoutDuplicates(t *testing.T) {
	withStreamCommand(t, func(ctx context.Context, name string, args []string, stdin string) (<-chan Token, <-chan Result) {
		if name != "claude" {
			t.Fatalf("unexpected command name %q", name)
		}
		rawTokens := make(chan Token)
		rawResults := make(chan Result, 1)
		go func() {
			rawTokens <- Token{Text: `{"type":"system","session_id":"claude-session"}`}
			rawTokens <- Token{Text: `{"type":"assistant","session_id":"claude-session","message":{"content":[{"type":"text","text":"hel"}]}}`}
			rawTokens <- Token{Text: `{"type":"assistant","session_id":"claude-session","message":{"content":[{"type":"text","text":"hello"}]}}`}
			rawTokens <- Token{Text: `{"type":"assistant","session_id":"claude-session","message":{"content":[{"type":"text","text":"hello"}]}}`}
			close(rawTokens)
			rawResults <- Result{}
			close(rawResults)
		}()
		return rawTokens, rawResults
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tokenCh, resultCh := NewClaudeAdapter().Invoke(ctx, Request{NewMessage: "prompt"})
	var collected strings.Builder
	for tok := range tokenCh {
		collected.WriteString(tok.Text)
	}

	if got, want := collected.String(), "hello"; got != want {
		t.Fatalf("claude text mismatch: got %q, want %q", got, want)
	}
	r := readOneResult(t, resultCh)
	if r.Err != nil {
		t.Fatalf("unexpected result error: %v", r.Err)
	}
	if r.SessionID != "claude-session" {
		t.Fatalf("expected session ID from stream, got %q", r.SessionID)
	}
}

func TestCLIAdapters_CancellationSendsTerminalResult(t *testing.T) {
	adapters := []Adapter{
		NewClaudeAdapter(),
		NewCodexAdapter(),
		NewGeminiAdapter(),
		NewOpencodeAdapter(),
	}

	for _, adapter := range adapters {
		t.Run(adapter.Name(), func(t *testing.T) {
			rawDelivered := make(chan struct{})
			withStreamCommand(t, func(ctx context.Context, name string, args []string, stdin string) (<-chan Token, <-chan Result) {
				rawTokens := make(chan Token)
				rawResults := make(chan Result, 1)
				go func() {
					rawTokens <- adapterRawToken(name)
					close(rawDelivered)
					<-ctx.Done()
					close(rawTokens)
					rawResults <- Result{}
					close(rawResults)
				}()
				return rawTokens, rawResults
			})

			ctx, cancel := context.WithCancel(context.Background())
			tokenCh, resultCh := adapter.Invoke(ctx, Request{NewMessage: "prompt"})
			select {
			case <-rawDelivered:
			case <-time.After(2 * time.Second):
				t.Fatal("adapter did not receive raw token")
			}
			cancel()

			for tok := range tokenCh {
				t.Fatalf("unexpected token after cancellation: %q", tok.Text)
			}

			r := readOneResult(t, resultCh)
			if !errors.Is(r.Err, context.Canceled) {
				t.Fatalf("expected context cancellation error, got %v", r.Err)
			}
		})
	}
}

func TestCLIAdapters_PropagateExactlyOneTerminalResult(t *testing.T) {
	adapters := []Adapter{
		NewClaudeAdapter(),
		NewCodexAdapter(),
		NewGeminiAdapter(),
		NewOpencodeAdapter(),
	}

	for _, adapter := range adapters {
		t.Run(adapter.Name(), func(t *testing.T) {
			withStreamCommand(t, func(ctx context.Context, name string, args []string, stdin string) (<-chan Token, <-chan Result) {
				rawTokens := make(chan Token)
				rawResults := make(chan Result, 1)
				go func() {
					rawTokens <- adapterRawToken(name)
					close(rawTokens)
					rawResults <- Result{SessionID: "raw-session"}
					close(rawResults)
				}()
				return rawTokens, rawResults
			})

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			tokenCh, resultCh := adapter.Invoke(ctx, Request{NewMessage: "prompt"})
			for range tokenCh {
			}
			r := readOneResult(t, resultCh)
			if r.Err != nil {
				t.Fatalf("unexpected result error: %v", r.Err)
			}
		})
	}
}

// --- helpers ---

func TestRunStreamingHelperProcess(t *testing.T) {
	mode, ok := helperProcessMode()
	if !ok {
		return
	}

	switch mode {
	case "success":
		fmt.Fprintln(os.Stdout, "alpha")
		fmt.Fprintln(os.Stdout, "beta")
		os.Exit(0)
	case "sleep":
		fmt.Fprintln(os.Stdout, "ready")
		time.Sleep(10 * time.Second)
		os.Exit(0)
	case "exit-error":
		fmt.Fprintln(os.Stdout, "before failure")
		fmt.Fprint(os.Stderr, strings.Repeat("x", 600))
		fmt.Fprint(os.Stderr, "stderr-tail-marker")
		os.Exit(7)
	case "long-line":
		fmt.Fprint(os.Stdout, strings.Repeat("x", 70*1024))
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown helper mode %q", mode)
		os.Exit(2)
	}
}

func helperProcessArgs(mode string) []string {
	return []string{"-test.run=TestRunStreamingHelperProcess", "--", "--helper-mode", mode}
}

func helperProcessMode() (string, bool) {
	for i, arg := range os.Args {
		if arg == "--helper-mode" && i+1 < len(os.Args) {
			return os.Args[i+1], true
		}
	}
	return "", false
}

func withStreamCommand(t *testing.T, fn func(context.Context, string, []string, string) (<-chan Token, <-chan Result)) {
	t.Helper()
	previous := streamCommand
	streamCommand = fn
	t.Cleanup(func() {
		streamCommand = previous
	})
}

func adapterRawToken(name string) Token {
	if name == "claude" {
		return Token{Text: `{"type":"assistant","message":{"content":[{"type":"text","text":"chunk"}]}}`}
	}
	return Token{Text: "chunk"}
}

func readOneResult(t *testing.T, resultCh <-chan Result) Result {
	t.Helper()

	r, ok := <-resultCh
	if !ok {
		t.Fatal("result channel closed without sending Result")
	}
	if _, ok := <-resultCh; ok {
		t.Fatal("result channel sent more than one Result")
	}
	return r
}

func assertEqualArgs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("args length mismatch: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("args mismatch: got %v, want %v", got, want)
		}
	}
}

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
