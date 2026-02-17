package mcp

// GetToolSchemas returns the list of available MCP tools with their schemas.
func GetToolSchemas() []Tool {
	return []Tool{
		{
			Name:        "browser_launch",
			Description: "Launch a new browser session",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"headless": map[string]interface{}{
						"type":        "boolean",
						"description": "Run browser in headless mode (no visible window)",
						"default":     false,
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_navigate",
			Description: "Navigate to a URL in the browser",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to navigate to",
					},
				},
				"required":             []string{"url"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_click",
			Description: "Click an element by CSS selector. Waits for element to be visible, stable, and enabled.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the element to click",
					},
				},
				"required":             []string{"selector"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_type",
			Description: "Type text into an element by CSS selector. Waits for element to be visible, stable, enabled, and editable.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the element to type into",
					},
					"text": map[string]interface{}{
						"type":        "string",
						"description": "The text to type",
					},
				},
				"required":             []string{"selector", "text"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_screenshot",
			Description: "Capture a screenshot of the current page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Optional filename to save the screenshot (e.g., screenshot.png)",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_find",
			Description: "Find an element by CSS selector and return its info (tag, text, bounding box)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the element to find",
					},
				},
				"required":             []string{"selector"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_evaluate",
			Description: "Execute JavaScript in the browser to extract data, query the DOM, or inspect page state. Returns the evaluated result. Use this to get text content, attributes, element data, or any information from the page.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "JavaScript expression to evaluate",
					},
				},
				"required":             []string{"expression"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_quit",
			Description: "Close the browser session",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_get_html",
			Description: "Get the HTML content of the page or a specific element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for a specific element (optional, defaults to full page HTML)",
					},
					"outer": map[string]interface{}{
						"type":        "boolean",
						"description": "Return outerHTML instead of innerHTML (default: false)",
						"default":     false,
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_find_all",
			Description: "Find all elements matching a CSS selector and return their info (tag, text, bounding box)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector to match elements",
					},
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of elements to return (default: 10)",
						"default":     10,
					},
				},
				"required":             []string{"selector"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_wait",
			Description: "Wait for an element to reach a specified state (attached, visible, or hidden)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the element to wait for",
					},
					"state": map[string]interface{}{
						"type":        "string",
						"description": "State to wait for: \"attached\" (exists in DOM), \"visible\" (visible on page), or \"hidden\" (not found or not visible)",
						"enum":        []string{"attached", "visible", "hidden"},
						"default":     "attached",
					},
					"timeout": map[string]interface{}{
						"type":        "number",
						"description": "Timeout in milliseconds (default: 30000)",
						"default":     30000,
					},
				},
				"required":             []string{"selector"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_new_tab",
			Description: "Open a new browser tab, optionally navigating to a URL",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL to navigate to in the new tab (optional)",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_list_tabs",
			Description: "List all open browser tabs with their URLs",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_switch_tab",
			Description: "Switch to a browser tab by index or URL substring",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "number",
						"description": "Tab index (0-based) from browser_list_tabs",
					},
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL substring to match (alternative to index)",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_close_tab",
			Description: "Close a browser tab by index (default: current tab)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "number",
						"description": "Tab index to close (default: 0, the current tab)",
						"default":     0,
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_hover",
			Description: "Hover over an element by CSS selector",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the element to hover over",
					},
				},
				"required":             []string{"selector"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_select",
			Description: "Select an option in a <select> element by value",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the <select> element",
					},
					"value": map[string]interface{}{
						"type":        "string",
						"description": "The value to select",
					},
				},
				"required":             []string{"selector", "value"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_scroll",
			Description: "Scroll the page or a specific element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"direction": map[string]interface{}{
						"type":        "string",
						"description": "Scroll direction: up, down, left, right (default: down)",
						"enum":        []string{"up", "down", "left", "right"},
						"default":     "down",
					},
					"amount": map[string]interface{}{
						"type":        "number",
						"description": "Number of scroll increments (default: 3)",
						"default":     3,
					},
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for element to scroll to (optional, defaults to viewport center)",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_keys",
			Description: "Press a key or key combination (e.g., \"Enter\", \"Control+a\", \"Shift+Tab\")",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keys": map[string]interface{}{
						"type":        "string",
						"description": "Key or key combination to press (e.g., \"Enter\", \"Control+a\", \"Shift+ArrowDown\")",
					},
				},
				"required":             []string{"keys"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_get_text",
			Description: "Get the text content of the page or a specific element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for a specific element (optional, defaults to full page text)",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_get_url",
			Description: "Get the current page URL",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_get_title",
			Description: "Get the current page title",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			},
		},
		{
			Name:        "browser_a11y_tree",
			Description: "Get the accessibility tree of the current page. Returns a tree of ARIA roles, names, and states — useful for understanding page structure without visual rendering.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"interestingOnly": map[string]interface{}{
						"type":        "boolean",
						"description": "Filter out generic/presentation nodes that are not interesting for accessibility. Default: true",
						"default":     true,
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_install",
			Description: "Install a fake clock on the page, overriding Date, setTimeout, setInterval, requestAnimationFrame, and performance.now",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"time": map[string]interface{}{
						"type":        "number",
						"description": "Initial time as epoch milliseconds (optional)",
					},
					"timezone": map[string]interface{}{
						"type":        "string",
						"description": "IANA timezone ID to override (e.g. 'America/New_York', 'Europe/London')",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_fast_forward",
			Description: "Jump the fake clock forward by N milliseconds, firing each due timer at most once",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"ticks": map[string]interface{}{
						"type":        "number",
						"description": "Number of milliseconds to fast-forward",
					},
				},
				"required":             []string{"ticks"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_run_for",
			Description: "Advance the fake clock by N milliseconds, firing all time-related callbacks systematically",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"ticks": map[string]interface{}{
						"type":        "number",
						"description": "Number of milliseconds to advance",
					},
				},
				"required":             []string{"ticks"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_pause_at",
			Description: "Jump the fake clock to a specific time and pause — no timers fire until resumed or advanced",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"time": map[string]interface{}{
						"type":        "number",
						"description": "Time as epoch milliseconds to pause at",
					},
				},
				"required":             []string{"time"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_resume",
			Description: "Resume real-time progression from the current fake clock time",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_set_fixed_time",
			Description: "Freeze Date.now() at a specific value permanently. Timers still run.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"time": map[string]interface{}{
						"type":        "number",
						"description": "Time as epoch milliseconds to freeze at",
					},
				},
				"required":             []string{"time"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_set_system_time",
			Description: "Set Date.now() to a specific value without triggering any timers",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"time": map[string]interface{}{
						"type":        "number",
						"description": "Time as epoch milliseconds to set",
					},
				},
				"required":             []string{"time"},
				"additionalProperties": false,
			},
		},
		{
			Name:        "page_clock_set_timezone",
			Description: "Override the browser timezone. Pass an IANA timezone ID (e.g. 'America/New_York'), or empty string to reset to system default",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"timezone": map[string]interface{}{
						"type":        "string",
						"description": "IANA timezone ID (e.g. 'America/New_York', 'Europe/London', 'Asia/Tokyo'). Empty string resets to system default.",
					},
				},
				"required":             []string{"timezone"},
				"additionalProperties": false,
			},
		},
	}
}
