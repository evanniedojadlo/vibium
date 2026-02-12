# Plan: Make `make` Work on Windows

## Context

We're setting up Windows as a dev environment. The Makefile uses Unix commands (ln, pkill, sed, etc.) that don't exist on Windows. However, `winget install Git.Git` includes **Git Bash** which provides most of these tools (cp, rm, mkdir, cat, sleep, bash if/for, etc.). GNU Make can be configured to use Git Bash as its shell, so the fix is targeted — not a full rewrite.

The user wants to `make build` and `make test` on Windows.

## Changes

### 1. Makefile — add Windows support

**File: `Makefile`**

Add OS detection preamble at the top (before `.PHONY`):
```makefile
# Windows: use Git Bash as shell so Unix commands (cp, rm, mkdir, etc.) work
ifeq ($(OS),Windows_NT)
  SHELL := C:/Program Files/Git/usr/bin/bash
  .SHELLFLAGS := -c
  EXE := .exe
else
  EXE :=
endif
```

Then apply `$(EXE)` suffix everywhere the binary is referenced:
- `build-go`: `go build -o bin/clicker$(EXE)`, replace `ln -sf` with `cp`
- `install-browser`: `./clicker/bin/clicker$(EXE) install`
- `serve`: `./clicker/bin/clicker$(EXE) serve`
- Node_modules copy block: use `$(EXE)` on target path

Platform-conditional targets:
- `double-tap`: use `taskkill /F /IM chrome.exe` on Windows, `pkill` elsewhere
- `clean-cache`: use `$$LOCALAPPDATA/vibium/...` on Windows, `~/Library/Caches/...` and `~/.cache/...` elsewhere

### 2. Test files — fix binary path + hardcoded /tmp/

**New file: `tests/helpers.js`**
```javascript
const path = require('node:path');
const EXE = process.platform === 'win32' ? '.exe' : '';
const CLICKER = path.join(__dirname, '../clicker/bin/clicker') + EXE;
module.exports = { CLICKER };
```

**10 test files** — replace `const CLICKER = path.join(...)` with:
```javascript
const { CLICKER } = require('../helpers');
```

Files:
- `tests/cli/navigation.test.js` (also fix `/tmp/` to `os.tmpdir()`)
- `tests/cli/elements.test.js`
- `tests/cli/actionability.test.js`
- `tests/cli/page-reading.test.js`
- `tests/cli/input-tools.test.js`
- `tests/cli/tabs.test.js`
- `tests/cli/process.test.js`
- `tests/daemon/lifecycle.test.js`
- `tests/daemon/concurrency.test.js`
- `tests/mcp/server.test.js`

### 3. Windows setup doc — add Make install + build instructions

**File: `docs/local-dev-setup-x86-windows.md`**

- Add `winget install GnuWin32.Make` to dev tools section
- Add note: must add `C:\Program Files (x86)\GnuWin32\bin` to PATH
- Update Build and Test section to use `make build && make test`

## Out of scope

- `set-version` target (uses macOS-specific `sed -i ''` — maintainer-only task)
- `package-python` target (uses Python venv activation with Unix paths)
- Windows process cleanup verification in process tests (already returns empty Set gracefully)

## Verification

1. Run `make test` on macOS — confirm no regressions
2. Run `make build` then `make test` on Windows VM in Git Bash
