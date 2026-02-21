package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newMapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "map",
		Short: "Map interactive page elements with @refs",
		Example: `  vibium map
  # Lists interactive elements with refs like @e1, @e2
  # Use refs with other commands: vibium click @e1`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !oneshot {
				result, err := daemonCall("browser_map", map[string]interface{}{})
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			fmt.Fprintf(os.Stderr, "Error: map command requires daemon mode\n")
			os.Exit(1)
		},
	}
}
