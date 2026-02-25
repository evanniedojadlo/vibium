package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [url] [selector]",
		Short: "Find an element by CSS selector or semantic locator",
		Example: `  vibium find "a"
  # → @e1 [a] "More information..."

  vibium find https://example.com "a"
  # Navigate to URL first, then find

  vibium find --text "Sign In"
  # → @e1 [button] "Sign In"

  vibium find --label "Email"
  # → @e1 [input type="email"] placeholder="Email"

  vibium click @e1
  # Use the returned @ref to interact with the found element

  vibium find --role link
  # → @e1 [a] "More information..."

  vibium find --role heading --text "Example"
  # Find heading containing "Example"

  vibium find --placeholder "Search..."
  vibium find --testid "submit-btn"
  vibium find --xpath "//div[@class='main']"`,
		Args: cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Collect semantic flags
			semanticFlags := map[string]string{
				"role":        "",
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
		},
	}

	cmd.Flags().String("role", "", "Find element by ARIA role (e.g. button, link, textbox)")
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
