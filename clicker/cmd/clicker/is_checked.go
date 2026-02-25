package main

import (
	"github.com/spf13/cobra"
)

func newIsCheckedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "is-checked [selector]",
		Short: "Check if a checkbox or radio is checked",
		Example: `  vibium is-checked "input[type=checkbox]"
  # Prints true or false`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			selector := args[0]

			result, err := daemonCall("browser_is_checked", map[string]interface{}{"selector": selector})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
