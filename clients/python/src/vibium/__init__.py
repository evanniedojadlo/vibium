"""Vibium - Browser automation for AI agents and humans.

Usage (sync, default):
    from vibium import browser
    bro = browser.start()
    vibe = bro.new_page()
    vibe.go("https://example.com")
    bro.stop()

Usage (async):
    from vibium.async_api import browser
    bro = await browser.start()
    vibe = await bro.new_page()
    await vibe.go("https://example.com")
    await bro.stop()
"""

from .sync_api.browser import browser, Browser
from .sync_api.page import Page
from .sync_api.element import Element
from .sync_api.context import BrowserContext

__version__ = "26.3.11"
__all__ = ["browser", "Browser", "Page", "Element", "BrowserContext"]
