# Recording Browser Sessions

Record a timeline of screenshots, network requests, DOM snapshots, and action groups — then view it in Record Player.

---

## What You'll Learn

How to capture a recording of a browser session and view it as an interactive timeline.

---

## Quick Start

The fastest way to record a session:

```bash
vibium record start --screenshots

vibium go https://example.com
vibium click 'a'

vibium record stop -o record.zip
```

Open `record.zip` in [Record Player](https://player.vibium.dev) to see a timeline of screenshots and actions.

---

## Basic Recording

Recording is managed with `vibium record start` and `vibium record stop`. The daemon is automatically started when needed.

```bash
vibium record start --name my-session

vibium go https://example.com
vibium click 'a'

vibium record stop -o record.zip
```

`record stop -o` writes the recording zip directly to the given path.

Enable `--screenshots` and `--snapshots` for a more complete recording:

```bash
vibium record start --screenshots --snapshots
```

- **screenshots** — captures the page periodically (~100ms), creating a visual filmstrip. Identical frames are deduplicated.
- **snapshots** — captures the full HTML when the recording stops, so you can inspect the DOM in the viewer.

To reduce recording size, use JPEG format with a lower quality setting:

```bash
vibium record start --screenshots --format jpeg --quality 0.3
```

The default format is JPEG at 0.5 quality. Lowering `--quality` produces smaller files — useful for long-running recordings or CI where file size matters.

---

## Actions

Every vibium command (`click`, `fill`, `go`, etc.) is automatically recorded in the recording as an action marker. You don't need to wrap commands in groups to see them — they show up individually in the timeline.

```bash
vibium record start --screenshots

vibium go https://example.com
vibium click '#btn'
vibium fill '#input' 'hello'

vibium record stop -o record.zip
```

To also record the raw BiDi protocol commands sent to the browser (e.g. `input.performActions`, `script.callFunction`), enable `--bidi`:

```bash
vibium record start --screenshots --bidi
```

This is useful for debugging low-level protocol issues but makes recordings larger.

---

## Action Groups

Use `record group start` and `record group stop` to label sections of your recording. Groups show up as named spans in the timeline.

```bash
vibium record start --screenshots
vibium go https://example.com

vibium record group start 'fill login form'
vibium fill '#username' 'alice'
vibium fill '#password' 'secret'
vibium record group stop

vibium record group start 'submit'
vibium click 'button[type="submit"]'
vibium record group stop

vibium record stop -o record.zip
```

Groups can be nested:

```bash
vibium record group start 'checkout flow'

vibium record group start 'shipping'
# ... fill shipping form
vibium record group stop

vibium record group start 'payment'
# ... fill payment form
vibium record group stop

vibium record group stop
```

---

## Chunks

Chunks split a long recording into segments without stopping the recording. Each chunk produces its own zip.

```bash
vibium record start --screenshots

# First chunk: login
vibium go https://example.com/login
vibium fill '#username' 'alice'
vibium record chunk stop -o login.zip

# Second chunk: dashboard
vibium record chunk start --name dashboard
vibium go https://example.com/dashboard
vibium record chunk stop -o dashboard.zip

# Final stop
vibium record stop
```

---

## Viewing Recordings

Open a recording in [Record Player](https://player.vibium.dev):

1. Go to [player.vibium.dev](https://player.vibium.dev)
2. Drop your `record.zip` file onto the page

The viewer shows:
- **Timeline** — scrub through screenshots frame by frame
- **Actions** — see group markers from `record group start`/`record group stop`
- **Network** — waterfall of all HTTP requests
- **Snapshots** — inspect the DOM at capture time

---

## Reference

### CLI Flags

| Command | Flag | Description |
|---------|------|-------------|
| `record start` | `--screenshots` | Capture screenshots periodically |
| `record start` | `--snapshots` | Capture HTML snapshots |
| `record start` | `--bidi` | Record raw BiDi commands in the recording |
| `record start` | `--name NAME` | Name for the recording |
| `record start` | `--title TITLE` | Title shown in Record Player |
| `record start` | `--format FORMAT` | Screenshot format: `jpeg` or `png` (default: `jpeg`) |
| `record start` | `--quality N` | JPEG quality 0.0–1.0 (default: `0.5`) |
| `record stop` | `-o, --output PATH` | Output file path (default: `record.zip`) |
| `record chunk stop` | `-o, --output PATH` | Output file path for the chunk |
| `record chunk start` | `--name NAME` | Name for the chunk |

---

## Next Steps

- [Recording Format](../explanation/recording-format.md) — detailed spec of the zip structure
- [Getting Started](getting-started-js.md) — first steps with Vibium
