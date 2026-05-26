package cmd

import (
	"fmt"
	"os"

	"github.com/archsinit/ralph-go/internal/config"
	"github.com/spf13/cobra"
)

var loopCmd = &cobra.Command{
	Use:   "loop",
	Short: "Run in loop mode",
	Long:  "Execute tasks in loop mode with continuous iteration.",
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
		fmt.Println("loop: not implemented")
	},
}
