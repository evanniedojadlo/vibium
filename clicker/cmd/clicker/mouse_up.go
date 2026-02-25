package main

import (
	"github.com/spf13/cobra"
)

func newMouseUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mouse-up",
		Short: "Release a mouse button",
		Example: `  vibium mouse-up
  # Release left mouse button

  vibium mouse-up --button 2
  # Release right mouse button`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			button, _ := cmd.Flags().GetInt("button")

			result, err := daemonCall("browser_mouse_up", map[string]interface{}{"button": float64(button)})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Int("button", 0, "Mouse button (0=left, 1=middle, 2=right)")
	return cmd
}
