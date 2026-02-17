# Trace File Format

## What Tracing Does

Tracing captures a timeline of everything that happens during a browser session — screenshots, network requests, DOM snapshots, and action groups — and packages it into a single zip file. The output is compatible with Playwright's trace viewer (`npx playwright show-trace trace.zip`).

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

**`before`** / **`after`** — Action group boundaries from `startGroup()` / `stopGroup()`.

```json
{"type":"before","callId":"group@3","apiName":"login flow","class":"Tracing","method":"group","params":{"name":"login flow"},"wallTime":1708000000300,"startTime":1708000000300}
{"type":"after","callId":"group-end@7","endTime":1708000000500}
```

Groups are nestable. The `callId` is a synthetic identifier (`group@<index>` / `group-end@<index>`) based on the event's position in the trace.

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

```
JS Client                    Go Engine (proxy)                   Browser
    │                              │                                │
    │  vibium:tracing.start        │                                │
    ├─────────────────────────────>│  Create TraceRecorder           │
    │                              │  Start screenshot goroutine     │
    │                              │        │                        │
    │  vibium:page.navigate        │        │  captureScreenshot     │
    ├─────────────────────────────>│        ├───────────────────────>│
    │                              │        │<──── PNG base64 ──────│
    │                              │        │  SHA1 → store in      │
    │                              │        │  resources map         │
    │                              │        │                        │
    │                       routeBrowserToClient                     │
    │                              │<──── BiDi events ──────────────│
    │                              │  RecordBidiEvent()              │
    │                              │  (network → .network,           │
    │                              │   other → .trace)               │
    │                              │                                 │
    │  vibium:tracing.stop         │                                 │
    ├─────────────────────────────>│  Stop screenshot goroutine      │
    │                              │  Capture final DOM snapshot     │
    │                              │  Build zip from:                │
    │                              │    events → <n>-trace.trace     │
    │                              │    network → <n>-trace.network  │
    │                              │    resources → resources/       │
    │<──── zip (base64 or file) ──│                                 │
```

The trace recorder hooks into the existing `routeBrowserToClient` message loop. Every BiDi event from the browser passes through this function — the recorder gets a copy before the event is forwarded to the JS client. This means tracing adds no extra browser round-trips for event collection; only the periodic screenshot captures generate additional traffic.

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

```bash
npx playwright show-trace trace.zip
```

The Playwright trace viewer reads the zip and renders a timeline with screenshots, actions, network waterfall, and DOM snapshots. Vibium traces include screenshots and network data; the action timeline shows group markers from `startGroup()`/`stopGroup()`.
