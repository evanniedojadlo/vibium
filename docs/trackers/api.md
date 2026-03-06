# Vibium API

161 commands across 22 categories, tracked across 6 implementation targets.

**Legend:** ✅ Done · 🟡 Partial · ⬜ Not started · — N/A

---

## Navigation

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.go(url)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.back()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.forward()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.reload()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.url()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.title()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.content()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Pages & Contexts

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `browser.page()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `browser.newPage()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `browser.newContext()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `context.newPage()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `browser.pages()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.close()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `browser.start(url)` | ✅ | ✅ | ✅ | ✅ | — | ✅ |
| `browser.stop()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `browser.onPage(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `browser.onPopup(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `browser.removeAllListeners(event?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.bringToFront()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.close()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Element Finding

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.find('css')` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.find({role, text, …})` | ✅ | ✅ | ✅ | ✅ | 🟡 | ✅ |
| `page.findAll('css')` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.findAll({…})` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `el.find('css')` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `el.find({…})` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |

## Selector Strategies

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `find({role: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({text: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({label: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({placeholder: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({alt: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({title: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({testid: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({xpath: '…'})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `find({role, text}) combo` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |

## Element Interaction

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `el.click()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.dblclick()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.fill(value)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.type(text)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.press(key)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.clear()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `el.check()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.uncheck()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.selectOption(val)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.setFiles(paths)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.hover()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.focus()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.highlight()` | ⬜ | ⬜ | ⬜ | ⬜ | ✅ | ✅ |
| `el.dragTo(target)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.tap()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `el.scrollIntoView()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.dispatchEvent(type)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |

## Element State

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `el.text()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.innerText()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.html()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.value()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.attr(name)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.bounds()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.isVisible()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.isHidden()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `el.isEnabled()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.isChecked()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.isEditable()` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `el.screenshot()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |

## Keyboard & Mouse

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.keyboard.press(key)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.keyboard.down(key)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.keyboard.up(key)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.keyboard.type(text)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `page.mouse.click(x,y)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.mouse.move(x,y)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.mouse.down()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.mouse.up()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.mouse.wheel(dx,dy)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.scroll(dir,amt)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.touch.tap(x,y)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |

## Network Interception

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.route(pattern, handler)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `route.fulfill(response)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `route.continue(overrides?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `route.abort(reason?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.onRequest(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.onResponse(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.setHeaders(headers)` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `page.unroute(pattern)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.removeAllListeners(event?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.onWebSocket(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |

## Request & Response

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `request.url()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `request.method()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `request.headers()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `request.postData()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `response.status()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `response.headers()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `response.body()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `response.json()` | ✅ | ✅ | ✅ | ✅ | — | — |

## Dialogs

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onDialog(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `dialog.accept(text?)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `dialog.dismiss()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `dialog.message()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `dialog.type()` | ✅ | ✅ | ✅ | ✅ | — | — |

## Screenshots & PDF

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.screenshot()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.screenshot({fullPage})` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.screenshot({clip})` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `page.pdf()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Cookies & Storage

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `context.cookies(urls?)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.setCookies()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.clearCookies()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.storageState()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.addInitScript()` | ✅ | ✅ | ✅ | ✅ | — | — |

## Emulation

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.setViewport(size)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.viewport()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.emulateMedia(opts)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.setContent(html)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.setGeolocation()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.window()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.setWindow(opts)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Frames

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.frame(nameOrUrl)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.frames()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.mainFrame()` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| Frames have full Page API | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |

## Accessibility

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.a11yTree()` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `el.role()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `el.label()` | ✅ | ✅ | ✅ | ✅ | — | — |

## Console & Errors

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onConsole(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.onError(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.consoleMessages()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.errors()` | ✅ | ✅ | ✅ | ✅ | — | — |

## Waiting

### Capture — set up before the action

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.capture.response(pat, fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.capture.request(pat, fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.capture.navigation(fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.capture.event(name, fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.capture.download(fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `page.capture.dialog(fn?)` | ✅ | ✅ | ✅ | ✅ | — | — |

### Wait Until — poll after the cause

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.waitUntil.url(pat)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.waitUntil.loaded(state?)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.waitUntil(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `el.waitUntil(state)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.wait(ms)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Downloads & Files

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onDownload(fn)` | ✅ | ✅ | ✅ | ✅ | — | — |
| `download.path()` | ✅ | ✅ | ✅ | ✅ | — | — |
| `download.saveAs(path)` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `el.setFiles(paths)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

## Clock

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.clock.install(opts?)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.fastForward(ms)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.runFor(ms)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.pauseAt(time)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.resume()` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.setFixedTime(time)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.setSystemTime(time)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| `page.clock.setTimezone(tz)` | ✅ | ✅ | ✅ | ✅ | ✅ | — |

## Tracing

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `context.tracing.start(opts)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.tracing.stop(opts)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `context.tracing.startChunk(opts)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `context.tracing.stopChunk(opts)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `context.tracing.startGroup(name)` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |
| `context.tracing.stopGroup()` | ✅ | ✅ | ✅ | ✅ | — | ⬜ |

## Evaluation

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.evaluate(expr)` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `page.addScript(src)` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `page.addStyle(src)` | ✅ | ✅ | ✅ | ✅ | ⬜ | ⬜ |
| `page.expose(name, fn)` | ✅ | ✅ | ✅ | ✅ | — | — |

## AI-Native Methods

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.check(claim)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.do(action)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.do(action, {data})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
