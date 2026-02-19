"""Sync Page wrapper."""

from __future__ import annotations

from typing import Any, Callable, Dict, List, Optional, Union, TYPE_CHECKING

from .._types import A11yNode
from .element import Element
from .element_list import ElementList
from .clock import Clock
from .route import Route

if TYPE_CHECKING:
    from .._sync_base import _EventLoopThread
    from ..async_api.page import Page as AsyncPage


class Keyboard:
    """Sync keyboard input."""

    def __init__(self, async_keyboard: Any, loop_thread: _EventLoopThread) -> None:
        self._async = async_keyboard
        self._loop = loop_thread

    def press(self, key: str) -> None:
        self._loop.run(self._async.press(key))

    def down(self, key: str) -> None:
        self._loop.run(self._async.down(key))

    def up(self, key: str) -> None:
        self._loop.run(self._async.up(key))

    def type(self, text: str) -> None:
        self._loop.run(self._async.type(text))


class Mouse:
    """Sync mouse input."""

    def __init__(self, async_mouse: Any, loop_thread: _EventLoopThread) -> None:
        self._async = async_mouse
        self._loop = loop_thread

    def click(self, x: float, y: float) -> None:
        self._loop.run(self._async.click(x, y))

    def move(self, x: float, y: float) -> None:
        self._loop.run(self._async.move(x, y))

    def down(self) -> None:
        self._loop.run(self._async.down())

    def up(self) -> None:
        self._loop.run(self._async.up())

    def wheel(self, delta_x: float, delta_y: float) -> None:
        self._loop.run(self._async.wheel(delta_x, delta_y))


class Touch:
    """Sync touch input."""

    def __init__(self, async_touch: Any, loop_thread: _EventLoopThread) -> None:
        self._async = async_touch
        self._loop = loop_thread

    def tap(self, x: float, y: float) -> None:
        self._loop.run(self._async.tap(x, y))


class Page:
    """Synchronous wrapper for async Page."""

    def __init__(self, async_page: AsyncPage, loop_thread: _EventLoopThread) -> None:
        self._async = async_page
        self._loop = loop_thread

        self.keyboard = Keyboard(async_page.keyboard, loop_thread)
        self.mouse = Mouse(async_page.mouse, loop_thread)
        self.touch = Touch(async_page.touch, loop_thread)
        self.clock = Clock(async_page.clock, loop_thread)

        # Sync event state
        self._console_messages: List[Dict[str, str]] = []
        self._errors: List[Dict[str, str]] = []

    @property
    def id(self) -> str:
        return self._async.id

    # --- Navigation ---

    def go(self, url: str) -> None:
        self._loop.run(self._async.go(url))

    def back(self) -> None:
        self._loop.run(self._async.back())

    def forward(self) -> None:
        self._loop.run(self._async.forward())

    def reload(self) -> None:
        self._loop.run(self._async.reload())

    # --- Info ---

    def url(self) -> str:
        return self._loop.run(self._async.url())

    def title(self) -> str:
        return self._loop.run(self._async.title())

    def content(self) -> str:
        return self._loop.run(self._async.content())

    # --- Finding ---

    def find(
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
        async_el = self._loop.run(self._async.find(
            selector, role=role, text=text, label=label, placeholder=placeholder,
            alt=alt, title=title, testid=testid, xpath=xpath, near=near, timeout=timeout,
        ))
        return Element(async_el, self._loop)

    def find_all(
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
        async_list = self._loop.run(self._async.find_all(
            selector, role=role, text=text, label=label, placeholder=placeholder,
            alt=alt, title=title, testid=testid, xpath=xpath, near=near, timeout=timeout,
        ))
        return ElementList(async_list, self._loop)

    def wait_for(
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
        async_el = self._loop.run(self._async.wait_for(
            selector, role=role, text=text, label=label, placeholder=placeholder,
            alt=alt, title=title, testid=testid, xpath=xpath, near=near, timeout=timeout,
        ))
        return Element(async_el, self._loop)

    # --- Waiting ---

    def wait(self, ms: int) -> None:
        self._loop.run(self._async.wait(ms))

    def wait_for_url(self, pattern: str, timeout: Optional[int] = None) -> None:
        self._loop.run(self._async.wait_for_url(pattern, timeout))

    def wait_for_load(self, state: Optional[str] = None, timeout: Optional[int] = None) -> None:
        self._loop.run(self._async.wait_for_load(state, timeout))

    def wait_for_function(self, fn: str, timeout: Optional[int] = None) -> Any:
        return self._loop.run(self._async.wait_for_function(fn, timeout))

    # --- Screenshots & PDF ---

    def screenshot(
        self,
        full_page: Optional[bool] = None,
        clip: Optional[Dict[str, Any]] = None,
    ) -> bytes:
        return self._loop.run(self._async.screenshot(full_page=full_page, clip=clip))

    def pdf(self) -> bytes:
        return self._loop.run(self._async.pdf())

    # --- Evaluation ---

    def eval(self, expression: str) -> Any:
        return self._loop.run(self._async.eval(expression))

    def evaluate(self, script: str) -> Any:
        """Execute a JS script (multi-statement, use 'return' for values)."""
        return self._loop.run(self._async.evaluate(script))

    def eval_handle(self, expression: str) -> str:
        return self._loop.run(self._async.eval_handle(expression))

    def add_script(self, source: str) -> None:
        self._loop.run(self._async.add_script(source))

    def add_style(self, source: str) -> None:
        self._loop.run(self._async.add_style(source))

    def expose(self, name: str, fn: str) -> None:
        self._loop.run(self._async.expose(name, fn))

    # --- Emulation ---

    def set_viewport(self, size: Dict[str, int]) -> None:
        self._loop.run(self._async.set_viewport(size))

    def viewport(self) -> Dict[str, int]:
        return self._loop.run(self._async.viewport())

    def emulate_media(self, **opts: Any) -> None:
        self._loop.run(self._async.emulate_media(**opts))

    def set_content(self, html: str) -> None:
        self._loop.run(self._async.set_content(html))

    def set_geolocation(self, coords: Dict[str, float]) -> None:
        self._loop.run(self._async.set_geolocation(coords))

    def set_window(self, **options: Any) -> None:
        self._loop.run(self._async.set_window(**options))

    def window(self) -> Dict[str, Any]:
        return self._loop.run(self._async.window())

    # --- Accessibility ---

    def a11y_tree(
        self,
        interesting_only: Optional[bool] = None,
        root: Optional[str] = None,
    ) -> A11yNode:
        return self._loop.run(self._async.a11y_tree(interesting_only, root))

    # --- Frames ---

    def frames(self) -> List[Page]:
        async_frames = self._loop.run(self._async.frames())
        return [Page(f, self._loop) for f in async_frames]

    def frame(self, name_or_url: str) -> Optional[Page]:
        async_frame = self._loop.run(self._async.frame(name_or_url))
        if async_frame is None:
            return None
        return Page(async_frame, self._loop)

    def main_frame(self) -> Page:
        return self

    # --- Lifecycle ---

    def bring_to_front(self) -> None:
        self._loop.run(self._async.bring_to_front())

    def close(self) -> None:
        self._loop.run(self._async.close())

    # --- Network ---

    def route(
        self,
        pattern: str,
        action: Union[str, Dict[str, Any], Callable[[Route], None]],
    ) -> None:
        """Intercept network requests.

        action can be:
          - 'continue' — pass through
          - 'abort' — block the request
          - dict — static fulfill ({status, body, headers})
          - callable — handler function receiving Route
        """
        if isinstance(action, str):
            if action == "abort":
                async def _abort_handler(async_route: Any) -> None:
                    await async_route.abort()
                self._loop.run(self._async.route(pattern, _abort_handler))
            else:  # 'continue'
                async def _continue_handler(async_route: Any) -> None:
                    await async_route.continue_()
                self._loop.run(self._async.route(pattern, _continue_handler))
        elif isinstance(action, dict):
            fulfill_opts = action

            async def _fulfill_handler(async_route: Any) -> None:
                await async_route.fulfill(**fulfill_opts)
            self._loop.run(self._async.route(pattern, _fulfill_handler))
        else:
            # Callable handler — use sync decision pattern
            def _sync_callback(async_route: Any) -> None:
                sync_route = Route(async_route)
                action(sync_route)
                decision = sync_route._decision
                import asyncio
                if decision["action"] == "fulfill":
                    opts = {k: v for k, v in decision.items() if k != "action" and v is not None}
                    asyncio.ensure_future(async_route.fulfill(**opts))
                elif decision["action"] == "abort":
                    asyncio.ensure_future(async_route.abort())
                else:
                    opts = {k: v for k, v in decision.items() if k != "action" and v is not None}
                    asyncio.ensure_future(async_route.continue_(**opts))

            self._loop.run(self._async.route(pattern, _sync_callback))

    def unroute(self, pattern: str) -> None:
        self._loop.run(self._async.unroute(pattern))

    def set_headers(self, headers: Dict[str, str]) -> None:
        self._loop.run(self._async.set_headers(headers))

    def wait_for_request(
        self,
        pattern: str,
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """Wait for a request matching a URL pattern. Returns dict."""
        req = self._loop.run(self._async.wait_for_request(pattern, timeout))
        return {
            "url": req.url(),
            "method": req.method(),
            "headers": req.headers(),
            "post_data": None,
        }

    def wait_for_response(
        self,
        pattern: str,
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """Wait for a response matching a URL pattern. Returns dict."""
        resp = self._loop.run(self._async.wait_for_response(pattern, timeout))
        return {
            "url": resp.url(),
            "status": resp.status(),
            "headers": resp.headers(),
        }

    # --- Events ---

    def on_dialog(
        self,
        action: Union[str, Callable],
    ) -> None:
        """Handle browser dialogs.

        action can be:
          - 'accept' — auto-accept
          - 'dismiss' — auto-dismiss
          - callable — handler function receiving Dialog (sync)
        """
        from .dialog import Dialog as SyncDialog

        if isinstance(action, str):
            async def _simple_handler(dialog: Any) -> None:
                if action == "accept":
                    await dialog.accept()
                else:
                    await dialog.dismiss()
            self._async.on_dialog(_simple_handler)
        else:
            def _sync_callback(dialog: Any) -> None:
                sync_dialog = SyncDialog(dialog)
                action(sync_dialog)
                decision = sync_dialog._decision
                import asyncio
                if decision["action"] == "accept":
                    asyncio.ensure_future(dialog.accept(decision.get("prompt_text")))
                else:
                    asyncio.ensure_future(dialog.dismiss())

            self._async.on_dialog(_sync_callback)

    def on_console(self, mode: str = "collect") -> None:
        """Start collecting console messages. Retrieve with console_messages()."""
        def _collector(msg: Any) -> None:
            self._console_messages.append({"type": msg.type(), "text": msg.text()})
        self._async.on_console(_collector)

    def console_messages(self) -> List[Dict[str, str]]:
        """Return collected console messages."""
        return list(self._console_messages)

    def on_error(self, mode: str = "collect") -> None:
        """Start collecting page errors. Retrieve with errors()."""
        def _collector(error: Exception) -> None:
            self._errors.append({"message": str(error)})
        self._async.on_error(_collector)

    def errors(self) -> List[Dict[str, str]]:
        """Return collected errors."""
        return list(self._errors)

    def remove_all_listeners(self, event: Optional[str] = None) -> None:
        self._async.remove_all_listeners(event)
        if not event or event == "console":
            self._console_messages.clear()
        if not event or event == "error":
            self._errors.clear()
