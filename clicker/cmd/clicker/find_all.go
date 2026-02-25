package main

import (
	"github.com/spf13/cobra"
)

func newFindAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-all [selector]",
		Short: "Find all elements matching a CSS selector",
		Example: `  vibium find-all "a"
  # → @e1 [a] "Home"  @e2 [a] "About"  ...

  vibium find-all "a" --limit 5
  # Limit results to 5 elements`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			limit, _ := cmd.Flags().GetInt("limit")

			var selector string
			if len(args) == 2 {
				// find-all <url> <selector> — navigate first
				_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
				if err != nil {
					printError(err)
					return
				}
				selector = args[1]
			} else {
				selector = args[0]
			}

			toolArgs := map[string]interface{}{
				"selector": selector,
				"limit":    float64(limit),
			}

			result, err := daemonCall("browser_find_all", toolArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Int("limit", 10, "Maximum number of elements to return")
	return cmd
}
