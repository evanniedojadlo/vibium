package main

import (
	"github.com/spf13/cobra"
)

func newWaitForURLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait-for-url [pattern]",
		Short: "Wait until the page URL contains a substring",
		Example: `  vibium wait-for-url "/dashboard"
  # Wait until URL contains "/dashboard"

  vibium wait-for-url "success" --timeout 10000
  # Wait up to 10 seconds`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pattern := args[0]
			timeout, _ := cmd.Flags().GetInt("timeout")

			toolArgs := map[string]interface{}{"pattern": pattern}
			if cmd.Flags().Changed("timeout") {
				toolArgs["timeout"] = float64(timeout)
			}

			result, err := daemonCall("browser_wait_for_url", toolArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Int("timeout", 30000, "Timeout in milliseconds")
	return cmd
}
