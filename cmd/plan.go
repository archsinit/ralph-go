package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run in plan mode",
	Long:  "Execute tasks in plan mode with step-by-step orchestration.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan: not implemented")
	},
}
