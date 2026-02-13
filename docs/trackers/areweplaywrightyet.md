# Are We Playwright Yet?

Vibium's Playwright-equivalent API coverage. 155 commands across 23 categories, tracked across 6 implementation targets.

**Legend:** ✅ Done · 🟡 Partial · ⬜ Not started · — N/A

---

## Navigation (9 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.go(url)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.back()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.forward()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.reload()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.url()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.title()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.waitForURL(pattern)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.waitForLoad(state?)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.content()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Pages & Contexts (10 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `browser.newPage()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `browser.newContext()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.newPage()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `browser.pages()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.close()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `browser.close()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `browser.onPage(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `browser.onPopup(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.bringToFront()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.close()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Element Finding (6 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.find('css')` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.find({role, text, …})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.findAll('css')` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.findAll({…})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.find('css')` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.find({…})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Selector Strategies (10 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `find({role: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({text: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({label: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({placeholder: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({alt: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({title: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({testid: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({xpath: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({near: '…'})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `find({role, text}) combo` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Locator Chaining & Filtering (8 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `el.first()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.last()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.nth(index)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.count()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.filter({hasText})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.filter({has})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.or(other)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.and(other)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Element Interaction (16 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `el.click()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.dblclick()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.fill(value)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.type(text)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.press(key)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.clear()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.check()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.uncheck()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.selectOption(val)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.setFiles(paths)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.hover()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.focus()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.dragTo(target)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.tap()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.scrollIntoView()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.dispatchEvent(type)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Element State (14 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `el.text()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.innerText()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.html()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.value()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.attr(name)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.bounds()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.isVisible()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.isHidden()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.isEnabled()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.isChecked()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.isEditable()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.eval(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.screenshot()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.waitFor({state})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Keyboard & Mouse (10 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.keyboard.press(key)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.keyboard.down(key)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.keyboard.up(key)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.keyboard.type(text)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.mouse.click(x,y)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.mouse.move(x,y)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.mouse.down()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.mouse.up()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.mouse.wheel(dx,dy)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.touch.tap(x,y)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Network Interception (11 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.route(pattern, handler)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `route.fulfill(response)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `route.continue(overrides?)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `route.abort(reason?)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.onRequest(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.onResponse(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.setHeaders(headers)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.waitForRequest(pat)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.waitForResponse(pat)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.routeWebSocket(pat)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.onWebSocket(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## Request & Response (8 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `request.url()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `request.method()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `request.headers()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `request.postData()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `response.status()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `response.headers()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `response.body()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `response.json()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## Dialogs (5 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onDialog(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `dialog.accept(text?)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `dialog.dismiss()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `dialog.message()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `dialog.type()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Screenshots & PDF (4 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.screenshot()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.screenshot({fullPage})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.screenshot({clip})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.pdf()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Cookies & Storage (5 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `context.cookies(urls?)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.setCookies()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.clearCookies()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.storageState()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `context.addInitScript()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## Emulation (6 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.setViewport(size)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.viewport()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.emulateMedia(opts)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.setContent(html)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.setGeolocation()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.grantPermissions()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Frames (4 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.frame(nameOrUrl)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.frames()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.mainFrame()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| Frames have full Page API | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Accessibility (4 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.a11yTree()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.a11yAudit()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.role()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `el.label()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Console, Errors & Workers (3 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onConsole(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.onError(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.workers()` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Waiting (5 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.waitFor(selector)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.wait(ms)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.waitForFunction(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.waitForEvent(name)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.pause()` | ⬜ | ⬜ | ⬜ | ⬜ | — | ⬜ |

## Downloads & Files (3 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.onDownload(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `download.saveAs(path)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.onFileChooser(fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## Clock (3 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.clock.install()` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.clock.fastForward(ms)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.clock.setFixedTime(t)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## Tracing (2 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.tracing.start(opts)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.tracing.stop(opts)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |

## Evaluation (5 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.eval(expr)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.evalHandle(expr)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |
| `page.addScript(src)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.addStyle(src)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.expose(name, fn)` | ⬜ | ⬜ | ⬜ | ⬜ | — | — |

## AI-Native Methods (4 commands)

| Command | JS async | JS sync | PY async | PY sync | MCP | CLI |
|---------|----------|---------|----------|---------|-----|-----|
| `page.check(claim)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.check(claim, {near})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.do(action)` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
| `page.do(action, {data})` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ |
