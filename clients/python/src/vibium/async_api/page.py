"""Async Page class — the main browser automation interface."""

from __future__ import annotations

import base64
import fnmatch
import re
from typing import Any, Callable, Dict, List, Optional, Union, TYPE_CHECKING

from .._types import A11yNode, BoundingBox, ElementInfo
from .element import Element
from .element_list import ElementList
from .clock import Clock
from .route import Route
from .network import Request, Response
from .dialog import Dialog
from .console import ConsoleMessage
from .download import Download
from .websocket_info import WebSocketInfo

if TYPE_CHECKING:
    from ..client import BiDiClient


def _match_pattern(pattern: str, url: str) -> bool:
    """Match a URL against a glob-like pattern."""
    if pattern == "**":
        return True
    if "*" in pattern:
        return fnmatch.fnmatch(url, pattern)
    return pattern in url


class Keyboard:
    """Page-level keyboard input."""

    def __init__(self, client: BiDiClient, context_id: str) -> None:
        self._client = client
        self._context_id = context_id

    async def press(self, key: str) -> None:
        await self._client.send("vibium:keyboard.press", {"context": self._context_id, "key": key})

    async def down(self, key: str) -> None:
        await self._client.send("vibium:keyboard.down", {"context": self._context_id, "key": key})

    async def up(self, key: str) -> None:
        await self._client.send("vibium:keyboard.up", {"context": self._context_id, "key": key})

    async def type(self, text: str) -> None:
        await self._client.send("vibium:keyboard.type", {"context": self._context_id, "text": text})


class Mouse:
    """Page-level mouse input."""

    def __init__(self, client: BiDiClient, context_id: str) -> None:
        self._client = client
        self._context_id = context_id

    async def click(self, x: float, y: float) -> None:
        await self._client.send("vibium:mouse.click", {"context": self._context_id, "x": x, "y": y})

    async def move(self, x: float, y: float) -> None:
        await self._client.send("vibium:mouse.move", {"context": self._context_id, "x": x, "y": y})

    async def down(self) -> None:
        await self._client.send("vibium:mouse.down", {"context": self._context_id})

    async def up(self) -> None:
        await self._client.send("vibium:mouse.up", {"context": self._context_id})

    async def wheel(self, delta_x: float, delta_y: float) -> None:
        await self._client.send("vibium:mouse.wheel", {
            "context": self._context_id, "x": 0, "y": 0,
            "deltaX": delta_x, "deltaY": delta_y,
        })


class Touch:
    """Page-level touch input."""

    def __init__(self, client: BiDiClient, context_id: str) -> None:
        self._client = client
        self._context_id = context_id

    async def tap(self, x: float, y: float) -> None:
        await self._client.send("vibium:touch.tap", {"context": self._context_id, "x": x, "y": y})


class Page:
    """Async page automation interface."""

    def __init__(self, client: BiDiClient, context_id: str) -> None:
        self._client = client
        self._context_id = context_id

        self.keyboard = Keyboard(client, context_id)
        self.mouse = Mouse(client, context_id)
        self.touch = Touch(client, context_id)
        self.clock = Clock(client, context_id)

        # Event state
        self._routes: List[Dict[str, Any]] = []
        self._request_callbacks: List[Callable] = []
        self._response_callbacks: List[Callable] = []
        self._dialog_callbacks: List[Callable] = []
        self._console_callbacks: List[Callable] = []
        self._error_callbacks: List[Callable] = []
        self._download_callbacks: List[Callable] = []
        self._pending_downloads: Dict[str, Download] = {}
        self._ws_callbacks: List[Callable] = []
        self._ws_connections: Dict[int, WebSocketInfo] = {}
        self._intercept_id: Optional[str] = None
        self._data_collector_id: Optional[str] = None

        # Register event handler
        self._event_handler = self._handle_event
        self._client.on_event(self._event_handler)

    @property
    def id(self) -> str:
        return self._context_id

    # --- Navigation ---

    async def go(self, url: str) -> None:
        """Navigate to a URL."""
        await self._client.send("vibium:page.navigate", {"context": self._context_id, "url": url})

    async def back(self) -> None:
        await self._client.send("vibium:page.back", {"context": self._context_id})

    async def forward(self) -> None:
        await self._client.send("vibium:page.forward", {"context": self._context_id})

    async def reload(self) -> None:
        await self._client.send("vibium:page.reload", {"context": self._context_id})

    # --- Info ---

    async def url(self) -> str:
        result = await self._client.send("vibium:page.url", {"context": self._context_id})
        return result["url"]

    async def title(self) -> str:
        result = await self._client.send("vibium:page.title", {"context": self._context_id})
        return result["title"]

    async def content(self) -> str:
        result = await self._client.send("vibium:page.content", {"context": self._context_id})
        return result["content"]

    # --- Finding ---

    async def find(
        self,
        selector: Optional[str] = None,
        /,
        *,
        role: Optional[str] = None,
        text: Optional[str] = None,
        label: Optional[str] = None,
        placeholder: Optional[str] = None,
        alt: Optional[str] = None,
        title: Optional[str] = None,
        testid: Optional[str] = None,
        xpath: Optional[str] = None,
        near: Optional[str] = None,
        timeout: Optional[int] = None,
    ) -> Element:
        """Find an element by CSS selector or semantic options."""
        params: Dict[str, Any] = {"context": self._context_id, "timeout": timeout}
        if selector is not None:
            params["selector"] = selector
        else:
            for key, val in [("role", role), ("text", text), ("label", label),
                             ("placeholder", placeholder), ("alt", alt), ("title", title),
                             ("testid", testid), ("xpath", xpath), ("near", near)]:
                if val is not None:
                    params[key] = val

        result = await self._client.send("vibium:find", params)
        info = ElementInfo(tag=result["tag"], text=result["text"], box=BoundingBox(**result["box"]))
        sel_str = selector or ""
        sel_params = {"selector": selector} if selector else {
            k: v for k, v in params.items() if k not in ("context", "timeout")
        }
        return Element(self._client, self._context_id, sel_str, info, None, sel_params)

    async def find_all(
        self,
        selector: Optional[str] = None,
        /,
        *,
        role: Optional[str] = None,
        text: Optional[str] = None,
        label: Optional[str] = None,
        placeholder: Optional[str] = None,
        alt: Optional[str] = None,
        title: Optional[str] = None,
        testid: Optional[str] = None,
        xpath: Optional[str] = None,
        near: Optional[str] = None,
        timeout: Optional[int] = None,
    ) -> ElementList:
        """Find all elements matching a selector or semantic options."""
        params: Dict[str, Any] = {"context": self._context_id, "timeout": timeout}
        if selector is not None:
            params["selector"] = selector
        else:
            for key, val in [("role", role), ("text", text), ("label", label),
                             ("placeholder", placeholder), ("alt", alt), ("title", title),
                             ("testid", testid), ("xpath", xpath), ("near", near)]:
                if val is not None:
                    params[key] = val

        result = await self._client.send("vibium:findAll", params)
        sel_str = selector or ""
        sel_params = {"selector": selector} if selector else {
            k: v for k, v in params.items() if k not in ("context", "timeout")
        }
        elements = []
        for el in result["elements"]:
            info = ElementInfo(tag=el["tag"], text=el["text"], box=BoundingBox(**el["box"]))
            elements.append(Element(self._client, self._context_id, sel_str, info, el.get("index"), sel_params))
        return ElementList(self._client, self._context_id, selector or params, elements)

    # --- Waiting ---

    async def wait_for(
        self,
        selector: Optional[str] = None,
        /,
        *,
        role: Optional[str] = None,
        text: Optional[str] = None,
        label: Optional[str] = None,
        placeholder: Optional[str] = None,
        alt: Optional[str] = None,
        title: Optional[str] = None,
        testid: Optional[str] = None,
        xpath: Optional[str] = None,
        near: Optional[str] = None,
        timeout: Optional[int] = None,
    ) -> Element:
        """Wait for a selector to appear on the page."""
        params: Dict[str, Any] = {"context": self._context_id, "timeout": timeout}
        if selector is not None:
            params["selector"] = selector
        else:
            for key, val in [("role", role), ("text", text), ("label", label),
                             ("placeholder", placeholder), ("alt", alt), ("title", title),
                             ("testid", testid), ("xpath", xpath), ("near", near)]:
                if val is not None:
                    params[key] = val

        result = await self._client.send("vibium:page.waitFor", params)
        info = ElementInfo(tag=result["tag"], text=result["text"], box=BoundingBox(**result["box"]))
        sel_str = selector or ""
        sel_params = {"selector": selector} if selector else {
            k: v for k, v in params.items() if k not in ("context", "timeout")
        }
        return Element(self._client, self._context_id, sel_str, info, None, sel_params)

    async def wait(self, ms: int) -> None:
        """Wait for a fixed amount of time (milliseconds)."""
        await self._client.send("vibium:page.wait", {"context": self._context_id, "ms": ms})

    async def wait_for_url(self, pattern: str, timeout: Optional[int] = None) -> None:
        """Wait until the page URL matches a pattern."""
        await self._client.send("vibium:page.waitForURL", {
            "context": self._context_id, "pattern": pattern, "timeout": timeout,
        })

    async def wait_for_load(self, state: Optional[str] = None, timeout: Optional[int] = None) -> None:
        """Wait until the page reaches a load state."""
        await self._client.send("vibium:page.waitForLoad", {
            "context": self._context_id, "state": state, "timeout": timeout,
        })

    async def wait_for_function(self, fn: str, timeout: Optional[int] = None) -> Any:
        """Wait until a function returns a truthy value."""
        result = await self._client.send("vibium:page.waitForFunction", {
            "context": self._context_id, "fn": fn, "timeout": timeout,
        })
        return result["value"]

    # --- Screenshots & PDF ---

    async def screenshot(
        self,
        full_page: Optional[bool] = None,
        clip: Optional[Dict[str, Any]] = None,
    ) -> bytes:
        """Take a screenshot. Returns PNG bytes."""
        result = await self._client.send("vibium:page.screenshot", {
            "context": self._context_id,
            "fullPage": full_page,
            "clip": clip,
        })
        return base64.b64decode(result["data"])

    async def pdf(self) -> bytes:
        """Print the page to PDF. Returns PDF bytes. Only works in headless mode."""
        result = await self._client.send("vibium:page.pdf", {"context": self._context_id})
        return base64.b64decode(result["data"])

    # --- Evaluation ---

    async def eval(self, expression: str) -> Any:
        """Evaluate a JS expression and return the deserialized value."""
        result = await self._client.send("vibium:page.eval", {
            "context": self._context_id, "expression": expression,
        })
        return result["value"]

    async def evaluate(self, script: str) -> Any:
        """Execute a JS script (multi-statement, use 'return' for values)."""
        result = await self._client.send("script.callFunction", {
            "functionDeclaration": f"() => {{ {script} }}",
            "target": {"context": self._context_id},
            "arguments": [],
            "awaitPromise": True,
            "resultOwnership": "root",
        })
        return result.get("result", {}).get("value")

    async def eval_handle(self, expression: str) -> str:
        """Evaluate a JS expression and return a handle ID."""
        result = await self._client.send("vibium:page.evalHandle", {
            "context": self._context_id, "expression": expression,
        })
        return result["handle"]

    async def add_script(self, source: str) -> None:
        """Inject a script into the page. Pass a URL or inline JavaScript."""
        is_url = source.startswith("http://") or source.startswith("https://") or source.startswith("//")
        params: Dict[str, Any] = {"context": self._context_id}
        if is_url:
            params["url"] = source
        else:
            params["content"] = source
        await self._client.send("vibium:page.addScript", params)

    async def add_style(self, source: str) -> None:
        """Inject a stylesheet into the page. Pass a URL or inline CSS."""
        is_url = source.startswith("http://") or source.startswith("https://") or source.startswith("//")
        params: Dict[str, Any] = {"context": self._context_id}
        if is_url:
            params["url"] = source
        else:
            params["content"] = source
        await self._client.send("vibium:page.addStyle", params)

    async def expose(self, name: str, fn: str) -> None:
        """Expose a function on window."""
        await self._client.send("vibium:page.expose", {
            "context": self._context_id, "name": name, "fn": fn,
        })

    # --- Emulation ---

    async def set_viewport(self, size: Dict[str, int]) -> None:
        """Set the viewport size. size: {width, height}."""
        await self._client.send("vibium:page.setViewport", {
            "context": self._context_id, "width": size["width"], "height": size["height"],
        })

    async def viewport(self) -> Dict[str, int]:
        """Get the current viewport size."""
        return await self._client.send("vibium:page.viewport", {"context": self._context_id})

    async def emulate_media(self, **opts: Any) -> None:
        """Override CSS media features."""
        await self._client.send("vibium:page.emulateMedia", {
            "context": self._context_id, **opts,
        })

    async def set_content(self, html: str) -> None:
        """Replace the page HTML content."""
        await self._client.send("vibium:page.setContent", {"context": self._context_id, "html": html})

    async def set_geolocation(self, coords: Dict[str, float]) -> None:
        """Override the browser's geolocation."""
        await self._client.send("vibium:page.setGeolocation", {
            "context": self._context_id, **coords,
        })

    async def set_window(self, **options: Any) -> None:
        """Set the OS browser window size, position, or state."""
        await self._client.send("vibium:page.setWindow", options)

    async def window(self) -> Dict[str, Any]:
        """Get the current OS browser window state and dimensions."""
        return await self._client.send("vibium:page.window", {})

    # --- Accessibility ---

    async def a11y_tree(
        self,
        interesting_only: Optional[bool] = None,
        root: Optional[str] = None,
    ) -> A11yNode:
        """Get the accessibility tree for the page."""
        params: Dict[str, Any] = {"context": self._context_id}
        if interesting_only is not None:
            params["interestingOnly"] = interesting_only
        if root is not None:
            params["root"] = root
        result = await self._client.send("vibium:page.a11yTree", params)
        return result["tree"]

    # --- Frames ---

    async def frames(self) -> List[Page]:
        """Get all child frames of this page."""
        result = await self._client.send("vibium:page.frames", {"context": self._context_id})
        return [Page(self._client, f["context"]) for f in result["frames"]]

    async def frame(self, name_or_url: str) -> Optional[Page]:
        """Find a frame by name or URL substring."""
        result = await self._client.send("vibium:page.frame", {
            "context": self._context_id, "nameOrUrl": name_or_url,
        })
        if not result or not result.get("context"):
            return None
        return Page(self._client, result["context"])

    def main_frame(self) -> Page:
        """Returns this page — the page IS its own main frame."""
        return self

    # --- Lifecycle ---

    async def bring_to_front(self) -> None:
        await self._client.send("browsingContext.activate", {"context": self._context_id})

    async def close(self) -> None:
        await self._client.send("browsingContext.close", {"context": self._context_id})

    # --- Network Interception ---

    async def route(self, pattern: str, handler: Callable[[Route], Any]) -> None:
        """Intercept network requests matching a URL pattern."""
        if self._intercept_id is None:
            result = await self._client.send("vibium:page.route", {"context": self._context_id})
            self._intercept_id = result["intercept"]

        self._ensure_data_collector()
        self._routes.append({"pattern": pattern, "handler": handler, "interceptId": self._intercept_id})

    async def unroute(self, pattern: str) -> None:
        """Remove a previously registered route."""
        self._routes = [r for r in self._routes if r["pattern"] != pattern]
        if not self._routes and self._intercept_id:
            await self._client.send("network.removeIntercept", {"intercept": self._intercept_id})
            self._intercept_id = None

    def on_request(self, fn: Callable[[Request], None]) -> None:
        """Register a callback for every outgoing request."""
        self._ensure_data_collector()
        self._request_callbacks.append(fn)

    def on_response(self, fn: Callable[[Response], None]) -> None:
        """Register a callback for every completed response."""
        self._ensure_data_collector()
        self._response_callbacks.append(fn)

    async def set_headers(self, headers: Dict[str, str]) -> None:
        """Set extra HTTP headers for all requests in this page."""
        result = await self._client.send("vibium:page.setHeaders", {
            "context": self._context_id, "headers": headers,
        })

        def _header_handler(route: Route) -> None:
            merged = {**route.request.headers(), **headers}
            import asyncio
            asyncio.ensure_future(route.continue_(headers=merged))

        self._routes.append({
            "pattern": "**",
            "handler": _header_handler,
            "interceptId": result["intercept"],
        })

    def on_web_socket(self, fn: Callable[[WebSocketInfo], None]) -> None:
        """Listen for WebSocket connections opened by the page."""
        is_first = len(self._ws_callbacks) == 0
        self._ws_callbacks.append(fn)
        if is_first:
            import asyncio
            asyncio.ensure_future(
                self._client.send("vibium:page.onWebSocket", {"context": self._context_id})
            )

    async def wait_for_request(self, pattern: str, timeout: Optional[int] = None) -> Request:
        """Wait for a request matching a URL pattern."""
        import asyncio
        timeout_ms = timeout or 30000
        future: asyncio.Future = asyncio.get_event_loop().create_future()

        def handler(request: Request) -> None:
            if _match_pattern(pattern, request.url()):
                self._request_callbacks.remove(handler)
                if not future.done():
                    future.set_result(request)

        self._ensure_data_collector()
        self._request_callbacks.append(handler)

        try:
            return await asyncio.wait_for(future, timeout=timeout_ms / 1000)
        except asyncio.TimeoutError:
            if handler in self._request_callbacks:
                self._request_callbacks.remove(handler)
            raise TimeoutError(f"Timeout waiting for request matching '{pattern}'")

    async def wait_for_response(self, pattern: str, timeout: Optional[int] = None) -> Response:
        """Wait for a response matching a URL pattern."""
        import asyncio
        timeout_ms = timeout or 30000
        future: asyncio.Future = asyncio.get_event_loop().create_future()

        def handler(response: Response) -> None:
            if _match_pattern(pattern, response.url()):
                self._response_callbacks.remove(handler)
                if not future.done():
                    future.set_result(response)

        self._ensure_data_collector()
        self._response_callbacks.append(handler)

        try:
            return await asyncio.wait_for(future, timeout=timeout_ms / 1000)
        except asyncio.TimeoutError:
            if handler in self._response_callbacks:
                self._response_callbacks.remove(handler)
            raise TimeoutError(f"Timeout waiting for response matching '{pattern}'")

    # --- Dialog Handling ---

    def on_dialog(self, handler: Callable[[Dialog], Any]) -> None:
        self._dialog_callbacks.append(handler)

    def on_console(self, handler: Callable[[ConsoleMessage], None]) -> None:
        self._console_callbacks.append(handler)

    def on_error(self, handler: Callable[[Exception], None]) -> None:
        self._error_callbacks.append(handler)

    def on_download(self, handler: Callable[[Download], None]) -> None:
        self._download_callbacks.append(handler)

    def remove_all_listeners(self, event: Optional[str] = None) -> None:
        """Remove all listeners for a given event, or all events."""
        if not event or event == "request":
            self._request_callbacks.clear()
        if not event or event == "response":
            self._response_callbacks.clear()
        if not event or event == "dialog":
            self._dialog_callbacks.clear()
        if not event or event == "console":
            self._console_callbacks.clear()
        if not event or event == "error":
            self._error_callbacks.clear()
        if not event or event == "download":
            self._download_callbacks.clear()
        if not event or event == "websocket":
            self._ws_callbacks.clear()
        if (not self._request_callbacks and not self._response_callbacks and not self._routes):
            self._teardown_data_collector()

    # --- Internal Event Handling ---

    def _ensure_data_collector(self) -> None:
        if self._data_collector_id is not None:
            return
        self._data_collector_id = "pending"
        import asyncio

        async def _setup() -> None:
            try:
                result = await self._client.send(
                    "network.addDataCollector",
                    {"dataTypes": ["request", "response"], "maxEncodedDataSize": 10 * 1024 * 1024},
                )
                self._data_collector_id = result["collector"]
            except Exception:
                self._data_collector_id = None

        asyncio.ensure_future(_setup())

    def _teardown_data_collector(self) -> None:
        cid = self._data_collector_id
        if not cid or cid == "pending":
            self._data_collector_id = None
            return
        self._data_collector_id = None
        import asyncio
        asyncio.ensure_future(
            self._client.send("network.removeDataCollector", {"collector": cid})
        )

    def _handle_event(self, event: Dict[str, Any]) -> None:
        """Dispatch a BiDi event to the appropriate handler."""
        params = event.get("params", {})
        event_context = params.get("context")

        # Filter events to this page's context
        if event_context and event_context != self._context_id:
            # log.entryAdded uses source.context
            method = event.get("method", "")
            if method == "log.entryAdded":
                source = params.get("source", {})
                if source.get("context") != self._context_id:
                    return
            else:
                return

        method = event.get("method", "")

        if method == "network.beforeRequestSent":
            self._handle_before_request_sent(params)
        elif method == "network.responseCompleted":
            self._handle_response_completed(params)
        elif method == "browsingContext.userPromptOpened":
            self._handle_user_prompt_opened(params)
        elif method == "browsingContext.downloadWillBegin":
            self._handle_download_will_begin(params)
        elif method == "browsingContext.downloadEnd":
            self._handle_download_completed(params)
        elif method == "log.entryAdded":
            self._handle_log_entry_added(params)
        elif method == "vibium:ws.created":
            self._handle_ws_created(params)
        elif method == "vibium:ws.message":
            self._handle_ws_message(params)
        elif method == "vibium:ws.closed":
            self._handle_ws_closed(params)

    def _handle_before_request_sent(self, params: Dict[str, Any]) -> None:
        is_blocked = params.get("isBlocked", False)
        request_data = params.get("request", {})
        request_id = request_data.get("request", "")

        if is_blocked and request_id:
            request_url = request_data.get("url", "")
            req = Request(params, self._client)

            for route_entry in self._routes:
                if _match_pattern(route_entry["pattern"], request_url):
                    route = Route(self._client, request_id, req)
                    try:
                        result = route_entry["handler"](route)
                        if hasattr(result, "__await__"):
                            import asyncio
                            asyncio.ensure_future(result)
                    except Exception:
                        pass
                    return

            # No matching route — auto-continue
            import asyncio
            asyncio.ensure_future(
                self._client.send("network.continueRequest", {"request": request_id})
            )
        else:
            req = Request(params, self._client)
            for cb in self._request_callbacks:
                cb(req)

    def _handle_response_completed(self, params: Dict[str, Any]) -> None:
        resp = Response(params, self._client)
        for cb in self._response_callbacks:
            cb(resp)

    def _handle_user_prompt_opened(self, params: Dict[str, Any]) -> None:
        dialog = Dialog(self._client, self._context_id, params)

        if self._dialog_callbacks:
            for cb in self._dialog_callbacks:
                try:
                    result = cb(dialog)
                    if hasattr(result, "__await__"):
                        import asyncio
                        asyncio.ensure_future(result)
                except Exception:
                    pass
        else:
            # Auto-dismiss if no handler registered
            import asyncio
            asyncio.ensure_future(dialog.dismiss())

    def _handle_log_entry_added(self, params: Dict[str, Any]) -> None:
        entry_type = params.get("type", "")
        if entry_type == "console":
            msg = ConsoleMessage(params)
            for cb in self._console_callbacks:
                cb(msg)
        elif entry_type == "javascript":
            text = params.get("text", "Unknown error")
            error = RuntimeError(text)
            for cb in self._error_callbacks:
                cb(error)

    def _handle_download_will_begin(self, params: Dict[str, Any]) -> None:
        url = params.get("url", "")
        filename = params.get("suggestedFilename", "")
        navigation = params.get("navigation", "")

        download = Download(self._client, url, filename)
        if navigation:
            self._pending_downloads[navigation] = download

        for cb in self._download_callbacks:
            cb(download)

    def _handle_download_completed(self, params: Dict[str, Any]) -> None:
        navigation = params.get("navigation", "")
        status = params.get("status", "complete")
        filepath = params.get("filepath")

        download = self._pending_downloads.pop(navigation, None)
        if download:
            download._complete(status, filepath)

    def _handle_ws_created(self, params: Dict[str, Any]) -> None:
        ws_id = params.get("id", 0)
        url = params.get("url", "")
        ws = WebSocketInfo(url)
        self._ws_connections[ws_id] = ws
        for cb in self._ws_callbacks:
            cb(ws)

    def _handle_ws_message(self, params: Dict[str, Any]) -> None:
        ws_id = params.get("id", 0)
        data = params.get("data", "")
        direction = params.get("direction", "received")
        ws = self._ws_connections.get(ws_id)
        if ws:
            ws._emit_message(data, direction)

    def _handle_ws_closed(self, params: Dict[str, Any]) -> None:
        ws_id = params.get("id", 0)
        code = params.get("code")
        reason = params.get("reason")
        ws = self._ws_connections.pop(ws_id, None)
        if ws:
            ws._emit_close(code, reason)
