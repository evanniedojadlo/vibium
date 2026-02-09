package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/features"
	"github.com/vibium/clicker/internal/process"
)

func newClickCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "click [url] [selector]",
		Short: "Click an element (optionally navigate to URL first)",
		Example: `  clicker click "a"
  # Clicks on current page (daemon mode)

  clicker click https://example.com "a"
  # Navigates to URL first, then clicks

  clicker click https://example.com "a" --timeout 5s
  # Custom timeout for actionability checks`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Daemon mode
			if !oneshot {
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
				timeout, _ := cmd.Flags().GetDuration("timeout")

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

				// Wait for element to be actionable (Visible, Stable, ReceivesEvents, Enabled)
				fmt.Printf("Waiting for element to be actionable: %s\n", selector)
				opts := features.WaitOptions{Timeout: timeout}
				if err := features.WaitForClick(client, "", selector, opts); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Clicking element: %s\n", selector)
				err = client.ClickElement("", selector)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error clicking: %v\n", err)
					os.Exit(1)
				}

				// TODO: Replace sleep with proper navigation wait (poll URL change or listen for BiDi events)
				fmt.Println("Waiting for navigation...")
				time.Sleep(1 * time.Second)

				// Get current URL after click
				currentURL, err := client.GetCurrentURL()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting URL: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Click complete! Current URL: %s\n", currentURL)
			})
		},
	}
	cmd.Flags().Duration("timeout", features.DefaultTimeout, "Timeout for actionability checks (e.g., 5s, 30s)")
	return cmd
}
