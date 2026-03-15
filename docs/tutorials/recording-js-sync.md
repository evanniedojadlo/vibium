# Recording Browser Sessions

Record a timeline of screenshots, network requests, DOM snapshots, and action groups — then view it in Record Player.

---

## What You'll Learn

How to capture a recording of a browser session and view it as an interactive timeline.

---

## Quick Start

The fastest way to record a session — use `page.context` to access recording without creating an explicit context:

```javascript
const { browser } = require('vibium/sync')

const bro = browser.start()
const vibe = bro.page()

vibe.context.recording.start({ screenshots: true })

vibe.go('https://example.com')
vibe.find('a').click()

vibe.context.recording.stop({ path: 'record.zip' })
bro.stop()
```

Open `record.zip` in [Record Player](https://player.vibium.dev) to see a timeline of screenshots and actions.

---

## Basic Recording

Recording lives on `BrowserContext`, not `Page`. The Quick Start above uses `page.context` as a shortcut — under the hood, every page belongs to a context, and `page.context` gives you direct access to it. This is equivalent to creating an explicit context:

```javascript
const { browser } = require('vibium/sync')

const bro = browser.start()
const ctx = bro.newContext()
const vibe = ctx.newPage()

ctx.recording.start({ name: 'my-session' })

vibe.go('https://example.com')
vibe.find('a').click()

const zip = ctx.recording.stop()
require('fs').writeFileSync('record.zip', zip)

bro.stop()
```

Use an explicit context when you need multiple pages in the same recording, or when you want to configure context options (viewport, locale, etc.). Use `page.context` when you just want to record a single page quickly.

`stop()` returns a `Buffer` containing the recording zip. You can also pass a `path` to write the file directly:

```javascript
ctx.recording.stop({ path: 'record.zip' })
```

Enable `screenshots` and `snapshots` for a more complete recording:

```javascript
ctx.recording.start({ screenshots: true, snapshots: true })
```

- **screenshots** — captures the page periodically (~100ms), creating a visual filmstrip. Identical frames are deduplicated.
- **snapshots** — captures the full HTML when the recording stops, so you can inspect the DOM in the viewer.

To reduce recording size, use JPEG format with a lower quality setting:

```javascript
ctx.recording.start({
  screenshots: true,
  format: 'jpeg',
  quality: 0.3,
})
```

The default format is JPEG at 0.5 quality. Lowering `quality` produces smaller files — useful for long-running recordings or CI where file size matters.

---

## Actions

Every vibium command (`click`, `fill`, `navigate`, etc.) is automatically recorded in the recording as an action marker. You don't need to wrap commands in groups to see them — they show up individually in the timeline.

```javascript
ctx.recording.start({ screenshots: true })

vibe.go('https://example.com')       // recorded as Page.navigate
vibe.find('#btn').click()             // recorded as Element.click
vibe.find('#input').fill('hello')     // recorded as Element.fill

ctx.recording.stop({ path: 'record.zip' })
```

To also record the raw BiDi protocol commands sent to the browser (e.g. `input.performActions`, `script.callFunction`), enable `bidi`:

```javascript
ctx.recording.start({ screenshots: true, bidi: true })
```

This is useful for debugging low-level protocol issues but makes recordings larger.

---

## Action Groups

Use `startGroup()` and `stopGroup()` to label sections of your recording. Groups show up as named spans in the timeline.

```javascript
ctx.recording.start({ screenshots: true })
vibe.go('https://example.com')

ctx.recording.startGroup('fill login form')
vibe.find('#username').fill('alice')
vibe.find('#password').fill('secret')
ctx.recording.stopGroup()

ctx.recording.startGroup('submit')
vibe.find('button[type="submit"]').click()
ctx.recording.stopGroup()

ctx.recording.stop({ path: 'record.zip' })
```

Groups can be nested:

```javascript
ctx.recording.startGroup('checkout flow')

  ctx.recording.startGroup('shipping')
  // ... fill shipping form
  ctx.recording.stopGroup()

  ctx.recording.startGroup('payment')
  // ... fill payment form
  ctx.recording.stopGroup()

ctx.recording.stopGroup()
```

---

## Chunks

Chunks split a long recording into segments without stopping the recording. Each chunk produces its own zip.

```javascript
ctx.recording.start({ screenshots: true })

// First chunk: login
vibe.go('https://example.com/login')
vibe.find('#username').fill('alice')
ctx.recording.stopChunk({ path: 'login.zip' })

// Second chunk: dashboard
ctx.recording.startChunk({ name: 'dashboard' })
vibe.go('https://example.com/dashboard')
ctx.recording.stopChunk({ path: 'dashboard.zip' })

// Final stop
ctx.recording.stop()
```

---

## Viewing Recordings

Open a recording in [Record Player](https://player.vibium.dev):

1. Go to [player.vibium.dev](https://player.vibium.dev)
2. Drop your `record.zip` file onto the page

The viewer shows:
- **Timeline** — scrub through screenshots frame by frame
- **Actions** — see group markers from `startGroup()`/`stopGroup()`
- **Network** — waterfall of all HTTP requests
- **Snapshots** — inspect the DOM at capture time

---

## Reference

### start() Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `name` | string | `"record"` | Name for the recording |
| `title` | string | — | Title shown in Record Player |
| `screenshots` | boolean | `false` | Capture screenshots (~100ms interval) |
| `snapshots` | boolean | `false` | Capture DOM snapshots on stop |
| `sources` | boolean | `false` | Reserved for future use |
| `bidi` | boolean | `false` | Record raw BiDi commands in the recording |
| `format` | `'jpeg'` \| `'png'` | `'jpeg'` | Screenshot image format |
| `quality` | number | `0.5` | JPEG quality 0.0–1.0 (ignored for PNG) |

### stop() / stopChunk() Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `path` | string | — | File path to write the zip to |

When `path` is omitted, the zip data is returned as a `Buffer`.

---

## Next Steps

- [Recording Format](../explanation/recording-format.md) — detailed spec of the zip structure
- [Getting Started](getting-started-js.md) — first steps with Vibium
