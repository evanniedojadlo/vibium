# Clicker Daemon Implementation Plan

## Context

Every CLI command today is one-shot: launch Chrome, execute, tear down. This makes multi-step workflows impossible and forces every agent to pay browser startup cost per command. The daemon keeps Chrome alive across commands, enabling session persistence through the same MCP tool handlers already in `internal/mcp/handlers.go`.

The design doc is at `docs/designs/clicker-daemon-design.md`. This plan covers Phase 1 (daemon core) broken into 9 independently committable steps. Phases 2-4 (new tools, HTTP transport, MCP bridge) come later.

---

## Step 1: Extract CLI commands from main.go into separate files

**Why:** `main.go` is 802 lines with all 14 commands inline. Every daemon step touches command files. Splitting now prevents merge conflicts later.

**Files created:**
- `cmd/clicker/navigate.go` — `newNavigateCmd() *cobra.Command`
- `cmd/clicker/screenshot.go` — `newScreenshotCmd()`
- `cmd/clicker/eval.go` — `newEvalCmd()`
- `cmd/clicker/find.go` — `newFindCmd()`
- `cmd/clicker/click.go` — `newClickCmd()`
- `cmd/clicker/type_cmd.go` — `newTypeCmd()`
- `cmd/clicker/check_actionable.go` — `newCheckActionableCmd()`
- `cmd/clicker/serve.go` — `newServeCmd()`
- `cmd/clicker/mcp_cmd.go` — `newMCPCmd()`
- `cmd/clicker/helpers.go` — `doWaitOpen()`, `waitAndClose()`, `printCheck()`
- `cmd/clicker/diagnostics.go` — version, paths, install, launch-test, ws-test, bidi-test

**Files modified:**
- `cmd/clicker/main.go` — Slim to `main()`, root command, global flags, `AddCommand()` calls (~50 lines)

**Pattern:** Each file returns a `*cobra.Command`. Global vars (`headless`, `waitOpen`, etc.) stay in `main.go` — accessible from all files since they're in `package main`.

**Verify:** `make test` passes. `clicker --help` output identical. Every subcommand works as before. Pure code move, zero behavior change.

---

## Step 2: Add socket/PID path helpers and daemon package skeleton

**Why:** Establish filesystem conventions (socket path, PID file, lock file) and the `IsRunning()` check that everything else depends on.

**Files created:**
- `internal/daemon/status.go` — `IsRunning() bool` (checks PID file + socket connectivity)
- `internal/daemon/pidfile.go` — `WritePID()`, `ReadPID()`, `RemovePID()`, `CleanStale()`

**Files modified:**
- `internal/paths/paths.go` — Add `GetSocketPath()`, `GetDaemonDir()`

**Socket paths** (reuses existing `GetCacheDir()`):
| Platform | Socket | PID |
|----------|--------|-----|
| macOS | `~/Library/Caches/vibium/clicker.sock` | `~/Library/Caches/vibium/clicker.pid` |
| Linux | `~/.cache/vibium/clicker.sock` | `~/.cache/vibium/clicker.pid` |
| Windows | `\\.\pipe\vibium-clicker` (named pipe) | `%LOCALAPPDATA%\vibium\clicker.pid` |

**`IsRunning()` logic:**
1. Read PID from file
2. Check if process exists (`os.FindProcess` + signal 0 on Unix)
3. Check if socket is connectable (`net.DialTimeout`)
4. Return true only if both pass

**Verify:** `cd clicker && go test ./internal/daemon/...` — test `IsRunning` returns false when no files exist, `WritePID`/`ReadPID` round-trips, `CleanStale` removes dead PID files.

---

## Step 3: Implement daemon socket listener and JSON-RPC router

**Why:** This is the core — the daemon process that listens on a socket and routes requests to the existing `mcp.Handlers`.

**Files created:**
- `internal/daemon/daemon.go` — `Daemon` struct with `New()`, `Run(ctx)`, `Shutdown()`
- `internal/daemon/router.go` — `handleConnection()`, JSON-RPC routing
- `internal/daemon/listener_unix.go` — `listen()` via `net.Listen("unix", path)`
- `internal/daemon/listener_windows.go` — stub (returns error with TODO)

**Key design:**
```go
type Daemon struct {
    listener     net.Listener
    handlers     *mcp.Handlers
    mu           sync.Mutex
    version      string
    startTime    time.Time
    lastActivity time.Time
    idleTimeout  time.Duration
}
```

**Routing:** The daemon duplicates the ~20-line `route()` switch from `mcp/server.go` (adding `daemon/status` and `daemon/shutdown`). This is intentional — the daemon needs additional methods and the existing server is tightly coupled to stdio.

**Concurrency:** Existing `mcp.Handlers` is single-client (designed for stdio). The daemon wraps all `handlers.Call()` in `d.mu.Lock()` to serialize. Multiple CLI connections wait their turn.

**Verify:** Go test — start daemon on temp socket, send `daemon/status` via `net.Dial`, verify response. Send `daemon/shutdown`, verify clean exit. No Chrome needed.

---

## Step 4: Implement socket client (CLI → daemon transport)

**Why:** CLI commands need a simple function to send a request to the daemon and get a response.

**Files created:**
- `internal/daemon/client.go` — `Call()`, `Status()`, `Shutdown()`

**API:**
```go
func Call(toolName string, args map[string]interface{}) (*mcp.ToolsCallResult, error)
func Status() (*StatusResult, error)
func Shutdown() error
```

**Implementation:** `net.DialTimeout` → write JSON-RPC + newline → read response line → parse → close. One connection per CLI invocation (no persistent connections needed).

**Timeouts:** 2s dial timeout, 60s read timeout (browser ops can be slow).

**Verify:** Go test — start daemon from Step 3, use client functions, verify round-trip. `make test` still passes.

---

## Step 5: Add `--json` flag and lazy browser launch

**Why:** Two prerequisites for daemon-mode CLI. `--json` gives agents structured output. Lazy launch means the daemon starts cheap (no Chrome) and launches Chrome on first tool call.

### 5a: `--json` flag

**Files created:**
- `cmd/clicker/output.go` — `printResult()`, `printError()`, `printJSON()`

**Files modified:**
- `cmd/clicker/main.go` — Add `jsonOutput bool` global + `--json` persistent flag

**Output envelope:**
```json
{"ok":true,"result":"Example Domain"}
{"ok":false,"error":"No element found matching: .nope"}
```

### 5b: Lazy browser launch

**Files modified:**
- `internal/mcp/handlers.go`:
  - Add `headless bool` field to `Handlers`
  - Change `NewHandlers(screenshotDir string)` → `NewHandlers(screenshotDir string, headless bool)`
  - Change `ensureBrowser()` to auto-call `browserLaunch()` instead of returning error
  - Add no-op guard to `browserLaunch()` when browser is already running
- `cmd/clicker/mcp_cmd.go` — Update `NewHandlers` call to pass `headless: false`

**Verify:** `clicker navigate https://example.com --json` outputs JSON envelope. MCP server auto-launches browser on first `browser_navigate` (no need for `browser_launch` first). All `make test` passes — update MCP test that expects "Call browser_launch first" error.

---

## Step 6: Add `clicker daemon start|stop|status` commands

**Why:** Proves the full daemon + client stack works end-to-end. Power users get explicit control.

**Files created:**
- `cmd/clicker/daemon_cmd.go` — `newDaemonCmd()` with subcommands

**Commands:**
- `clicker daemon start` — Foreground, listens on socket, writes PID file
- `clicker daemon start -d` — Daemonize (re-exec as detached child, poll socket, exit)
- `clicker daemon stop` — Send `daemon/shutdown` via client
- `clicker daemon status` — Send `daemon/status`, print info (or "not running")
- `--idle-timeout 30m` flag on start
- `--headless` passed through to handlers

**Daemonization:** Go has no `fork()`. Use `os.StartProcess` with `Setsid: true` (Unix) to detach. Parent polls socket availability (max 5s), prints PID, exits.

**Verify:** Manual — `clicker daemon start -d`, `clicker daemon status` shows running, `clicker daemon stop` shuts it down cleanly. Socket and PID files removed.

---

## Step 7: Add daemon awareness to existing CLI commands + auto-start

**Why:** The main feature — CLI commands transparently use the daemon. Browser persists between commands.

**Files created:**
- `cmd/clicker/daemon_client.go` — `daemonCall()` helper with auto-start logic

**Files modified:**
- `cmd/clicker/main.go` — Add `--oneshot` persistent flag, `VIBIUM_ONESHOT` env var support
- `cmd/clicker/navigate.go` — Add daemon path before existing oneshot path
- `cmd/clicker/screenshot.go` — Same pattern
- `cmd/clicker/find.go` — Same
- `cmd/clicker/click.go` — Same
- `cmd/clicker/type_cmd.go` — Same
- `cmd/clicker/eval.go` — Same

**Pattern for each command:**
```go
if !oneshot {
    result, err := daemonCall("browser_navigate", map[string]any{"url": url})
    printResult(result)
    return
}
// ... existing oneshot code unchanged ...
```

**Auto-start logic** (in `daemonCall`):
1. Try `daemon.Call()` → success? Return result.
2. Connection refused? Spawn `clicker daemon start --_internal` as detached process.
3. Poll socket (backoff, max 5s).
4. Retry `daemon.Call()`.

**CLI argument adaptation:** Current commands take `[url] [selector]`. In daemon mode, URL is optional — if first arg looks like a URL (`http://` or `https://`), navigate first then execute. If it's just a selector, execute directly on current page. In oneshot mode: unchanged.

**Existing test compatibility:** Add `VIBIUM_ONESHOT=1` env var check. Update `Makefile` test targets to set this. All existing tests run in oneshot mode — zero test file changes.

**Verify:** `make test` passes (via `VIBIUM_ONESHOT=1`). Manual: `clicker navigate https://example.com` → auto-starts daemon → navigates. `clicker find "h1"` → reuses session. `clicker navigate https://google.com --oneshot` → old behavior.

---

## Step 8: Idle timeout and graceful shutdown

**Why:** Forgotten daemons shouldn't waste resources. Clean shutdown prevents orphaned Chrome processes.

**Files modified:**
- `internal/daemon/daemon.go` — Add idle timeout goroutine, signal handling, `Shutdown()` cleanup sequence

**Idle timeout:** Goroutine checks `time.Since(lastActivity)` every minute. `touchActivity()` called on every incoming request.

**Shutdown sequence:**
1. Cancel context (stops accept loop)
2. `handlers.Close()` (closes BiDi connection, terminates Chrome)
3. Close listener
4. `os.Remove(socketPath)` (Unix socket cleanup)
5. Remove PID file

**Signal handling:** Daemon command installs its own `SIGINT`/`SIGTERM` handler that calls `Shutdown()`. Does NOT use `process.WithCleanup()` (that's for oneshot commands).

**Verify:** `clicker daemon start --idle-timeout 5s` → wait 5s → `clicker daemon status` shows "not running". Socket/PID files cleaned up. `kill <daemon-pid>` → clean Chrome shutdown.

---

## Step 9: Integration tests

**Why:** Verify the full daemon lifecycle through the CLI binary.

**Files created:**
- `tests/daemon/lifecycle.test.js` — Start/stop, navigate+find across commands, auto-start, idle timeout, stale PID cleanup, `--oneshot` bypass
- `tests/daemon/concurrency.test.js` — Rapid sequential commands, error recovery

**Files modified:**
- `Makefile` — Add `test-daemon` target, add to `test` target

**Test pattern** (matches existing style using Node.js `node:test`):
```javascript
const result = execSync('clicker navigate https://example.com --json');
const parsed = JSON.parse(result);
assert(parsed.ok === true);
```

**Verify:** `make test-daemon` passes. `make test` passes (includes daemon + existing tests).

---

## Summary

```
Step 1: Extract main.go          (pure refactor, no behavior change)
Step 2: Paths + PID helpers       (new internal/daemon/ skeleton)
Step 3: Daemon listener + router  (core daemon loop)
Step 4: Socket client             (CLI → daemon transport)
Step 5: --json + lazy launch      (prerequisites for daemon CLI)
Step 6: daemon start/stop/status  (explicit daemon management)
Step 7: Daemon-aware CLI          (the main feature)
Step 8: Idle timeout + shutdown   (production readiness)
Step 9: Integration tests         (verify everything)
```

Each step is independently committable and keeps `make test` passing. Steps 1-2 have no dependencies on each other. Steps 3-9 are sequential.

## Key files

| File | Role |
|------|------|
| `cmd/clicker/main.go` | 802-line monolith → slim entry point (Step 1) |
| `internal/mcp/handlers.go` | 8 tool handlers reused by daemon (Step 5 modifies) |
| `internal/mcp/server.go` | JSON-RPC types + routing pattern (daemon mirrors) |
| `internal/paths/paths.go` | Platform paths — add socket path (Step 2) |
| `internal/daemon/` | New package — daemon, client, status, PID (Steps 2-4) |
