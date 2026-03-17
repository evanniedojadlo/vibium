# Client Implementation Guide

> **Draft:** This is a work-in-progress draft that may be used to generate client libraries for additional languages in the future.

Reference for implementing Vibium clients in new languages (Java, C#, Ruby, Kotlin, Swift, Rust, Go, Nim, etc.).

Use the **JS client** (`clients/javascript/`) and **Python client** (`clients/python/`) as reference implementations.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Class Hierarchy](#class-hierarchy)
3. [Command Reference](#command-reference)
4. [Naming Conventions](#naming-conventions)
5. [Error Types](#error-types)
6. [Async / Sync Patterns](#async--sync-patterns)
7. [Reserved Keyword Handling](#reserved-keyword-handling)
8. [Aliases](#aliases)
9. [Binary Discovery](#binary-discovery)
10. [Testing Checklist](#testing-checklist)

---

## Architecture Overview

```
┌────────────────┐  stdin/stdout  ┌─────────────┐    BiDi/WS     ┌─────────┐
│ Client (JS,    │◄──────────────►│   vibium    │◄──────────────►│ Chrome  │
│ Python, etc.)  │ ndjson (pipes) │   binary    │ WebDriver BiDi │ browser │
└────────────────┘                └─────────────┘                └─────────┘
```

1. Client spawns the `vibium pipe` command as a subprocess
2. Client communicates via newline-delimited JSON over stdin/stdout
3. The binary sends a `vibium:lifecycle.ready` signal on stdout once the browser is launched
4. `vibium:` extension commands are handled by the binary; standard BiDi commands are forwarded to Chrome

### Message Format

**Request** (client → vibium):
```json
{"id": 1, "method": "vibium:page.navigate", "params": {"context": "ctx-1", "url": "https://example.com"}}
```

**Success response** (vibium → client):
```json
{"id": 1, "type": "success", "result": {}}
```

**Error response** (vibium → client):
```json
{"id": 1, "type": "error", "error": "timeout", "message": "Timeout after 30000ms waiting for '#btn'"}
```

**Event** (vibium → client, no `id`):
```json
{"method": "browsingContext.load", "params": {"context": "ctx-1", "url": "https://example.com"}}
```

---

## Class Hierarchy

All clients must implement these classes:

```
Browser                  ← manages browser lifecycle
├── .context             ← default BrowserContext (property)
├── .keyboard            ← (accessed via Page)
├── .mouse               ← (accessed via Page)
└── .touch               ← (accessed via Page)

BrowserContext            ← cookie/storage isolation boundary
├── .recording           ← Recording (property)
└── newPage()            ← creates Page

Page                      ← a browser tab
├── .keyboard            ← Keyboard (property)
├── .mouse               ← Mouse (property)
├── .touch               ← Touch (property)
├── .clock               ← Clock (property)
├── .context             ← back-reference to BrowserContext
├── find() / findAll()   ← returns Element(s)
├── route()              ← creates Route via callback
├── onDialog()           ← creates Dialog via callback
├── onConsole()          ← creates ConsoleMessage via callback
├── onDownload()         ← creates Download via callback
├── onRequest()          ← creates Request via callback
├── onResponse()         ← creates Response via callback
└── onWebSocket()        ← creates WebSocketInfo via callback

Element                   ← a resolved DOM element
├── click/fill/type/...  ← interaction methods
├── text/html/value/...  ← state query methods
└── find() / findAll()   ← scoped element search

Keyboard                  ← page-level keyboard input
Mouse                     ← page-level mouse input
Touch                     ← page-level touch input
Clock                     ← fake timer control
Recording                 ← trace recording control
Route                     ← network interception handler
  └── .request           ← Request (property)
Dialog                    ← browser dialog (alert/confirm/prompt)
Request                   ← network request info
Response                  ← network response info
Download                  ← file download handle
ConsoleMessage            ← console.log() message
WebSocketInfo             ← WebSocket connection info
```

### Data Types

These should be structured types (interfaces/structs), not raw dicts:

| Type | Fields |
|---|---|
| `Cookie` | `name`, `value`, `domain`, `path`, `size`, `httpOnly`, `secure`, `sameSite`, `expiry?` |
| `SetCookieParam` | `name`, `value`, `domain?`, `url?`, `path?`, `httpOnly?`, `secure?`, `sameSite?`, `expiry?` |
| `StorageState` | `cookies: Cookie[]`, `origins: OriginState[]` |
| `OriginState` | `origin`, `localStorage: {name, value}[]`, `sessionStorage: {name, value}[]` |
| `BoundingBox` | `x`, `y`, `width`, `height` |
| `ElementInfo` | `tag`, `text`, `box: BoundingBox` |
| `A11yNode` | `role`, `name?`, `value?`, `description?`, `disabled?`, `expanded?`, `focused?`, `checked?`, `pressed?`, `selected?`, `level?`, `multiselectable?`, `children?: A11yNode[]` |
| `ScreenshotOptions` | `fullPage?`, `clip?: {x, y, width, height}` |
| `FindOptions` | `timeout?` |

---

## Command Reference

All extension commands use the `vibium:` prefix. Standard WebDriver BiDi commands (e.g., `browsingContext.getTree`, `session.subscribe`) are forwarded directly to Chrome.

### Browser

| Wire Command | Description | JS | Python |
|---|---|---|---|
| *binary launch + WebSocket connect* | Launch browser and connect | `browser.start(opts?)` | `browser.start(opts?)` |
| `vibium:browser.page` | Get the default page | `page()` | `page()` |
| `vibium:browser.newPage` | Create a new page | `newPage()` | `new_page()` |
| `vibium:browser.newContext` | Create a new browser context | `newContext()` | `new_context()` |
| `vibium:browser.pages` | List all pages | `pages()` | `pages()` |
| `vibium:browser.stop` | Stop the browser | `stop()` | `stop()` |
| *client-side event listener* | Listen for new page events | `onPage(cb)` | `on_page(cb)` |
| *client-side event listener* | Listen for popup events | `onPopup(cb)` | `on_popup(cb)` |
| *client-side* | Remove all event listeners | `removeAllListeners(ev?)` | `remove_all_listeners(ev?)` |

### Page

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:page.navigate` | Navigate to a URL | `go(url)` | `go(url)` |
| `vibium:page.back` | Go back | `back()` | `back()` |
| `vibium:page.forward` | Go forward | `forward()` | `forward()` |
| `vibium:page.reload` | Reload page | `reload()` | `reload()` |
| `vibium:page.url` | Get current URL | `url()` | `url()` |
| `vibium:page.title` | Get page title | `title()` | `title()` |
| `vibium:page.content` | Get page HTML | `content()` | `content()` |
| `vibium:page.find` | Find a single element | `find(sel, opts?)` | `find(sel, **opts)` |
| `vibium:page.findAll` | Find all matching elements | `findAll(sel, opts?)` | `find_all(sel, **opts)` |
| `vibium:page.screenshot` | Take a page screenshot | `screenshot(opts?)` | `screenshot(opts?)` |
| `vibium:page.pdf` | Generate PDF | `pdf()` | `pdf()` |
| `vibium:page.eval` | Evaluate JavaScript | `evaluate(expr)` | `evaluate(expr)` |
| `vibium:page.addScript` | Add a script tag | `addScript(src)` | `add_script(src)` |
| `vibium:page.addStyle` | Add a style tag | `addStyle(src)` | `add_style(src)` |
| `vibium:page.expose` | Expose a function to the page | `expose(name, fn)` | `expose(name, fn)` |
| `vibium:page.wait` | Wait for a duration | `wait(ms)` | `wait(ms)` |
| `vibium:page.waitFor` | Wait for a selector | `waitFor(sel, opts?)` | `wait_for(sel, **opts)` |
| `vibium:page.waitForFunction` | Wait for a JS function to return truthy | `waitForFunction(fn, opts?)` | `wait_for_function(fn, **opts)` |
| `vibium:page.waitForURL` | Wait for URL to match | `waitForURL(url, opts?)` | `wait_for_url(url, **opts)` |
| `vibium:page.waitForLoad` | Wait for page load | `waitForLoad(opts?)` | `wait_for_load(**opts)` |
| `vibium:page.scroll` | Scroll the page | `scroll(dir?, amt?, sel?)` | `scroll(dir?, amt?, sel?)` |
| `vibium:page.setViewport` | Set viewport size | `setViewport(size)` | `set_viewport(size)` |
| `vibium:page.viewport` | Get viewport size | `viewport()` | `viewport()` |
| `vibium:page.emulateMedia` | Override CSS media features | `emulateMedia(opts)` | `emulate_media(**opts)` |
| `vibium:page.setContent` | Set page HTML | `setContent(html)` | `set_content(html)` |
| `vibium:page.setGeolocation` | Override geolocation | `setGeolocation(coords)` | `set_geolocation(coords)` |
| `vibium:page.setWindow` | Set window size/position | `setWindow(opts)` | `set_window(**opts)` |
| `vibium:page.window` | Get window info | `window()` | `window()` |
| `vibium:page.a11yTree` | Get accessibility tree | `a11yTree(opts?)` | `a11y_tree(opts?)` |
| `vibium:page.frames` | List all frames | `frames()` | `frames()` |
| `vibium:page.frame` | Get a frame by name/URL | `frame(nameOrUrl)` | `frame(name_or_url)` |
| *returns self (top frame)* | Get the main frame | `mainFrame()` | `main_frame()` |
| `browsingContext.activate` | Bring page to front | `bringToFront()` | `bring_to_front()` |
| `browsingContext.close` | Close the page | `close()` | `close()` |
| `vibium:page.route` | Register a route handler | `route(pattern, handler)` | `route(pattern, handler)` |
| `network.removeIntercept` | Remove a route handler | `unroute(pattern)` | `unroute(pattern)` |
| `vibium:page.setHeaders` | Set extra HTTP headers | `setHeaders(headers)` | `set_headers(headers)` |
| *client-side event listener* | Listen for requests | `onRequest(fn)` | `on_request(fn)` |
| *client-side event listener* | Listen for responses | `onResponse(fn)` | `on_response(fn)` |
| *client-side event listener* | Listen for dialogs | `onDialog(fn)` | `on_dialog(fn)` |
| *client-side event listener* | Listen for console messages | `onConsole(fn)` | `on_console(fn)` |
| *client-side event listener* | Listen for page errors | `onError(fn)` | `on_error(fn)` |
| *client-side event listener* | Listen for downloads | `onDownload(fn)` | `on_download(fn)` |
| `vibium:page.onWebSocket` | Subscribe to WebSocket events | `onWebSocket(fn)` | `on_web_socket(fn)` |
| *client-side* | Remove all event listeners | `removeAllListeners(ev?)` | `remove_all_listeners(ev?)` |

### Element

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:element.click` | Click an element | `click(opts?)` | `click(timeout?)` |
| `vibium:element.dblclick` | Double-click an element | `dblclick(opts?)` | `dblclick(timeout?)` |
| `vibium:element.fill` | Fill an input field | `fill(value, opts?)` | `fill(value, timeout?)` |
| `vibium:element.type` | Type text character by character | `type(text, opts?)` | `type(text, timeout?)` |
| `vibium:element.press` | Press a key on a focused element | `press(key, opts?)` | `press(key, timeout?)` |
| `vibium:element.clear` | Clear an input field | `clear(opts?)` | `clear(timeout?)` |
| `vibium:element.check` | Check a checkbox | `check(opts?)` | `check(timeout?)` |
| `vibium:element.uncheck` | Uncheck a checkbox | `uncheck(opts?)` | `uncheck(timeout?)` |
| `vibium:element.selectOption` | Select a dropdown option | `selectOption(val, opts?)` | `select_option(val, timeout?)` |
| `vibium:element.hover` | Hover over an element | `hover(opts?)` | `hover(timeout?)` |
| `vibium:element.focus` | Focus an element | `focus(opts?)` | `focus(timeout?)` |
| `vibium:element.dragTo` | Drag an element to a target | `dragTo(target, opts?)` | `drag_to(target, timeout?)` |
| `vibium:element.tap` | Tap an element (touch) | `tap(opts?)` | `tap(timeout?)` |
| `vibium:element.scrollIntoView` | Scroll element into view | `scrollIntoView(opts?)` | `scroll_into_view(timeout?)` |
| `vibium:element.dispatchEvent` | Dispatch a DOM event on an element | `dispatchEvent(type, init?)` | `dispatch_event(type, init?)` |
| `vibium:element.setFiles` | Set files on a file input | `setFiles(files, opts?)` | `set_files(files, timeout?)` |
| `vibium:element.text` | Get element text content | `text()` | `text()` |
| `vibium:element.innerText` | Get element inner text | `innerText()` | `inner_text()` |
| `vibium:element.html` | Get element outer HTML | `html()` | `html()` |
| `vibium:element.value` | Get input element value | `value()` | `value()` |
| `vibium:element.attr` | Get element attribute | `attr(name)` | `attr(name)` |
| `vibium:element.bounds` | Get element bounding box | `bounds()` | `bounds()` |
| `vibium:element.isVisible` | Check if element is visible | `isVisible()` | `is_visible()` |
| `vibium:element.isHidden` | Check if element is hidden | `isHidden()` | `is_hidden()` |
| `vibium:element.isEnabled` | Check if element is enabled | `isEnabled()` | `is_enabled()` |
| `vibium:element.isChecked` | Check if element is checked | `isChecked()` | `is_checked()` |
| `vibium:element.isEditable` | Check if element is editable | `isEditable()` | `is_editable()` |
| `vibium:element.role` | Get element ARIA role | `role()` | `role()` |
| `vibium:element.label` | Get element accessible label | `label()` | `label()` |
| `vibium:element.screenshot` | Screenshot an element | `screenshot()` | `screenshot()` |
| `vibium:element.waitFor` | Wait for element state | `waitUntil(state?, opts?)` | `wait_until(state?, timeout?)` |
| `vibium:element.find` | Find a single element (scoped) | `find(sel, opts?)` | `find(sel, **opts)` |
| `vibium:element.findAll` | Find all matching elements (scoped) | `findAll(sel, opts?)` | `find_all(sel, **opts)` |

### BrowserContext

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:context.newPage` | Create a page in a context | `newPage()` | `new_page()` |
| `browser.removeUserContext` | Close the context | `close()` | `close()` |
| `vibium:context.cookies` | Get cookies | `cookies(urls?)` | `cookies(urls?)` |
| `vibium:context.setCookies` | Set cookies | `setCookies(cookies)` | `set_cookies(cookies)` |
| `vibium:context.clearCookies` | Clear cookies | `clearCookies()` | `clear_cookies()` |
| `vibium:context.storage` | Get storage state | `storage()` | `storage()` |
| `vibium:context.setStorage` | Set storage state | `setStorage(state)` | `set_storage(state)` |
| `vibium:context.clearStorage` | Clear all storage | `clearStorage()` | `clear_storage()` |
| `vibium:context.addInitScript` | Add an init script | `addInitScript(script)` | `add_init_script(script)` |

### Keyboard

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:keyboard.press` | Press a key | `press(key)` | `press(key)` |
| `vibium:keyboard.down` | Key down | `down(key)` | `down(key)` |
| `vibium:keyboard.up` | Key up | `up(key)` | `up(key)` |
| `vibium:keyboard.type` | Type text | `type(text)` | `type(text)` |

### Mouse

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:mouse.click` | Click at coordinates | `click(x, y, opts?)` | `click(x, y, **opts)` |
| `vibium:mouse.move` | Move mouse | `move(x, y, opts?)` | `move(x, y, **opts)` |
| `vibium:mouse.down` | Mouse button down | `down(opts?)` | `down(**opts)` |
| `vibium:mouse.up` | Mouse button up | `up(opts?)` | `up(**opts)` |
| `vibium:mouse.wheel` | Scroll mouse wheel | `wheel(dx, dy)` | `wheel(dx, dy)` |

### Touch

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:touch.tap` | Tap at coordinates | `tap(x, y)` | `tap(x, y)` |

### Clock

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:clock.install` | Install fake timers | `install(opts?)` | `install(time?, timezone?)` |
| `vibium:clock.fastForward` | Fast-forward time | `fastForward(ticks)` | `fast_forward(ticks)` |
| `vibium:clock.runFor` | Run timers for a duration | `runFor(ticks)` | `run_for(ticks)` |
| `vibium:clock.pauseAt` | Pause clock at a time | `pauseAt(time)` | `pause_at(time)` |
| `vibium:clock.resume` | Resume clock | `resume()` | `resume()` |
| `vibium:clock.setFixedTime` | Set fixed fake time | `setFixedTime(time)` | `set_fixed_time(time)` |
| `vibium:clock.setSystemTime` | Set system time | `setSystemTime(time)` | `set_system_time(time)` |
| `vibium:clock.setTimezone` | Set timezone | `setTimezone(tz)` | `set_timezone(tz)` |

### Recording

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:recording.start` | Start recording | `start(opts?)` | `start(opts?)` |
| `vibium:recording.stop` | Stop recording, return trace | `stop(opts?)` | `stop(path?)` |
| `vibium:recording.startChunk` | Start a recording chunk | `startChunk(opts?)` | `start_chunk(opts?)` |
| `vibium:recording.stopChunk` | Stop a recording chunk | `stopChunk(opts?)` | `stop_chunk(path?)` |
| `vibium:recording.startGroup` | Start a logical group | `startGroup(name, opts?)` | `start_group(name, location?)` |
| `vibium:recording.stopGroup` | Stop a logical group | `stopGroup()` | `stop_group()` |

### Route

| Wire Command | Description | JS | Python |
|---|---|---|---|
| — | Request being intercepted | `.request` (property) | *passed via callback args* |
| `vibium:network.fulfill` | Fulfill an intercepted request | `fulfill(resp?)` | `fulfill(status?, headers?, ...)` |
| `vibium:network.continue` | Continue an intercepted request | `continue(overrides?)` | `continue_(overrides?)` |
| `vibium:network.abort` | Abort an intercepted request | `abort()` | `abort()` |

### Dialog

| Wire Command | Description | JS | Python |
|---|---|---|---|
| *from event data* | Get dialog message | `message()` | `message()` |
| *from event data* | Get dialog type | `type()` | `type()` |
| *from event data* | Get dialog default value | `defaultValue()` | `default_value()` |
| `browsingContext.handleUserPrompt` | Accept the dialog | `accept(promptText?)` | `accept(prompt_text?)` |
| `browsingContext.handleUserPrompt` | Dismiss the dialog | `dismiss()` | `dismiss()` |

### Download

| Wire Command | Description | JS | Python |
|---|---|---|---|
| `vibium:download.saveAs` | Save a download to path | `saveAs(path)` | `save_as(path)` |
| *from event data* | Get download URL | `url()` | `url()` |
| *from event data* | Get download filename | `filename()` | `filename()` |
| *from event data* | Get download path | `path()` | `path()` |

### Request / Response / ConsoleMessage / WebSocketInfo

These are lightweight data classes constructed from events. See the JS or Python source for their exact fields.

---

## Naming Conventions

### Method Names

| Convention | JS | Python | Java/Kotlin | C# | Ruby | Rust | Go |
|---|---|---|---|---|---|---|---|
| Multi-word methods | `camelCase` | `snake_case` | `camelCase` | `PascalCase` | `snake_case` | `snake_case` | `PascalCase` |
| Boolean queries | `isVisible()` | `is_visible()` | `isVisible()` | `IsVisible()` | `visible?` | `is_visible()` | `IsVisible()` |
| Setters | `setViewport()` | `set_viewport()` | `setViewport()` | `SetViewport()` | `set_viewport` / `viewport=` | `set_viewport()` | `SetViewport()` |
| Event handlers | `onDialog(fn)` | `on_dialog(fn)` | `onDialog(fn)` | `OnDialog(fn)` | `on_dialog(&block)` | `on_dialog(fn)` | `OnDialog(fn)` |

### Wire → Client Mapping

The wire protocol uses `camelCase`. Each language converts to its idiomatic style:

```
vibium:page.setViewport  →  JS: setViewport()   Python: set_viewport()   Ruby: set_viewport
vibium:element.isVisible →  JS: isVisible()     Python: is_visible()     Ruby: visible?
vibium:page.a11yTree     →  JS: a11yTree()      Python: a11y_tree()      Ruby: a11y_tree
```

### Parameter Names

Wire parameters are `camelCase`. Convert to language idioms:

```
Wire: {"colorScheme": "dark", "reducedMotion": "reduce"}
JS:   colorScheme: "dark", reducedMotion: "reduce"     (same as wire)
Py:   color_scheme="dark", reduced_motion="reduce"     (snake_case)
Ruby: color_scheme: "dark", reduced_motion: "reduce"   (snake_case)
```

**Important:** Always convert at the client boundary. Never leak wire-protocol casing to users (see [#91](https://github.com/VibiumDev/vibium/issues/91)).

---

## Error Types

Every client must define these error types:

| Error | When Thrown |
|---|---|
| `ConnectionError` | WebSocket connection to vibium binary failed |
| `TimeoutError` | Element wait or `waitForFunction` timed out |
| `ElementNotFoundError` | Selector matched no elements |
| `BrowserCrashedError` | Browser process died unexpectedly |

### Wire Error Detection

The wire protocol returns errors in this format:

```json
{"id": 1, "type": "error", "error": "timeout", "message": "Timeout after 30000ms waiting for '#btn'"}
```

Map the `error` field to structured error types:
- `"timeout"` → `TimeoutError`
- Messages containing `"not found"` or `"no elements"` → `ElementNotFoundError`
- WebSocket close with no response → `BrowserCrashedError`
- WebSocket connection failure → `ConnectionError`

### Language-Specific Names

Some languages have built-in `TimeoutError` or `ConnectionError`. Use prefixed names to avoid conflicts:

| Language | Timeout | Connection |
|---|---|---|
| JS/TS | `TimeoutError` | `ConnectionError` |
| Python | `VibiumTimeoutError` | `VibiumConnectionError` |
| Java | `VibiumTimeoutException` | `VibiumConnectionException` |
| C# | `VibiumTimeoutException` | `VibiumConnectionException` |
| Ruby | `TimeoutError` | `ConnectionError` (namespaced under `Vibium::`) |
| Rust | `Error::Timeout` | `Error::Connection` (enum variants) |
| Go | `ErrTimeout` | `ErrConnection` (sentinel errors) |

---

## Async / Sync Patterns

### Every client must have an async API

The wire protocol is inherently async (WebSocket messages). The primary API should be async.

### Sync wrappers are optional but recommended

For scripting and REPL use, a sync wrapper dramatically improves the getting-started experience.

| Language | Async Pattern | Sync Pattern |
|---|---|---|
| JS/TS | `async/await` (native) | Separate `*Sync` classes |
| Python | `async/await` | Separate `sync_api/` module (blocks on event loop) |
| Java | `CompletableFuture<T>` | Blocking `.get()` wrappers |
| Kotlin | `suspend fun` (coroutines) | `runBlocking { }` wrappers |
| C# | `Task<T>` / `async` | `.GetAwaiter().GetResult()` wrappers |
| Ruby | Not needed (GIL) | Primary API is sync; use threads for events |
| Rust | `async fn` (tokio/async-std) | `block_on()` wrappers |
| Go | Goroutines (inherently concurrent) | Primary API is sync with channels for events |
| Swift | `async/await` (structured concurrency) | Sync wrappers with `DispatchSemaphore` |

### Event Handling

Events (`onDialog`, `onRequest`, etc.) are received as WebSocket messages with no `id`. The client must:

1. Parse incoming messages
2. If `type` is `"success"` or `"error"` → match to pending request by `id`
3. If `method` is present (event) → dispatch to registered listeners

---

## Reserved Keyword Handling

Some method names conflict with language reserved words. Here's how to handle them:

| Wire Method | Conflict | Resolution |
|---|---|---|
| `vibium:network.continue` | `continue` is reserved in most languages | Python: `continue_()`, Java: `doContinue()`, Ruby: `continue_request`, C#: `Continue()` (C# allows PascalCase), Rust: `r#continue()` or `continue_()`, Go: `Continue()` |

### General Rules

1. **Append underscore** (Python, Ruby): `continue_()`, `import_()`
2. **Prefix with `do`** (Java, Kotlin): `doContinue()`
3. **Raw identifier** (Rust): `r#continue()`
4. **PascalCase avoids most conflicts** (C#, Go)

---

## Aliases

The JS client provides some aliases for Playwright compatibility and discoverability. New clients should include these:

| Primary | Alias | Reason |
|---|---|---|
| `attr(name)` | `getAttribute(name)` | Playwright compat |
| `bounds()` | `boundingBox()` | Playwright compat |
| `go(url)` | — | Short and memorable; `navigate` is the wire name |
| `waitUntil(state)` | — | Maps to `vibium:element.waitFor` on wire |

### Which to Include

- **Always include the primary name** (shorter, Vibium-native)
- **Include Playwright aliases** for `getAttribute` and `boundingBox` — many users come from Playwright
- **Do not** alias everything — keep the API surface small

---

## Binary Discovery

Each client needs to find and launch the `vibium` binary. The resolution order:

1. **Environment variable** `VIBIUM_BIN_PATH` — highest priority
2. **PATH lookup** — `which vibium` / `where vibium`
3. **npm-installed binary** — check `node_modules/.bin/vibium`
4. **Known install locations** — platform-specific defaults

### Reference

- JS: `clients/javascript/src/clicker/binary.ts` → `getVibiumBinPath()`
- Python: `clients/python/src/vibium/binary.py` → `find_vibium_bin()`

---

## Testing Checklist

Before releasing a new client, verify:

- [ ] `browser.start()` launches a visible browser
- [ ] `browser.start(headless=True)` launches headless
- [ ] `page.go(url)` navigates and waits for load
- [ ] `page.find("selector")` returns an Element
- [ ] `element.click()` performs a click
- [ ] `element.fill("text")` fills an input
- [ ] `page.screenshot()` returns image bytes
- [ ] `page.evaluate("1 + 1")` returns `2`
- [ ] `context.cookies()` / `setCookies()` round-trips
- [ ] `page.route()` intercepts and can fulfill requests
- [ ] `page.onDialog()` handles alert/confirm/prompt
- [ ] Error types are raised (timeout, element not found)
- [ ] `browser.stop()` cleanly shuts down
- [ ] Binary discovery works via `VIBIUM_BIN_PATH` and PATH
- [ ] Sync wrapper works (if provided)

Run the existing test suite against your client:

```bash
make test  # runs CLI + JS + MCP + Python tests
```
