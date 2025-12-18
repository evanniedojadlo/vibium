# Contributing to Vibium

## Development Environment

We recommend developing inside a macOS VM to limit the blast radius of AI-assisted tools like Claude Code. See [docs/local-dev-setup.md](docs/local-dev-setup.md) for complete VM setup instructions.

If you prefer to develop directly on your host machine, follow the steps below.

---

## Prerequisites

- Go 1.21+
- Node.js 18+

---

## Clone and Build

```bash
git clone https://github.com/VibiumDev/vibium.git
cd vibium
make
```

This installs npm dependencies and builds both clicker and the JS client.

---

## Available Make Targets

```bash
make             # Build everything (default)
make build-go    # Build clicker binary
make build-js    # Build JS client
make deps        # Install npm dependencies
make serve       # Start proxy server on :9515
make clean       # Clean binaries and JS dist
make clean-cache # Clean cached Chrome for Testing
make clean-all   # Clean everything
make help        # Show all targets
```

---

## Using the JS Client

After building, you can test the JS client in a Node REPL:

```bash
cd clients/javascript && node --experimental-repl-await
```

```javascript
const { browser } = await import('./dist/index.mjs')
const vibe = await browser.launch({ headless: false })
await vibe.go('https://example.com')
const shot = await vibe.screenshot()
require('fs').writeFileSync('test.png', shot)
await vibe.quit()
```

---

## Using Clicker

Clicker is the Go binary at the heart of Vibium. It handles browser lifecycle, WebDriver BiDi protocol, and will eventually expose an MCP server for AI agents.

Long-term, clicker runs silently in the background — called by JS/TS, Python, or Java client libraries. Most users won't interact with it directly.

For now, the CLI is a development and testing aid. It lets you verify browser automation works before the client libraries are built on top.

After building, the binary is at `./clicker/bin/clicker`.

### Setup

```bash
./clicker/bin/clicker install   # Download Chrome for Testing + chromedriver
./clicker/bin/clicker paths     # Show browser and cache paths
./clicker/bin/clicker version   # Show version
```

### Browser Commands

```bash
# Navigate to a URL
./clicker/bin/clicker navigate https://example.com

# Take a screenshot
./clicker/bin/clicker screenshot https://example.com -o shot.png

# Evaluate JavaScript
./clicker/bin/clicker eval https://example.com "document.title"

# Find an element
./clicker/bin/clicker find https://example.com "a"

# Click an element
./clicker/bin/clicker click https://example.com "a"

# Type into an element
./clicker/bin/clicker type https://the-internet.herokuapp.com/inputs "input" "12345"
```

### Useful Flags

```bash
--headed          # Show the browser window (headless by default)
--wait-open 5     # Wait 5 seconds after navigation for page to load
--wait-close 3    # Keep browser open 3 seconds before closing
```

Example:

```bash
./clicker/bin/clicker screenshot https://example.com --headed --wait-close 5 -o shot.png
```

---

## Submitting Changes

- **Team members**: push directly to `VibiumDev/vibium`
- **External contributors**: fork the repo, push to your fork, then open a PR to `VibiumDev/vibium`

See [docs/local-dev-setup.md](docs/local-dev-setup.md) for details on the fork-based workflow.
