package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	msgs := []Message{
		{Author: "user", Role: "user", Text: "hello", TS: time.Now().UTC().Truncate(time.Millisecond)},
		{Author: "claude", Role: "assistant", Text: "hi there", TS: time.Now().UTC().Truncate(time.Millisecond)},
		{Author: "user", Role: "user", Text: "bye", TS: time.Now().UTC().Truncate(time.Millisecond)},
	}

	for _, m := range msgs {
		if err := s.Append(m); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if len(s2.Messages) != len(msgs) {
		t.Fatalf("got %d messages, want %d", len(s2.Messages), len(msgs))
	}
	for i, want := range msgs {
		got := s2.Messages[i]
		if got.Author != want.Author || got.Role != want.Role || got.Text != want.Text {
			t.Errorf("message[%d] mismatch: got %+v, want %+v", i, got, want)
		}
		if got.Seq != i {
			t.Errorf("message[%d] Seq=%d, want %d", i, got.Seq, i)
		}
	}
}

func TestAgentSessionRoundTrip(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if err := s.SetAgentSession("claude", "sess-abc123"); err != nil {
		t.Fatalf("SetAgentSession: %v", err)
	}

	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if got := s2.ResumeID("claude"); got != "sess-abc123" {
		t.Errorf("ResumeID=%q, want %q", got, "sess-abc123")
	}
	if got := s2.ResumeID("unknown"); got != "" {
		t.Errorf("ResumeID for unknown agent=%q, want empty", got)
	}
}

func TestCorruptTrailingLine(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if err := s.Append(Message{Author: "user", Role: "user", Text: "first"}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if err := s.Append(Message{Author: "claude", Role: "assistant", Text: "second"}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// Append a corrupt/truncated line to the transcript file.
	f, err := os.OpenFile(filepath.Join(dir, transcriptFile), os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("open transcript for corruption: %v", err)
	}
	f.WriteString(`{"seq":2,"author":"corrupt","role":"user","text":"trun`)
	f.Close()

	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen after corruption: %v", err)
	}

	if len(s2.Messages) != 2 {
		t.Errorf("got %d messages after corrupt line, want 2", len(s2.Messages))
	}
}

func TestResumedSequenceAssignment(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// Append messages (they will get seq 0, 1, 2)
	msgs := []Message{
		{Author: "user", Role: "user", Text: "first"},
		{Author: "claude", Role: "assistant", Text: "second"},
		{Author: "user", Role: "user", Text: "third"},
	}
	for _, m := range msgs {
		if err := s.Append(m); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	// Reopen and verify nextSeq is correct (should be 3, not len(Messages))
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	// Append a new message and verify it gets seq 3
	if err := s2.Append(Message{Author: "user", Role: "user", Text: "fourth"}); err != nil {
		t.Fatalf("Append after resume: %v", err)
	}

	if s2.Messages[3].Seq != 3 {
		t.Errorf("resumed message got Seq=%d, want 3", s2.Messages[3].Seq)
	}

	// Reopen again and verify all messages are present with correct sequences
	s3, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen after second append: %v", err)
	}

	if len(s3.Messages) != 4 {
		t.Fatalf("got %d messages after resume, want 4", len(s3.Messages))
	}
	for i := 0; i < 4; i++ {
		if s3.Messages[i].Seq != i {
			t.Errorf("message[%d] Seq=%d, want %d", i, s3.Messages[i].Seq, i)
		}
	}
}

func TestNonTrailingCorruptLine(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// Write two valid messages
	if err := s.Append(Message{Author: "user", Role: "user", Text: "first"}); err != nil {
		t.Fatalf("Append first: %v", err)
	}
	if err := s.Append(Message{Author: "claude", Role: "assistant", Text: "second"}); err != nil {
		t.Fatalf("Append second: %v", err)
	}

	// Manually append a corrupt line in the middle and a valid line after it
	f, err := os.OpenFile(filepath.Join(dir, transcriptFile), os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("open transcript: %v", err)
	}
	f.WriteString(`not valid json at all` + "\n")
	f.WriteString(`{"seq":3,"author":"valid","role":"user","text":"after corrupt"}` + "\n")
	f.Close()

	// Reopening should error because of the non-trailing corrupt line
	_, err = Open(dir)
	if err == nil {
		t.Error("Open should error on non-trailing corrupt line, got nil")
	}
}

func TestAppendMemoryConsistency(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// After successful append, verify the message is in memory
	if err := s.Append(Message{Author: "user", Role: "user", Text: "test"}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	if len(s.Messages) != 1 {
		t.Errorf("after append, got %d messages in memory, want 1", len(s.Messages))
	}
	if s.Messages[0].Text != "test" {
		t.Errorf("message text=%q, want 'test'", s.Messages[0].Text)
	}

	// Reopen and verify the message persisted correctly
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if len(s2.Messages) != 1 {
		t.Errorf("after reopen, got %d messages, want 1", len(s2.Messages))
	}
	if s2.Messages[0].Text != "test" {
		t.Errorf("persisted message text=%q, want 'test'", s2.Messages[0].Text)
	}
}

func TestNonZeroBasedSequences(t *testing.T) {
	dir := t.TempDir()

	// Manually write messages with non-contiguous sequences to the file
	f, err := os.OpenFile(filepath.Join(dir, transcriptFile), os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("create transcript: %v", err)
	}
	f.WriteString(`{"seq":5,"author":"user","role":"user","text":"five"}` + "\n")
	f.WriteString(`{"seq":10,"author":"claude","role":"assistant","text":"ten"}` + "\n")
	f.Close()

	// Reopen and verify nextSeq is set to last valid seq+1 = 11
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if len(s2.Messages) != 2 {
		t.Fatalf("got %d messages, want 2", len(s2.Messages))
	}

	// Append a new message and verify it gets seq 11
	if err := s2.Append(Message{Author: "user", Role: "user", Text: "eleven"}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	if s2.Messages[2].Seq != 11 {
		t.Errorf("new message got Seq=%d, want 11", s2.Messages[2].Seq)
	}
}

func TestNonMonotonicSequences(t *testing.T) {
	dir := t.TempDir()

	// Write messages in non-monotonic order: 10, then 5
	f, err := os.OpenFile(filepath.Join(dir, transcriptFile), os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("create transcript: %v", err)
	}
	f.WriteString(`{"seq":10,"author":"claude","role":"assistant","text":"ten"}` + "\n")
	f.WriteString(`{"seq":5,"author":"user","role":"user","text":"five"}` + "\n")
	f.Close()

	// Reopen: nextSeq should be last valid message's seq + 1 = 6, not max(seq) + 1 = 11
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if len(s.Messages) != 2 {
		t.Fatalf("got %d messages, want 2", len(s.Messages))
	}

	// The next appended message should get seq 6 (not 11)
	if err := s.Append(Message{Author: "user", Role: "user", Text: "six"}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	if s.Messages[2].Seq != 6 {
		t.Errorf("new message got Seq=%d, want 6 (last valid seq was 5)", s.Messages[2].Seq)
	}
}

func TestAgentSessionMemoryConsistency(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// Set agent session
	if err := s.SetAgentSession("claude", "sess-abc123"); err != nil {
		t.Fatalf("SetAgentSession: %v", err)
	}

	// Verify it's in memory
	if got := s.ResumeID("claude"); got != "sess-abc123" {
		t.Errorf("ResumeID in memory=%q, want %q", got, "sess-abc123")
	}

	// Reopen and verify it persisted
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if got := s2.ResumeID("claude"); got != "sess-abc123" {
		t.Errorf("ResumeID after reopen=%q, want %q", got, "sess-abc123")
	}
}

func TestLargeMessageHandling(t *testing.T) {
	dir := t.TempDir()

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// Create a large message (500KB) to test scanner buffer
	largeText := ""
	for i := 0; i < 10000; i++ {
		largeText += "This is a long message that repeats many times to test large message handling. "
	}

	// Append the large message
	if err := s.Append(Message{Author: "user", Role: "user", Text: largeText}); err != nil {
		t.Fatalf("Append large message: %v", err)
	}

	if len(s.Messages) != 1 {
		t.Fatalf("after append, got %d messages, want 1", len(s.Messages))
	}

	// Reopen and verify large message loaded correctly
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reOpen: %v", err)
	}

	if len(s2.Messages) != 1 {
		t.Fatalf("after reopen, got %d messages, want 1", len(s2.Messages))
	}

	if s2.Messages[0].Text != largeText {
		t.Error("large message text not preserved correctly")
	}
}
