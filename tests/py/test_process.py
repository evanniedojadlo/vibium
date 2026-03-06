"""Process tests — cleanup, multiple sessions (3 tests)."""

import pytest


def test_sync_cleanup(test_server):
    """Sync browser closes cleanly."""
    from vibium import browser
    bro = browser.start(headless=True)
    vibe = bro.page()
    vibe.go(test_server)
    assert vibe.title() == "Test App"
    bro.stop()
    # After close, launching another should work
    bro2 = browser.start(headless=True)
    vibe2 = bro2.page()
    vibe2.go(test_server)
    assert vibe2.title() == "Test App"
    bro2.stop()


async def test_async_cleanup(test_server):
    """Async browser closes cleanly."""
    from vibium.async_api import browser
    bro = await browser.start(headless=True)
    vibe = await bro.page()
    await vibe.go(test_server)
    assert await vibe.title() == "Test App"
    await bro.stop()
    # After close, launching another should work
    bro2 = await browser.start(headless=True)
    vibe2 = await bro2.page()
    await vibe2.go(test_server)
    assert await vibe2.title() == "Test App"
    await bro2.stop()


def test_multiple_sessions(test_server):
    """Multiple browser sessions can coexist."""
    from vibium import browser
    bro1 = browser.start(headless=True)
    bro2 = browser.start(headless=True)
    try:
        v1 = bro1.page()
        v2 = bro2.page()
        v1.go(test_server)
        v2.go(test_server + "/subpage")
        assert v1.title() == "Test App"
        assert v2.title() == "Subpage"
    finally:
        bro1.stop()
        bro2.stop()
