package main

import (
	"github.com/spf13/cobra"
)

func newURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "url",
		Short: "Get the current page URL",
		Example: `  vibium url
  # Prints: https://example.com`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_get_url", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
