# day 8: elements & sync api

the js client now has element interactions and a synchronous api.

---

## what shipped

**element class (async)**
- `element.click()` - click the element
- `element.type(text)` - type into the element
- `element.text()` - get text content
- `element.getAttribute(name)` - get attribute value
- `element.boundingBox()` - get {x, y, width, height}

**vibe.find(selector)**
- `await vibe.find('a')` returns an Element instance

**sync api**
- `browserSync.launch()` - blocking browser launch
- `vibe.go(url)`, `vibe.find(selector)`, `vibe.screenshot()`, `vibe.quit()`
- uses Worker thread + Atomics.wait under the hood
- proper cleanup on Ctrl+C (no zombie processes)

---

## try it

```bash
make
cd clients/javascript && node --experimental-repl-await
```

**async api:**
```javascript
const { browser } = await import('./dist/index.mjs')
const vibe = await browser.launch({ headless: false })
await vibe.go('https://example.com')
const link = await vibe.find('a')
console.log(await link.text())  // "Learn more"
await link.click()
await vibe.quit()
```

**sync api:**
```javascript
const { browserSync } = require('./dist')
const vibe = browserSync.launch({ headless: false })
vibe.go('https://example.com')
const link = vibe.find('a')
console.log(link.text())  // "Learn more"
link.click()
vibe.quit()
```

---

day 8 of 14. six days left.

*december 19, 2025*
