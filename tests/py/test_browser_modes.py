"""Browser mode tests — headless, headed, default (3 sync tests)."""

from vibium import browser


def test_headless(test_server):
    bro = browser.start(headless=True)
    try:
        vibe = bro.page()
        vibe.go(test_server)
        assert vibe.title() == "Test App"
    finally:
        bro.stop()


def test_headed(test_server):
    bro = browser.start(headless=False)
    try:
        vibe = bro.page()
        vibe.go(test_server)
        assert vibe.title() == "Test App"
    finally:
        bro.stop()


def test_default_visible(test_server):
    """Default launch() is not headless (browser visible)."""
    bro = browser.start()
    try:
        vibe = bro.page()
        vibe.go(test_server)
        assert vibe.title() == "Test App"
    finally:
        bro.stop()
