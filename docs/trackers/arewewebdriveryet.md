# Are We WebDriver Yet?

Vibium's coverage of the classic [WebDriver](https://w3c.github.io/webdriver/) protocol (W3C, pre-BiDi). Maps every HTTP endpoint to its Vibium equivalent.

**Legend:** ✅ Done · 🟡 Partial · ⬜ Not started

---

## Session (3 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 1 | New Session | `POST /session` | `browser.start(caps?)` | ✅ |
| 2 | Delete Session | `DELETE /session/{id}` | `browser.stop()` | ✅ |
| 3 | Status | `GET /status` | `browser.status()` | ⬜ |

## Timeouts (2 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 4 | Get Timeouts | `GET /session/{id}/timeouts` | `browser.timeouts()` | ⬜ |
| 5 | Set Timeouts | `POST /session/{id}/timeouts` | `browser.setTimeouts(t)` | ⬜ |

## Navigation (6 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 6 | Navigate To | `POST /session/{id}/url` | `page.go(url)` | ✅ |
| 7 | Get Current URL | `GET /session/{id}/url` | `page.url()` | ✅ |
| 8 | Back | `POST /session/{id}/back` | `page.back()` | ✅ |
| 9 | Forward | `POST /session/{id}/forward` | `page.forward()` | ✅ |
| 10 | Refresh | `POST /session/{id}/refresh` | `page.reload()` | ✅ |
| 11 | Get Title | `GET /session/{id}/title` | `page.title()` | ✅ |

## Contexts (12 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 12 | Get Window Handle | `GET /session/{id}/window` | `page.id` | ⬜ |
| 13 | Close Window | `DELETE /session/{id}/window` | `page.close()` | ✅ |
| 14 | Switch To Window | `POST /session/{id}/window` | N/A (not needed) | ⬜ |
| 15 | Get Window Handles | `GET /session/{id}/window/handles` | `browser.pages()` | ✅ |
| 16 | New Window | `POST /session/{id}/window/new` | `browser.newPage()` | ✅ |
| 17 | Switch To Frame | `POST /session/{id}/frame` | `page.frame(ref)` | ✅ |
| 18 | Switch To Parent Frame | `POST /session/{id}/frame/parent` | `frame.parent()` | ⬜ |
| 19 | Get Window Rect | `GET /session/{id}/window/rect` | `page.window()` | ✅ |
| 20 | Set Window Rect | `POST /session/{id}/window/rect` | `page.setWindow()` | ✅ |
| 21 | Maximize Window | `POST /session/{id}/window/maximize` | `page.setWindow()` | ✅ |
| 22 | Minimize Window | `POST /session/{id}/window/minimize` | `page.setWindow()` | ✅ |
| 23 | Fullscreen Window | `POST /session/{id}/window/fullscreen` | `page.setWindow()` | ✅ |

## Elements (7 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 24 | Find Element | `POST /session/{id}/element` | `page.find(sel)` | ✅ |
| 25 | Find Elements | `POST /session/{id}/elements` | `page.findAll(sel)` | ✅ |
| 26 | Find Element From Element | `POST /session/{id}/element/{eid}/element` | `el.find(sel)` | ✅ |
| 27 | Find Elements From Element | `POST /session/{id}/element/{eid}/elements` | `el.findAll(sel)` | ⬜ |
| 28 | Find Element From Shadow Root | `POST /session/{id}/shadow/{sid}/element` | — | ⬜ |
| 29 | Find Elements From Shadow Root | `POST /session/{id}/shadow/{sid}/elements` | — | ⬜ |
| 30 | Get Active Element | `GET /session/{id}/element/active` | `page.activeElement()` | ⬜ |

## Element State (11 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 31 | Get Element Shadow Root | `GET /session/{id}/element/{eid}/shadow` | — | ⬜ |
| 32 | Is Element Selected | `GET /session/{id}/element/{eid}/selected` | `el.isChecked()` | ✅ |
| 33 | Get Element Attribute | `GET /session/{id}/element/{eid}/attribute/{name}` | `el.attr(name)` | ✅ |
| 34 | Get Element Property | `GET /session/{id}/element/{eid}/property/{name}` | `page.evaluate(fn, el)` | ✅ |
| 35 | Get Element CSS Value | `GET /session/{id}/element/{eid}/css/{name}` | `page.evaluate(fn, el)` | ✅ |
| 36 | Get Element Text | `GET /session/{id}/element/{eid}/text` | `el.text()` | ✅ |
| 37 | Get Element Tag Name | `GET /session/{id}/element/{eid}/name` | `page.evaluate(fn, el)` | ✅ |
| 38 | Get Element Rect | `GET /session/{id}/element/{eid}/rect` | `el.bounds()` | ✅ |
| 39 | Is Element Enabled | `GET /session/{id}/element/{eid}/enabled` | `el.isEnabled()` | ✅ |
| 40 | Get Computed Role | `GET /session/{id}/element/{eid}/computedrole` | `el.role()` | ✅ |
| 41 | Get Computed Label | `GET /session/{id}/element/{eid}/computedlabel` | `el.label()` | ✅ |

## Element Interaction (3 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 42 | Element Click | `POST /session/{id}/element/{eid}/click` | `el.click()` | ✅ |
| 43 | Element Clear | `POST /session/{id}/element/{eid}/clear` | `el.clear()` | ✅ |
| 44 | Element Send Keys | `POST /session/{id}/element/{eid}/value` | `el.type(text)` | ✅ |

## Document (3 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 45 | Get Page Source | `GET /session/{id}/source` | `page.content()` | ✅ |
| 46 | Execute Script | `POST /session/{id}/execute/sync` | `page.evaluate(expr)` | ✅ |
| 47 | Execute Async Script | `POST /session/{id}/execute/async` | `page.evaluate(asyncExpr)` | ✅ |

## Cookies (5 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 48 | Get All Cookies | `GET /session/{id}/cookie` | `context.cookies()` | ✅ |
| 49 | Get Named Cookie | `GET /session/{id}/cookie/{name}` | `context.cookies({name})` | ✅ |
| 50 | Add Cookie | `POST /session/{id}/cookie` | `context.setCookies([c])` | ✅ |
| 51 | Delete Cookie | `DELETE /session/{id}/cookie/{name}` | `context.clearCookies({name})` | ✅ |
| 52 | Delete All Cookies | `DELETE /session/{id}/cookie` | `context.clearCookies()` | ✅ |

## Actions (2 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 53 | Perform Actions | `POST /session/{id}/actions` | `page.keyboard.* / page.mouse.*` | ✅ |
| 54 | Release Actions | `DELETE /session/{id}/actions` | (automatic) | ⬜ |

## User Prompts (4 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 55 | Dismiss Alert | `POST /session/{id}/alert/dismiss` | `dialog.dismiss()` | ✅ |
| 56 | Accept Alert | `POST /session/{id}/alert/accept` | `dialog.accept()` | ✅ |
| 57 | Get Alert Text | `GET /session/{id}/alert/text` | `dialog.message()` | ✅ |
| 58 | Send Alert Text | `POST /session/{id}/alert/text` | `dialog.accept(text)` | ✅ |

## Screen Capture (2 commands)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 59 | Take Screenshot | `GET /session/{id}/screenshot` | `page.screenshot()` | ✅ |
| 60 | Take Element Screenshot | `GET /session/{id}/element/{eid}/screenshot` | `el.screenshot()` | ✅ |

## Print (1 command)

| # | WebDriver | Endpoint | Vibium | Status |
|---|-----------|----------|--------|--------|
| 61 | Print Page | `POST /session/{id}/print` | `page.pdf()` | ✅ |

---

**Total: 61 commands** · ✅ 49 · 🟡 0 · ⬜ 12
