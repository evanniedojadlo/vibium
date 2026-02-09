package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/process"
)

func newFindAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-all [selector]",
		Short: "Find all elements matching a CSS selector",
		Example: `  clicker find-all "a"
  # Find all links on current page (daemon mode)

  clicker find-all "a" --limit 5
  # Limit results to 5 elements`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			limit, _ := cmd.Flags().GetInt("limit")

			// Daemon mode
			if !oneshot {
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
				return
			}

			// Oneshot mode — requires URL + selector
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Error: requires [url] [selector] in oneshot mode\n")
				os.Exit(1)
			}
			url := args[0]
			selector := args[1]
			process.WithCleanup(func() {
				fmt.Println("Launching browser...")
				launchResult, err := browser.Launch(browser.LaunchOptions{Headless: headless})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error launching browser: %v\n", err)
					os.Exit(1)
				}
				defer waitAndClose(launchResult)

				conn, err := bidi.Connect(launchResult.WebSocketURL)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error connecting: %v\n", err)
					os.Exit(1)
				}
				defer conn.Close()

				client := bidi.NewClient(conn)

				_, err = client.Navigate("", url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error navigating: %v\n", err)
					os.Exit(1)
				}

				doWaitOpen()

				elements, err := client.FindAllElements("", selector, limit)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error finding elements: %v\n", err)
					os.Exit(1)
				}

				for i, el := range elements {
					fmt.Printf("[%d] tag=%s, text=\"%s\", box={x:%.0f, y:%.0f, w:%.0f, h:%.0f}\n",
						i, el.Tag, el.Text, el.Box.X, el.Box.Y, el.Box.Width, el.Box.Height)
				}
				if len(elements) == 0 {
					fmt.Println("No elements found")
				}
			})
		},
	}
	cmd.Flags().Int("limit", 10, "Maximum number of elements to return")
	return cmd
}
