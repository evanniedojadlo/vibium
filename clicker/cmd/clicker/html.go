package main

import (
	"github.com/spf13/cobra"
)

func newHTMLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "html [selector]",
		Short: "Get HTML content of the page or an element",
		Example: `  vibium html
  # Get full page HTML

  vibium html "div.content"
  # Get innerHTML of a specific element

  vibium html "div.content" --outer
  # Get outerHTML of a specific element

  vibium html https://example.com "h1"
  # Navigate then get element HTML`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			outer, _ := cmd.Flags().GetBool("outer")

			toolArgs := map[string]interface{}{}
			if outer {
				toolArgs["outer"] = true
			}
			if len(args) == 2 {
				// html <url> <selector> — navigate first
				_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
				if err != nil {
					printError(err)
					return
				}
				toolArgs["selector"] = args[1]
			} else if len(args) == 1 {
				if isURL(args[0]) {
					// html <url> — navigate then get full page HTML
					_, err := daemonCall("browser_navigate", map[string]interface{}{"url": args[0]})
					if err != nil {
						printError(err)
						return
					}
				} else {
					// html <selector> — get element HTML on current page
					toolArgs["selector"] = args[0]
				}
			}

			result, err := daemonCall("browser_get_html", toolArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Bool("outer", false, "Return outerHTML instead of innerHTML")
	return cmd
}
