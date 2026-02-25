package main

import (
	"github.com/spf13/cobra"
)

func newDblClickCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dblclick [selector]",
		Short: "Double-click an element",
		Example: `  vibium dblclick "td.cell"
  # Double-click to edit a table cell

  vibium dblclick @e2
  # Double-click element from map`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			selector := args[0]

			result, err := daemonCall("browser_dblclick", map[string]interface{}{"selector": selector})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
