package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newEvalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval [url] [expression]",
		Short: "Evaluate a JavaScript expression (optionally navigate to URL first)",
		Example: `  vibium eval "document.title"
  # Evaluates on current page

  vibium eval https://example.com "document.title"
  # Navigates to URL first, then evaluates

  echo 'document.title' | vibium eval --stdin
  # Read expression from stdin (avoids shell quoting issues)`,
		Args: cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			useStdin, _ := cmd.Flags().GetBool("stdin")

			var expression string

			if useStdin {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					printError(fmt.Errorf("failed to read stdin: %w", err))
					return
				}
				expression = strings.TrimSpace(string(data))
			} else if len(args) == 2 {
				// eval <url> <expression> — navigate first
				_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
				if err != nil {
					printError(err)
					return
				}
				expression = args[1]
			} else if len(args) == 1 {
				// eval <expression> — current page
				expression = args[0]
			} else {
				fmt.Fprintf(os.Stderr, "Error: expression is required (use args or --stdin)\n")
				os.Exit(1)
			}

			// Evaluate
			result, err := daemonCall("browser_evaluate", map[string]interface{}{"expression": expression})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Bool("stdin", false, "Read expression from stdin")
	return cmd
}
