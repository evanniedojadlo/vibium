package main

import (
	"github.com/spf13/cobra"
)

func newNavigateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go [url]",
		Short: "Go to a URL and print page info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]

			result, err := daemonCall("browser_navigate", map[string]interface{}{"url": url})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
