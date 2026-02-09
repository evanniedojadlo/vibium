package bidi

import (
	"encoding/json"
	"fmt"

	errs "github.com/vibium/clicker/internal/errors"
)

// ElementInfo contains information about a found element.
type ElementInfo struct {
	SharedID string  `json:"sharedId"`
	Tag      string  `json:"tag"`
	Text     string  `json:"text"`
	Box      BoxInfo `json:"box"`
}

// BoxInfo contains bounding box coordinates.
type BoxInfo struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// FindElement finds an element by CSS selector and returns its info.
// If context is empty, it uses the first available context.
func (c *Client) FindElement(context, selector string) (*ElementInfo, error) {
	// If no context provided, get the first one from the tree
	if context == "" {
		tree, err := c.GetTree()
		if err != nil {
			return nil, fmt.Errorf("failed to get browsing context: %w", err)
		}
		if len(tree.Contexts) == 0 {
			return nil, fmt.Errorf("no browsing contexts available")
		}
		context = tree.Contexts[0].Context
	}

	// JavaScript to find element and extract info as JSON string
	// We return a JSON string to avoid BiDi's complex object serialization
	script := `
		(selector) => {
			const el = document.querySelector(selector);
			if (!el) return null;
			const rect = el.getBoundingClientRect();
			return JSON.stringify({
				tag: el.tagName.toLowerCase(),
				text: (el.textContent || '').trim().substring(0, 100),
				box: {
					x: rect.x,
					y: rect.y,
					width: rect.width,
					height: rect.height
				}
			});
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

	msg, err := c.SendCommand("script.callFunction", params)
	if err != nil {
		return nil, err
	}

	// Parse the result
	var callResult struct {
		Type   string          `json:"type"`
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(msg.Result, &callResult); err != nil {
		return nil, fmt.Errorf("failed to parse script.callFunction result: %w", err)
	}

	if callResult.Type == "exception" {
		return nil, fmt.Errorf("script exception: %s", string(callResult.Result))
	}

	// Parse the remote value (string containing JSON)
	var remoteValue struct {
		Type  string `json:"type"`
		Value string `json:"value,omitempty"`
	}

	if err := json.Unmarshal(callResult.Result, &remoteValue); err != nil {
		return nil, fmt.Errorf("failed to parse remote value: %w", err)
	}

	// Check if element was found
	if remoteValue.Type == "null" {
		return nil, &errs.ElementNotFoundError{Selector: selector, Context: context}
	}

	// Parse the JSON string value
	var info ElementInfo
	if err := json.Unmarshal([]byte(remoteValue.Value), &info); err != nil {
		return nil, fmt.Errorf("failed to parse element info: %w", err)
	}

	return &info, nil
}

// GetElementCenter returns the center coordinates of an element's bounding box.
func (info *ElementInfo) GetCenter() (float64, float64) {
	return info.Box.X + info.Box.Width/2, info.Box.Y + info.Box.Height/2
}

// FindAllElements finds all elements matching a CSS selector and returns their info.
// If context is empty, it uses the first available context.
// limit caps the number of results returned (0 = no limit).
func (c *Client) FindAllElements(context, selector string, limit int) ([]ElementInfo, error) {
	if context == "" {
		tree, err := c.GetTree()
		if err != nil {
			return nil, fmt.Errorf("failed to get browsing context: %w", err)
		}
		if len(tree.Contexts) == 0 {
			return nil, fmt.Errorf("no browsing contexts available")
		}
		context = tree.Contexts[0].Context
	}

	script := `
		(selector, limit) => {
			const els = document.querySelectorAll(selector);
			const results = [];
			const max = limit > 0 ? Math.min(els.length, limit) : els.length;
			for (let i = 0; i < max; i++) {
				const el = els[i];
				const rect = el.getBoundingClientRect();
				results.push({
					tag: el.tagName.toLowerCase(),
					text: (el.textContent || '').trim().substring(0, 100),
					box: {
						x: rect.x,
						y: rect.y,
						width: rect.width,
						height: rect.height
					}
				});
			}
			return JSON.stringify(results);
		}
	`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": context},
		"arguments": []map[string]interface{}{
			{"type": "string", "value": selector},
			{"type": "number", "value": limit},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	msg, err := c.SendCommand("script.callFunction", params)
	if err != nil {
		return nil, err
	}

	var callResult struct {
		Type   string          `json:"type"`
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(msg.Result, &callResult); err != nil {
		return nil, fmt.Errorf("failed to parse script.callFunction result: %w", err)
	}

	if callResult.Type == "exception" {
		return nil, fmt.Errorf("script exception: %s", string(callResult.Result))
	}

	var remoteValue struct {
		Type  string `json:"type"`
		Value string `json:"value,omitempty"`
	}

	if err := json.Unmarshal(callResult.Result, &remoteValue); err != nil {
		return nil, fmt.Errorf("failed to parse remote value: %w", err)
	}

	var elements []ElementInfo
	if err := json.Unmarshal([]byte(remoteValue.Value), &elements); err != nil {
		return nil, fmt.Errorf("failed to parse elements: %w", err)
	}

	return elements, nil
}
