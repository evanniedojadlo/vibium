---
name: vibe-check
description: Browser automation for AI agents. Use when the user needs to navigate websites, read page content, fill forms, click elements, take screenshots, or manage browser tabs.
---

# Vibium Browser Automation — CLI Reference

The `vibium` CLI automates Chrome via the command line. The browser auto-launches on first use (daemon mode keeps it running between commands).

```
vibium go <url> → vibium map → vibium click @e1 → vibium map
```

## Core Workflow

Every browser automation follows this pattern:

1. **Navigate**: `vibium go <url>`
2. **Map**: `vibium map` (get element refs like `@e1`, `@e2`)
3. **Interact**: Use refs to click, fill, select — e.g. `vibium click @e1`
4. **Re-map**: After navigation or DOM changes, get fresh refs with `vibium map`

## Binary Resolution

Before running any commands, resolve the `vibium` binary path once:

1. Try `vibium` directly (works if globally installed via `npm install -g vibium`)
2. Fall back to `./clicker/bin/vibium` (dev environment, in project root)
3. Fall back to `./node_modules/.bin/vibium` (local npm install)

Run `vibium --help` (or the resolved path) to confirm. Use the resolved path for all subsequent commands.

**Windows note:** Use forward slashes in paths (e.g. `./clicker/bin/vibium.exe`) and quote paths containing spaces.

## Commands

### Discovery
- `vibium map` — map interactive elements with @refs (recommended before interacting)
- `vibium diff map` — compare current vs last map (see what changed)

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
- `vibium eval "<js>"` — run JavaScript and print result (`--stdin` to read from stdin)
- `vibium count "<selector>"` — count matching elements
- `vibium screenshot -o file.png` — capture screenshot (`--full-page`, `--annotate`)
- `vibium a11y-tree` — accessibility tree (`--everything` for all nodes)

### Interaction
- `vibium click "<selector>"` — click an element (also accepts `@ref` from map)
- `vibium dblclick "<selector>"` — double-click an element
- `vibium type "<selector>" "<text>"` — type into an input (appends to existing value)
- `vibium fill "<selector>" "<text>"` — clear field and type new text (replaces value)
- `vibium press <key> [selector]` — press a key on element or focused element
- `vibium focus "<selector>"` — focus an element
- `vibium hover "<selector>"` — hover over an element
- `vibium scroll [direction]` — scroll page (`--amount N`, `--selector`)
- `vibium scroll-into-view "<selector>"` — scroll element into view (centered)
- `vibium keys "<combo>"` — press keys (Enter, Control+a, Shift+Tab)
- `vibium select "<selector>" "<value>"` — pick a dropdown option
- `vibium check "<selector>"` — check a checkbox/radio (idempotent)
- `vibium uncheck "<selector>"` — uncheck a checkbox (idempotent)

### Mouse Primitives
- `vibium mouse-click [x] [y]` — click at coordinates or current position (`--button 0|1|2`)
- `vibium mouse-move <x> <y>` — move mouse to coordinates
- `vibium mouse-down` — press mouse button (`--button 0|1|2`)
- `vibium mouse-up` — release mouse button (`--button 0|1|2`)
- `vibium drag "<source>" "<target>"` — drag from one element to another

### Element State
- `vibium value "<selector>"` — get input/textarea/select value
- `vibium attr "<selector>" "<attribute>"` — get HTML attribute value
- `vibium is-visible "<selector>"` — check if element is visible (true/false)
- `vibium is-enabled "<selector>"` — check if element is enabled (true/false)
- `vibium is-checked "<selector>"` — check if checkbox/radio is checked (true/false)

### Waiting
- `vibium wait "<selector>"` — wait for element (`--state visible|hidden|attached`, `--timeout ms`)
- `vibium wait-for-url "<pattern>"` — wait until URL contains substring (`--timeout ms`)
- `vibium wait-for-load` — wait until page is fully loaded (`--timeout ms`)
- `vibium wait-for-text "<text>"` — wait until text appears on page (`--timeout ms`)
- `vibium wait-for-fn "<expression>"` — wait until JS expression returns truthy (`--timeout ms`)
- `vibium sleep <ms>` — pause execution (max 30000ms)

### Capture
- `vibium screenshot -o file.png` — capture screenshot (`--full-page`, `--annotate`)
- `vibium pdf -o file.pdf` — save page as PDF

### Dialogs
- `vibium dialog accept [text]` — accept dialog (optionally with prompt text)
- `vibium dialog dismiss` — dismiss dialog

### Emulation
- `vibium set-viewport <width> <height>` — set viewport size (`--dpr` for device pixel ratio)
- `vibium viewport` — get current viewport dimensions
- `vibium window` — get OS browser window dimensions and state
- `vibium set-window <width> <height> [x] [y]` — set window size and position (`--state`)
- `vibium emulate-media` — override CSS media features (`--color-scheme`, `--reduced-motion`, `--forced-colors`, `--contrast`, `--media`)
- `vibium set-geolocation <lat> <lng>` — override geolocation (`--accuracy`)
- `vibium set-content "<html>"` — replace page HTML (`--stdin` to read from stdin)

### Frames
- `vibium frames` — list all iframes on the page
- `vibium frame "<nameOrUrl>"` — find a frame by name or URL substring

### File Upload
- `vibium upload "<selector>" <files...>` — set files on input[type=file]

### Tracing
- `vibium trace start` — start recording (`--screenshots`, `--snapshots`, `--name`)
- `vibium trace stop` — stop recording and save ZIP (`-o path`)

### Cookies
- `vibium cookies` — list all cookies
- `vibium cookies set <name> <value>` — set a cookie
- `vibium cookies clear` — clear all cookies

### Storage State
- `vibium storage-state` — export cookies + localStorage + sessionStorage (`-o state.json`)
- `vibium restore-storage <path>` — restore state from JSON file

### Downloads
- `vibium download set-dir <path>` — set download directory

### Tabs
- `vibium tabs` — list open tabs
- `vibium tab-new [url]` — open new tab
- `vibium tab-switch <index|url>` — switch tab
- `vibium tab-close [index]` — close tab

### Debug
- `vibium highlight "<selector>"` — highlight element visually (3 seconds)

### Session
- `vibium quit` — close the browser (daemon keeps running)
- `vibium close` — alias for quit
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

**Ref-based workflow (recommended for AI):**
```sh
vibium go https://example.com
vibium map
vibium click @e1
vibium map  # re-map after interaction
```

**Verify action worked:**
```sh
vibium map
vibium click @e3
vibium diff map  # see what changed
```

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

**Save as PDF:**
```sh
vibium go https://example.com
vibium pdf -o page.pdf
```

**Annotated screenshot:**
```sh
vibium screenshot -o annotated.png --annotate
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

## Ref Lifecycle

Refs (`@e1`, `@e2`) are invalidated when the page changes. Always re-map after:
- Clicking links or buttons that navigate
- Form submissions
- Dynamic content loading (dropdowns, modals)

## Tips

- All click/type/hover/fill actions auto-wait for the element to be actionable
- All selector arguments also accept `@ref` from `vibium map`
- Use `vibium map` before interacting to discover interactive elements
- Use `vibium fill` to replace a field's value, `vibium type` to append to it
- Use `vibium find` to inspect an element before interacting
- Use `vibium find-by-role` for semantic element lookup (more reliable than CSS selectors)
- Use `vibium a11y-tree` to understand page structure without visual rendering
- Use `vibium text "<selector>"` to read specific sections
- `vibium eval` is the escape hatch for complex DOM queries
- `vibium check`/`vibium uncheck` are idempotent — safe to call without checking state first
- Screenshots save to the current directory by default (`-o` to change)
