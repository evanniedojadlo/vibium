# Remote Browser Control

Run Chrome on one machine, control it from another.

---

## Server (the machine with the browser)

Install vibium (this downloads Chrome + chromedriver automatically):

```bash
npm install -g vibium
```

Find the chromedriver path and start it:

```bash
vibium paths
# Chromedriver: /Users/you/.cache/vibium/.../chromedriver

$(vibium paths | grep Chromedriver | cut -d' ' -f2) --port=9515 --allowed-ips=""
```

---

## Client (your dev machine)

### JavaScript

```javascript
const { browser } = require('vibium/sync')

const bro = browser.connect('ws://your-server:9515/session')
const page = bro.page()

page.go('https://example.com')
console.log(page.title())        // "Example Domain"
console.log(page.find('h1').text())  // "Example Domain"

bro.close()
```

### Python

```python
from vibium.sync_api import browser

bro = browser.connect("ws://your-server:9515/session")
page = bro.page()

page.go("https://example.com")
print(page.title())          # "Example Domain"
print(page.find("h1").text())    # "Example Domain"

bro.close()
```

---

## With Authentication

If your endpoint requires auth headers (e.g. a cloud browser provider):

```javascript
const bro = browser.connect('wss://cloud.example.com/bidi', {
  headers: { 'Authorization': 'Bearer my-token' }
})
```

```python
bro = browser.connect("wss://cloud.example.com/bidi", headers={
    "Authorization": "Bearer my-token",
})
```

---

## How It Works

```
Client machine                    Server machine
┌──────────┐   stdin/stdout   ┌─────────┐   WebSocket   ┌─────────────┐
│ your code│ ──── pipes ────► │ vibium  │ ────────────► │ chromedriver│
└──────────┘                  └─────────┘               └──────┬──────┘
                                                               │
                                                        ┌──────▼──────┐
                                                        │   Chrome    │
                                                        └─────────────┘
```

`browser.connect()` starts a local vibium process that proxies to the remote chromedriver. All vibium features (auto-wait, screenshots, tracing) work over remote connections.
