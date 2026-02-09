# Contributing to Vibium

## Development Environment

We recommend developing inside a VM to limit the blast radius of AI-assisted tools like Claude Code. See the setup guide for your platform:
- [macOS](docs/local-dev-setup-mac.md)
- [Linux x86](docs/local-dev-setup-x86-linux.md)
- [Windows x86](docs/local-dev-setup-x86-windows.md)

If you prefer to develop directly on your host machine, follow the steps below.

---

## Prerequisites

- Go 1.21+
- Node.js 18+
- Python 3.9+ (for Python client development)
- GitHub CLI (optional, for managing issues/PRs from terminal):
  - macOS: `brew install gh`
  - Linux: `sudo apt install gh` or `sudo dnf install gh`
  - Windows: `winget install GitHub.cli`

---

## Clone and Build

```bash
git clone https://github.com/VibiumDev/vibium.git
cd vibium
make
make test
```

This installs npm dependencies, builds clicker and the JS client, downloads Chrome for Testing (if needed), and runs the test suite.

---

## Available Make Targets

### Build

```bash
make                       # Build everything (default)
make build-go              # Build clicker binary
make build-js              # Build JS client
make build-go-all          # Cross-compile clicker for all platforms
```

### Package

```bash
make package               # Build all packages (npm + Python)
make package-js            # Build npm packages only
make package-python        # Build Python wheels only
```

### Test

```bash
make test                  # Run all tests (auto-installs Chrome for Testing)
make test-cli              # Run CLI tests only
make test-js               # Run JS library tests only
make test-mcp              # Run MCP server tests only
make test-python           # Run Python client tests
make test-daemon           # Run daemon lifecycle tests
```

### Other

```bash
make install-browser       # Install Chrome for Testing
make deps                  # Install npm dependencies
make serve                 # Start proxy server on :9515
make double-tap            # Kill zombie Chrome/chromedriver processes
make get-version           # Show current version
make set-version VERSION=x.x.x  # Set version across all packages
```

### Clean

```bash
make clean                 # Clean binaries and JS dist
make clean-go              # Clean clicker binaries
make clean-js              # Clean JS client dist
make clean-npm-packages    # Clean built npm packages
make clean-python-packages # Clean Python packages
make clean-packages        # Clean all packages (npm + Python)
make clean-cache           # Clean cached Chrome for Testing
make clean-all             # Clean everything
```

---

## Using the JS Client

After building, you can test the JS client in a Node REPL:

```bash
cd clients/javascript && node --experimental-repl-await
```

```javascript
// Option 1: require (REPL-friendly)
const { browserSync } = require('./dist')

// Option 2: dynamic import (REPL with --experimental-repl-await)
const { browser } = await import('./dist/index.mjs')

// Option 3: static import (in .mjs files)
import { browser } from './dist/index.mjs'
```

Sync example:

```javascript
const { browserSync } = require('./dist')
const vibe = browserSync.launch()
vibe.go('https://example.com')

const el = vibe.find('h1')
console.log(el.text())

// Execute JavaScript
const title = vibe.evaluate('document.title')
console.log('Page title:', title)

const shot = vibe.screenshot()
require('fs').writeFileSync('test.png', shot)
vibe.quit()
```

Async example:

```javascript
const { browser } = await import('./dist/index.mjs')
const vibe = await browser.launch()
await vibe.go('https://example.com')

const el = await vibe.find('h1')
console.log(await el.text())

// Execute JavaScript
const title = await vibe.evaluate('document.title')
console.log('Page title:', title)

const shot = await vibe.screenshot()
require('fs').writeFileSync('test.png', shot)
await vibe.quit()
```

---

## Using the Python Client

The Python client provides both sync and async APIs.

### Setup

For local development, use a virtual environment:

```bash
cd clients/python
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -e .           # Editable install - code changes take effect immediately
```

Or install from PyPI:

```bash
pip install vibium
```

### Sync Example

```python
from vibium import browser_sync

vibe = browser_sync.launch()
vibe.go("https://example.com")

el = vibe.find("h1")
print(el.text())

# Execute JavaScript
title = vibe.evaluate("document.title")
print(f"Page title: {title}")

with open("screenshot.png", "wb") as f:
    f.write(vibe.screenshot())

vibe.quit()
```

### Async Example

```python
import asyncio
from vibium import browser

async def main():
    vibe = await browser.launch()
    await vibe.go("https://example.com")

    el = await vibe.find("h1")
    print(await el.text())

    # Execute JavaScript
    title = await vibe.evaluate("document.title")
    print(f"Page title: {title}")

    with open("screenshot.png", "wb") as f:
        f.write(await vibe.screenshot())

    await vibe.quit()

asyncio.run(main())
```

---

## Using Clicker

Clicker is the Go binary at the heart of Vibium. It handles browser lifecycle, WebDriver BiDi protocol, and exposes an MCP server for AI agents.

Long-term, clicker runs silently in the background — called by client libraries (JS/TS, Python, etc.). Most users won't interact with it directly.

For now, the CLI is a development and testing aid. It lets you verify browser automation works before the client libraries are built on top.

After building, the binary is at `./clicker/bin/clicker`.

### Setup

```bash
./clicker/bin/clicker install   # Download Chrome for Testing + chromedriver
./clicker/bin/clicker paths     # Show browser and cache paths
./clicker/bin/clicker version   # Show version
```

### Browser Commands

By default, clicker runs in **daemon mode** — the browser stays open between commands:

```bash
# Navigate to a URL
./clicker/bin/clicker navigate https://example.com

# Interact with the current page (no URL needed)
./clicker/bin/clicker find "h1"
./clicker/bin/clicker click "a"
./clicker/bin/clicker type "input" "hello"
./clicker/bin/clicker eval "document.title"
./clicker/bin/clicker screenshot -o shot.png

# You can also provide a URL to navigate first
./clicker/bin/clicker find https://example.com "a"
./clicker/bin/clicker screenshot https://example.com -o shot.png
```

Use `--oneshot` to launch a fresh browser for each command (the old behavior):

```bash
./clicker/bin/clicker navigate https://example.com --oneshot
```

### Useful Flags

```bash
--headless        # Hide the browser window (visible by default)
--oneshot          # Launch a fresh browser per command (no daemon)
--json             # Output results as JSON
--wait-open 5     # Wait 5 seconds after navigation for page to load
--wait-close 3    # Keep browser open 3 seconds before closing (oneshot only)
```

### Daemon Management

```bash
./clicker/bin/clicker daemon start    # Start daemon in foreground
./clicker/bin/clicker daemon start -d # Start daemon in background
./clicker/bin/clicker daemon status   # Show daemon status
./clicker/bin/clicker daemon stop     # Stop the daemon
```

The daemon auto-starts on the first command, so you rarely need to manage it manually.

---

## Using the MCP Server

Clicker includes an MCP (Model Context Protocol) server for AI agent integration.

### Available Tools

| Tool | Description |
|------|-------------|
| `browser_launch` | Start a browser session |
| `browser_navigate` | Go to a URL |
| `browser_click` | Click an element by CSS selector |
| `browser_type` | Type into an element |
| `browser_screenshot` | Capture the page |
| `browser_find` | Find element info |
| `browser_evaluate` | Execute JavaScript to extract data or inspect page state |
| `browser_quit` | Close the browser |
| `browser_get_text` | Get text content of page or element |
| `browser_get_url` | Get the current page URL |
| `browser_get_title` | Get the current page title |
| `browser_get_html` | Get HTML content of page or element |
| `browser_find_all` | Find all elements matching a CSS selector |
| `browser_wait` | Wait for element to reach a state (attached/visible/hidden) |
| `browser_hover` | Hover over an element |
| `browser_select` | Select an option in a `<select>` element |
| `browser_scroll` | Scroll the page or an element |
| `browser_keys` | Press a key or key combination |
| `browser_new_tab` | Open a new browser tab |
| `browser_list_tabs` | List all open browser tabs |
| `browser_switch_tab` | Switch to a tab by index or URL |
| `browser_close_tab` | Close a browser tab |

### Running the MCP Server

```bash
# Run directly (for testing)
./clicker/bin/clicker mcp

# With custom screenshot directory
./clicker/bin/clicker mcp --screenshot-dir ./screenshots

# Disable screenshot file saving (inline base64 only)
./clicker/bin/clicker mcp --screenshot-dir ""
```

### Configuring with Claude Code

```bash
claude mcp add vibium -- clicker mcp
```

### Testing with JSON-RPC

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./clicker/bin/clicker mcp
```

---

## Debugging

For low-level debugging tools and troubleshooting tips, see [docs/how-to-guides/debugging.md](docs/how-to-guides/debugging.md).

---

## Submitting Changes

- **Team members**: push directly to `VibiumDev/vibium`
- **External contributors**: fork the repo, push to your fork, then open a PR to `VibiumDev/vibium`

See [docs/local-dev-setup-mac.md](docs/local-dev-setup-mac.md) for details on the fork-based workflow.
