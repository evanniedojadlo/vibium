---
name: vibe-check
description: Browser automation for AI agents. Use when the user needs to navigate websites, read page content, fill forms, click elements, take screenshots, or manage browser tabs.
---

# Vibium Browser Automation — CLI Reference

The `vibium` CLI automates Chrome via the command line. The browser auto-launches on first use (daemon mode keeps it running between commands).

```
vibium go <url> → vibium text → vibium screenshot -o shot.png
```

## Binary Resolution

Before running any commands, resolve the `vibium` binary path once:

1. Try `vibium` directly (works if globally installed via `npm install -g vibium`)
2. Fall back to `./clicker/bin/vibium` (dev environment, in project root)
3. Fall back to `./node_modules/.bin/vibium` (local npm install)

Run `vibium --help` (or the resolved path) to confirm. Use the resolved path for all subsequent commands.

**Windows note:** Use forward slashes in paths (e.g. `./clicker/bin/vibium.exe`) and quote paths containing spaces.

## Commands

### Navigation
- `vibium go <url>` — go to a page
- `vibium back` — go back in history
- `vibium forward` — go forward in history
- `vibium reload` — reload the current page
- `vibium url` — print current URL
- `vibium title` — print page title

### Reading Content
- `vibium text` — get all page text
- `vibium text "<selector>"` — get text of a specific element
- `vibium html` — get page HTML (use `--outer` for outerHTML)
- `vibium find "<selector>"` — element info (tag, text, bounding box)
- `vibium find-all "<selector>"` — all matching elements (`--limit N`)
- `vibium find-by-role` — find element by ARIA role/name (`--role`, `--name`, `--selector`, `--timeout`)
- `vibium eval "<js>"` — run JavaScript and print result
- `vibium screenshot -o file.png` — capture screenshot (`--full-page` for entire document)
- `vibium a11y-tree` — accessibility tree (`--everything` for all nodes)

### Interaction
- `vibium click "<selector>"` — click an element
- `vibium type "<selector>" "<text>"` — type into an input (appends to existing value)
- `vibium fill "<selector>" "<text>"` — clear field and type new text (replaces value)
- `vibium press <key> [selector]` — press a key on element or focused element
- `vibium hover "<selector>"` — hover over an element
- `vibium scroll [direction]` — scroll page (`--amount N`, `--selector`)
- `vibium scroll-into-view "<selector>"` — scroll element into view (centered)
- `vibium keys "<combo>"` — press keys (Enter, Control+a, Shift+Tab)
- `vibium select "<selector>" "<value>"` — pick a dropdown option
- `vibium check "<selector>"` — check a checkbox/radio (idempotent)
- `vibium uncheck "<selector>"` — uncheck a checkbox (idempotent)

### Element State
- `vibium value "<selector>"` — get input/textarea/select value
- `vibium attr "<selector>" "<attribute>"` — get HTML attribute value
- `vibium is-visible "<selector>"` — check if element is visible (true/false)

### Waiting
- `vibium wait "<selector>"` — wait for element (`--state visible|hidden|attached`, `--timeout ms`)
- `vibium wait-for-url "<pattern>"` — wait until URL contains substring (`--timeout ms`)
- `vibium wait-for-load` — wait until page is fully loaded (`--timeout ms`)
- `vibium sleep <ms>` — pause execution (max 30000ms)

### Tabs
- `vibium tabs` — list open tabs
- `vibium tab-new [url]` — open new tab
- `vibium tab-switch <index|url>` — switch tab
- `vibium tab-close [index]` — close tab

### Session
- `vibium quit` — close the browser (daemon keeps running)
- `vibium daemon start` — start background browser
- `vibium daemon status` — check if running
- `vibium daemon stop` — stop daemon

## Global Flags

| Flag | Description |
|------|-------------|
| `--headless` | Hide browser window |
| `--json` | Output as JSON |
| `--oneshot` | One-shot mode (no daemon) |
| `-v, --verbose` | Debug logging |
| `--wait-open N` | Wait N seconds after navigation |
| `--wait-close N` | Keep browser open N seconds before closing |

## Daemon vs Oneshot

By default, commands connect to a **daemon** — a background process that keeps the browser alive between commands. This is fast and lets you chain commands against the same page.

Use `--oneshot` (or `VIBIUM_ONESHOT=1`) to launch a fresh browser for each command, then tear it down. Useful for CI or one-off scripts.

## Common Patterns

**Read a page:**
```sh
vibium go https://example.com
vibium text
```

**Fill a form:**
```sh
vibium go https://example.com/login
vibium fill "input[name=email]" "user@example.com"
vibium fill "input[name=password]" "secret"
vibium click "button[type=submit]"
vibium wait-for-url "/dashboard"
```

**Check page structure without rendering:**
```sh
vibium go https://example.com
vibium a11y-tree
```

**Extract structured data:**
```sh
vibium go https://example.com
vibium eval "JSON.stringify([...document.querySelectorAll('a')].map(a => a.href))"
```

**Inspect an element:**
```sh
vibium attr "a" "href"
vibium value "input[name=email]"
vibium is-visible ".modal"
```

**Multi-tab workflow:**
```sh
vibium tab-new https://docs.example.com
vibium text "h1"
vibium tab-switch 0
```

**Full-page screenshot:**
```sh
vibium screenshot -o full.png --full-page
```

## Tips

- All click/type/hover/fill actions auto-wait for the element to be actionable
- Use `vibium fill` to replace a field's value, `vibium type` to append to it
- Use `vibium find` to inspect an element before interacting
- Use `vibium find-by-role` for semantic element lookup (more reliable than CSS selectors)
- Use `vibium a11y-tree` to understand page structure without visual rendering
- Use `vibium text "<selector>"` to read specific sections
- `vibium eval` is the escape hatch for complex DOM queries
- `vibium check`/`vibium uncheck` are idempotent — safe to call without checking state first
- Screenshots save to the current directory by default (`-o` to change)
