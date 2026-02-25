package main

import (
	"github.com/spf13/cobra"
)

func newCookiesCmd() *cobra.Command {
	cookiesCmd := &cobra.Command{
		Use:   "cookies",
		Short: "Manage browser cookies",
		Example: `  vibium cookies
  # List all cookies`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_get_cookies", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}

	setCmd := &cobra.Command{
		Use:   "set [name] [value]",
		Short: "Set a cookie",
		Example: `  vibium cookies set "session" "abc123"
  # Set a cookie with name and value`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_set_cookie", map[string]interface{}{
				"name":  args[0],
				"value": args[1],
			})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cookies",
		Example: `  vibium cookies clear
  # Delete all cookies`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := daemonCall("browser_delete_cookies", map[string]interface{}{})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}

	cookiesCmd.AddCommand(setCmd)
	cookiesCmd.AddCommand(clearCmd)
	return cookiesCmd
}
