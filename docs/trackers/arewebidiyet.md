# Are We BiDi Yet?

Vibium's coverage of the [WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) spec. Maps every BiDi command and event to its Vibium equivalent.

**Legend:** ✅ Done · 🟡 Partial · ⬜ Not started

---

## Object Model

| BiDi Concept | Vibium | Notes |
|--------------|--------|-------|
| Session | `browser` (Browser) | `browser.start()` creates browser process |
| User Context | `context` (Context) | `browser.newContext()` — isolated cookies/storage |
| Browsing Context | `page` (Page) | `browser.newPage()` or `context.newPage()` |
| Element / Node | `Element` | `page.find()` returns Element |
| Realm | (internal) | Script execution contexts |
| Network Intercept | `route` (Route) | `page.route()` returns Route |

```javascript
import { browser } from 'vibium'

const bro = await browser.start()     // ≈ session.new
const ctx = await bro.newContext()     // ≈ browser.createUserContext
const page = await ctx.newPage()      // ≈ browsingContext.create
await page.go('https://example.com')  // ≈ browsingContext.navigate
await page.find('#btn').click()       // ≈ browsingContext.locateNodes + input.performActions
await bro.stop()                      // ≈ session.end
```

---

## Commands (63)

### session (5 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 1 | `session.status` | `browser.status()` | ⬜ |
| 2 | `session.new` | `browser.start(caps?)` | ✅ |
| 3 | `session.end` | `browser.stop()` | ✅ |
| 4 | `session.subscribe` | (internal) | ⬜ |
| 5 | `session.unsubscribe` | (internal) | ⬜ |

### browser (7 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 6 | `browser.close` | `browser.stop()` | ✅ |
| 7 | `browser.createUserContext` | `browser.newContext()` | ✅ |
| 8 | `browser.getClientWindows` | `browser.windows()` | ⬜ |
| 9 | `browser.getUserContexts` | `browser.contexts()` | ⬜ |
| 10 | `browser.removeUserContext` | `context.close()` | ✅ |
| 11 | `browser.setClientWindowState` | `page.setWindow()` | ✅ |
| 12 | `browser.setDownloadBehavior` | `page.setDownloadBehavior()` | ⬜ |

### browsingContext (13 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 13 | `browsingContext.activate` | `page.bringToFront()` | ✅ |
| 14 | `browsingContext.captureScreenshot` | `page.screenshot()` | ✅ |
| 15 | `browsingContext.close` | `page.close()` | ✅ |
| 16 | `browsingContext.create` | `browser.newPage()` | ✅ |
| 17 | `browsingContext.getTree` | `browser.pages()` / `page.frames()` | ✅ |
| 18 | `browsingContext.handleUserPrompt` | `dialog.accept()` / `dialog.dismiss()` | ✅ |
| 19 | `browsingContext.locateNodes` | `page.find()` / `page.findAll()` | ✅ |
| 20 | `browsingContext.navigate` | `page.go(url)` | ✅ |
| 21 | `browsingContext.print` | `page.pdf()` | ✅ |
| 22 | `browsingContext.reload` | `page.reload()` | ✅ |
| 23 | `browsingContext.setBypassCSP` | `page.setBypassCSP()` | ⬜ |
| 24 | `browsingContext.setViewport` | `page.setViewport()` | ✅ |
| 25 | `browsingContext.traverseHistory` | `page.back()` / `page.forward()` | ✅ |

### emulation (11 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 26 | `emulation.setForcedColorsModeThemeOverride` | — | ⬜ |
| 27 | `emulation.setGeolocationOverride` | `page.setGeolocation()` | ✅ |
| 28 | `emulation.setLocaleOverride` | — | ⬜ |
| 29 | `emulation.setNetworkConditions` | — | ⬜ |
| 30 | `emulation.setScreenOrientationOverride` | — | ⬜ |
| 31 | `emulation.setScreenSettingsOverride` | — | ⬜ |
| 32 | `emulation.setScriptingEnabled` | — | ⬜ |
| 33 | `emulation.setScrollbarTypeOverride` | — | ⬜ |
| 34 | `emulation.setTimezoneOverride` | — | ⬜ |
| 35 | `emulation.setTouchOverride` | — | ⬜ |
| 36 | `emulation.setUserAgentOverride` | — | ⬜ |

### network (13 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 37 | `network.addDataCollector` | — | ⬜ |
| 38 | `network.addIntercept` | `page.route()` | ✅ |
| 39 | `network.continueRequest` | `route.continue()` | ✅ |
| 40 | `network.continueResponse` | `route.continue()` | ✅ |
| 41 | `network.continueWithAuth` | `route.authenticate()` | ⬜ |
| 42 | `network.disownData` | — | ⬜ |
| 43 | `network.failRequest` | `route.abort()` | ✅ |
| 44 | `network.getData` | — | ⬜ |
| 45 | `network.provideResponse` | `route.fulfill()` | ✅ |
| 46 | `network.removeDataCollector` | — | ⬜ |
| 47 | `network.removeIntercept` | `page.unroute()` | ✅ |
| 48 | `network.setCacheBehavior` | — | ⬜ |
| 49 | `network.setExtraHeaders` | `page.setHeaders()` | ✅ |

### script (6 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 50 | `script.addPreloadScript` | `context.addInitScript()` | ✅ |
| 51 | `script.callFunction` | `page.evaluate()` | ✅ |
| 52 | `script.disown` | (internal) | ⬜ |
| 53 | `script.evaluate` | `page.evaluate()` | ✅ |
| 54 | `script.getRealms` | — | ⬜ |
| 55 | `script.removePreloadScript` | — | ⬜ |

### input (3 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 56 | `input.performActions` | `page.keyboard.*` / `page.mouse.*` | ✅ |
| 57 | `input.releaseActions` | (automatic) | ⬜ |
| 58 | `input.setFiles` | `el.setFiles()` | ✅ |

### storage (3 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 59 | `storage.deleteCookies` | `context.clearCookies()` | ✅ |
| 60 | `storage.getCookies` | `context.cookies()` | ✅ |
| 61 | `storage.setCookie` | `context.setCookies([c])` | ✅ |

### webExtension (2 commands)

| # | BiDi Command | Vibium | Status |
|---|-------------|--------|--------|
| 62 | `webExtension.install` | — | ⬜ |
| 63 | `webExtension.uninstall` | — | ⬜ |

---

## Events (24)

### browsingContext (14 events)

| # | BiDi Event | Vibium | Status |
|---|-----------|--------|--------|
| 64 | `browsingContext.contextCreated` | `browser.onPage(fn)` | ✅ |
| 65 | `browsingContext.contextDestroyed` | — | ⬜ |
| 66 | `browsingContext.domContentLoaded` | — | ⬜ |
| 67 | `browsingContext.downloadEnd` | — | ⬜ |
| 68 | `browsingContext.downloadWillBegin` | `page.onDownload(fn)` | ✅ |
| 69 | `browsingContext.fragmentNavigated` | — | ⬜ |
| 70 | `browsingContext.historyUpdated` | — | ⬜ |
| 71 | `browsingContext.load` | — | ⬜ |
| 72 | `browsingContext.navigationAborted` | — | ⬜ |
| 73 | `browsingContext.navigationCommitted` | — | ⬜ |
| 74 | `browsingContext.navigationFailed` | — | ⬜ |
| 75 | `browsingContext.navigationStarted` | — | ⬜ |
| 76 | `browsingContext.userPromptClosed` | — | ⬜ |
| 77 | `browsingContext.userPromptOpened` | `page.onDialog(fn)` | ✅ |

### input (1 event)

| # | BiDi Event | Vibium | Status |
|---|-----------|--------|--------|
| 78 | `input.fileDialogOpened` | — | ⬜ |

### log (1 event)

| # | BiDi Event | Vibium | Status |
|---|-----------|--------|--------|
| 79 | `log.entryAdded` | `page.onConsole(fn)` / `page.onError(fn)` | ✅ |

### network (5 events)

| # | BiDi Event | Vibium | Status |
|---|-----------|--------|--------|
| 80 | `network.authRequired` | (via route) | ⬜ |
| 81 | `network.beforeRequestSent` | `page.onRequest(fn)` | ✅ |
| 82 | `network.fetchError` | — | ⬜ |
| 83 | `network.responseCompleted` | `page.onResponse(fn)` | ✅ |
| 84 | `network.responseStarted` | — | ⬜ |

### script (3 events)

| # | BiDi Event | Vibium | Status |
|---|-----------|--------|--------|
| 85 | `script.message` | — | ⬜ |
| 86 | `script.realmCreated` | — | ⬜ |
| 87 | `script.realmDestroyed` | — | ⬜ |

---

**Total: 63 commands + 24 events = 87** · ✅ 40 · 🟡 0 · ⬜ 47
