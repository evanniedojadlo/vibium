# Plan: README Improvements

## 1. Remove hardcoded tool counts from prose

Tool counts go stale constantly (81 → 83 → 85...). Replace with qualitative language.

| File | Line | Current | Replacement |
|------|------|---------|-------------|
| `README.md` | 11 | "learns 81 browser automation tools instantly" | "learns the full browser automation toolkit instantly" |

Leave alone (historical/tracking docs — counts are correct for their snapshot):
- `docs/updates/2026-02-28-26.2-release.md` — release notes are frozen snapshots
- `docs/trackers/api.md` — tracking doc, counts are the point
- `tests/mcp/server.test.js` — test assertion, must stay exact

## 2. Add hero GIF

Record a short (~10s) terminal GIF showing:
```
vibium go https://var.parts && vibium map && vibium click @e1 && vibium diff map
```

Place at top of README, right after the tagline and before "Why Vibium?". Use a simple `![demo](docs/assets/demo.gif)`.

Tool: `vhs` (Charm CLI) or `asciinema` + `agg` for terminal recording, or screen capture of the actual browser + terminal side by side.

## 3. Move architecture diagram lower

Move the `## Architecture` section below `## Language APIs`. Most visitors want "what does it do" → "how do I install it" → "how do I use it" before internals. New order:

1. Tagline + hero GIF
2. Why Vibium?
3. Agent Setup (+ CLI Quick Reference + MCP)
4. Language APIs
5. Architecture ← moved here
6. Platform Support
7. Contributing
8. Roadmap
9. License

## 4. Explain `npx skills add`

The `skills` CLI isn't standard — first-time readers will be confused. Add a one-liner:

```markdown
> `skills` is the [Skills CLI](https://github.com/anthropics/skills) for managing agent skills. Install with `npm install -g skills`.
```

Or if `skills` isn't a real standalone tool, clarify what it is and link to docs.

## 5. Add inline quick start

Add a 3-line "try it now" block right after "New here?" so people can copy-paste without reading a tutorial:

```markdown
**Try it now:**
```bash
npm install -g vibium
vibium go https://var.parts && vibium map
```

## 6. Trim install plumbing from Language APIs

Move cache paths and `VIBIUM_SKIP_BROWSER_DOWNLOAD` to a setup/install doc (or a `<details>` collapse). Keep the Language APIs section focused on the API itself.

## 7. Add badges

Add standard OSS badges at the top: npm version, license, CI status.

```markdown
[![npm](https://img.shields.io/npm/v/vibium)](https://www.npmjs.com/package/vibium)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
```

---

## Priority

| # | Change | Effort | Impact |
|---|--------|--------|--------|
| 1 | Remove hardcoded tool counts | 5 min | Prevents recurring staleness |
| 2 | Hero GIF | 30 min | Biggest visual impact |
| 3 | Move architecture lower | 5 min | Better information hierarchy |
| 4 | Explain `skills` CLI | 5 min | Unblocks confused newcomers |
| 5 | Inline quick start | 5 min | Reduces time to "aha" |
| 6 | Trim install plumbing | 10 min | Cleaner README |
| 7 | Badges | 5 min | Polish |
