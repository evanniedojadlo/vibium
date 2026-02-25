package main

import (
	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/proxy"
)

func newClickCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "click [url] [selector]",
		Short: "Click an element (optionally navigate to URL first)",
		Example: `  vibium click "a"
  # Clicks on current page

  vibium click https://example.com "a"
  # Navigates to URL first, then clicks

  vibium click https://example.com "a" --timeout 5s
  # Custom timeout for actionability checks`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var selector string
			if len(args) == 2 {
				// click <url> <selector> — navigate first
				_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
				if err != nil {
					printError(err)
					return
				}
				selector = args[1]
			} else {
				// click <selector> — current page
				selector = args[0]
			}

			// Click element
			result, err := daemonCall("browser_click", map[string]interface{}{"selector": selector})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Duration("timeout", proxy.DefaultTimeout, "Timeout for actionability checks (e.g., 5s, 30s)")
	return cmd
}
