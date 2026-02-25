package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/process"
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
			process.WithCleanup(func() {
				url := args[0]
				selector := args[1]

				fmt.Println("Launching browser...")
				launchResult, err := browser.Launch(browser.LaunchOptions{Headless: headless})
				if err != nil {
					fatalExit("Error launching browser: %v", err)
				}
				defer waitAndClose(launchResult)

				fmt.Println("Connecting to BiDi...")
				conn, err := bidi.Connect(launchResult.WebSocketURL)
				if err != nil {
					fatalExit("Error connecting: %v", err)
				}
				defer conn.Close()

				client := bidi.NewClient(conn)

				fmt.Printf("Navigating to %s...\n", url)
				_, err = client.Navigate("", url)
				if err != nil {
					fatalExit("Error navigating: %v", err)
				}

				doWaitOpen()

				fmt.Printf("\nChecking actionability for selector: %s\n", selector)

				result, err := checkAllActionability(client, "", selector)
				if err != nil {
					fatalExit("Error: %v", err)
				}

				printCheck("Visible", result.Visible)
				printCheck("Stable", result.Stable)
				printCheck("ReceivesEvents", result.ReceivesEvents)
				printCheck("Enabled", result.Enabled)
				printCheck("Editable", result.Editable)
			})
		},
	}
}

type actionabilityResult struct {
	Visible        bool `json:"visible"`
	Stable         bool `json:"stable"`
	ReceivesEvents bool `json:"receivesEvents"`
	Enabled        bool `json:"enabled"`
	Editable       bool `json:"editable"`
}

func checkAllActionability(client *bidi.Client, context, selector string) (*actionabilityResult, error) {
	if context == "" {
		tree, err := client.GetTree()
		if err != nil {
			return nil, fmt.Errorf("failed to get browsing context: %w", err)
		}
		if len(tree.Contexts) == 0 {
			return nil, fmt.Errorf("no browsing contexts available")
		}
		context = tree.Contexts[0].Context
	}

	script := `
		(selector) => {
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
		}
	`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": context},
		"arguments": []map[string]interface{}{
			{"type": "string", "value": selector},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	msg, err := client.SendCommand("script.callFunction", params)
	if err != nil {
		return nil, err
	}

	var callResult struct {
		Type   string          `json:"type"`
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(msg.Result, &callResult); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}
	if callResult.Type == "exception" {
		return nil, fmt.Errorf("script exception: %s", string(callResult.Result))
	}

	var remoteValue struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(callResult.Result, &remoteValue); err != nil {
		return nil, fmt.Errorf("failed to parse remote value: %w", err)
	}

	var result actionabilityResult
	if err := json.Unmarshal([]byte(remoteValue.Value), &result); err != nil {
		return nil, fmt.Errorf("failed to parse actionability result: %w", err)
	}

	return &result, nil
}
