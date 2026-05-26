package session

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	transcriptFile = "transcript.jsonl"
	agentsFile     = "agents.json"
)

// Message is one turn in the chatroom transcript.
type Message struct {
	Seq    int       `json:"seq"`
	Author string    `json:"author"`
	Role   string    `json:"role"`
	Text   string    `json:"text"`
	TS     time.Time `json:"ts"`
}

// Session holds the in-memory chatroom state for one ralph session directory.
type Session struct {
	// Dir is the on-disk directory that backs this session.
	Dir string

	// Messages is the ordered transcript.
	Messages []Message

	// nextSeq is the sequence number to assign to the next appended message.
	nextSeq int

	// AgentSessions maps agent name to the CLI session ID returned by the agent adapter.
	AgentSessions map[string]string
}

// Open opens (or creates) the session directory at dir and loads any existing state.
func Open(dir string) (*Session, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	s := &Session{
		Dir:           dir,
		AgentSessions: make(map[string]string),
	}

	if err := s.loadMessages(); err != nil {
		return nil, err
	}
	if err := s.loadAgents(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Session) loadMessages() error {
	path := filepath.Join(s.Dir, transcriptFile)
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var m Message
		if err := json.Unmarshal(line, &m); err != nil {
			// Check if there are more lines after this one.
			if scanner.Scan() {
				// Non-trailing corrupt line is a hard error.
				return fmt.Errorf("corrupt line %d in transcript: %w", lineNum, err)
			}
			// Trailing corrupt line is tolerated (likely truncated on crash).
			break
		}
		s.Messages = append(s.Messages, m)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Track the next sequence number based on the highest sequence in loaded messages.
	s.nextSeq = 0
	for _, m := range s.Messages {
		if m.Seq >= s.nextSeq {
			s.nextSeq = m.Seq + 1
		}
	}

	return nil
}

// ResumeID returns the CLI session ID stored for the named agent, or "" if none.
func (s *Session) ResumeID(agent string) string {
	return s.AgentSessions[agent]
}

// Append assigns a sequence number, timestamps the message, persists it to
// transcript.jsonl, and flushes to disk before returning.
func (s *Session) Append(m Message) error {
	// Assign sequence number and timestamp before writing to disk.
	m.Seq = s.nextSeq
	if m.TS.IsZero() {
		m.TS = time.Now().UTC()
	}

	path := filepath.Join(s.Dir, transcriptFile)
	isNewFile := s.nextSeq == 0 && len(s.Messages) == 0

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(m); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	// fsync the directory after creating a new transcript file for durability.
	if isNewFile {
		if d, err := os.Open(s.Dir); err == nil {
			defer d.Close()
			syscall.Fsync(int(d.Fd()))
		}
	}

	// Only update in-memory state after successful write and sync to disk.
	s.Messages = append(s.Messages, m)
	s.nextSeq++

	return nil
}

// SetAgentSession stores the CLI session ID for an agent and rewrites agents.json atomically.
func (s *Session) SetAgentSession(agent, id string) error {
	s.AgentSessions[agent] = id

	data, err := json.Marshal(s.AgentSessions)
	if err != nil {
		return err
	}

	dir := s.Dir
	tmp, err := os.CreateTemp(dir, "agents-*.json.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	target := filepath.Join(dir, agentsFile)
	if err := os.Rename(tmpName, target); err != nil {
		os.Remove(tmpName)
		return err
	}

	// fsync the directory to ensure the rename is durable (where supported).
	if d, err := os.Open(dir); err == nil {
		defer d.Close()
		syscall.Fsync(int(d.Fd()))
	}

	return nil
}

func (s *Session) loadAgents() error {
	path := filepath.Join(s.Dir, agentsFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.AgentSessions)
}
