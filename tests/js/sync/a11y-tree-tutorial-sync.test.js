/**
 * Tests that verify the sync code examples in docs/tutorials/a11y-tree-js.md are correct.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../../clients/javascript/dist/sync');

describe('A11y Tree Tutorial (JS Sync)', () => {

  // Basic tree structure — sync
  test('a11yTree() returns tree with role and children (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <h1>Welcome</h1>
        <label for="user">Username</label>
        <input id="user" type="text" />
        <button aria-label="Sign in">Log In</button>
      `);

      const tree = vibe.a11yTree();

      assert.strictEqual(tree.role, 'WebArea');
      assert.ok(Array.isArray(tree.children), 'tree should have children');
    } finally {
      bro.close();
    }
  });

  // Semantic find + click with aria-label — sync
  test('find({ role, label }) + click with aria-label (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <div id="result">not clicked</div>
        <button aria-label="Sign in" onclick="document.getElementById('result').textContent='signed in'">Log In</button>
      `);

      vibe.find({ role: 'button', label: 'Sign in' }).click();

      assert.strictEqual(vibe.find('#result').text(), 'signed in');
    } finally {
      bro.close();
    }
  });

  // CSS find + fill — sync
  test('CSS find + fill works for inputs (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <label for="user">Username</label>
        <input id="user" type="text" />
      `);

      vibe.find('#user').fill('alice');

      assert.strictEqual(vibe.find('#user').value(), 'alice');
    } finally {
      bro.close();
    }
  });

  // Semantic find + click with text — sync
  test('find({ role, text }) + click for button text (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <div id="result">waiting</div>
        <button onclick="document.getElementById('result').textContent='done'">Submit</button>
      `);

      vibe.find({ role: 'button', text: 'Submit' }).click();

      assert.strictEqual(vibe.find('#result').text(), 'done');
    } finally {
      bro.close();
    }
  });

  // CSS find + text() — sync
  test('CSS find + text() works (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent('<h1>Welcome</h1>');

      assert.strictEqual(vibe.find('h1').text(), 'Welcome');
    } finally {
      bro.close();
    }
  });

  // Scoping with root — sync
  test('a11yTree({ root }) scopes to selector (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <h1>Title</h1>
        <nav><a href="/a">Link A</a></nav>
      `);

      const navTree = vibe.a11yTree({ root: 'nav' });

      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }

      const roles = collectRoles(navTree);
      assert.ok(roles.includes('link'), 'nav tree should include links');
      assert.ok(!roles.includes('heading'), 'should not include heading outside root');
    } finally {
      bro.close();
    }
  });

  // interestingOnly — sync
  test('interestingOnly: false shows generic nodes (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent('<div><span>hello</span></div>');

      const fullTree = vibe.a11yTree({ interestingOnly: false });

      function collectRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) roles.push(...collectRoles(child));
        }
        return roles;
      }

      assert.ok(collectRoles(fullTree).includes('generic'), 'should show generic nodes');
    } finally {
      bro.close();
    }
  });

  // Tree data flows into find() — sync
  test('tree data flows into find(): discover and click (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <div id="result">not clicked</div>
        <button aria-label="Sign in" onclick="document.getElementById('result').textContent='signed in'">Log In</button>
      `);

      const tree = vibe.a11yTree();

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

      vibe.find({ role: 'button', label: btn.name }).click();
      assert.strictEqual(vibe.find('#result').text(), 'signed in');
    } finally {
      bro.close();
    }
  });

  // Checkbox state from tree — sync
  test('tree state drives action: check unchecked checkbox (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent('<input type="checkbox" aria-label="Remember me" />');

      function findByRole(node, role) {
        if (node.role === role) return node;
        for (const child of node.children || []) {
          const found = findByRole(child, role);
          if (found) return found;
        }
        return null;
      }

      const tree = vibe.a11yTree();
      const checkbox = findByRole(tree, 'checkbox');
      assert.strictEqual(checkbox.checked, false);

      if (!checkbox.checked) {
        vibe.find({ role: 'checkbox', label: checkbox.name }).click();
      }

      const tree2 = vibe.a11yTree();
      const checkbox2 = findByRole(tree2, 'checkbox');
      assert.strictEqual(checkbox2.checked, true);
    } finally {
      bro.close();
    }
  });

  // Practical workflow with tree data — sync
  test('practical workflow: tree data → fill → tree-driven click (sync)', () => {
    const bro = browser.launch({ headless: true });
    try {
      const vibe = bro.page();
      vibe.setContent(`
        <h1>Welcome</h1>
        <label for="user">Username</label>
        <input id="user" type="text" />
        <button aria-label="Sign in">Log In</button>
      `);

      const tree = vibe.a11yTree();
      assert.strictEqual(tree.role, 'WebArea');

      function findByRole(node, role) {
        if (node.role === role) return node;
        for (const child of node.children || []) {
          const found = findByRole(child, role);
          if (found) return found;
        }
        return null;
      }
      const btn = findByRole(tree, 'button');

      vibe.find('#user').fill('alice');
      vibe.find({ role: 'button', label: btn.name }).click();

      assert.strictEqual(vibe.find('h1').text(), 'Welcome');
    } finally {
      bro.close();
    }
  });
});
