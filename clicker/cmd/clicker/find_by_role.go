package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newFindByRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-by-role",
		Short: "Find an element by ARIA role and accessible name",
		Example: `  vibium find-by-role --role button --name "Submit"
  # Find a button with text "Submit"

  vibium find-by-role --role link --name "Learn more"
  # Find a link containing "Learn more"

  vibium find-by-role --role textbox
  # Find the first textbox on the page`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			role, _ := cmd.Flags().GetString("role")
			name, _ := cmd.Flags().GetString("name")
			selector, _ := cmd.Flags().GetString("selector")
			timeout, _ := cmd.Flags().GetInt("timeout")

			if !oneshot {
				toolArgs := map[string]interface{}{}
				if role != "" {
					toolArgs["role"] = role
				}
				if name != "" {
					toolArgs["name"] = name
				}
				if selector != "" {
					toolArgs["selector"] = selector
				}
				if cmd.Flags().Changed("timeout") {
					toolArgs["timeout"] = float64(timeout)
				}

				result, err := daemonCall("browser_find_by_role", toolArgs)
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			fmt.Fprintf(os.Stderr, "Error: find-by-role command requires daemon mode\n")
			os.Exit(1)
		},
	}
	cmd.Flags().String("role", "", "ARIA role to match (e.g., button, link, textbox)")
	cmd.Flags().String("name", "", "Accessible name to match (substring)")
	cmd.Flags().String("selector", "", "Additional CSS selector to narrow results")
	cmd.Flags().Int("timeout", 30000, "Timeout in milliseconds")
	return cmd
}
