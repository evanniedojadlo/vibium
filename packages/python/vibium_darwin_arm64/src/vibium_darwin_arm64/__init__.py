"""Vibium clicker binary for macOS ARM64."""

from pathlib import Path

__version__ = "0.1.3"

def get_binary_path() -> str:
    """Get the path to the clicker binary."""
    return str(Path(__file__).parent / "bin" / "clicker")
