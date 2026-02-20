/**
 * Tests that verify the code examples in docs/tutorials/a11y-tree-js.md are correct.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../../clients/javascript/dist');

describe('A11y Tree Tutorial (JS Async)', () => {

  // "Getting the Tree" section — basic tree structure
  test('a11yTree() returns tree with role and children', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <h1>Welcome</h1>
        <label for="user">Username</label>
        <input id="user" type="text" />
        <button aria-label="Sign in">Log In</button>
      `);

      const tree = await vibe.a11yTree();

      assert.strictEqual(tree.role, 'WebArea');
      assert.ok(Array.isArray(tree.children), 'tree should have children');

      // Tree should contain heading, textbox, and button roles
      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }
      const roles = collectRoles(tree);
      assert.ok(roles.includes('heading'), 'should have heading');
      assert.ok(roles.includes('textbox'), 'should have textbox');
      assert.ok(roles.includes('button'), 'should have button');
    } finally {
      await bro.close();
    }
  });

  // Tree shows names from aria-label and <label for>
  test('a11yTree() shows names from explicit labels', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <label for="user">Username</label>
        <input id="user" type="text" />
        <button aria-label="Sign in">Log In</button>
      `);

      const tree = await vibe.a11yTree();

      function findByRole(node, role) {
        if (node.role === role) return node;
        if (node.children) {
          for (const child of node.children) {
            const found = findByRole(child, role);
            if (found) return found;
          }
        }
        return null;
      }

      const textbox = findByRole(tree, 'textbox');
      assert.ok(textbox, 'tree should have textbox');
      assert.strictEqual(textbox.name, 'Username');

      const button = findByRole(tree, 'button');
      assert.ok(button, 'tree should have button');
      assert.strictEqual(button.name, 'Sign in');
    } finally {
      await bro.close();
    }
  });

  // "From Tree to Action" — semantic find + click with aria-label
  test('find({ role, label }) + click works with aria-label', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <div id="result">not clicked</div>
        <button aria-label="Sign in" onclick="document.getElementById('result').textContent='signed in'">Log In</button>
      `);

      await vibe.find({ role: 'button', label: 'Sign in' }).click();

      const result = await vibe.find('#result');
      assert.strictEqual(await result.text(), 'signed in');
    } finally {
      await bro.close();
    }
  });

  // CSS find + fill (fill requires CSS selector)
  test('CSS find + fill works for inputs', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <label for="user">Username</label>
        <input id="user" type="text" />
      `);

      await vibe.find('#user').fill('alice');

      const input = await vibe.find('#user');
      assert.strictEqual(await input.value(), 'alice');
    } finally {
      await bro.close();
    }
  });

  // semantic find + click with text content
  test('find({ role, text }) + click works for button text', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <div id="result">waiting</div>
        <button onclick="document.getElementById('result').textContent='done'">Submit</button>
      `);

      await vibe.find({ role: 'button', text: 'Submit' }).click();

      const result = await vibe.find('#result');
      assert.strictEqual(await result.text(), 'done');
    } finally {
      await bro.close();
    }
  });

  // CSS find + text() for reading state
  test('CSS find + text() works for reading state', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent('<h1>Welcome</h1>');

      const heading = await vibe.find('h1');
      assert.strictEqual(await heading.text(), 'Welcome');
    } finally {
      await bro.close();
    }
  });

  // "Scoping with root"
  test('a11yTree({ root }) scopes to CSS selector', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <h1>Title</h1>
        <nav><a href="/a">Link A</a><a href="/b">Link B</a></nav>
      `);

      const navTree = await vibe.a11yTree({ root: 'nav' });

      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }

      const roles = collectRoles(navTree);
      assert.ok(roles.includes('link'), 'nav tree should include links');
      assert.ok(!roles.includes('heading'), 'nav tree should not include heading outside root');
    } finally {
      await bro.close();
    }
  });

  // "Filtering with interestingOnly"
  test('interestingOnly: false includes generic nodes', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent('<div><span>hello</span></div>');

      const fullTree = await vibe.a11yTree({ interestingOnly: false });

      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }

      assert.ok(collectRoles(fullTree).includes('generic'), 'should show generic nodes');
    } finally {
      await bro.close();
    }
  });

  test('default interestingOnly filters generic nodes', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent('<div><span>hello</span></div>');

      const tree = await vibe.a11yTree();

      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }

      assert.ok(!collectRoles(tree).includes('generic'), 'should filter generic nodes');
    } finally {
      await bro.close();
    }
  });

  // Tree shows checked state
  test('a11yTree() captures checkbox state', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <label><input type="checkbox" checked /> Remember me</label>
      `);

      const tree = await vibe.a11yTree();

      function findByRole(node, role) {
        if (node.role === role) return node;
        if (node.children) {
          for (const child of node.children) {
            const found = findByRole(child, role);
            if (found) return found;
          }
        }
        return null;
      }

      const checkbox = findByRole(tree, 'checkbox');
      assert.ok(checkbox, 'tree should contain a checkbox');
      assert.strictEqual(checkbox.checked, true, 'checkbox should be checked');
    } finally {
      await bro.close();
    }
  });

  // "Using tree data in code" — tree name flows into find()
  test('tree data flows into find(): discover button name, then click', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <div id="result">not clicked</div>
        <button aria-label="Sign in" onclick="document.getElementById('result').textContent='signed in'">Log In</button>
      `);

      const tree = await vibe.a11yTree();

      function findByRole(node, role) {
        if (node.role === role) return node;
        for (const child of node.children || []) {
          const found = findByRole(child, role);
          if (found) return found;
        }
        return null;
      }

      const btn = findByRole(tree, 'button');
      assert.ok(btn, 'tree should contain a button');
      assert.strictEqual(btn.name, 'Sign in');

      // Use the name from the tree to click
      await vibe.find({ role: 'button', label: btn.name }).click();

      const result = await vibe.find('#result');
      assert.strictEqual(await result.text(), 'signed in');
    } finally {
      await bro.close();
    }
  });

  // "Using tree data in code" — checkbox state drives conditional action
  test('tree state drives action: check unchecked checkbox', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <input type="checkbox" aria-label="Remember me" />
      `);

      const tree = await vibe.a11yTree();

      function findByRole(node, role) {
        if (node.role === role) return node;
        for (const child of node.children || []) {
          const found = findByRole(child, role);
          if (found) return found;
        }
        return null;
      }

      const checkbox = findByRole(tree, 'checkbox');
      assert.ok(checkbox, 'tree should contain a checkbox');
      assert.strictEqual(checkbox.checked, false, 'should start unchecked');

      // Use tree state to decide whether to click
      if (!checkbox.checked) {
        await vibe.find({ role: 'checkbox', label: checkbox.name }).click();
      }

      // Verify it's now checked
      const tree2 = await vibe.a11yTree();
      const checkbox2 = findByRole(tree2, 'checkbox');
      assert.strictEqual(checkbox2.checked, true, 'should now be checked');
    } finally {
      await bro.close();
    }
  });

  // "Practical Workflow" — full flow using tree data
  test('practical workflow: tree data → CSS fill → tree-driven click → CSS read', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.setContent(`
        <h1>Welcome</h1>
        <label for="user">Username</label>
        <input id="user" type="text" />
        <button aria-label="Sign in" onclick="document.getElementById('user').value='submitted'">Log In</button>
      `);

      // 1. Inspect the tree
      const tree = await vibe.a11yTree();
      assert.strictEqual(tree.role, 'WebArea');

      // 2. Find button in tree
      function findByRole(node, role) {
        if (node.role === role) return node;
        for (const child of node.children || []) {
          const found = findByRole(child, role);
          if (found) return found;
        }
        return null;
      }
      const btn = findByRole(tree, 'button');
      assert.strictEqual(btn.name, 'Sign in');

      // 3. Fill with CSS, click using name from tree
      await vibe.find('#user').fill('alice');
      await vibe.find({ role: 'button', label: btn.name }).click();

      // 4. Read state with CSS
      const heading = await vibe.find('h1');
      assert.strictEqual(await heading.text(), 'Welcome');
    } finally {
      await bro.close();
    }
  });
});
