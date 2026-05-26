package cmd

import (
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "ralph-go",
	Short: "ralph-go: multi-agent orchestration",
	Long:  "A tool for orchestrating multi-agent workflows with plan and loop modes.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./ralph.toml", "Path to config file")
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(loopCmd)
}
