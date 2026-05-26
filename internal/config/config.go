package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Duration wraps time.Duration to parse human-readable strings from TOML (e.g., "15m").
type Duration time.Duration

// UnmarshalText parses a duration string like "15m", "1h", "30s".
func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	*d = Duration(parsed)
	return nil
}

// Config is the top-level TOML configuration for ralph-go.
type Config struct {
	Agents []Agent    `toml:"agent"`
	Plan   PlanConfig `toml:"plan"`
	Loop   LoopConfig `toml:"loop"`
	Paths  Paths      `toml:"paths"`
}

// Agent defines a named LLM agent (claude, gemini, etc).
type Agent struct {
	Name         string `toml:"name"`
	CLI          string `toml:"cli"`
	SystemPrompt string `toml:"system_prompt"`
	Enabled      bool   `toml:"enabled"`
}

// PlanConfig holds settings for plan mode.
type PlanConfig struct {
	TurnOrder []string `toml:"turn_order"`
	PlanAgent string   `toml:"plan_agent"`
}

// LoopConfig holds settings for loop mode.
type LoopConfig struct {
	Executor    string     `toml:"executor"`
	Reviewer    string     `toml:"reviewer"`
	MaxRetries  int        `toml:"max_retries"`
	TaskTimeout Duration   `toml:"task_timeout"`
	Git         GitConfig  `toml:"git"`
	Ntfy        NtfyConfig `toml:"ntfy"`
}

// GitConfig holds git integration settings.
type GitConfig struct {
	CommitPrefix string `toml:"commit_prefix"`
}

// NtfyConfig holds ntfy.sh notification settings.
type NtfyConfig struct {
	Server string `toml:"server"`
	Topic  string `toml:"topic"`
	Token  string `toml:"token"`
}

// Paths holds directory paths for sessions, logs, and output.
type Paths struct {
	SessionDir string `toml:"session_dir"`
	LogDir     string `toml:"log_dir"`
	OutDir     string `toml:"out_dir"`
}

// Known CLI agent types.
var KnownCLIs = map[string]bool{
	"claude":   true,
	"codex":    true,
	"gemini":   true,
	"opencode": true,
}

// Load reads and decodes a TOML config file into Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", path, err)
	}

	return &cfg, nil
}

// Validate checks semantic correctness of the config and returns aggregated errors.
func (c *Config) Validate() error {
	var errs []string

	// Check agent constraints
	seenNames := make(map[string]bool)
	agentsByName := make(map[string]*Agent)
	for i, agent := range c.Agents {
		if agent.Name == "" {
			errs = append(errs, fmt.Sprintf("agent[%d]: name is empty", i))
		} else if seenNames[agent.Name] {
			errs = append(errs, fmt.Sprintf("agent[%d]: duplicate name %q", i, agent.Name))
		} else {
			seenNames[agent.Name] = true
			agentsByName[agent.Name] = &c.Agents[i]
		}

		if !KnownCLIs[agent.CLI] {
			errs = append(errs, fmt.Sprintf("agent[%d] %q: CLI %q not in known set (claude, codex, gemini, opencode)", i, agent.Name, agent.CLI))
		}
	}

	// Check plan config
	if c.Plan.PlanAgent == "" {
		errs = append(errs, "plan.plan_agent: must be set")
	} else {
		agent, exists := agentsByName[c.Plan.PlanAgent]
		if !exists {
			errs = append(errs, fmt.Sprintf("plan.plan_agent: agent %q not found", c.Plan.PlanAgent))
		} else if !agent.Enabled {
			errs = append(errs, fmt.Sprintf("plan.plan_agent: agent %q is not enabled", c.Plan.PlanAgent))
		}
	}

	if len(c.Plan.TurnOrder) == 0 {
		errs = append(errs, "plan.turn_order: must contain at least one entry")
	}

	for _, name := range c.Plan.TurnOrder {
		if name == "user" {
			continue // literal "user" is allowed
		}
		agent, exists := agentsByName[name]
		if !exists {
			errs = append(errs, fmt.Sprintf("plan.turn_order: agent %q not found", name))
		} else if !agent.Enabled {
			errs = append(errs, fmt.Sprintf("plan.turn_order: agent %q is not enabled", name))
		}
	}

	// Check loop config
	if c.Loop.Executor == "" {
		errs = append(errs, "loop.executor: must be set")
	} else {
		agent, exists := agentsByName[c.Loop.Executor]
		if !exists {
			errs = append(errs, fmt.Sprintf("loop.executor: agent %q not found", c.Loop.Executor))
		} else if !agent.Enabled {
			errs = append(errs, fmt.Sprintf("loop.executor: agent %q is not enabled", c.Loop.Executor))
		}
	}

	if c.Loop.Reviewer == "" {
		errs = append(errs, "loop.reviewer: must be set")
	} else {
		agent, exists := agentsByName[c.Loop.Reviewer]
		if !exists {
			errs = append(errs, fmt.Sprintf("loop.reviewer: agent %q not found", c.Loop.Reviewer))
		} else if !agent.Enabled {
			errs = append(errs, fmt.Sprintf("loop.reviewer: agent %q is not enabled", c.Loop.Reviewer))
		}
	}

	if c.Loop.MaxRetries < 0 {
		errs = append(errs, fmt.Sprintf("loop.max_retries: must be >= 0, got %d", c.Loop.MaxRetries))
	}

	if c.Loop.TaskTimeout <= 0 {
		errs = append(errs, fmt.Sprintf("loop.task_timeout: must be > 0, got %v", time.Duration(c.Loop.TaskTimeout)))
	}

	if c.Loop.Ntfy.Topic == "" {
		errs = append(errs, "loop.ntfy.topic: must be non-empty")
	}

	if len(errs) > 0 {
		var msg string
		for i, err := range errs {
			if i > 0 {
				msg += "\n"
			}
			msg += err
		}
		return fmt.Errorf("config validation failed:\n%s", msg)
	}

	return nil
}
