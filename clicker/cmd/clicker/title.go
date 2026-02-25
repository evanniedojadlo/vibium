package main

import (
	"github.com/spf13/cobra"
)

func newTitleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "title",
		Short: "Get the current page title",
		Example: `  vibium title
  # Prints: Example Domain`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_get_title", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
