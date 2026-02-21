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
	cmd := &cobra.Command{
		Use:   "find [url] [selector]",
		Short: "Find an element by CSS selector or semantic locator",
		Example: `  vibium find "a"
  # Find by CSS selector on current page

  vibium find https://example.com "a"
  # Navigate to URL first, then find

  vibium find --text "Sign In"
  # Find element containing text "Sign In"

  vibium find --label "Email"
  # Find input associated with label "Email"

  vibium find --placeholder "Search..."
  # Find element with placeholder

  vibium find --testid "submit-btn"
  # Find element by data-testid attribute

  vibium find --xpath "//div[@class='main']"
  # Find element by XPath expression`,
		Args: cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Collect semantic flags
			semanticFlags := map[string]string{
				"text":        "",
				"label":       "",
				"placeholder": "",
				"testid":      "",
				"xpath":       "",
				"alt":         "",
				"title":       "",
			}
			hasSemantic := false
			for key := range semanticFlags {
				val, _ := cmd.Flags().GetString(key)
				if val != "" {
					semanticFlags[key] = val
					hasSemantic = true
				}
			}

			// Daemon mode
			if !oneshot {
				toolArgs := map[string]interface{}{}

				if hasSemantic {
					// Semantic find — no positional selector required
					// But allow optional URL as first positional arg
					if len(args) >= 1 {
						// Check if first arg looks like a URL
						if isURL(args[0]) {
							_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
							if err != nil {
								printError(err)
								return
							}
						}
					}
					for key, val := range semanticFlags {
						if val != "" {
							toolArgs[key] = val
						}
					}
				} else {
					// CSS selector find (original behavior)
					if len(args) == 0 {
						fmt.Fprintf(os.Stderr, "Error: requires a CSS selector or semantic flag (--text, --label, etc.)\n")
						os.Exit(1)
					}
					if len(args) == 2 {
						_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
						if err != nil {
							printError(err)
							return
						}
						toolArgs["selector"] = args[1]
					} else {
						toolArgs["selector"] = args[0]
					}
				}

				result, err := daemonCall("browser_find", toolArgs)
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

	cmd.Flags().String("text", "", "Find element containing this text")
	cmd.Flags().String("label", "", "Find input by associated label text")
	cmd.Flags().String("placeholder", "", "Find element by placeholder attribute")
	cmd.Flags().String("testid", "", "Find element by data-testid attribute")
	cmd.Flags().String("xpath", "", "Find element by XPath expression")
	cmd.Flags().String("alt", "", "Find element by alt attribute")
	cmd.Flags().String("title", "", "Find element by title attribute")

	return cmd
}

// isURL returns true if the string looks like a URL (starts with http:// or https://).
func isURL(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}
