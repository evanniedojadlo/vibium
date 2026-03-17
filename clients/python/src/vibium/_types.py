"""Shared types for the Vibium Python client."""

from __future__ import annotations

from dataclasses import dataclass
from typing import List, Optional, TypedDict, Union


@dataclass
class BoundingBox:
    """Bounding box of an element."""

    x: float
    y: float
    width: float
    height: float


@dataclass
class ElementInfo:
    """Information about an element returned from find commands."""

    tag: str
    text: str
    box: BoundingBox


class Cookie(TypedDict, total=False):
    """Browser cookie."""

    name: str
    value: str
    domain: str
    path: str
    size: int
    httpOnly: bool
    secure: bool
    sameSite: str
    expiry: int


class SetCookieParam(TypedDict, total=False):
    """Parameters for setting a cookie."""

    name: str
    value: str
    domain: str
    url: str
    path: str
    httpOnly: bool
    secure: bool
    sameSite: str
    expiry: int


class _StorageEntry(TypedDict):
    name: str
    value: str


class OriginState(TypedDict):
    """Storage state for a single origin."""

    origin: str
    localStorage: List[_StorageEntry]
    sessionStorage: List[_StorageEntry]


class StorageState(TypedDict):
    """Full browser storage state."""

    cookies: List[Cookie]
    origins: List[OriginState]


class A11yNode(TypedDict, total=False):
    """Accessibility tree node."""

    role: str
    name: str
    value: Union[str, int, float]
    description: str
    disabled: bool
    expanded: bool
    focused: bool
    checked: Union[bool, str]
    pressed: Union[bool, str]
    selected: bool
    required: bool
    readonly: bool
    level: int
    valuemin: int
    valuemax: int
    children: List[A11yNode]


class ViewportSize(TypedDict):
    """Browser viewport dimensions."""

    width: int
    height: int


class WindowInfo(TypedDict):
    """Browser window state and dimensions."""

    state: str
    width: int
    height: int
    x: int
    y: int
