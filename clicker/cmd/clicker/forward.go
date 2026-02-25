package main

import (
	"github.com/spf13/cobra"
)

func newForwardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "forward",
		Short: "Navigate forward in browser history",
		Example: `  vibium forward
  # Go forward one page (like clicking the forward button)`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_forward", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
