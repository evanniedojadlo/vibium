"""Parses annotated markdown tutorials and provides test runners.

Annotations (same format as the JS tutorial runner):
    <!-- helpers -->              — next code block defines shared helper functions
    <!-- test: async "name" -->   — next code block is an async test
    <!-- test: sync "name" -->    — next code block is a sync test
"""

import re
import textwrap
from pathlib import Path

# tests/py/helpers/ -> project root
PROJECT_ROOT = Path(__file__).resolve().parents[3]


def extract_blocks(md_path):
    """Parse markdown and extract annotated ``python`` code blocks."""
    content = (PROJECT_ROOT / md_path).read_text()
    lines = content.split("\n")
    blocks = []

    pending = None
    in_code_block = False
    is_annotated = False
    code_lines = []

    for line in lines:
        if not in_code_block:
            if re.match(r"<!--\s*helpers\s*-->", line):
                pending = {"type": "helpers"}
                continue

            m = re.match(r'<!--\s*test:\s*(async|sync)\s+"([^"]+)"\s*-->', line)
            if m:
                pending = {"type": "test", "mode": m.group(1), "name": m.group(2)}
                continue

            if re.match(r"^```python\s*$", line):
                in_code_block = True
                if pending:
                    is_annotated = True
                    code_lines = []
                else:
                    is_annotated = False
                continue
        else:
            if re.match(r"^```\s*$", line):
                in_code_block = False
                if is_annotated and pending:
                    blocks.append({**pending, "code": "\n".join(code_lines)})
                    pending = None
                continue
            if is_annotated:
                code_lines.append(line)

    return blocks


def get_tutorial_tests(md_path, mode):
    """Return ``[(name, helpers, code), ...]`` for *pytest.mark.parametrize*."""
    blocks = extract_blocks(md_path)
    helpers = ""
    tests = []

    for block in blocks:
        if block["type"] == "helpers":
            helpers += block["code"] + "\n"
            continue
        if block.get("mode") != mode:
            continue
        tests.append((block["name"], helpers, block["code"]))

    return tests


async def run_async_block(code, vibe):
    """Exec a block of async Python code with *vibe* as the page object."""
    indented = textwrap.indent(code, "    ")
    wrapped = f"async def _run(vibe):\n{indented}\n"
    ns = {}
    exec(compile(wrapped, "<tutorial>", "exec"), ns)
    await ns["_run"](vibe)


def run_sync_block(code, vibe):
    """Exec a block of sync Python code with *vibe* as the page object."""
    exec(compile(code, "<tutorial>", "exec"), {"vibe": vibe})
