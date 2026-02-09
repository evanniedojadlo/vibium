package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/process"
)

func newFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find [url] [selector]",
		Short: "Find an element by CSS selector (optionally navigate to URL first)",
		Example: `  clicker find "a"
  # Finds on current page (daemon mode)

  clicker find https://example.com "a"
  # Navigates to URL first, then finds
  # Prints: tag=A, text="Learn more", box={x,y,w,h}`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Daemon mode
			if !oneshot {
				var selector string
				if len(args) == 2 {
					// find <url> <selector> — navigate first
					_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
					if err != nil {
						printError(err)
						return
					}
					selector = args[1]
				} else {
					// find <selector> — current page
					selector = args[0]
				}

				// Find element
				result, err := daemonCall("browser_find", map[string]interface{}{"selector": selector})
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			// Oneshot mode (original behavior) — requires URL + selector
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

				fmt.Printf("Finding element: %s\n", selector)
				info, err := client.FindElement("", selector)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error finding element: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Found: tag=%s, text=\"%s\", box={x:%.0f, y:%.0f, w:%.0f, h:%.0f}\n",
					info.Tag, info.Text, info.Box.X, info.Box.Y, info.Box.Width, info.Box.Height)
			})
		},
	}
}
