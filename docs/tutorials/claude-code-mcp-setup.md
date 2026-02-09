# Setting Up Vibium MCP in Claude Code

This tutorial covers how to manage MCP (Model Context Protocol) servers in Claude Code, specifically for Vibium browser automation.

## Prerequisites

- Claude Code CLI installed:
  - macOS: `brew install --cask claude-code`
  - npm: `npm install -g @anthropic/claude-code`
- Vibium installed (`npm install vibium`) or clicker binary built locally

## List Installed MCP Servers

To see all currently configured MCP servers:

```bash
claude mcp list
```

Example output (after adding vibium):
```
Checking MCP server health...

vibium: npx -y vibium mcp - ✓ Connected
```

If no MCP servers are configured, the list will be empty.

## Add Vibium MCP

### Option 1: Using npx (Recommended)

```bash
claude mcp add vibium -- npx -y vibium mcp
```

### Option 2: Using Local Binary

If you built clicker locally:

```bash
claude mcp add vibium -- /path/to/clicker mcp
```

For example, from the vibium repo root:

```bash
claude mcp add vibium -- ./clicker/bin/clicker mcp
```

### Option 2b: Custom Screenshot Directory

By default, screenshots are saved to:
- macOS: `~/Pictures/Vibium/`
- Linux: `~/Pictures/Vibium/`
- Windows: `%USERPROFILE%\Pictures\Vibium\`

To use a different directory:

```bash
claude mcp add vibium -- ./clicker/bin/clicker mcp --screenshot-dir ./screenshots
```

To disable file saving (base64 inline only):

```bash
claude mcp add vibium -- ./clicker/bin/clicker mcp --screenshot-dir ""
```

### Option 3: Using Absolute Path

```bash
claude mcp add vibium -- $HOME/Projects/vibium/clicker/bin/clicker mcp
```

### Verify Installation

After adding, verify it appears in the list:

```bash
claude mcp list
```

You should see:
```
vibium: /path/to/clicker mcp
```

## Remove Vibium MCP

To remove the Vibium MCP server:

```bash
claude mcp remove vibium
```

Verify removal:

```bash
claude mcp list
```

Vibium should no longer appear in the output.

## How Claude Discovers MCP Tools

When Claude Code starts, it connects to each configured MCP server and performs a discovery handshake:

**Step 1: Initialize** - Establish the connection and exchange capabilities

```
→ {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"capabilities": {}}}
← {"jsonrpc": "2.0", "id": 1, "result": {"capabilities": {"tools": {}}, "serverInfo": {"name": "vibium", "version": "0.1.0"}}}
```

**Step 2: List Tools** - Get available tools with their schemas

```
→ {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
← {"jsonrpc": "2.0", "id": 2, "result": {"tools": [
    {"name": "browser_launch", "description": "Launch a new browser session", "inputSchema": {...}},
    {"name": "browser_navigate", "description": "Navigate to a URL", "inputSchema": {...}},
    ...
  ]}}
```

The `inputSchema` (JSON Schema) tells Claude:
- What parameters each tool accepts
- Which parameters are required vs optional
- Parameter types and descriptions

You can inspect exactly what Claude learns:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./clicker/bin/clicker mcp | jq .result.tools
```

**Important:** Tool discovery happens **on startup**. After adding or modifying an MCP server, you must start a new Claude Code session for changes to take effect.

## Testing the Integration

Once added, start a new Claude Code session and ask Claude to use browser automation:

```
> Take a screenshot of https://example.com
```

Claude will use the Vibium MCP tools:
1. `browser_launch` - Start the browser
2. `browser_navigate` - Go to the URL
3. `browser_screenshot` - Capture the page
4. `browser_quit` - Close the browser

## Available MCP Tools (22)

| Tool | Description |
|------|-------------|
| `browser_launch` | Start a browser session (visible by default) |
| `browser_navigate` | Navigate to a URL |
| `browser_get_url` | Get current page URL |
| `browser_get_title` | Get current page title |
| `browser_get_text` | Get text content (page or element) |
| `browser_get_html` | Get HTML content (innerHTML or outerHTML) |
| `browser_find` | Find element info (tag, text, bounding box) |
| `browser_find_all` | Find all matching elements |
| `browser_evaluate` | Execute JavaScript in the browser |
| `browser_screenshot` | Capture a screenshot |
| `browser_click` | Click an element by CSS selector |
| `browser_type` | Type text into an element |
| `browser_hover` | Hover over an element |
| `browser_scroll` | Scroll page or element |
| `browser_keys` | Press key or key combination |
| `browser_select` | Select dropdown option |
| `browser_wait` | Wait for element state |
| `browser_new_tab` | Open a new tab |
| `browser_list_tabs` | List all open tabs |
| `browser_switch_tab` | Switch to tab by index or URL |
| `browser_close_tab` | Close a tab |
| `browser_quit` | Close the browser session |

## Troubleshooting

### MCP server not responding

Check that the clicker binary exists and is executable:

```bash
./clicker/bin/clicker mcp --help
```

### Browser fails to launch

Ensure Chrome for Testing is installed:

```bash
./clicker/bin/clicker install
```

### View MCP server logs

Run clicker directly to see any error output. You can test the full flow by sending JSON-RPC messages to stdin:

**Initialize the connection:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./clicker/bin/clicker mcp
```

Expected response:
```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"vibium","version":"0.1.0"}}}
```

**Test a full browser session (headed):**

Create a test script to send multiple commands:

```bash
cat << 'EOF' | ./clicker/bin/clicker mcp
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"browser_launch","arguments":{"headless":false}}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"browser_navigate","arguments":{"url":"https://example.com"}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"browser_quit","arguments":{}}}
EOF
```

Expected output (one JSON response per line):
```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"vibium","version":"0.1.0"}}}
{"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"Browser launched (headless: false)"}]}}
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"Navigated to https://example.com/"}]}}
{"jsonrpc":"2.0","id":4,"result":{"content":[{"type":"text","text":"Browser session closed"}]}}
```

**Individual commands for reference:**

| Action | JSON-RPC Message |
|--------|------------------|
| Initialize | `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}` |
| Launch (headed) | `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"browser_launch","arguments":{"headless":false}}}` |
| Navigate | `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"browser_navigate","arguments":{"url":"https://example.com"}}}` |
| Screenshot | `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"browser_screenshot","arguments":{}}}` |
| Screenshot (to file) | `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"browser_screenshot","arguments":{"filename":"page.png"}}}` |
| Quit | `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"browser_quit","arguments":{}}}` |
