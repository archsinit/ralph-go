package agent

import (
	"fmt"

	"github.com/archsinit/ralph-go/internal/config"
)

// New returns the Adapter for the given CLI name.
// Returns an error for unknown CLI names.
func New(cli string, _ config.Agent) (Adapter, error) {
	switch cli {
	case "claude":
		return NewClaudeAdapter(), nil
	case "codex":
		return NewCodexAdapter(), nil
	case "gemini":
		return NewGeminiAdapter(), nil
	case "opencode":
		return NewOpencodeAdapter(), nil
	case "echo":
		return NewEchoAdapter(), nil
	default:
		return nil, fmt.Errorf("unknown agent CLI %q", cli)
	}
}
