# Recording Browser Sessions

Record a timeline of screenshots, network requests, DOM snapshots, and action groups — then view it in Record Player.

---

## What You'll Learn

How to capture a recording of a browser session and view it as an interactive timeline.

---

## Quick Start

The fastest way to record a session — use `page.context` to access recording without creating an explicit context:

```python
import asyncio
from vibium.async_api import browser

async def main():
    bro = await browser.start()
    vibe = await bro.page()

    await vibe.context.recording.start(screenshots=True)

    await vibe.go('https://example.com')
    await vibe.find('a').click()

    await vibe.context.recording.stop(path='record.zip')
    await bro.stop()

asyncio.run(main())
```

Open `record.zip` in [Record Player](https://player.vibium.dev) to see a timeline of screenshots and actions.

---

## Basic Recording

Recording lives on `BrowserContext`, not `Page`. The Quick Start above uses `page.context` as a shortcut — under the hood, every page belongs to a context, and `page.context` gives you direct access to it. This is equivalent to creating an explicit context:

```python
import asyncio
from vibium.async_api import browser

async def main():
    bro = await browser.start()
    ctx = await bro.new_context()
    vibe = await ctx.new_page()

    await ctx.recording.start(name='my-session')

    await vibe.go('https://example.com')
    await vibe.find('a').click()

    data = await ctx.recording.stop()
    with open('record.zip', 'wb') as f:
        f.write(data)

    await bro.stop()

asyncio.run(main())
```

Use an explicit context when you need multiple pages in the same recording, or when you want to configure context options (viewport, locale, etc.). Use `page.context` when you just want to record a single page quickly.

`stop()` returns `bytes` containing the recording zip. You can also pass a `path` to write the file directly:

```python
await ctx.recording.stop(path='record.zip')
```

Enable `screenshots` and `snapshots` for a more complete recording:

```python
await ctx.recording.start(screenshots=True, snapshots=True)
```

- **screenshots** — captures the page periodically (~100ms), creating a visual filmstrip. Identical frames are deduplicated.
- **snapshots** — captures the full HTML when the recording stops, so you can inspect the DOM in the viewer.

To reduce recording size, use JPEG format with a lower quality setting:

```python
await ctx.recording.start(
    screenshots=True,
    format='jpeg',
    quality=0.3,
)
```

The default format is JPEG at 0.5 quality. Lowering `quality` produces smaller files — useful for long-running recordings or CI where file size matters.

---

## Actions

Every vibium command (`click`, `fill`, `navigate`, etc.) is automatically recorded in the recording as an action marker. You don't need to wrap commands in groups to see them — they show up individually in the timeline.

```python
await ctx.recording.start(screenshots=True)

await vibe.go('https://example.com')       # recorded as Page.navigate
await vibe.find('#btn').click()             # recorded as Element.click
await vibe.find('#input').fill('hello')     # recorded as Element.fill

await ctx.recording.stop(path='record.zip')
```

To also record the raw BiDi protocol commands sent to the browser (e.g. `input.performActions`, `script.callFunction`), enable `bidi`:

```python
await ctx.recording.start(screenshots=True, bidi=True)
```

This is useful for debugging low-level protocol issues but makes recordings larger.

---

## Action Groups

Use `start_group()` and `stop_group()` to label sections of your recording. Groups show up as named spans in the timeline.

```python
await ctx.recording.start(screenshots=True)
await vibe.go('https://example.com')

await ctx.recording.start_group('fill login form')
await vibe.find('#username').fill('alice')
await vibe.find('#password').fill('secret')
await ctx.recording.stop_group()

await ctx.recording.start_group('submit')
await vibe.find('button[type="submit"]').click()
await ctx.recording.stop_group()

await ctx.recording.stop(path='record.zip')
```

Groups can be nested:

```python
await ctx.recording.start_group('checkout flow')

await ctx.recording.start_group('shipping')
# ... fill shipping form
await ctx.recording.stop_group()

await ctx.recording.start_group('payment')
# ... fill payment form
await ctx.recording.stop_group()

await ctx.recording.stop_group()
```

---

## Chunks

Chunks split a long recording into segments without stopping the recording. Each chunk produces its own zip.

```python
await ctx.recording.start(screenshots=True)

# First chunk: login
await vibe.go('https://example.com/login')
await vibe.find('#username').fill('alice')
await ctx.recording.stop_chunk(path='login.zip')

# Second chunk: dashboard
await ctx.recording.start_chunk(name='dashboard')
await vibe.go('https://example.com/dashboard')
await ctx.recording.stop_chunk(path='dashboard.zip')

# Final stop
await ctx.recording.stop()
```

---

## Viewing Recordings

Open a recording in [Record Player](https://player.vibium.dev):

1. Go to [player.vibium.dev](https://player.vibium.dev)
2. Drop your `record.zip` file onto the page

The viewer shows:
- **Timeline** — scrub through screenshots frame by frame
- **Actions** — see group markers from `start_group()`/`stop_group()`
- **Network** — waterfall of all HTTP requests
- **Snapshots** — inspect the DOM at capture time

---

## Reference

### start() Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `name` | str | `"record"` | Name for the recording |
| `title` | str | — | Title shown in Record Player |
| `screenshots` | bool | `False` | Capture screenshots (~100ms interval) |
| `snapshots` | bool | `False` | Capture DOM snapshots on stop |
| `sources` | bool | `False` | Reserved for future use |
| `bidi` | bool | `False` | Record raw BiDi commands in the recording |
| `format` | `'jpeg'` \| `'png'` | `'jpeg'` | Screenshot image format |
| `quality` | float | `0.5` | JPEG quality 0.0–1.0 (ignored for PNG) |

### stop() / stop_chunk() Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `path` | str | — | File path to write the zip to |

When `path` is omitted, the zip data is returned as `bytes`.

---

## Next Steps

- [Recording Format](../explanation/recording-format.md) — detailed spec of the zip structure
- [Getting Started](getting-started-python.md) — first steps with Vibium
