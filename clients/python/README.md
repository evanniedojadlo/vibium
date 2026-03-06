# Vibium Python Client

Browser automation for AI agents and humans.

## Installation

```bash
pip install vibium
```

Chrome downloads automatically on first use. To install Chrome ahead of time:

```bash
vibium install
```

## Quick Start

```python
from vibium import browser

bro = browser.start()
vibe = bro.page()
vibe.go("https://example.com")

# Take a screenshot
png = vibe.screenshot()
with open("screenshot.png", "wb") as f:
    f.write(png)

# Find and click a link
link = vibe.find("a")
print(link.text())
link.click()

bro.stop()
```

## Async API

```python
import asyncio
from vibium.async_api import browser

async def main():
    bro = await browser.start()
    vibe = await bro.page()
    await vibe.go("https://example.com")

    link = await vibe.find("a")
    await link.click()

    await bro.stop()

asyncio.run(main())
```

## CLI

```bash
vibium install   # Download Chrome for Testing
vibium version   # Show version
```

## Requirements

- Python 3.9+

## License

Apache-2.0
