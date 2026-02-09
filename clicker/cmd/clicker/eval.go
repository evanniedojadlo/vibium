package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/process"
)

func newEvalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "eval [url] [expression]",
		Short: "Evaluate a JavaScript expression (optionally navigate to URL first)",
		Example: `  clicker eval "document.title"
  # Evaluates on current page (daemon mode)

  clicker eval https://example.com "document.title"
  # Navigates to URL first, then evaluates
  # Prints: Example Domain`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Daemon mode
			if !oneshot {
				var expression string
				if len(args) == 2 {
					// eval <url> <expression> — navigate first
					_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
					if err != nil {
						printError(err)
						return
					}
					expression = args[1]
				} else {
					// eval <expression> — current page
					expression = args[0]
				}

				// Evaluate
				result, err := daemonCall("browser_evaluate", map[string]interface{}{"expression": expression})
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			// Oneshot mode (original behavior) — requires URL + expression
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Error: requires [url] [expression] in oneshot mode\n")
				os.Exit(1)
			}
			url := args[0]
			expression := args[1]
			process.WithCleanup(func() {
				fmt.Println("Launching browser...")
				launchResult, err := browser.Launch(browser.LaunchOptions{Headless: headless})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error launching browser: %v\n", err)
					os.Exit(1)
				}
				defer waitAndClose(launchResult)

				fmt.Println("Connecting to BiDi...")
				conn, err := bidi.Connect(launchResult.WebSocketURL)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error connecting: %v\n", err)
					os.Exit(1)
				}
				defer conn.Close()

				client := bidi.NewClient(conn)

				fmt.Printf("Navigating to %s...\n", url)
				_, err = client.Navigate("", url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error navigating: %v\n", err)
					os.Exit(1)
				}

				doWaitOpen()

				fmt.Printf("Evaluating: %s\n", expression)
				result, err := client.Evaluate("", expression)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error evaluating: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Result: %v\n", result)
			})
		},
	}
}
