# Finding Elements with the Accessibility Tree

Learn how to inspect page structure with `a11yTree()` and use the results to find and interact with elements.

---

## What You'll Learn

- How to get a page's accessibility tree
- How to use tree output to build selectors for `find()`
- `everything` and `root` options
- When to use semantic selectors vs CSS selectors

---

## Getting the Tree

`a11yTree()` returns the page's accessibility tree — a structured view of every element as assistive technology sees it.

**Sync:**

```javascript
const { browser } = require('vibium/sync')

const bro = browser.launch()
const vibe = bro.page()
vibe.go('https://example.com')

const tree = vibe.a11yTree()
console.log(JSON.stringify(tree, null, 2))

bro.close()
```

**Async:**

```javascript
const { browser } = require('vibium')

async function main() {
  const bro = await browser.launch()
  const vibe = await bro.page()
  await vibe.go('https://example.com')

  const tree = await vibe.a11yTree()
  console.log(JSON.stringify(tree, null, 2))

  await bro.close()
}

main()
```

Here's an example of what the tree looks like for a login form:

```json
{
  "role": "WebArea",
  "name": "Login",
  "children": [
    { "role": "heading", "level": 1 },
    { "role": "textbox", "name": "Username" },
    { "role": "textbox", "name": "Password" },
    { "role": "checkbox", "name": "Remember me", "checked": false },
    { "role": "button", "name": "Sign in" }
  ]
}
```

Each node has a `role` (what the element is). Nodes with explicit labels (`aria-label`, `<label for>`, etc.) also have a `name`. Some nodes have state properties like `checked`, `disabled`, or `expanded`.

---

## From Tree to Action

The accessibility tree tells you what's on the page. You can then use `find()` to locate and interact with those elements.

### Semantic selectors

`find()` accepts semantic options that correspond to what you see in the tree:

```javascript
// Tree shows: { role: "button", name: "Sign in" }
// The name comes from aria-label, so use label:
vibe.find({ role: 'button', label: 'Sign in' }).click()

// For buttons/links where the name comes from text content, use text:
vibe.find({ role: 'button', text: 'Submit' }).click()
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

```javascript
const heading = vibe.find('h1')
console.log(heading.text()) // Read text content

const input = vibe.find('#username')
console.log(input.value()) // Read input value
```

### Using tree data in code

You can read the tree programmatically and use its data to drive actions — useful for scripts and AI agents that discover page structure at runtime.

```javascript
const tree = vibe.a11yTree()

// Walk the tree to find a node by role
function findByRole(node, role) {
  if (node.role === role) return node
  for (const child of node.children || []) {
    const found = findByRole(child, role)
    if (found) return found
  }
  return null
}

// Discover the button's name from the tree, then click it
const btn = findByRole(tree, 'button')
console.log(btn.name) // "Sign in"
vibe.find({ role: 'button', label: btn.name }).click()
```

The tree also exposes element state. For example, you can check whether a checkbox is already checked before clicking it:

```javascript
const checkbox = findByRole(tree, 'checkbox')

if (!checkbox.checked) {
  vibe.find({ role: 'checkbox', label: checkbox.name }).click()
}
```

---

## Scoping with `root`

On complex pages, the full tree can be large. Use `root` to inspect just one section:

```javascript
// Only get the tree for the nav element
const navTree = vibe.a11yTree({ root: 'nav' })

// Only get the tree for a specific element
const formTree = vibe.a11yTree({ root: '#login-form' })
```

The `root` parameter accepts a CSS selector. The tree will only include that element and its descendants.

---

## Filtering with `everything`

By default, `a11yTree()` hides generic container nodes (divs, spans with no semantic role). This keeps the output focused on meaningful elements.

Set `everything: true` to see all nodes:

```javascript
// Default: only semantic elements
const tree = vibe.a11yTree()

// Show all nodes including generic containers
const fullTree = vibe.a11yTree({ everything: true })
```

**When to use `everything: true`:**
- Debugging layout issues where you need to see the full DOM structure
- When elements you expect aren't appearing in the default tree

**When to keep the default:**
- Most of the time — the filtered tree is much easier to read
- When looking for interactive elements (buttons, links, inputs)

---

## Practical Workflow

Here's the full pattern: inspect the tree, then use what you learn to find and interact with elements.

**Sync:**

```javascript
const { browser } = require('vibium/sync')

const bro = browser.launch()
const vibe = bro.page()

vibe.setContent(`
  <h1>Welcome</h1>
  <label for="user">Username</label>
  <input id="user" type="text" />
  <button aria-label="Sign in">Log In</button>
`)

// 1. Inspect the tree to understand the page
const tree = vibe.a11yTree()

// 2. Find the button in the tree and read its name
function findByRole(node, role) {
  if (node.role === role) return node
  for (const child of node.children || []) {
    const found = findByRole(child, role)
    if (found) return found
  }
  return null
}
const btn = findByRole(tree, 'button')
console.log(`Found: ${btn.role} "${btn.name}"`) // Found: button "Sign in"

// 3. Fill inputs using CSS selectors
vibe.find('#user').fill('alice')

// 4. Click using the name discovered from the tree
vibe.find({ role: 'button', label: btn.name }).click()

// 5. Read state using CSS selectors
console.log('Heading:', vibe.find('h1').text())

bro.close()
```

**Async:**

```javascript
const { browser } = require('vibium')

async function main() {
  const bro = await browser.launch()
  const vibe = await bro.page()

  await vibe.setContent(`
    <h1>Welcome</h1>
    <label for="user">Username</label>
    <input id="user" type="text" />
    <button aria-label="Sign in">Log In</button>
  `)

  // 1. Inspect the tree to understand the page
  const tree = await vibe.a11yTree()

  // 2. Find the button in the tree and read its name
  function findByRole(node, role) {
    if (node.role === role) return node
    for (const child of node.children || []) {
      const found = findByRole(child, role)
      if (found) return found
    }
    return null
  }
  const btn = findByRole(tree, 'button')
  console.log(`Found: ${btn.role} "${btn.name}"`) // Found: button "Sign in"

  // 3. Fill inputs using CSS selectors
  await vibe.find('#user').fill('alice')

  // 4. Click using the name discovered from the tree
  await vibe.find({ role: 'button', label: btn.name }).click()

  // 5. Read state using CSS selectors
  console.log('Heading:', await vibe.find('h1').text())

  await bro.close()
}

main()
```

---

## Reference

### a11yTree() Node Fields

| Field | Type | Description |
|---|---|---|
| `role` | string | ARIA role (e.g. "button", "link", "heading") |
| `name` | string | Accessible name (from aria-label, `<label>`, etc.) |
| `value` | string \| number | Current value (inputs, sliders) |
| `description` | string | Accessible description |
| `children` | A11yNode[] | Child nodes |
| `disabled` | boolean | Whether the element is disabled |
| `checked` | boolean \| 'mixed' | Checkbox/radio state |
| `pressed` | boolean \| 'mixed' | Toggle button state |
| `selected` | boolean | Whether the element is selected |
| `expanded` | boolean | Whether a collapsible is open |
| `focused` | boolean | Whether the element has focus |
| `required` | boolean | Whether the field is required |
| `readonly` | boolean | Whether the field is read-only |
| `level` | number | Heading level (1-6) |
| `valuemin` | number | Minimum value (sliders, spinbuttons) |
| `valuemax` | number | Maximum value (sliders, spinbuttons) |

### find() Selector Options

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

```javascript
// WRONG: "Submit" comes from text content, not aria-label
vibe.find({ role: 'button', label: 'Submit' })

// RIGHT: use text for text-content-derived names
vibe.find({ role: 'button', text: 'Submit' })
```

### Tree is too large

Use `root` to scope to a section:

```javascript
const navTree = vibe.a11yTree({ root: 'nav' })
```

### Tree node has no `name`

The `name` field only appears when the element has an explicit accessible name (via `aria-label`, `<label>`, `alt`, `placeholder`, or `title`). Elements whose name comes only from text content (like `<h1>Title</h1>`) may not show `name` in the tree. The tree still shows their `role`, which you can use with `find()`.

### Everything shows as "generic"

Elements without semantic roles (plain divs, spans) appear as "generic" in the tree. This usually means the page lacks proper semantic HTML or ARIA attributes.

---

## Next Steps

- [Getting Started](getting-started.md) — First steps with Vibium (JavaScript)
- [Accessibility Tree (Python)](a11y-tree-python.md) — This same tutorial in Python
