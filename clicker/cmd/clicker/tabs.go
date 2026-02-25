package main

import (
	"github.com/spf13/cobra"
)

func newTabsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tabs",
		Short: "List all open browser tabs",
		Example: `  vibium tabs
  # [0] https://example.com
  # [1] https://google.com`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_list_tabs", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
