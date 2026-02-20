"""Sync network event tests — on_request, on_response (6 sync tests)."""

import time


def test_on_request(sync_browser, test_server):
    """on_request fires for fetch requests."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    urls = []
    vibe.on_request(lambda req: urls.append(req.url()))
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert any("api/data" in u for u in urls)


def test_on_request_method_headers(sync_browser, test_server):
    """on_request captures method and headers."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    captured = {}

    def handler(req):
        if "api/data" in req.url():
            captured["method"] = req.method()
            captured["headers"] = req.headers()

    vibe.on_request(handler)
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert captured.get("method") == "GET"
    assert isinstance(captured.get("headers"), dict)


def test_on_response(sync_browser, test_server):
    """on_response fires for fetch requests."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    statuses = []
    vibe.on_response(lambda resp: statuses.append(resp.status()))
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert 200 in statuses


def test_on_response_url_status(sync_browser, test_server):
    """on_response captures url, status, headers."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    captured = {}

    def handler(resp):
        if "api/data" in resp.url():
            captured["url"] = resp.url()
            captured["status"] = resp.status()
            captured["headers"] = resp.headers()

    vibe.on_response(handler)
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert "api/data" in captured.get("url", "")
    assert captured.get("status") == 200
    assert isinstance(captured.get("headers"), dict)


def test_remove_request_listeners(sync_browser, test_server):
    """remove_all_listeners('request') stops on_request callbacks."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    urls = []
    vibe.on_request(lambda req: urls.append(req.url()))
    vibe.remove_all_listeners("request")
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert len(urls) == 0


def test_remove_response_listeners(sync_browser, test_server):
    """remove_all_listeners('response') stops on_response callbacks."""
    vibe = sync_browser.new_page()
    vibe.go(test_server + "/fetch")
    statuses = []
    vibe.on_response(lambda resp: statuses.append(resp.status()))
    vibe.remove_all_listeners("response")
    vibe.eval("doFetch()")
    vibe.wait(500)
    assert len(statuses) == 0
