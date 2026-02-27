"""Pipe client for communicating with the vibium binary via stdin/stdout."""

from __future__ import annotations

import asyncio
import json
from typing import Any, Callable, Dict, List, Optional, TYPE_CHECKING

if TYPE_CHECKING:
    from .binary import VibiumProcess


class BiDiError(Exception):
    """Raised when a BiDi command fails."""

    def __init__(self, error: str, message: str):
        self.error = error
        self.message = message
        super().__init__(f"{error}: {message}")


class BiDiClient:
    """Pipe client for BiDi protocol with event dispatch."""

    def __init__(
        self,
        process: VibiumProcess,
    ):
        self._process = process
        self._stdin = process._process.stdin
        self._stdout = process._process.stdout
        self._next_id = 1
        self._pending: Dict[int, asyncio.Future] = {}
        self._receiver_task: Optional[asyncio.Task] = None
        self._event_handlers: List[Callable[[Dict[str, Any]], None]] = []

    @classmethod
    async def connect(cls, process: VibiumProcess) -> BiDiClient:
        """Create a BiDiClient from a VibiumProcess with pipe streams."""
        client = cls(process)
        client._receiver_task = asyncio.create_task(client._receive_loop())
        # Replay any events that arrived before the ready signal
        if hasattr(process, "_pre_ready_lines"):
            for line in process._pre_ready_lines:
                client._dispatch_message(line)
            del process._pre_ready_lines
        return client

    def _dispatch_message(self, line: str) -> None:
        """Parse and dispatch a single message line."""
        try:
            data = json.loads(line)
        except (json.JSONDecodeError, ValueError):
            return
        msg_id = data.get("id")
        if msg_id is not None and msg_id in self._pending:
            future = self._pending[msg_id]
            if not future.done():
                future.set_result(data)
        elif msg_id is None and "method" in data:
            for handler in self._event_handlers:
                try:
                    handler(data)
                except Exception:
                    pass

    def on_event(self, handler: Callable[[Dict[str, Any]], None]) -> None:
        """Register an event handler for messages without an id (events)."""
        self._event_handlers.append(handler)

    def remove_event_handler(self, handler: Callable[[Dict[str, Any]], None]) -> None:
        """Remove a previously registered event handler."""
        try:
            self._event_handlers.remove(handler)
        except ValueError:
            pass

    async def _receive_loop(self) -> None:
        """Background task to receive and dispatch messages from stdout."""
        try:
            while True:
                line_bytes = await self._stdout.readline()  # type: ignore[union-attr]
                if not line_bytes:
                    break  # EOF — process exited
                line = line_bytes.decode().strip()
                if not line:
                    continue
                self._dispatch_message(line)
        except (asyncio.CancelledError, OSError):
            pass
        finally:
            for future in self._pending.values():
                if not future.done():
                    future.set_exception(ConnectionError("Connection closed"))

    async def send(self, method: str, params: Optional[Dict[str, Any]] = None, timeout: float = 60) -> Any:
        """Send a command and wait for the response."""
        msg_id = self._next_id
        self._next_id += 1

        command = {
            "id": msg_id,
            "method": method,
            "params": params or {},
        }

        future: asyncio.Future = asyncio.get_running_loop().create_future()
        self._pending[msg_id] = future

        try:
            line = json.dumps(command) + "\n"
            self._stdin.write(line.encode())  # type: ignore[union-attr]
            await self._stdin.drain()  # type: ignore[union-attr]
            try:
                response = await asyncio.wait_for(future, timeout=timeout)
            except asyncio.TimeoutError:
                raise TimeoutError(f"Command '{method}' timed out after {timeout}s")

            if response.get("type") == "error":
                raise BiDiError(
                    response.get("error", "unknown"),
                    response.get("message", "Unknown error"),
                )

            return response.get("result")
        finally:
            self._pending.pop(msg_id, None)

    async def close(self) -> None:
        """Close the pipe connection."""
        if self._receiver_task:
            self._receiver_task.cancel()
            try:
                await self._receiver_task
            except asyncio.CancelledError:
                pass

        # Close stdin to signal the pipe process
        if self._stdin and not self._stdin.is_closing():
            self._stdin.close()
