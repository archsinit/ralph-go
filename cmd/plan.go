package cmd

import (
	"fmt"
	"os"

	"github.com/archsinit/ralph-go/internal/config"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run in plan mode",
	Long:  "Execute tasks in plan mode with step-by-step orchestration.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := cfg.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		_ = cfg // TODO: use config
		fmt.Println("plan: not implemented")
	},
}
