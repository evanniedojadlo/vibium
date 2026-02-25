package main

import (
	"github.com/spf13/cobra"
)

func newViewportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "viewport",
		Short: "Get the current viewport dimensions",
		Example: `  vibium viewport
  # {"width":1280,"height":720,"devicePixelRatio":1}`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_get_viewport", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
