package mcp

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
	"github.com/vibium/clicker/internal/features"
	"github.com/vibium/clicker/internal/log"
)

// Handlers manages browser session state and executes tool calls.
type Handlers struct {
	launchResult  *browser.LaunchResult
	client        *bidi.Client
	conn          *bidi.Connection
	screenshotDir string
	headless      bool
}

// NewHandlers creates a new Handlers instance.
// screenshotDir specifies where screenshots are saved. If empty, file saving is disabled.
// headless controls whether the browser is launched in headless mode.
func NewHandlers(screenshotDir string, headless bool) *Handlers {
	return &Handlers{
		screenshotDir: screenshotDir,
		headless:      headless,
	}
}

// Call executes a tool by name with the given arguments.
func (h *Handlers) Call(name string, args map[string]interface{}) (*ToolsCallResult, error) {
	log.Debug("tool call", "name", name, "args", args)

	switch name {
	case "browser_launch":
		return h.browserLaunch(args)
	case "browser_navigate":
		return h.browserNavigate(args)
	case "browser_click":
		return h.browserClick(args)
	case "browser_type":
		return h.browserType(args)
	case "browser_screenshot":
		return h.browserScreenshot(args)
	case "browser_find":
		return h.browserFind(args)
	case "browser_evaluate":
		return h.browserEvaluate(args)
	case "browser_quit":
		return h.browserQuit(args)
	case "browser_get_text":
		return h.browserGetText(args)
	case "browser_get_url":
		return h.browserGetURL(args)
	case "browser_get_title":
		return h.browserGetTitle(args)
	case "browser_get_html":
		return h.browserGetHTML(args)
	case "browser_find_all":
		return h.browserFindAll(args)
	case "browser_wait":
		return h.browserWait(args)
	case "browser_hover":
		return h.browserHover(args)
	case "browser_select":
		return h.browserSelect(args)
	case "browser_scroll":
		return h.browserScroll(args)
	case "browser_keys":
		return h.browserKeys(args)
	case "browser_new_tab":
		return h.browserNewTab(args)
	case "browser_list_tabs":
		return h.browserListTabs(args)
	case "browser_switch_tab":
		return h.browserSwitchTab(args)
	case "browser_close_tab":
		return h.browserCloseTab(args)
	case "browser_a11y_tree":
		return h.browserA11yTree(args)
	case "page_clock_install":
		return h.pageClockInstall(args)
	case "page_clock_fast_forward":
		return h.pageClockFastForward(args)
	case "page_clock_run_for":
		return h.pageClockRunFor(args)
	case "page_clock_pause_at":
		return h.pageClockPauseAt(args)
	case "page_clock_resume":
		return h.pageClockResume(args)
	case "page_clock_set_fixed_time":
		return h.pageClockSetFixedTime(args)
	case "page_clock_set_system_time":
		return h.pageClockSetSystemTime(args)
	case "page_clock_set_timezone":
		return h.pageClockSetTimezone(args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// Close cleans up any active browser sessions.
func (h *Handlers) Close() {
	if h.conn != nil {
		h.conn.Close()
		h.conn = nil
	}
	if h.launchResult != nil {
		h.launchResult.Close()
		h.launchResult = nil
	}
	h.client = nil
}

// browserLaunch launches a new browser session.
func (h *Handlers) browserLaunch(args map[string]interface{}) (*ToolsCallResult, error) {
	// If browser is already running, return success (no-op)
	if h.client != nil {
		return &ToolsCallResult{
			Content: []Content{{
				Type: "text",
				Text: "Browser already running",
			}},
		}, nil
	}

	// Parse options — per-call headless overrides the default
	useHeadless := h.headless
	if val, ok := args["headless"].(bool); ok {
		useHeadless = val
	}

	// Launch browser
	launchResult, err := browser.Launch(browser.LaunchOptions{Headless: useHeadless})
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to BiDi
	conn, err := bidi.Connect(launchResult.WebSocketURL)
	if err != nil {
		launchResult.Close()
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	h.launchResult = launchResult
	h.conn = conn
	h.client = bidi.NewClient(conn)

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Browser launched (headless: %v)", useHeadless),
		}},
	}, nil
}

// browserNavigate navigates to a URL.
func (h *Handlers) browserNavigate(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	url, ok := args["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("url is required")
	}

	result, err := h.client.Navigate("", url)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Navigated to %s", result.URL),
		}},
	}, nil
}

// browserClick clicks an element.
func (h *Handlers) browserClick(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	// Wait for element to be actionable
	opts := features.DefaultWaitOptions()
	if err := features.WaitForClick(h.client, "", selector, opts); err != nil {
		return nil, err
	}

	// Click the element
	if err := h.client.ClickElement("", selector); err != nil {
		return nil, fmt.Errorf("failed to click: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Clicked element: %s", selector),
		}},
	}, nil
}

// browserType types text into an element.
func (h *Handlers) browserType(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	text, ok := args["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text is required")
	}

	// Wait for element to be actionable
	opts := features.DefaultWaitOptions()
	if err := features.WaitForType(h.client, "", selector, opts); err != nil {
		return nil, err
	}

	// Type into the element
	if err := h.client.TypeIntoElement("", selector, text); err != nil {
		return nil, fmt.Errorf("failed to type: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Typed into element: %s", selector),
		}},
	}, nil
}

// browserScreenshot captures a screenshot.
func (h *Handlers) browserScreenshot(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	base64Data, err := h.client.CaptureScreenshot("")
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// If filename provided, save to file (only if screenshotDir is configured)
	if filename, ok := args["filename"].(string); ok && filename != "" {
		if h.screenshotDir == "" {
			return nil, fmt.Errorf("screenshot file saving is disabled (use --screenshot-dir to enable)")
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(h.screenshotDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create screenshot directory: %w", err)
		}

		// Use only the basename to prevent path traversal
		safeName := filepath.Base(filename)
		fullPath := filepath.Join(h.screenshotDir, safeName)

		pngData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode screenshot: %w", err)
		}
		if err := os.WriteFile(fullPath, pngData, 0644); err != nil {
			return nil, fmt.Errorf("failed to save screenshot: %w", err)
		}
		return &ToolsCallResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Screenshot saved to %s", fullPath),
			}},
		}, nil
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type:     "image",
			Data:     base64Data,
			MimeType: "image/png",
		}},
	}, nil
}

// browserFind finds an element and returns its info.
func (h *Handlers) browserFind(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	info, err := h.client.FindElement("", selector)
	if err != nil {
		return nil, err
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("tag=%s, text=\"%s\", box={x:%.0f, y:%.0f, w:%.0f, h:%.0f}",
				info.Tag, info.Text, info.Box.X, info.Box.Y, info.Box.Width, info.Box.Height),
		}},
	}, nil
}

// browserEvaluate executes JavaScript code in the browser.
func (h *Handlers) browserEvaluate(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	expression, ok := args["expression"].(string)
	if !ok || expression == "" {
		return nil, fmt.Errorf("expression is required")
	}

	result, err := h.client.Evaluate("", expression)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate: %w", err)
	}

	// Format result as string
	var resultText string
	switch v := result.(type) {
	case string:
		resultText = v
	case nil:
		resultText = "null"
	default:
		resultText = fmt.Sprintf("%v", v)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: resultText,
		}},
	}, nil
}

// browserQuit closes the browser session.
func (h *Handlers) browserQuit(args map[string]interface{}) (*ToolsCallResult, error) {
	if h.launchResult == nil {
		return &ToolsCallResult{
			Content: []Content{{
				Type: "text",
				Text: "No browser session to close",
			}},
		}, nil
	}

	h.Close()

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: "Browser session closed",
		}},
	}, nil
}

// browserNewTab creates a new browser tab.
func (h *Handlers) browserNewTab(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	url, _ := args["url"].(string)

	contextID, err := h.client.CreateTab(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create tab: %w", err)
	}

	msg := "New tab opened"
	if url != "" {
		msg = fmt.Sprintf("New tab opened and navigated to %s", url)
	}
	_ = contextID

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: msg,
		}},
	}, nil
}

// browserListTabs lists all open browser tabs.
func (h *Handlers) browserListTabs(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	tree, err := h.client.GetTree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tabs: %w", err)
	}

	var text string
	for i, ctx := range tree.Contexts {
		text += fmt.Sprintf("[%d] %s\n", i, ctx.URL)
	}
	if text == "" {
		text = "No tabs open"
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: text,
		}},
	}, nil
}

// browserSwitchTab switches to a tab by index or URL substring.
func (h *Handlers) browserSwitchTab(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	tree, err := h.client.GetTree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tabs: %w", err)
	}

	var contextID string

	// Try index first
	if idx, ok := args["index"].(float64); ok {
		i := int(idx)
		if i < 0 || i >= len(tree.Contexts) {
			return nil, fmt.Errorf("tab index %d out of range (0-%d)", i, len(tree.Contexts)-1)
		}
		contextID = tree.Contexts[i].Context
	} else if url, ok := args["url"].(string); ok && url != "" {
		// Search by URL substring
		for _, ctx := range tree.Contexts {
			if containsSubstring(ctx.URL, url) {
				contextID = ctx.Context
				break
			}
		}
		if contextID == "" {
			return nil, fmt.Errorf("no tab matching URL %q", url)
		}
	} else {
		return nil, fmt.Errorf("index or url is required")
	}

	if err := h.client.ActivateTab(contextID); err != nil {
		return nil, err
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Switched to tab: %s", contextID),
		}},
	}, nil
}

// browserCloseTab closes a tab by index (default: current tab).
func (h *Handlers) browserCloseTab(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	tree, err := h.client.GetTree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tabs: %w", err)
	}

	if len(tree.Contexts) == 0 {
		return nil, fmt.Errorf("no tabs open")
	}

	idx := 0
	if i, ok := args["index"].(float64); ok {
		idx = int(i)
	}

	if idx < 0 || idx >= len(tree.Contexts) {
		return nil, fmt.Errorf("tab index %d out of range (0-%d)", idx, len(tree.Contexts)-1)
	}

	contextID := tree.Contexts[idx].Context
	if err := h.client.CloseTab(contextID); err != nil {
		return nil, err
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Closed tab %d", idx),
		}},
	}, nil
}

// browserA11yTree returns the accessibility tree of the current page.
func (h *Handlers) browserA11yTree(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	interestingOnly := true
	if val, ok := args["interestingOnly"].(bool); ok {
		interestingOnly = val
	}

	script := a11yTreeMCPScript()
	result, err := h.client.CallFunction("", script, []interface{}{interestingOnly})
	if err != nil {
		return nil, fmt.Errorf("failed to get accessibility tree: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("%v", result),
		}},
	}, nil
}

// a11yTreeMCPScript returns the JS function for the MCP a11y tree tool.
func a11yTreeMCPScript() string {
	return `(interestingOnly) => {
		const IMPLICIT_ROLES = {
			A: (el) => el.hasAttribute('href') ? 'link' : '',
			AREA: (el) => el.hasAttribute('href') ? 'link' : '',
			ARTICLE: () => 'article',
			ASIDE: () => 'complementary',
			BUTTON: () => 'button',
			DETAILS: () => 'group',
			DIALOG: () => 'dialog',
			FOOTER: () => 'contentinfo',
			FORM: () => 'form',
			H1: () => 'heading', H2: () => 'heading', H3: () => 'heading',
			H4: () => 'heading', H5: () => 'heading', H6: () => 'heading',
			HEADER: () => 'banner',
			HR: () => 'separator',
			IMG: (el) => el.getAttribute('alt') ? 'img' : 'presentation',
			INPUT: (el) => {
				const t = (el.getAttribute('type') || 'text').toLowerCase();
				const map = {button:'button',checkbox:'checkbox',image:'button',
					number:'spinbutton',radio:'radio',range:'slider',
					reset:'button',search:'searchbox',submit:'button',text:'textbox',
					email:'textbox',tel:'textbox',url:'textbox',password:'textbox'};
				return map[t] || 'textbox';
			},
			LI: () => 'listitem',
			MAIN: () => 'main',
			MENU: () => 'list',
			NAV: () => 'navigation',
			OL: () => 'list',
			OPTION: () => 'option',
			OUTPUT: () => 'status',
			PROGRESS: () => 'progressbar',
			SECTION: () => 'region',
			SELECT: (el) => el.hasAttribute('multiple') ? 'listbox' : 'combobox',
			SUMMARY: () => 'button',
			TABLE: () => 'table',
			TBODY: () => 'rowgroup', THEAD: () => 'rowgroup', TFOOT: () => 'rowgroup',
			TD: () => 'cell',
			TEXTAREA: () => 'textbox',
			TH: () => 'columnheader',
			TR: () => 'row',
			UL: () => 'list',
		};

		function getRole(el) {
			if (typeof el.computedRole === 'string' && el.computedRole !== '') return el.computedRole;
			const explicit = el.getAttribute('role');
			if (explicit) return explicit.toLowerCase();
			const fn = IMPLICIT_ROLES[el.tagName];
			return fn ? fn(el) : 'generic';
		}

		function getName(el) {
			if (typeof el.computedName === 'string') return el.computedName;
			const ariaLabel = el.getAttribute('aria-label');
			if (ariaLabel) return ariaLabel;
			const labelledBy = el.getAttribute('aria-labelledby');
			if (labelledBy) {
				const parts = labelledBy.split(/\\s+/).map(id => {
					const ref = document.getElementById(id);
					return ref ? (ref.textContent || '').trim() : '';
				}).filter(Boolean);
				if (parts.length) return parts.join(' ');
			}
			if (el.id) {
				const assocLabel = document.querySelector('label[for="' + el.id + '"]');
				if (assocLabel) return (assocLabel.textContent || '').trim();
			}
			const placeholder = el.getAttribute('placeholder');
			if (placeholder) return placeholder;
			const alt = el.getAttribute('alt');
			if (alt) return alt;
			const title = el.getAttribute('title');
			if (title) return title;
			return '';
		}

		function getChildren(el) {
			if (el.shadowRoot) return Array.from(el.shadowRoot.children);
			return Array.from(el.children);
		}

		function getHeadingLevel(el) {
			const tag = el.tagName;
			if (tag === 'H1') return 1;
			if (tag === 'H2') return 2;
			if (tag === 'H3') return 3;
			if (tag === 'H4') return 4;
			if (tag === 'H5') return 5;
			if (tag === 'H6') return 6;
			const level = el.getAttribute('aria-level');
			if (level) return parseInt(level, 10);
			return undefined;
		}

		function buildNode(el) {
			const role = getRole(el);
			const name = getName(el);
			const childNodes = [];
			for (const child of getChildren(el)) {
				if (child.nodeType !== 1) continue;
				const nodes = buildNode(child);
				if (nodes) {
					if (Array.isArray(nodes)) childNodes.push(...nodes);
					else childNodes.push(nodes);
				}
			}
			if (interestingOnly) {
				if (role === 'none' || role === 'presentation') return childNodes.length ? childNodes : null;
				if (role === 'generic' && !name) return childNodes.length ? childNodes : null;
			}
			const node = { role: role };
			if (name) node.name = name;
			if (el.hasAttribute('disabled') || el.disabled) node.disabled = true;
			if (el.hasAttribute('aria-expanded')) node.expanded = el.getAttribute('aria-expanded') === 'true';
			if (document.activeElement === el) node.focused = true;
			if (typeof el.checked === 'boolean' && (el.type === 'checkbox' || el.type === 'radio')) {
				node.checked = el.checked;
			} else if (el.hasAttribute('aria-checked')) {
				const v = el.getAttribute('aria-checked');
				node.checked = v === 'true' ? true : v === 'mixed' ? 'mixed' : false;
			}
			if (el.hasAttribute('aria-pressed')) {
				const v = el.getAttribute('aria-pressed');
				node.pressed = v === 'true' ? true : v === 'mixed' ? 'mixed' : false;
			}
			if (el.hasAttribute('aria-selected') && el.getAttribute('aria-selected') === 'true') node.selected = true;
			if (el.hasAttribute('required') || el.required) node.required = true;
			if (el.hasAttribute('readonly') || el.readOnly) node.readonly = true;
			const level = getHeadingLevel(el);
			if (level !== undefined) node.level = level;
			if (childNodes.length) node.children = childNodes;
			return node;
		}

		const children = [];
		for (const child of getChildren(document.body)) {
			if (child.nodeType !== 1) continue;
			const nodes = buildNode(child);
			if (nodes) {
				if (Array.isArray(nodes)) children.push(...nodes);
				else children.push(nodes);
			}
		}
		return JSON.stringify({ role: 'WebArea', name: document.title, children: children });
	}`
}

// containsSubstring checks if s contains substr (case-sensitive).
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && strings.Contains(s, substr)
}

// browserHover moves the mouse over an element.
func (h *Handlers) browserHover(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	info, err := h.client.FindElement("", selector)
	if err != nil {
		return nil, err
	}

	x, y := info.GetCenter()
	if err := h.client.MoveMouse("", x, y); err != nil {
		return nil, fmt.Errorf("failed to hover: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Hovered over element: %s", selector),
		}},
	}, nil
}

// browserSelect selects an option in a <select> element.
func (h *Handlers) browserSelect(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	value, ok := args["value"].(string)
	if !ok || value == "" {
		return nil, fmt.Errorf("value is required")
	}

	script := `(selector, value) => {
		const el = document.querySelector(selector);
		if (!el) return JSON.stringify({error: 'Element not found'});
		if (el.tagName.toLowerCase() !== 'select') return JSON.stringify({error: 'Element is not a <select>'});
		el.value = value;
		el.dispatchEvent(new Event('change', {bubbles: true}));
		return JSON.stringify({selected: el.value});
	}`

	result, err := h.client.CallFunction("", script, []interface{}{selector, value})
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Selected value %q in %s (result: %v)", value, selector, result),
		}},
	}, nil
}

// browserScroll scrolls the page or an element.
func (h *Handlers) browserScroll(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	direction := "down"
	if d, ok := args["direction"].(string); ok && d != "" {
		direction = d
	}

	amount := 3
	if a, ok := args["amount"].(float64); ok {
		amount = int(a)
	}

	// Determine scroll target coordinates
	x, y := 0, 0
	if selector, ok := args["selector"].(string); ok && selector != "" {
		info, err := h.client.FindElement("", selector)
		if err != nil {
			return nil, err
		}
		cx, cy := info.GetCenter()
		x, y = int(cx), int(cy)
	} else {
		// Viewport center
		result, err := h.client.Evaluate("", "JSON.stringify({w: window.innerWidth, h: window.innerHeight})")
		if err == nil && result != nil {
			x, y = 400, 300 // Reasonable fallback
		}
	}

	// Map direction to deltas (120 pixels per scroll "notch")
	deltaX, deltaY := 0, 0
	pixels := amount * 120
	switch direction {
	case "down":
		deltaY = pixels
	case "up":
		deltaY = -pixels
	case "right":
		deltaX = pixels
	case "left":
		deltaX = -pixels
	default:
		return nil, fmt.Errorf("invalid direction: %q (use up, down, left, right)", direction)
	}

	if err := h.client.ScrollWheel("", x, y, deltaX, deltaY); err != nil {
		return nil, fmt.Errorf("failed to scroll: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Scrolled %s by %d", direction, amount),
		}},
	}, nil
}

// browserKeys presses a key or key combination.
func (h *Handlers) browserKeys(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	keys, ok := args["keys"].(string)
	if !ok || keys == "" {
		return nil, fmt.Errorf("keys is required")
	}

	if err := h.client.PressKeyCombo("", keys); err != nil {
		return nil, fmt.Errorf("failed to press keys: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Pressed keys: %s", keys),
		}},
	}, nil
}

// browserGetHTML returns the HTML content of the page or an element.
func (h *Handlers) browserGetHTML(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	outer, _ := args["outer"].(bool)

	var expr string
	if selector, ok := args["selector"].(string); ok && selector != "" {
		if outer {
			expr = fmt.Sprintf(`document.querySelector(%q)?.outerHTML || ''`, selector)
		} else {
			expr = fmt.Sprintf(`document.querySelector(%q)?.innerHTML || ''`, selector)
		}
	} else {
		expr = `document.documentElement.outerHTML`
	}

	result, err := h.client.Evaluate("", expr)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML: %w", err)
	}

	html := ""
	if result != nil {
		html = fmt.Sprintf("%v", result)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: html,
		}},
	}, nil
}

// browserFindAll finds all elements matching a CSS selector.
func (h *Handlers) browserFindAll(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	elements, err := h.client.FindAllElements("", selector, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find elements: %w", err)
	}

	var text string
	for i, el := range elements {
		text += fmt.Sprintf("[%d] tag=%s, text=\"%s\", box={x:%.0f, y:%.0f, w:%.0f, h:%.0f}\n",
			i, el.Tag, el.Text, el.Box.X, el.Box.Y, el.Box.Width, el.Box.Height)
	}
	if text == "" {
		text = "No elements found"
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: text,
		}},
	}, nil
}

// browserWait waits for an element to reach a specified state.
func (h *Handlers) browserWait(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return nil, fmt.Errorf("selector is required")
	}

	state := "attached"
	if s, ok := args["state"].(string); ok && s != "" {
		state = s
	}

	timeout := features.DefaultTimeout
	if t, ok := args["timeout"].(float64); ok {
		timeout = time.Duration(t) * time.Millisecond
	}

	opts := features.WaitOptions{Timeout: timeout}

	switch state {
	case "attached":
		if err := features.WaitForSelector(h.client, "", selector, opts); err != nil {
			return nil, err
		}
	case "visible":
		if err := features.WaitForSelector(h.client, "", selector, opts); err != nil {
			return nil, err
		}
		visible, err := features.CheckVisible(h.client, "", selector)
		if err != nil {
			return nil, fmt.Errorf("visibility check failed: %w", err)
		}
		if !visible {
			return nil, fmt.Errorf("element %q found but not visible", selector)
		}
	case "hidden":
		if err := features.WaitForHidden(h.client, "", selector, opts); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid state: %q (use \"attached\", \"visible\", or \"hidden\")", state)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Element %q reached state: %s", selector, state),
		}},
	}, nil
}

// browserGetText returns the text content of the page or an element.
func (h *Handlers) browserGetText(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	var expr string
	if selector, ok := args["selector"].(string); ok && selector != "" {
		expr = fmt.Sprintf(`document.querySelector(%q)?.innerText || ''`, selector)
	} else {
		expr = `document.body.innerText`
	}

	result, err := h.client.Evaluate("", expr)
	if err != nil {
		return nil, fmt.Errorf("failed to get text: %w", err)
	}

	text := ""
	if result != nil {
		text = fmt.Sprintf("%v", result)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: text,
		}},
	}, nil
}

// browserGetURL returns the current page URL.
func (h *Handlers) browserGetURL(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	result, err := h.client.Evaluate("", "window.location.href")
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	url := ""
	if result != nil {
		url = fmt.Sprintf("%v", result)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: url,
		}},
	}, nil
}

// browserGetTitle returns the current page title.
func (h *Handlers) browserGetTitle(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	result, err := h.client.Evaluate("", "document.title")
	if err != nil {
		return nil, fmt.Errorf("failed to get title: %w", err)
	}

	title := ""
	if result != nil {
		title = fmt.Sprintf("%v", result)
	}

	return &ToolsCallResult{
		Content: []Content{{
			Type: "text",
			Text: title,
		}},
	}, nil
}

// pageClockInstall installs a fake clock on the page.
func (h *Handlers) pageClockInstall(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	_, err := h.client.CallFunction("", clockInstallScript(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to install clock: %w", err)
	}

	if timeVal, ok := args["time"].(float64); ok {
		script := fmt.Sprintf("() => { window.__vibiumClock.setSystemTime(%v); return 'ok'; }", timeVal)
		if _, err := h.client.CallFunction("", script, nil); err != nil {
			return nil, fmt.Errorf("failed to set initial time: %w", err)
		}
	}

	if tz, ok := args["timezone"].(string); ok && tz != "" {
		if err := h.setTimezoneOverride(tz); err != nil {
			return nil, fmt.Errorf("failed to set timezone: %w", err)
		}
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: "Clock installed"}},
	}, nil
}

// pageClockFastForward fast-forwards the fake clock.
func (h *Handlers) pageClockFastForward(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	ticks, ok := args["ticks"].(float64)
	if !ok {
		return nil, fmt.Errorf("ticks is required")
	}

	script := fmt.Sprintf("() => { window.__vibiumClock.fastForward(%v); return 'ok'; }", ticks)
	if _, err := h.client.CallFunction("", script, nil); err != nil {
		return nil, fmt.Errorf("clock.fastForward failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Fast-forwarded %v ms", ticks)}},
	}, nil
}

// pageClockRunFor advances the fake clock, firing all callbacks.
func (h *Handlers) pageClockRunFor(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	ticks, ok := args["ticks"].(float64)
	if !ok {
		return nil, fmt.Errorf("ticks is required")
	}

	script := fmt.Sprintf("() => { window.__vibiumClock.runFor(%v); return 'ok'; }", ticks)
	if _, err := h.client.CallFunction("", script, nil); err != nil {
		return nil, fmt.Errorf("clock.runFor failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Ran for %v ms", ticks)}},
	}, nil
}

// pageClockPauseAt pauses the fake clock at a specific time.
func (h *Handlers) pageClockPauseAt(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	timeVal, ok := args["time"].(float64)
	if !ok {
		return nil, fmt.Errorf("time is required")
	}

	script := fmt.Sprintf("() => { window.__vibiumClock.pauseAt(%v); return 'ok'; }", timeVal)
	if _, err := h.client.CallFunction("", script, nil); err != nil {
		return nil, fmt.Errorf("clock.pauseAt failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Paused at %v", timeVal)}},
	}, nil
}

// pageClockResume resumes real-time progression.
func (h *Handlers) pageClockResume(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	if _, err := h.client.CallFunction("", "() => { window.__vibiumClock.resume(); return 'ok'; }", nil); err != nil {
		return nil, fmt.Errorf("clock.resume failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: "Clock resumed"}},
	}, nil
}

// pageClockSetFixedTime freezes Date.now() at a value.
func (h *Handlers) pageClockSetFixedTime(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	timeVal, ok := args["time"].(float64)
	if !ok {
		return nil, fmt.Errorf("time is required")
	}

	script := fmt.Sprintf("() => { window.__vibiumClock.setFixedTime(%v); return 'ok'; }", timeVal)
	if _, err := h.client.CallFunction("", script, nil); err != nil {
		return nil, fmt.Errorf("clock.setFixedTime failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Fixed time set to %v", timeVal)}},
	}, nil
}

// pageClockSetSystemTime sets Date.now() without triggering timers.
func (h *Handlers) pageClockSetSystemTime(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	timeVal, ok := args["time"].(float64)
	if !ok {
		return nil, fmt.Errorf("time is required")
	}

	script := fmt.Sprintf("() => { window.__vibiumClock.setSystemTime(%v); return 'ok'; }", timeVal)
	if _, err := h.client.CallFunction("", script, nil); err != nil {
		return nil, fmt.Errorf("clock.setSystemTime failed: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("System time set to %v", timeVal)}},
	}, nil
}

// pageClockSetTimezone overrides or resets the browser timezone.
func (h *Handlers) pageClockSetTimezone(args map[string]interface{}) (*ToolsCallResult, error) {
	if err := h.ensureBrowser(); err != nil {
		return nil, err
	}

	tz, _ := args["timezone"].(string)

	if tz == "" {
		// Reset to default
		if err := h.clearTimezoneOverride(); err != nil {
			return nil, fmt.Errorf("failed to clear timezone: %w", err)
		}
		return &ToolsCallResult{
			Content: []Content{{Type: "text", Text: "Timezone reset to system default"}},
		}, nil
	}

	if err := h.setTimezoneOverride(tz); err != nil {
		return nil, fmt.Errorf("failed to set timezone: %w", err)
	}

	return &ToolsCallResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Timezone set to %s", tz)}},
	}, nil
}

// setTimezoneOverride uses BiDi emulation.setTimezoneOverride.
func (h *Handlers) setTimezoneOverride(timezone string) error {
	tree, err := h.client.GetTree()
	if err != nil {
		return fmt.Errorf("failed to get browsing context: %w", err)
	}
	if len(tree.Contexts) == 0 {
		return fmt.Errorf("no browsing contexts available")
	}

	_, err = h.client.SendCommand("emulation.setTimezoneOverride", map[string]interface{}{
		"timezone": timezone,
		"contexts": []interface{}{tree.Contexts[0].Context},
	})
	return err
}

// clearTimezoneOverride resets the browser timezone to the system default.
func (h *Handlers) clearTimezoneOverride() error {
	tree, err := h.client.GetTree()
	if err != nil {
		return fmt.Errorf("failed to get browsing context: %w", err)
	}
	if len(tree.Contexts) == 0 {
		return fmt.Errorf("no browsing contexts available")
	}

	_, err = h.client.SendCommand("emulation.setTimezoneOverride", map[string]interface{}{
		"timezone": nil,
		"contexts": []interface{}{tree.Contexts[0].Context},
	})
	return err
}

// clockInstallScript returns the JS that installs the fake clock on the page.
// This is the same script used by the proxy handlers (defined separately to avoid circular imports).
func clockInstallScript() string {
	return `() => {
	if (window.__vibiumClock) return 'already_installed';

	const OrigDate = Date;
	const origSetTimeout = setTimeout;
	const origClearTimeout = clearTimeout;
	const origSetInterval = setInterval;
	const origClearInterval = clearInterval;
	const origRAF = requestAnimationFrame;
	const origCAF = cancelAnimationFrame;
	const origPerfNow = performance.now.bind(performance);

	let currentTime = OrigDate.now();
	let fixedTime = null;
	let paused = false;
	let nextId = 1;
	let resumeTimer = null;
	const timers = new Map();

	class FakeDate extends OrigDate {
		constructor(...args) {
			if (args.length === 0) {
				super(fixedTime !== null ? fixedTime : currentTime);
			} else {
				super(...args);
			}
		}
		static now() {
			return fixedTime !== null ? fixedTime : currentTime;
		}
		static parse(s) { return OrigDate.parse(s); }
		static UTC(...args) { return OrigDate.UTC(...args); }
	}

	function fakeSetTimeout(fn, delay, ...args) {
		if (typeof fn !== 'function') return 0;
		const id = nextId++;
		timers.set(id, {
			callback: fn, args: args,
			triggerTime: currentTime + (delay || 0),
			interval: 0, type: 'timeout'
		});
		return id;
	}

	function fakeSetInterval(fn, delay, ...args) {
		if (typeof fn !== 'function') return 0;
		const id = nextId++;
		timers.set(id, {
			callback: fn, args: args,
			triggerTime: currentTime + (delay || 0),
			interval: delay || 0, type: 'interval'
		});
		return id;
	}

	function fakeClearTimeout(id) { timers.delete(id); }
	function fakeClearInterval(id) { timers.delete(id); }

	let rafId = 1;
	const rafCallbacks = new Map();
	function fakeRAF(fn) { const id = rafId++; rafCallbacks.set(id, fn); return id; }
	function fakeCAF(id) { rafCallbacks.delete(id); }

	const startPerfTime = origPerfNow();
	const startCurrentTime = currentTime;
	function fakePerfNow() { return startPerfTime + (currentTime - startCurrentTime); }

	window.Date = FakeDate;
	window.setTimeout = fakeSetTimeout;
	window.setInterval = fakeSetInterval;
	window.clearTimeout = fakeClearTimeout;
	window.clearInterval = fakeClearInterval;
	window.requestAnimationFrame = fakeRAF;
	window.cancelAnimationFrame = fakeCAF;
	performance.now = fakePerfNow;

	function getDueTimers(upTo) {
		const due = [];
		for (const [id, t] of timers) {
			if (t.triggerTime <= upTo) due.push([id, t]);
		}
		due.sort((a, b) => a[1].triggerTime - b[1].triggerTime);
		return due;
	}

	function fireRAFs() {
		const cbs = Array.from(rafCallbacks.entries());
		rafCallbacks.clear();
		for (const [, fn] of cbs) { try { fn(currentTime); } catch (e) {} }
	}

	const clock = {
		fastForward(ms) {
			const target = currentTime + ms;
			currentTime = target;
			const due = getDueTimers(target);
			for (const [id, t] of due) {
				timers.delete(id);
				try { t.callback(...t.args); } catch (e) {}
			}
			fireRAFs();
		},
		runFor(ms) {
			const target = currentTime + ms;
			while (currentTime < target) {
				let earliest = null;
				let earliestId = null;
				for (const [id, t] of timers) {
					if (t.triggerTime <= target && (!earliest || t.triggerTime < earliest.triggerTime)) {
						earliest = t; earliestId = id;
					}
				}
				if (!earliest || earliest.triggerTime > target) { currentTime = target; break; }
				currentTime = earliest.triggerTime;
				if (earliest.type === 'interval' && earliest.interval > 0) {
					earliest.triggerTime = currentTime + earliest.interval;
				} else { timers.delete(earliestId); }
				try { earliest.callback(...earliest.args); } catch (e) {}
			}
			currentTime = target;
			fireRAFs();
		},
		pauseAt(time) {
			currentTime = time; paused = true;
			if (resumeTimer) { origClearInterval(resumeTimer); resumeTimer = null; }
			const due = getDueTimers(time);
			for (const [id, t] of due) { timers.delete(id); try { t.callback(...t.args); } catch (e) {} }
		},
		resume() {
			if (resumeTimer) return;
			paused = false;
			let lastReal = OrigDate.now();
			resumeTimer = origSetInterval(() => {
				const now = OrigDate.now();
				const delta = now - lastReal;
				lastReal = now;
				currentTime += delta;
				const due = getDueTimers(currentTime);
				for (const [id, t] of due) {
					if (t.type === 'interval' && t.interval > 0) { t.triggerTime = currentTime + t.interval; }
					else { timers.delete(id); }
					try { t.callback(...t.args); } catch (e) {}
				}
				fireRAFs();
			}, 16);
		},
		setFixedTime(time) { fixedTime = time; },
		setSystemTime(time) { currentTime = time; fixedTime = null; }
	};

	window.__vibiumClock = clock;
	return 'installed';
}`
}

// ensureBrowser checks that a browser session is active.
// If no browser is running, it auto-launches one (lazy launch).
func (h *Handlers) ensureBrowser() error {
	if h.client == nil {
		_, err := h.browserLaunch(map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("auto-launch failed: %w", err)
		}
	}
	return nil
}
