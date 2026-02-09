# Getting Started with Vibium MCP

Let AI control your browser. This guide shows you how to add Vibium to your AI coding assistant.

---

## What You'll Get

After setup, you can ask your AI assistant things like:
- "Take a screenshot of https://example.com"
- "Go to Hacker News and find the top story"
- "Fill out this form and click submit"

The AI will control a real browser to do it.

---

## Prerequisites

Install one of the supported AI coding assistants:

- **Claude Code:** [claude.ai/download](https://claude.ai/download)
- **Gemini CLI:** [github.com/google-gemini/gemini-cli](https://github.com/google-gemini/gemini-cli)

---

## Quick Setup

**Claude Code:**
```bash
claude mcp add vibium -- npx -y vibium mcp
```

**Gemini CLI:**
```bash
gemini mcp add vibium npx -y vibium mcp
```

That's it. Chrome downloads automatically on first use.

---

## Try It

Restart your AI assistant, then ask:

```
Take a screenshot of https://example.com
```

You'll see:
1. A Chrome window open
2. The page load
3. The AI respond with the screenshot

Screenshots are saved to `~/Pictures/Vibium/` (macOS/Linux) or `Pictures\Vibium\` (Windows).

---

## Available Tools

Once Vibium is added, your AI can use these tools:

| Tool | What It Does |
|------|--------------|
| `browser_launch` | Opens Chrome (visible by default) |
| `browser_navigate` | Goes to a URL |
| `browser_find` | Finds an element by CSS selector |
| `browser_click` | Clicks an element |
| `browser_type` | Types text into an element |
| `browser_screenshot` | Captures the page |
| `browser_quit` | Closes the browser |

You don't need to know these — just ask in plain English.

---

## Detailed Setup Guides

For troubleshooting, advanced options, and platform-specific instructions:

- [Claude Code MCP Setup](claude-code-mcp-setup.md)
- [Gemini CLI MCP Setup](gemini-cli-mcp-setup.md)

---

## Next Steps

**Use the JS API directly:**
See [Getting Started Tutorial](getting-started.md) for programmatic control.

**Learn more about MCP:**
[Model Context Protocol docs](https://modelcontextprotocol.io)
