package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/archsinit/ralph-go/internal/agent"
	"github.com/archsinit/ralph-go/internal/config"
	"github.com/archsinit/ralph-go/internal/orchestrator"
	"github.com/archsinit/ralph-go/internal/session"
	"github.com/archsinit/ralph-go/internal/tui"
	"github.com/spf13/cobra"
)

var (
	planSessionFlag string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run in plan mode",
	Long:  "Execute tasks in plan mode with step-by-step orchestration.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPlan(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	planCmd.Flags().StringVar(&planSessionFlag, "session", "", "Session directory to resume (default: create new)")
	// Note: planCmd is registered in root.go's init(), not here, to avoid duplicate registration
}

func runPlan() error {
	// Load and validate config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// Determine session directory
	sessionDir := planSessionFlag
	if sessionDir == "" {
		sessionDir = cfg.Paths.SessionDir
		if sessionDir == "" {
			// Use a default based on config path
			configDir := filepath.Dir(configPath)
			sessionDir = filepath.Join(configDir, "session")
		}
	}

	// Open session
	sess, err := session.Open(sessionDir)
	if err != nil {
		return fmt.Errorf("open session: %w", err)
	}

	// Build adapters from config
	adapters := make(map[string]agent.Adapter)
	for _, cfgAgent := range cfg.Agents {
		if !cfgAgent.Enabled {
			continue
		}
		adapter, err := agent.New(cfgAgent.CLI, cfgAgent)
		if err != nil {
			return fmt.Errorf("create adapter %q: %w", cfgAgent.Name, err)
		}
		adapters[cfgAgent.Name] = adapter
	}

	// Create turn order from config
	turnOrder := orchestrator.NewTurnOrder(cfg.Plan.TurnOrder)

	// Create TUI with initial messages
	initialMsgs := make([]tui.Message, len(sess.Messages))
	for i, m := range sess.Messages {
		initialMsgs[i] = tui.Message{
			Author: m.Author,
			Text:   m.Text,
		}
	}

	// Create bridge before TUI so we can pass submit and cancel callbacks
	bridge := orchestrator.NewTUIBridge(nil)

	// Run TUI with submit and cancel callbacks (returns immediately now)
	tuiHandle, err := tui.Run(
		tui.WithInitialMessages(initialMsgs),
		tui.WithSubmitCallback(bridge.GetSubmitCallback()),
		tui.WithCancelCallback(bridge.GetCancelCallback()),
	)
	if err != nil {
		return fmt.Errorf("run TUI: %w", err)
	}

	// Update bridge with TUI handle
	bridge.SetHandle(tuiHandle)

	// Create orchestrator engine
	engine := orchestrator.NewEngine(cfg, sess, turnOrder, adapters, bridge)

	// Run the engine in a goroutine so we can wait for both engine and TUI
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engineErr := make(chan error, 1)
	go func() {
		if err := engine.Run(ctx); err != nil && err != context.Canceled {
			engineErr <- fmt.Errorf("engine: %w", err)
		} else {
			engineErr <- nil
		}
	}()

	// Wait for engine to finish
	engineExitErr := <-engineErr

	// Always quit the TUI regardless of engine exit status
	tuiHandle.Quit()
	if err := tuiHandle.Wait(); err != nil {
		// Log TUI error but don't fail if engine already succeeded or quit normally
		if engineExitErr == nil || engineExitErr == context.Canceled {
			// Engine quit cleanly (via /quit or normal exit), TUI error is secondary
			return nil
		}
		return fmt.Errorf("TUI: %w", err)
	}

	// Report engine error if there was one (but not if it was normal context.Canceled)
	if engineExitErr != nil && engineExitErr != context.Canceled {
		return engineExitErr
	}

	return nil
}
