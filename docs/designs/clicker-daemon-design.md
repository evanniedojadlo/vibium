# Clicker Daemon Design Document

## Status: Draft
## Date: February 2026

---

## 1. Problem

Every `clicker` CLI command today is one-shot: launch Chrome → execute action → tear down. This means:

- **No session continuity.** An agent that wants to navigate, then find an element, then click, then read text must pay the browser startup cost each time and loses all state between commands.
- **MCP tools are isolated.** `browser_launch` starts a browser, but the MCP server running in stdio mode already manages its own lifecycle — the CLI and MCP are two separate worlds that can't share a session.
- **Impossible to build multi-step workflows.** The CLI is useless for anything beyond single-action screenshots. Agents need a living browser they can interact with over time.
- **MCP locks out most agents.** MCP requires explicit client support — the agent framework has to implement the protocol. A CLI tool that persists state and outputs JSON works with *every* agent that can shell out. That's all of them.

## 2. Proposal

Add a daemon mode to `clicker` that runs the existing MCP server logic as a long-lived background process, reachable via Unix socket. The existing `clicker mcp` (inline, stdio) continues to work exactly as it does today — no changes, no breakage. The daemon is a new, additive capability that enables two things: session persistence across CLI commands, and a universal agent skill interface via the CLI.

**Three interfaces, one daemon, one browser session:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        clicker daemon                           │
│                     (JSON-RPC 2.0 / MCP)                        │
│                            │                                    │
│                        Chrome (BiDi)                            │
└───────▲─────────────────▲──────────────────────▲────────────────┘
        │                 │                      │
   Unix socket       Unix socket            Unix socket
        │                 │                      │
  ┌─────┴──────┐   ┌──────┴──────┐   ┌──────────┴──────────┐
  │  CLI       │   │ MCP bridge  │   │  HTTP bridge        │
  │            │   │             │   │                     │
  │ clicker    │   │ clicker mcp │   │ clicker http        │
  │ navigate   │   │ --connect   │   │                     │
  │ find       │   │             │   │                     │
  │ click ...  │   │ (stdio)     │   │ (Streamable HTTP)   │
  └─────┬──────┘   └──────┬──────┘   └──────────┬──────────┘
        │                 │                      │
   Any agent         MCP clients           Remote agents
   (bash/shell)    (Claude Code)          (multi-machine)
```

- **MCP** (`clicker mcp`, existing) is the native integration for agents that support it. Rich, bidirectional, proper tool schemas. Unchanged.
- **CLI + daemon** is the universal integration. Any agent that can run shell commands gets browser automation. A skill file tells the agent the commands exist; `--json` makes output machine-parseable.
- **HTTP** (future) enables remote and multi-agent setups.

The daemon auto-starts transparently. Users and agents never need to think about it:

```bash
# First command: no daemon running → auto-starts one in background
clicker navigate https://example.com

# Subsequent commands: daemon already running → reuses session
clicker find "h1"
clicker click "a.login"
clicker screenshot -o page.png

# 30 minutes of inactivity → daemon quietly exits
```

The daemon reuses the same MCP tool handlers that `clicker mcp` already uses. No new protocol. No rewrite. Just a new transport (socket) and a new lifecycle (long-lived process).

## 3. Design Principles

1. **Don't break `clicker mcp`.** The existing inline MCP server is the zero-config path that works with Claude Code today. It continues to work identically — owns its own browser, manages its own lifecycle, stdio in/out. No changes.

2. **One protocol.** The daemon speaks JSON-RPC 2.0 (same as MCP). CLI commands targeting the daemon send `tools/call` messages over the socket. No REST, no gRPC, no custom protocol.

3. **Same tool handlers, different transport.** The daemon imports and runs the exact same tool handler functions that `clicker mcp` already uses. The only difference is where requests come from (socket vs stdin) and where responses go (socket vs stdout).

4. **Auto-start, auto-cleanup.** The daemon starts transparently on first CLI command and exits after an idle timeout (default 30 minutes). Users and agents never manage daemon lifecycle. `clicker daemon start|stop` exists for power users who want explicit control.

5. **The CLI is a first-class agent interface.** The daemon + CLI combination is not a lesser version of MCP — it's the most portable version. Any agent that can shell out gets browser automation. The CLI must be designed for agents: predictable argument patterns, `--json` output on every command, composable single-purpose commands, self-describing help text.

6. **Skill-file compatible.** The CLI surface should be documentable in a single skill file that any agent framework can consume. The skill file is the agent's "SDK" — no client library, no protocol integration, just a prompt that says "here are the commands you can run."

## 4. Architecture

### 4.1 Daemon Process

The daemon is a long-lived process that:

- Listens on a Unix domain socket (TCP on Windows)
- Accepts JSON-RPC 2.0 requests
- Manages browser sessions and contexts
- Handles cleanup on shutdown

**Socket location:**

| Platform | Path |
|----------|------|
| Linux    | `~/.cache/vibium/clicker.sock` |
| macOS    | `~/Library/Caches/vibium/clicker.sock` |
| Windows  | `\\.\pipe\vibium-clicker` (named pipe) |

**Bookkeeping files** (same directory as socket):

- `clicker.pid` — daemon process ID
- `clicker.lock` — advisory lock to prevent double-start

### 4.2 Session and Context Model

```
daemon
  └── session (1 Chrome process)
        ├── context "default"  ← implicit, always exists
        ├── context "task-1"   ← isolated browsing context
        └── context "task-2"
```

**V1 scope:** One session, one default context. This covers the overwhelmingly common case: a single agent driving a single browser. Named sessions and multiple contexts are deferred to V2 (see §8).

A session maps to one Chrome process. A context maps to a WebDriver BiDi browsing context (similar to an incognito window — isolated cookies, storage, history).

### 4.3 Protocol: JSON-RPC 2.0 Over Unix Socket

All communication uses newline-delimited JSON-RPC 2.0 messages, identical to MCP's wire format.

**Request (CLI → daemon):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "browser_navigate",
    "arguments": { "url": "https://example.com" }
  }
}
```

**Response (daemon → CLI):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      { "type": "text", "text": "Navigated to https://example.com" }
    ]
  }
}
```

This is the exact same shape as an MCP `tools/call`. The daemon doesn't know or care whether the request came from the CLI, stdio bridge, or HTTP bridge.

**Daemon management methods** (non-MCP, daemon-specific):

| Method | Description |
|--------|-------------|
| `daemon/status` | Returns daemon version, uptime, session info |
| `daemon/shutdown` | Graceful shutdown |

These use the same JSON-RPC envelope but are routed internally rather than to the MCP tool handler.

### 4.4 Transport Adapters

**`clicker mcp` (existing, unchanged):**

Runs the MCP server inline, exactly as today. Owns its own browser process. Communicates over stdin/stdout. This is what `claude mcp add vibium -- npx -y vibium` uses. **No changes.**

**`clicker mcp --connect` (new, optional):**

Instead of running the MCP server inline, bridges stdin/stdout to a running daemon's socket. Useful when you want Claude Code to share a browser session with CLI commands.

```
Claude Code ──stdin──► clicker mcp --connect ──socket──► daemon
Claude Code ◄─stdout── clicker mcp --connect ◄─socket── daemon
```

If no daemon is running, prints an error and suggests `clicker daemon start`.

**CLI adapter** (enhanced, existing commands):

CLI commands auto-start the daemon if it isn't running. The user never manages the daemon directly.

```
User types:    clicker navigate https://example.com
Binary does:   1. Check if daemon socket exists
               2a. YES → connect to socket, send tools/call, print result
               2b. NO  → spawn daemon in background, wait for socket (~100ms),
                          send tools/call, print result
```

The daemon starts with a default idle timeout (30 minutes). If no commands arrive for that duration, it shuts down cleanly and removes its socket/PID files. The next CLI command auto-starts a fresh daemon.

A `--oneshot` flag bypasses the daemon entirely for CI or scripting contexts where you want the old launch-execute-teardown behavior.

**HTTP adapter** (`clicker http`, new, Phase 3):

Bridges HTTP to a running daemon's socket. Implements MCP Streamable HTTP transport. Enables remote agents and multi-agent setups.

```
Remote Agent ──HTTP POST──► clicker http ──socket──► daemon
Remote Agent ◄─HTTP resp─── clicker http ◄─socket── daemon
```

### 4.5 CLI Command Mapping

Every CLI subcommand maps to exactly one MCP tool call. The CLI binary is a thin arg-parser that formats JSON-RPC.

| CLI Command | MCP Tool | Arguments | Status |
|-------------|----------|-----------|--------|
| `clicker navigate <url>` | `browser_navigate` | `{url}` | ✅ Exists |
| `clicker screenshot [-o file]` | `browser_screenshot` | `{path?}` | ✅ Exists |
| `clicker find <selector>` | `browser_find` | `{selector}` | ✅ Exists |
| `clicker click <selector>` | `browser_click` | `{selector}` | ✅ Exists |
| `clicker type <selector> <text>` | `browser_type` | `{selector, text}` | ✅ Exists |
| `clicker eval <js>` | `browser_eval` | `{script}` | ✅ Exists |
| `clicker text [selector]` | `browser_get_text` | `{selector?}` | New |
| `clicker url` | `browser_get_url` | `{}` | New |
| `clicker title` | `browser_get_title` | `{}` | New |
| `clicker scroll [selector] [direction]` | `browser_scroll` | `{selector?, direction?}` | New |
| `clicker wait <selector> [timeout]` | `browser_wait` | `{selector, timeout?}` | New |
| `clicker tabs` | `browser_list_tabs` | `{}` | New |
| `clicker keys <keys>` | `browser_keys` | `{keys}` | New |
| `clicker quit` | `browser_quit` | `{}` | ✅ Exists |

**Output modes:**
- Default: human-readable (just the text content, for interactive use)
- `--json`: structured envelope (for agents and scripting)

```bash
$ clicker title
Example Domain

$ clicker title --json
{"ok":true,"result":"Example Domain"}

$ clicker find "h1" --json
{"ok":true,"result":{"selector":"h1","text":"Example Domain","tag":"h1","visible":true}}

$ clicker find ".nope" --json
{"ok":false,"error":"No element found matching: .nope"}
```

### 4.6 Daemon Lifecycle

**Starting:**

```bash
# Explicit
clicker daemon start            # foreground
clicker daemon start -d         # background (daemonize)

# Implicit (auto-start)
clicker navigate https://x.com  # no daemon running → starts one, then navigates
```

Auto-start behavior:

1. CLI command connects to socket.
2. Connection refused → no daemon running.
3. Start daemon as background process.
4. Wait for socket to become available (poll with backoff, max 5s).
5. Send the original command.
6. If `--oneshot` flag is set, skip auto-start and use legacy one-shot mode.

**Stopping:**

```bash
clicker daemon stop     # sends daemon/shutdown, waits for clean exit
clicker daemon restart  # stop + start
```

On `daemon/shutdown`:
1. Close all browser contexts.
2. Terminate Chrome process.
3. Remove socket and PID files.
4. Exit.

**Status:**

```bash
$ clicker daemon status
vibium clicker daemon v0.2.0
status:   running
pid:      48291
uptime:   1h 23m
socket:   ~/.cache/vibium/clicker.sock
chrome:   pid 48305, Chrome for Testing 131.0.6778.85
session:  1 context(s), current URL: https://example.com
```

### 4.7 Browser Lifecycle Within the Daemon

The daemon manages Chrome lazily:

1. **Daemon starts** → no Chrome process yet.
2. **First tool call arrives** (e.g., `browser_navigate`) → daemon launches Chrome, creates default context, executes the command.
3. **`browser_quit` received** → Chrome process terminates, session state cleared. Next tool call launches a fresh Chrome.
4. **`browser_launch` received while Chrome is running** → no-op (returns current session info). This maintains compatibility with the existing MCP flow where agents call `browser_launch` first.

This means the daemon itself is extremely cheap when idle — just a process listening on a socket.

### 4.8 Crash Recovery

**Chrome crashes:**
- Daemon detects broken BiDi WebSocket.
- Next tool call returns an error: `{"error": {"code": -32000, "message": "Browser crashed. Send any browser command to restart."}}`.
- Next command auto-launches a new Chrome (session state is lost).

**Daemon crashes:**
- Socket file and PID file become stale.
- Next CLI command detects stale PID (process doesn't exist), cleans up files, auto-starts a new daemon.

**Idle timeout** (optional, off by default):
- `clicker daemon start --idle-timeout 30m`
- If no JSON-RPC messages received for the timeout duration, daemon shuts down cleanly.
- Useful for CI environments and forgotten daemons.

## 5. New MCP Tools

The daemon architecture makes it trivial to add new tools. These are the tools to add alongside the daemon work, since they require session persistence to be useful.

### 5.1 Page Reading

| Tool | Description | Arguments | Status |
|------|-------------|-----------|--------|
| `browser_eval` | Execute JavaScript, return result | `{script: string}` | ✅ Exists |
| `browser_get_text` | Get text content of page or element | `{selector?: string}` | |
| `browser_get_html` | Get HTML of page or element | `{selector?: string, outer?: bool}` | |
| `browser_get_url` | Get current URL | `{}` | |
| `browser_get_title` | Get page title | `{}` | |

### 5.2 Interaction

| Tool | Description | Arguments | Status |
|------|-------------|-----------|--------|
| `browser_scroll` | Scroll page or element | `{selector?: string, direction?: "up"\|"down"\|"left"\|"right", amount?: number}` | |
| `browser_hover` | Hover over element | `{selector: string}` | |
| `browser_select` | Select option in dropdown | `{selector: string, value: string}` | |
| `browser_keys` | Send keyboard input | `{keys: string}` (e.g., `"Enter"`, `"Control+a"`) | |
| `browser_wait` | Wait for element to appear | `{selector: string, timeout?: number, state?: "visible"\|"attached"\|"hidden"}` | |
| `browser_find_all` | Find all matching elements | `{selector: string, limit?: number}` | |

### 5.3 Tab Management

| Tool | Description | Arguments | Status |
|------|-------------|-----------|--------|
| `browser_new_tab` | Open new tab, optionally navigate | `{url?: string}` | |
| `browser_list_tabs` | List open tabs with URLs | `{}` | |
| `browser_switch_tab` | Switch to tab by index or URL match | `{index?: number, url?: string}` | |
| `browser_close_tab` | Close a tab | `{index?: number}` | |

## 6. CLI as Agent Skill Interface

### 6.1 Why This Matters

MCP is the right integration for agents that support it natively. But MCP requires the agent framework to implement the protocol — not all do, and new frameworks appear constantly. The CLI + daemon is the universal fallback: any agent that can run a shell command gets browser automation.

An "agent skill" is a prompt that tells an agent what tools it has. For CLI-based tools, the skill is just documentation of the commands. The agent reads the skill, shells out to `clicker`, parses the output, and acts on it. No SDK, no protocol integration, no client library.

This means the CLI surface design directly determines how good an agent can be at browser automation. A well-designed CLI is a well-designed skill.

### 6.2 CLI Design Guidelines for Agents

**Predictable.** Agents hallucinate flags when interfaces are inconsistent. Every command follows the same pattern:

```
clicker <command> [positional args] [--flags]
```

No subcommand nesting (`clicker browser navigate`), no ambiguous flag overloading. One command, one action.

**Self-describing.** An agent should be able to discover capabilities without a skill file:

```bash
$ clicker --help
Browser automation CLI. Browser persists between commands.

Commands:
  navigate <url>              Go to URL
  text [selector]             Get text content of page or element
  find <selector>             Find element by CSS selector
  find-all <selector>         Find all matching elements
  click <selector>            Click element
  type <selector> <text>      Type into input/textarea
  ...

$ clicker wait --help
Wait for an element to appear.

Usage: clicker wait <selector> [--timeout 5000] [--state visible]

Arguments:
  selector    CSS selector to wait for

Flags:
  --timeout   Max wait time in ms (default: 5000)
  --state     Wait condition: visible, attached, hidden (default: visible)
  --json      Output as JSON
```

**JSON-first for agents.** Every command supports `--json`. The output envelope is consistent across all commands so agents can parse reliably:

```bash
$ clicker title --json
{"ok":true,"result":"Example Domain"}

$ clicker find "h1" --json
{"ok":true,"result":{"selector":"h1","text":"Example Domain","tag":"h1","visible":true}}

$ clicker find ".nonexistent" --json
{"ok":false,"error":"No element found matching: .nonexistent"}
```

Human-readable output (no `--json`) is the default for interactive use — just the text content, no wrapper:

```bash
$ clicker title
Example Domain

$ clicker text "h1"
Example Domain
```

**Composable.** Commands do one thing. Agents chain them:

```bash
clicker navigate https://github.com/login
clicker type "#login_field" "username"
clicker type "#password" "hunter2"
clicker click "[type=submit]"
clicker wait ".dashboard"
clicker text ".dashboard h1"
```

Not: `clicker login --url https://github.com --user username --pass hunter2 --wait .dashboard --extract .dashboard h1`

**Non-zero exit codes on failure.** Agents check `$?`. If a command fails (element not found, navigation error, timeout), it exits non-zero and prints the error. This lets agents use simple conditional logic:

```bash
if clicker find ".error-message" --json 2>/dev/null; then
  # handle error state
fi
```

### 6.3 Reference Skill File

This is an example skill file that can be dropped into any agent framework. It's the entire integration — no SDK, no protocol, just a prompt:

````markdown
# Skill: Browser Automation (Vibium Clicker)

You have access to `clicker`, a browser automation CLI.
The browser launches automatically and persists between commands.
All commands support `--json` for structured output.

## Navigation
clicker navigate <url>              # Go to URL
clicker url                         # Get current URL
clicker title                       # Get page title
clicker scroll [up|down|left|right] # Scroll the page

## Reading
clicker text [selector]             # Get text content (page or element)
clicker find <selector>             # Find element, returns info
clicker find-all <selector>         # Find all matching elements
clicker eval <javascript>           # Execute JS, return result
clicker screenshot [-o file.png]    # Capture viewport

## Interaction
clicker click <selector>            # Click element
clicker type <selector> <text>      # Type into input/textarea
clicker keys <keys>                 # Keyboard: Enter, Tab, Escape, Control+a
clicker select <selector> <value>   # Select dropdown option
clicker hover <selector>            # Hover over element

## Waiting
clicker wait <selector>             # Wait for element (default 5s timeout)
clicker wait <selector> --timeout 10000 --state visible

## Tabs
clicker tabs                        # List open tabs
clicker tab-new [url]               # Open new tab
clicker tab-switch <index>          # Switch to tab
clicker tab-close [index]           # Close tab

## Session
clicker quit                        # Close browser (next command opens fresh one)
clicker daemon status               # Check daemon status

## Flags (available on all commands)
--json                              # Structured JSON output
--oneshot                           # Don't use daemon (one-shot mode)

## Tips
- CSS selectors: "#id", ".class", "tag", "[attr=value]", "tag.class"
- The browser is visible by default. You can watch it work.
- If a command fails, check the error with --json for details.
- `clicker eval` can run arbitrary JavaScript for anything the
  built-in commands don't cover.

## Example: Search GitHub
clicker navigate https://github.com
clicker type "[name=q]" "vibium"
clicker keys Enter
clicker wait ".repo-list"
clicker text ".repo-list" --json
````

This skill file works with Claude Code, Cursor, Windsurf, Codex, Gemini CLI, Aider, local models via ollama — anything that can be prompted and can run shell commands.

### 6.4 Generating the Skill File

The skill file should be generated from the CLI itself, not hand-maintained:

```bash
clicker skill > SKILL.md
```

This command introspects the registered commands, their arguments, and flags, and outputs a formatted skill file. When new tools are added, the skill file updates automatically. The generated output should be tuned for LLM consumption: concise, example-heavy, no ambiguity.

## 7. Migration Path

This is additive. Nothing breaks. The word "migration" is generous — it's really just "what gets built in what order."

### Phase 1: Daemon Core

Build the daemon process, socket listener, and session manager. The daemon imports the existing MCP tool handlers and exposes them over the socket.

**What changes:**
- New `clicker daemon start|stop|status|restart` subcommands.
- Existing CLI commands (`navigate`, `screenshot`, etc.) gain daemon awareness: if a daemon is running, use it; otherwise, one-shot as before.

**What doesn't change:**
- `clicker mcp` continues to work inline, exactly as today.
- `claude mcp add vibium -- npx -y vibium` continues to work.
- All existing tests pass without modification.

### Phase 2: New Tools

Add the tools from §5 to the shared tool handler set. They're available in both `clicker mcp` (inline) and the daemon simultaneously, since both use the same handlers. Corresponding CLI subcommands are added for daemon usage.

### Phase 3: HTTP Transport

Add `clicker http` adapter for Streamable HTTP MCP transport. Additive, no impact on existing users.

### Phase 4: `clicker mcp --connect`

Add the optional bridge mode so MCP clients like Claude Code can use a shared daemon session instead of inline mode.

### Phase 5: Multi-Session (Future)

Named sessions, multiple contexts. Deferred until there's demand.

## 8. Implementation Plan

### 7.1 Code Changes

The existing codebase is untouched. New code is added alongside it.

**New files:**

```
clicker/
  cmd/clicker/
    daemon.go             # NEW: daemon start/stop/status/restart subcommands
  internal/
    daemon/
      daemon.go           # NEW: main daemon loop, socket listener
      session.go          # NEW: browser session management (wraps existing BiDi code)
      router.go           # NEW: JSON-RPC method routing → existing tool handlers
    transport/
      socket.go           # NEW: Unix socket client (for CLI → daemon)
      http.go             # NEW (Phase 3): HTTP bridge → daemon
```

**Modified files (minimal changes):**

```
  cmd/clicker/
    navigate.go           # ADD: check for daemon, use socket if available
    screenshot.go         # ADD: check for daemon, use socket if available
    find.go               # ADD: same pattern
    click.go              # ADD: same pattern
    type.go               # ADD: same pattern
    eval.go               # ADD: same pattern
```

The pattern for modifying existing CLI commands is:

```go
func runNavigate(url string) {
    if daemon.IsRunning() {
        // Daemon exists → send via socket
        result := transport.Call("browser_navigate", map[string]any{"url": url})
        fmt.Println(result)
    } else if !oneshot {
        // No daemon → auto-start one, then send
        daemon.AutoStart()
        result := transport.Call("browser_navigate", map[string]any{"url": url})
        fmt.Println(result)
    } else {
        // --oneshot flag → existing one-shot behavior
        existingNavigate(url)
    }
}
```

**Shared tool handlers:**

The existing MCP tool handler functions (currently called from `clicker mcp`) need to be importable by the daemon's router. If they're already in a shared package, no changes needed. If they're embedded in the `mcp` command, extract them to a shared `internal/mcp/tools.go` — this is a move, not a rewrite.

### 7.2 Task Breakdown

**Phase 1 tasks (daemon core):**

1. **Implement `internal/daemon/daemon.go`**
   - Socket listener (Unix domain socket / Windows named pipe)
   - PID file and lock file management
   - Signal handling (SIGINT, SIGTERM → graceful shutdown)
   - Idle timeout (optional)

2. **Implement `internal/daemon/router.go`**
   - Accept JSON-RPC 2.0 messages from socket
   - Route `tools/call` to existing tool handlers
   - Route `daemon/*` to daemon management handlers
   - Return JSON-RPC responses

3. **Implement `internal/daemon/session.go`**
   - Lazy Chrome launch on first tool call
   - Session state (current URL, browsing context handle)
   - Chrome crash detection and recovery
   - Clean teardown on quit/shutdown

4. **Ensure tool handlers are importable**
   - If already in a shared package: no work.
   - If embedded in `clicker mcp` command: move to `internal/mcp/tools.go`. Pure code move, no logic changes.

5. **Implement `internal/transport/socket.go`**
   - Client-side socket connection for CLI → daemon
   - Connect, send JSON-RPC, receive response, disconnect

6. **Implement `cmd/clicker/daemon.go`**
   - `clicker daemon start [-d] [--idle-timeout]`
   - `clicker daemon stop`
   - `clicker daemon status`
   - `clicker daemon restart`

7. **Add daemon awareness to existing CLI commands**
   - Each command checks `daemon.IsRunning()`
   - If running: send via socket
   - If not running: auto-start daemon, then send via socket
   - `--oneshot` flag bypasses daemon entirely (legacy behavior)
   - Add `--json` output flag

**Phase 2 tasks (new tools + skill file):**

8. Add `browser_get_text`, `browser_get_html`, `browser_get_url`, `browser_get_title`
9. Add `browser_scroll`, `browser_hover`, `browser_select`, `browser_keys`
10. Add `browser_wait`
11. Add `browser_find_all`
12. Add `browser_new_tab`, `browser_list_tabs`, `browser_switch_tab`, `browser_close_tab`
13. Add corresponding CLI subcommands for each new tool
14. Implement `clicker skill` — auto-generate skill file from registered commands

Note: `browser_eval` already exists in both CLI and MCP.

**Phase 3 tasks (HTTP transport):**

15. Implement `internal/transport/http.go` — Streamable HTTP MCP bridge
16. Implement `clicker http [--port 9517]` subcommand

**Phase 4 tasks (`--connect` mode):**

17. Implement `clicker mcp --connect` — stdio bridge to daemon socket

### 7.3 Testing Strategy

- **Existing tests pass unchanged.** The inline `clicker mcp` path is untouched, so all current tests must continue to pass with zero modifications. This is the primary correctness check.
- **Daemon lifecycle tests:** start, stop, status, restart, stale PID cleanup, idle timeout.
- **Socket communication tests:** connect, send/receive JSON-RPC, handle disconnect, handle concurrent clients.
- **Tool-via-daemon tests:** each tool called through the daemon socket produces the same results as calling it through inline `clicker mcp`.
- **CLI daemon-awareness tests:** CLI commands use daemon when running, fall back to one-shot when not.
- **Transport bridge tests:** `clicker mcp --connect` round-trip, HTTP bridge round-trip.
- **Crash recovery tests:** kill Chrome mid-session, verify daemon reports error and recovers on next command.

## 9. Deferred

These are explicitly out of scope for this design but noted for future consideration.

- **Named sessions** (`--session work`): multiple Chrome instances managed by one daemon. Needed for multi-agent setups. Requires extending the session manager and adding session identifiers to every tool call.
- **Named contexts** (`--context task1`): multiple isolated browsing contexts within one session. Useful for parallel agent tasks that shouldn't share cookies/state.
- **MCP Resources**: exposing page URL, title, and DOM as MCP resources (read-only, no tool call needed). Requires MCP resource support in the daemon.
- **MCP Notifications**: push events from daemon to clients (page navigated, element appeared, console error). Requires bidirectional communication over the socket, which JSON-RPC supports but needs careful design for the stdio bridge.
- **Remote daemon**: daemon listening on a network interface instead of local socket, for running Chrome on a remote machine. Security implications need careful design (auth, TLS).

## 10. Open Questions

1. **Should `browser_launch` become a no-op?** Currently agents call `browser_launch` as their first action. With the daemon, Chrome launches lazily on first use. Should `browser_launch` still exist as an explicit "ensure browser is ready" signal, or should it be silently accepted and ignored?

2. **Socket vs. TCP for local communication.** Unix sockets are faster and don't require port allocation, but add platform-specific code. Named pipes on Windows behave differently. Is it worth supporting TCP as a fallback even on Unix for simplicity?

3. **What happens to `--headless` and `--wait-open`?** These are currently per-command flags. With the daemon, they become session-level configuration set at daemon start or first browser launch. How should per-command flags that are really session-level config be handled?

4. **Concurrent CLI commands.** Two terminal windows both send `clicker click` at the same time. The daemon processes them sequentially on the same browser — is that correct, or should there be explicit locking/queueing with user-visible feedback?

5. **`clicker mcp --connect` config story.** If a user wants Claude Code to use the daemon instead of inline mode, they'd need to change their MCP config to `clicker mcp --connect`. Should there be a helper command like `clicker mcp config --connect` that patches the Claude Code config? Or is this too niche to worry about early on?

6. **`--json` output envelope.** The design proposes a simple `{"ok":true,"result":...}` / `{"ok":false,"error":...}` envelope for CLI `--json` output. Should this match the MCP response shape (`{"content":[{"type":"text","text":...}]}`) for consistency, or is a simpler agent-friendly envelope better? MCP shape is more verbose but means one format everywhere.

7. **Skill file scope.** Should the `clicker skill` command generate a minimal skill (just commands and examples) or a comprehensive one (including tips, error handling patterns, CSS selector help)? Longer skill files consume more agent context. Could offer `clicker skill --minimal` and `clicker skill --full`.

8. **Default idle timeout value.** 30 minutes is proposed. Too short and agents doing long tasks lose their session. Too long and forgotten daemons waste resources. Should this be configurable per-agent via an environment variable (e.g., `VIBIUM_IDLE_TIMEOUT=1h`)?
