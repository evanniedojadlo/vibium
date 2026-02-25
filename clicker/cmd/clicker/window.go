package main

import (
	"github.com/spf13/cobra"
)

func newWindowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "window",
		Short: "Get the OS browser window dimensions and state",
		Example: `  vibium window
  # {"state":"normal","x":0,"y":25,"width":1280,"height":720}`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_get_window", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
