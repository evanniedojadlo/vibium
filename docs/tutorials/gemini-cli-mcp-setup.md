# Setting Up Vibium MCP in Gemini CLI

This tutorial covers how to configure Vibium as an MCP (Model Context Protocol) server in Google's Gemini CLI.

## Prerequisites

- Gemini CLI installed ([installation guide](https://github.com/google-gemini/gemini-cli))
- Google account (free tier: 60 requests/min, 1,000 requests/day — no credit card required)
- Node.js 18+ (for npx)

## Add Vibium MCP

### Option 1: Using the CLI (Recommended)

```bash
gemini mcp add vibium npx -y vibium mcp
```

### Option 2: Manual Configuration

Edit `~/.gemini/settings.json` (create it if it doesn't exist):

```json
{
  "mcpServers": {
    "vibium": {
      "command": "npx",
      "args": ["-y", "vibium", "mcp"]
    }
  }
}
```

For project-specific config, create `.gemini/settings.json` in your project directory instead.

### Option 3: Using Local Binary

If you built clicker locally:

```bash
gemini mcp add vibium /path/to/clicker mcp
```

Or in `settings.json`:

```json
{
  "mcpServers": {
    "vibium": {
      "command": "/path/to/clicker",
      "args": ["mcp"]
    }
  }
}
```

### Custom Screenshot Directory

By default, screenshots are saved to:
- macOS: `~/Pictures/Vibium/`
- Linux: `~/Pictures/Vibium/`
- Windows: `%USERPROFILE%\Pictures\Vibium\`

To use a different directory:

```json
{
  "mcpServers": {
    "vibium": {
      "command": "npx",
      "args": ["-y", "vibium", "mcp", "--screenshot-dir", "./screenshots"]
    }
  }
}
```

To disable file saving (base64 inline only):

```json
{
  "mcpServers": {
    "vibium": {
      "command": "npx",
      "args": ["-y", "vibium", "mcp", "--screenshot-dir", ""]
    }
  }
}
```

## Verify Installation

List configured MCP servers:

```bash
gemini mcp list
```

You should see `vibium` in the output.

## Testing the Integration

Start Gemini CLI and ask it to use browser automation:

```
> Take a screenshot of https://example.com
```

Gemini will use the Vibium MCP tools:
1. `browser_launch` - Start the browser
2. `browser_navigate` - Go to the URL
3. `browser_screenshot` - Capture the page
4. `browser_quit` - Close the browser

## Available MCP Tools

| Tool | Description |
|------|-------------|
| `browser_launch` | Start a browser session (visible by default) |
| `browser_navigate` | Navigate to a URL |
| `browser_click` | Click an element by CSS selector |
| `browser_type` | Type text into an element |
| `browser_screenshot` | Capture a screenshot |
| `browser_find` | Find element info (tag, text, bounding box) |
| `browser_quit` | Close the browser session |

## Troubleshooting

### "Client is not connected" error

This usually means the MCP server command isn't running correctly. Try:

1. **Verify npx works:**
   ```bash
   npx -y vibium mcp
   ```
   You should see "Vibium MCP server" output (press Ctrl+C to exit).

2. **Check Chrome for Testing is installed:**
   ```bash
   npx -y vibium mcp
   # Then in another terminal:
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | npx -y vibium mcp
   ```

3. **Use absolute paths** if relative paths aren't working.

### Browser fails to launch

The first time you run Vibium, it downloads Chrome for Testing. If this fails:

```bash
# Run directly to see error messages
npx -y vibium mcp
```

On macOS, if you see a Gatekeeper warning about chromedriver, this should be fixed in v0.1.5+.

### Test MCP server manually

Send JSON-RPC messages directly to verify the server works:

```bash
cat << 'EOF' | npx -y vibium mcp
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"browser_launch","arguments":{"headless":true}}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"browser_navigate","arguments":{"url":"https://example.com"}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"browser_screenshot","arguments":{}}}
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"browser_quit","arguments":{}}}
EOF
```

You should see JSON responses for each command.

### Enable debug mode

Run Gemini CLI with `--debug` for verbose output:

```bash
gemini --debug
```

Or press F12 in interactive mode to open the debug console.

## Remove Vibium MCP

```bash
gemini mcp remove vibium
```

Or manually remove the `vibium` entry from your `settings.json`.

## References

- [Gemini CLI MCP Documentation](https://geminicli.com/docs/tools/mcp-server/)
- [Gemini CLI Configuration](https://github.com/google-gemini/gemini-cli/blob/main/docs/get-started/configuration.md)
