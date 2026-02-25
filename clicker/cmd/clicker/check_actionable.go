package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newCheckActionableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-actionable [url] [selector]",
		Short: "Check actionability of an element (Visible, Stable, ReceivesEvents, Enabled, Editable)",
		Example: `  vibium check-actionable https://example.com "a"
  # Output:
  # Checking actionability for selector: a
  # ✓ Visible: true
  # ✓ Stable: true
  # ✓ ReceivesEvents: true
  # ✓ Enabled: true
  # ✗ Editable: false`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			selector := args[1]

			// Navigate to URL
			_, err := daemonCall("browser_navigate", map[string]interface{}{"url": url})
			if err != nil {
				printError(err)
				return
			}

			fmt.Printf("\nChecking actionability for selector: %s\n", selector)

			// Evaluate actionability script
			script := `(() => {
				const selector = ` + fmt.Sprintf("%q", selector) + `;
				const el = document.querySelector(selector);
				if (!el) return JSON.stringify({ error: 'element not found' });

				const rect = el.getBoundingClientRect();
				const style = window.getComputedStyle(el);
				const visible = rect.width > 0 && rect.height > 0 &&
					style.visibility !== 'hidden' && style.display !== 'none';

				const cx = rect.x + rect.width/2, cy = rect.y + rect.height/2;
				const hit = document.elementFromPoint(cx, cy);
				const receivesEvents = hit && (el === hit || el.contains(hit));

				let enabled = true;
				if (el.disabled === true) enabled = false;
				else if (el.getAttribute('aria-disabled') === 'true') enabled = false;
				else {
					const fs = el.closest('fieldset[disabled]');
					if (fs) { const legend = fs.querySelector('legend'); if (!legend || !legend.contains(el)) enabled = false; }
				}

				let editable = enabled && !el.readOnly && el.getAttribute('aria-readonly') !== 'true';
				if (editable) {
					const tag = el.tagName.toLowerCase();
					if (tag === 'input') {
						const t = (el.type || 'text').toLowerCase();
						editable = ['text','password','email','number','search','tel','url'].includes(t);
					} else if (tag !== 'textarea' && !el.isContentEditable) {
						editable = false;
					}
				}

				return JSON.stringify({ visible, stable: true, receivesEvents, enabled, editable });
			})()`

			result, err := daemonCall("browser_evaluate", map[string]interface{}{"expression": script})
			if err != nil {
				printError(err)
				return
			}

			// Parse the result — daemon returns the evaluated value as text
			resultText := ""
			if result != nil {
				for _, c := range result.Content {
					if c.Type == "text" {
						resultText = c.Text
						break
					}
				}
			}

			var actionResult struct {
				Visible        bool   `json:"visible"`
				Stable         bool   `json:"stable"`
				ReceivesEvents bool   `json:"receivesEvents"`
				Enabled        bool   `json:"enabled"`
				Editable       bool   `json:"editable"`
				Error          string `json:"error"`
			}
			if err := json.Unmarshal([]byte(resultText), &actionResult); err != nil {
				printError(fmt.Errorf("failed to parse actionability result: %w", err))
				return
			}
			if actionResult.Error != "" {
				printError(fmt.Errorf("%s", actionResult.Error))
				return
			}

			printCheck("Visible", actionResult.Visible)
			printCheck("Stable", actionResult.Stable)
			printCheck("ReceivesEvents", actionResult.ReceivesEvents)
			printCheck("Enabled", actionResult.Enabled)
			printCheck("Editable", actionResult.Editable)
		},
	}
}
