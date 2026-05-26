package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loopCmd = &cobra.Command{
	Use:   "loop",
	Short: "Run in loop mode",
	Long:  "Execute tasks in loop mode with continuous iteration.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("loop: not implemented")
	},
}
