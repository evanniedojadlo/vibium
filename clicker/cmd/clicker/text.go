package main

import (
	"github.com/spf13/cobra"
)

func newTextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "text [selector]",
		Short: "Get text content of the page or an element",
		Example: `  vibium text
  # Get all page text

  vibium text "h1"
  # Get text of a specific element

  vibium text https://example.com
  # Navigate then get all page text

  vibium text https://example.com "h1"
  # Navigate then get element text`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			toolArgs := map[string]interface{}{}
			if len(args) == 2 {
				// text <url> <selector> — navigate first
				_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
				if err != nil {
					printError(err)
					return
				}
				toolArgs["selector"] = args[1]
			} else if len(args) == 1 {
				if isURL(args[0]) {
					// text <url> — navigate then get all page text
					_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
					if err != nil {
						printError(err)
						return
					}
				} else {
					// text <selector> — get element text on current page
					toolArgs["selector"] = args[0]
				}
			}

			result, err := daemonCall("browser_get_text", toolArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
}
