# Trace File Format

## What Tracing Does

Tracing captures a timeline of everything that happens during a browser session — screenshots, network requests, DOM snapshots, and action groups — and packages it into a single zip file.

## Zip Structure

A trace zip contains three kinds of entries:

```
trace.zip
├── 0-trace.trace           # Main event timeline (newline-delimited JSON)
├── 0-trace.network         # Network events (newline-delimited JSON)
└── resources/
    ├── a1b2c3d4e5...png    # Screenshot frames (named by SHA1)
    └── f6a7b8c9d0...html   # DOM snapshots (named by SHA1)
```

The number prefix (`0-`) is the chunk index. Each `stopChunk()` call produces a zip with the current chunk's events. The first chunk is `0`, the next is `1`, etc.

## Event Files

### `<n>-trace.trace`

Newline-delimited JSON. Each line is a self-contained event object with a `type` field. Events appear in chronological order.

**`context-options`** — Always the first event. Metadata about the recording session.

```json
{"type":"context-options","browserName":"chromium","platform":"darwin","wallTime":1708000000000,"title":"my test","options":{},"sdkLanguage":"javascript"}
```

| Field | Type | Description |
|-------|------|-------------|
| `browserName` | string | Always `"chromium"` |
| `platform` | string | OS: `"darwin"`, `"linux"`, or `"windows"` |
| `wallTime` | number | Unix timestamp in milliseconds |
| `title` | string | From `start({ title })` or `start({ name })` |
| `options` | object | Browser context options (currently `{}`) |
| `sdkLanguage` | string | Always `"javascript"` |

**`screencast-frame`** — A screenshot captured during recording. References a PNG file in `resources/` by SHA1 hash.

```json
{"type":"screencast-frame","pageId":"ABCDEF123","sha1":"a1b2c3d4e5f6...","width":1280,"height":720,"timestamp":1708000000100}
```

| Field | Type | Description |
|-------|------|-------------|
| `pageId` | string | Browsing context ID of the captured page |
| `sha1` | string | Lowercase hex SHA1 of the PNG data |
| `width` | number | Screenshot width in pixels (read from PNG header) |
| `height` | number | Screenshot height in pixels |
| `timestamp` | number | Unix ms when the screenshot was taken |

When `screenshots: true` is set, a background goroutine captures a screenshot every ~100ms. Identical frames are deduplicated by SHA1 — if the page doesn't change, only one PNG is stored in `resources/`.

**`frame-snapshot`** — A DOM snapshot (full `document.documentElement.outerHTML`). References an HTML file in `resources/`.

```json
{"type":"frame-snapshot","pageId":"ABCDEF123","sha1":"f6a7b8c9d0...","frameUrl":"https://example.com","title":"","timestamp":1708000000200,"wallTime":1708000000200,"callId":"snapshot@5"}
```

Captured once when `stop()` or `stopChunk()` is called (if `snapshots: true`).

**`before`** / **`after`** — Paired markers that bracket an operation. There are three kinds: actions, action groups, and BiDi commands.

#### Actions (auto-recorded)

Every vibium command emits a `before`/`after` pair automatically — both mutations (`click`, `fill`, `navigate`) and read-only queries (`text`, `isVisible`, `getAttribute`). This matches Playwright's behavior of recording all SDK calls.

```json
{"type":"before","callId":"action@1","apiName":"Page.navigate","class":"Page","method":"vibium:page.navigate","params":{"url":"https://example.com"},"wallTime":1708000000300,"startTime":1708000000300}
{"type":"after","callId":"action-end@1","endTime":1708000000400}
{"type":"before","callId":"action@2","apiName":"Element.click","class":"Element","method":"vibium:click","params":{"selector":"#login"},"wallTime":1708000000500,"startTime":1708000000500}
{"type":"after","callId":"action-end@2","endTime":1708000000600}
{"type":"before","callId":"action@3","apiName":"Element.text","class":"Element","method":"vibium:el.text","params":{"selector":".result"},"wallTime":1708000000700,"startTime":1708000000700}
{"type":"after","callId":"action-end@3","endTime":1708000000750}
```

| Field | Type | Description |
|-------|------|-------------|
| `callId` | string | `action@<N>` for the start, `action-end@<N>` for the end. `N` is a monotonic counter shared across actions and BiDi commands. |
| `apiName` | string | Human-readable name like `Element.click`, `Page.navigate`, `Element.text`. Derived from the vibium method by `apiNameFromMethod()`. |
| `class` | string | `Element`, `Page`, `Browser`, `BrowserContext`, `Network`, `Dialog`, etc. |
| `method` | string | The raw vibium method (e.g., `vibium:click`, `vibium:el.text`). |
| `params` | object | The command parameters as sent by the client. |

The `dispatch()` wrapper in the router records these markers — every vibium command that goes through `dispatch()` gets traced. Tracing commands themselves (`tracing.start`, `tracing.stop`, etc.) are excluded since they control tracing.

#### Action groups (user-defined)

Groups are named spans from `startGroup()` / `stopGroup()`. They wrap multiple actions under a single label in the timeline.

```json
{"type":"before","callId":"group@3","apiName":"login flow","class":"Tracing","method":"group","params":{"name":"login flow"},"wallTime":1708000000300,"startTime":1708000000300}
{"type":"after","callId":"group-end@7","endTime":1708000000500}
```

Groups are nestable. The `callId` is `group@<index>` / `group-end@<index>` where `<index>` is the event's position in the trace array.

#### BiDi commands (opt-in)

When tracing is started with `bidi: true`, every raw BiDi command sent to the browser via `sendInternalCommand` is also recorded. This is useful for debugging low-level protocol issues.

```json
{"type":"before","callId":"bidi@4","apiName":"browsingContext.navigate","class":"BiDi","method":"browsingContext.navigate","params":{"context":"ABC123","url":"https://example.com"},"wallTime":1708000000350,"startTime":1708000000350}
{"type":"after","callId":"bidi-end@4","endTime":1708000000390}
```

| Field | Type | Description |
|-------|------|-------------|
| `callId` | string | `bidi@<N>` / `bidi-end@<N>`. Shares the same monotonic counter as actions. |
| `apiName` | string | The BiDi method name (e.g., `browsingContext.navigate`). |
| `class` | string | Always `"BiDi"`. |

BiDi markers nest inside action markers — a single vibium command (like `Page.navigate`) may produce several BiDi commands internally.

**`event`** — A BiDi browser event (context creation, dialog, log entry, etc.) recorded as-is.

```json
{"type":"event","method":"browsingContext.contextCreated","params":{...},"time":1708000000150,"class":"BrowserContext"}
```

These are all non-network BiDi events that flow through the router while tracing is active.

### `<n>-trace.network`

Newline-delimited JSON. Each line is a network event recorded from BiDi.

```json
{"type":"resource-snapshot","method":"network.responseCompleted","params":{...},"timestamp":1708000000400}
```

| Field | Type | Description |
|-------|------|-------------|
| `method` | string | BiDi method: `network.beforeRequestSent`, `network.responseCompleted`, or `network.fetchError` |
| `params` | object | Full BiDi event params (request URL, headers, status, timing, etc.) |
| `timestamp` | number | Unix ms when the event was recorded |

## Resources Directory

Binary assets referenced by SHA1 hash from the event files.

| Extension | Content | Source |
|-----------|---------|--------|
| `.png` | Screenshot frames | `screencast-frame` events |
| `.html` | DOM snapshots | `frame-snapshot` events |

The file name is the full lowercase hex SHA1 hash of the content (e.g., `a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0.png`). This provides natural deduplication — if two screenshots are pixel-identical, they share one file.

## How Recording Works

### Session lifecycle

```
Client                       Proxy                            Browser
  │                            │                                 │
  │  tracing.start             │                                 │
  ├───────────────────────────>│  Create TraceRecorder           │
  │                            │  Start screenshot goroutine     │
  │                            │                                 │
  │  page.navigate, click, …   │                                 │
  ├───────────────────────────>│  dispatch() (see below)         │
  │                            │                                 │
  │                            │<─── BiDi events ────────────────│
  │                            │  RecordBidiEvent()              │
  │                            │  (network → .network,           │
  │                            │   other → .trace)               │
  │                            │                                 │
  │  tracing.stop              │                                 │
  ├───────────────────────────>│  Stop screenshot goroutine      │
  │                            │  Capture final DOM snapshot     │
  │                            │  Build zip                      │
  │<──── zip (base64 or file) ─│                                 │
```

### Inside `dispatch()` for a single command

```
dispatch(session, cmd, handler)
  │
  ├── RecordAction(method, params)       // before marker
  │
  │   handler(session, cmd)
  │     │
  │     ├── sendInternalCommand ────────────────> Browser
  │     │   ├── [if bidi: true] RecordBidiCommand()
  │     │   │   ···wait for response···
  │     │   └── [if bidi: true] RecordBidiCommandEnd()
  │     │
  │     ├── sendInternalCommand ────────────────> Browser
  │     │   └── (same pattern, one per BiDi call)
  │     │
  │     └── sendSuccess / sendError
  │
  └── RecordActionEnd()                  // after marker
```

The `dispatch()` wrapper records action markers around every vibium command handler. Inside the handler, `sendInternalCommand` optionally records BiDi command markers when `bidi: true` was passed to `tracing.start`. The `routeBrowserToClient` loop independently records BiDi *events* (browser-initiated messages like context creation, network activity, log entries) — these are passive observations that require no extra round-trips. Only the periodic screenshot captures generate additional traffic.

## Chunks vs. Single Trace

By default, `start()` → `stop()` produces one zip covering the entire session.

Chunks split a trace into segments without stopping the recording:

```
start()  ──────────────────────────────────────────────  stop()
           │                    │                    │
      events A             events B             events C
           │                    │                    │
       stopChunk()         startChunk()          (final)
       → zip with A        stopChunk()
                           → zip with B
                                                → zip with C
```

Each chunk gets its own `context-options` header and chunk index. Resources (screenshots, snapshots) are shared across chunks — the `resources/` map is not cleared on `startChunk()`, so a stopChunk zip may contain resources referenced by earlier chunks too.

## Viewing Traces

Open a trace in [Vibium Trace](https://trace.vibium.dev):

1. Go to [trace.vibium.dev](https://trace.vibium.dev)
2. Drop your `trace.zip` file onto the page

The viewer renders a timeline with screenshots, actions, network waterfall, and DOM snapshots. Every vibium command appears as an individual action in the timeline (e.g., `Page.navigate`, `Element.click`, `Element.text`). Action groups from `startGroup()`/`stopGroup()` appear as labeled spans that wrap multiple actions. With `bidi: true`, raw BiDi commands are also visible as nested entries within their parent actions.
