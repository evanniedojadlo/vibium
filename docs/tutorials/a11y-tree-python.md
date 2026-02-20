# Finding Elements with the Accessibility Tree (Python)

Learn how to inspect page structure with `a11y_tree()` and use the results to find and interact with elements.

---

## What You'll Learn

- How to get a page's accessibility tree
- How to use tree output to build selectors for `find()`
- Filtering with `interesting_only` and scoping with `root`
- When to use semantic selectors vs CSS selectors

---

## Getting the Tree

`a11y_tree()` returns the page's accessibility tree — a structured view of every element as assistive technology sees it.

**Sync:**

```python
from vibium import browser

bro = browser.launch()
vibe = bro.page()
vibe.go("https://example.com")

tree = vibe.a11y_tree()
print(tree)

bro.close()
```

**Async:**

```python
import asyncio
from vibium.async_api import browser

async def main():
    bro = await browser.launch()
    vibe = await bro.page()
    await vibe.go("https://example.com")

    tree = await vibe.a11y_tree()
    print(tree)

    await bro.close()

asyncio.run(main())
```

Here's an example of what the tree looks like for a login form:

```python
{
    "role": "WebArea",
    "name": "Login",
    "children": [
        {"role": "heading", "level": 1},
        {"role": "textbox", "name": "Username"},
        {"role": "textbox", "name": "Password"},
        {"role": "checkbox", "name": "Remember me", "checked": False},
        {"role": "button", "name": "Sign in"}
    ]
}
```

Each node has a `role` (what the element is). Nodes with explicit labels (`aria-label`, `<label for>`, etc.) also have a `name`. Some nodes have state properties like `checked`, `disabled`, or `expanded`.

---

## From Tree to Action

The accessibility tree tells you what's on the page. You can then use `find()` to locate and interact with those elements.

### Semantic selectors

`find()` accepts semantic keyword arguments that correspond to what you see in the tree:

```python
# Tree shows: {"role": "button", "name": "Sign in"}
# The name comes from aria-label, so use label:
vibe.find(role="button", label="Sign in").click()

# For buttons/links where the name comes from text content, use text:
vibe.find(role="button", text="Submit").click()
```

**Which parameter maps to the tree's `name`?** It depends on the source:

| Source of the name | find() parameter |
|---|---|
| `aria-label` attribute | `label` |
| `aria-labelledby` reference | `label` |
| `<label for="id">` element | `label` |
| Visible text content | `text` |
| `alt` attribute (images) | `alt` |
| `placeholder` attribute | `placeholder` |
| `title` attribute | `title` |

### CSS selectors

CSS selectors always work for both finding and reading element state:

```python
heading = vibe.find("h1")
print(heading.text())  # Read text content

input_el = vibe.find("#username")
print(input_el.value())  # Read input value
```

### Using tree data in code

You can read the tree programmatically and use its data to drive actions — useful for scripts and AI agents that discover page structure at runtime.

```python
tree = vibe.a11y_tree()

# Walk the tree to find a node by role
def find_by_role(node, role):
    if node["role"] == role:
        return node
    for child in node.get("children", []):
        found = find_by_role(child, role)
        if found:
            return found
    return None

# Discover the button's name from the tree, then click it
btn = find_by_role(tree, "button")
print(btn["name"])  # "Sign in"
vibe.find(role="button", label=btn["name"]).click()
```

The tree also exposes element state. For example, you can check whether a checkbox is already checked before clicking it:

```python
checkbox = find_by_role(tree, "checkbox")

if not checkbox.get("checked"):
    vibe.find(role="checkbox", label=checkbox["name"]).click()
```

---

## Scoping with `root`

On complex pages, the full tree can be large. Use `root` to inspect just one section:

```python
# Only get the tree for the nav element
nav_tree = vibe.a11y_tree(root="nav")

# Only get the tree for a specific element
form_tree = vibe.a11y_tree(root="#login-form")
```

The `root` parameter accepts a CSS selector. The tree will only include that element and its descendants.

---

## Filtering with `interesting_only`

By default, `a11y_tree()` sets `interesting_only` to `True`, which hides generic container nodes (divs, spans with no semantic role). This keeps the output focused on meaningful elements.

Set it to `False` to see everything:

```python
# Default: only semantic elements
tree = vibe.a11y_tree()

# Show all nodes including generic containers
full_tree = vibe.a11y_tree(interesting_only=False)
```

**When to use `interesting_only=False`:**
- Debugging layout issues where you need to see the full DOM structure
- When elements you expect aren't appearing in the default tree

**When to keep the default (`True`):**
- Most of the time — the filtered tree is much easier to read
- When looking for interactive elements (buttons, links, inputs)

---

## Practical Workflow

Here's the full pattern: inspect the tree, then use what you learn to find and interact with elements.

**Sync:**

```python
from vibium import browser

bro = browser.launch()
vibe = bro.page()

vibe.set_content("""
  <h1>Welcome</h1>
  <label for="user">Username</label>
  <input id="user" type="text" />
  <button aria-label="Sign in">Log In</button>
""")

# 1. Inspect the tree to understand the page
tree = vibe.a11y_tree()

# 2. Find the button in the tree and read its name
def find_by_role(node, role):
    if node["role"] == role:
        return node
    for child in node.get("children", []):
        found = find_by_role(child, role)
        if found:
            return found
    return None

btn = find_by_role(tree, "button")
print(f'Found: {btn["role"]} "{btn["name"]}"')  # Found: button "Sign in"

# 3. Fill inputs using CSS selectors
vibe.find("#user").fill("alice")

# 4. Click using the name discovered from the tree
vibe.find(role="button", label=btn["name"]).click()

# 5. Read state using CSS selectors
print("Heading:", vibe.find("h1").text())

bro.close()
```

**Async:**

```python
import asyncio
from vibium.async_api import browser

async def main():
    bro = await browser.launch()
    vibe = await bro.page()

    await vibe.set_content("""
      <h1>Welcome</h1>
      <label for="user">Username</label>
      <input id="user" type="text" />
      <button aria-label="Sign in">Log In</button>
    """)

    # 1. Inspect the tree to understand the page
    tree = await vibe.a11y_tree()

    # 2. Find the button in the tree and read its name
    def find_by_role(node, role):
        if node["role"] == role:
            return node
        for child in node.get("children", []):
            found = find_by_role(child, role)
            if found:
                return found
        return None

    btn = find_by_role(tree, "button")
    print(f'Found: {btn["role"]} "{btn["name"]}"')

    # 3. Fill inputs using CSS selectors
    await vibe.find("#user").fill("alice")

    # 4. Click using the name discovered from the tree
    await vibe.find(role="button", label=btn["name"]).click()

    # 5. Read state using CSS selectors
    print("Heading:", await vibe.find("h1").text())

    await bro.close()

asyncio.run(main())
```

---

## Reference

### a11y_tree() Node Fields

The return value is a dict. Possible keys:

| Field | Type | Description |
|---|---|---|
| `role` | str | ARIA role (e.g. "button", "link", "heading") |
| `name` | str | Accessible name (from aria-label, `<label>`, etc.) |
| `value` | str \| int \| float | Current value (inputs, sliders) |
| `description` | str | Accessible description |
| `children` | list[dict] | Child nodes |
| `disabled` | bool | Whether the element is disabled |
| `checked` | bool \| "mixed" | Checkbox/radio state |
| `pressed` | bool \| "mixed" | Toggle button state |
| `selected` | bool | Whether the element is selected |
| `expanded` | bool | Whether a collapsible is open |
| `focused` | bool | Whether the element has focus |
| `required` | bool | Whether the field is required |
| `readonly` | bool | Whether the field is read-only |
| `level` | int | Heading level (1-6) |
| `valuemin` | int \| float | Minimum value (sliders, spinbuttons) |
| `valuemax` | int \| float | Maximum value (sliders, spinbuttons) |

### find() Keyword Arguments

| Parameter | What it matches |
|---|---|
| `role` | ARIA role |
| `text` | Visible text content (innerText) |
| `label` | Explicit label: `aria-label`, `aria-labelledby`, `<label for>` |
| `placeholder` | Placeholder attribute |
| `alt` | Alt attribute (images) |
| `title` | Title attribute |
| `testid` | `data-testid` attribute |
| `xpath` | XPath expression |
| `near` | CSS selector of a nearby element |
| `timeout` | Max wait time in ms |

### CSS vs Semantic Selectors

| Use CSS selectors when... | Use semantic selectors when... |
|---|---|
| Reading element state (text, value, etc.) | Interacting (click, fill, check) |
| You know the exact HTML structure | Finding by role and label (like a user would) |
| Targeting by class name or ID | The HTML structure might change |

---

## Troubleshooting

### "Element not found" with label

`label` only matches explicit labelling mechanisms (`aria-label`, `aria-labelledby`, `<label for>`). It does **not** match text content. If the element's name comes from its visible text (buttons, links, headings), use `text` instead:

```python
# WRONG: "Submit" comes from text content, not aria-label
vibe.find(role="button", label="Submit")

# RIGHT: use text for text-content-derived names
vibe.find(role="button", text="Submit")
```

### Tree is too large

Use `root` to scope to a section:

```python
nav_tree = vibe.a11y_tree(root="nav")
```

### Tree node has no `name`

The `name` field only appears when the element has an explicit accessible name (via `aria-label`, `<label>`, `alt`, `placeholder`, or `title`). Elements whose name comes only from text content (like `<h1>Title</h1>`) may not show `name` in the tree. The tree still shows their `role`, which you can use with `find()`.

### Everything shows as "generic"

Elements without semantic roles (plain divs, spans) appear as "generic" in the tree. This usually means the page lacks proper semantic HTML or ARIA attributes.

### Python-specific notes

- Method name is `a11y_tree()` (snake_case), not `a11yTree()`
- Option is `interesting_only` (snake_case), not `interestingOnly`
- The tree is a Python dict, not a JavaScript object — access fields with `tree["role"]`, not `tree.role`

---

## Next Steps

- [Getting Started (Python)](getting-started-python.md) — First steps with Vibium in Python
- [Accessibility Tree (JavaScript)](a11y-tree-js.md) — This same tutorial in JavaScript
