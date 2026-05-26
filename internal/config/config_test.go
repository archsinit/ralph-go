package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.toml")

	content := `
[[agent]]
name = "claude"
cli = "claude"
system_prompt = "You are Claude."
enabled = true

[plan]
turn_order = ["claude", "user"]
plan_agent = "claude"

[loop]
executor = "claude"
reviewer = "claude"
max_retries = 2
task_timeout = "15m"

[loop.ntfy]
server = "https://ntfy.sh"
topic = "test"
token = ""

[loop.git]
commit_prefix = ""

[paths]
session_dir = ".ralph/sessions"
log_dir = ".ralph/logs"
out_dir = "."
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatalf("config is nil")
	}
	if len(cfg.Agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(cfg.Agents))
	}
	if cfg.Agents[0].Name != "claude" {
		t.Errorf("expected agent name claude, got %s", cfg.Agents[0].Name)
	}
	if cfg.Plan.PlanAgent != "claude" {
		t.Errorf("expected plan_agent claude, got %s", cfg.Plan.PlanAgent)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	path := "/nonexistent/path/to/file.toml"
	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !contains(err.Error(), path) {
		t.Errorf("error should mention path %q, got: %v", path, err)
	}
}

func TestLoad_MalformedTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.toml")
	content := `
[[agent]]
name = "test"
cli = "unknown"
invalid syntax here [[[
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", SystemPrompt: "test", Enabled: true},
					{Name: "codex", CLI: "codex", SystemPrompt: "test", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude", "codex", "user"},
					PlanAgent: "claude",
				},
				Loop: LoopConfig{
					Executor:    "claude",
					Reviewer:    "codex",
					MaxRetries:  2,
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Server: "https://ntfy.sh",
						Topic:  "test",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate agent names",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
					{Name: "claude", CLI: "codex", Enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "empty agent name",
			config: &Config{
				Agents: []Agent{
					{Name: "", CLI: "claude", Enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown CLI",
			config: &Config{
				Agents: []Agent{
					{Name: "test", CLI: "unknown-cli", Enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "turn_order references missing agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude", "nonexistent", "user"},
				},
			},
			wantErr: true,
		},
		{
			name: "turn_order references disabled agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
					{Name: "codex", CLI: "codex", Enabled: false},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude", "codex", "user"},
				},
			},
			wantErr: true,
		},
		{
			name: "plan_agent references missing agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					PlanAgent: "nonexistent",
				},
			},
			wantErr: true,
		},
		{
			name: "plan_agent references disabled agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: false},
				},
				Plan: PlanConfig{
					PlanAgent: "claude",
				},
			},
			wantErr: true,
		},
		{
			name: "executor references missing agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Loop: LoopConfig{
					Executor:    "nonexistent",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Topic: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "executor references disabled agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: false},
				},
				Loop: LoopConfig{
					Executor:    "claude",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Topic: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "reviewer references missing agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Loop: LoopConfig{
					Reviewer:    "nonexistent",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Topic: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing plan_agent",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude"},
					PlanAgent: "",
				},
				Loop: LoopConfig{
					Executor:    "claude",
					Reviewer:    "claude",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy:        NtfyConfig{Topic: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty turn_order",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{},
					PlanAgent: "claude",
				},
				Loop: LoopConfig{
					Executor:    "claude",
					Reviewer:    "claude",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy:        NtfyConfig{Topic: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing loop executor",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude"},
					PlanAgent: "claude",
				},
				Loop: LoopConfig{
					Executor:    "",
					Reviewer:    "claude",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy:        NtfyConfig{Topic: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing loop reviewer",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Plan: PlanConfig{
					TurnOrder: []string{"claude"},
					PlanAgent: "claude",
				},
				Loop: LoopConfig{
					Executor:    "claude",
					Reviewer:    "",
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy:        NtfyConfig{Topic: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "negative max_retries",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Loop: LoopConfig{
					MaxRetries:  -1,
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Topic: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "zero task_timeout",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Loop: LoopConfig{
					TaskTimeout: Duration(0),
					Ntfy: NtfyConfig{
						Topic: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty ntfy topic",
			config: &Config{
				Agents: []Agent{
					{Name: "claude", CLI: "claude", Enabled: true},
				},
				Loop: LoopConfig{
					TaskTimeout: Duration(15 * time.Minute),
					Ntfy: NtfyConfig{
						Topic: "",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
