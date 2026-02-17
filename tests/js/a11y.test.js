/**
 * JS Library Tests: Accessibility (a11yTree, el.role, el.label)
 * Tests page.a11yTree(), el.role(), el.label().
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

// --- el.role() ---

describe('Element Accessibility: role()', () => {
  test('role() returns "link" for <a> element', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const link = await page.find('a');
      const role = await link.role();
      assert.strictEqual(role, 'link');
    } finally {
      await b.close();
    }
  });

  test('role() returns "heading" for <h1> element', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const h1 = await page.find('h1');
      const role = await h1.role();
      assert.strictEqual(role, 'heading');
    } finally {
      await b.close();
    }
  });

  test('role() reads explicit role attribute', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<div role="alert" id="msg">Error!</div>');

      const el = await page.find('#msg');
      const role = await el.role();
      assert.strictEqual(role, 'alert');
    } finally {
      await b.close();
    }
  });

  test('fluent: find().role() chains', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const role = await page.find('a').role();
      assert.strictEqual(role, 'link');
    } finally {
      await b.close();
    }
  });
});

// --- el.label() ---

describe('Element Accessibility: label()', () => {
  test('label() returns accessible name for a link', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const link = await page.find('a');
      const label = await link.label();
      assert.ok(label.length > 0, `label should not be empty, got: "${label}"`);
    } finally {
      await b.close();
    }
  });

  test('label() reads aria-label', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<button aria-label="Close dialog">X</button>');

      const btn = await page.find('button');
      const label = await btn.label();
      assert.strictEqual(label, 'Close dialog');
    } finally {
      await b.close();
    }
  });

  test('label() resolves aria-labelledby', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent(`
        <span id="lbl">Username</span>
        <input id="inp" aria-labelledby="lbl" />
      `);

      const input = await page.find('#inp');
      const label = await input.label();
      assert.strictEqual(label, 'Username');
    } finally {
      await b.close();
    }
  });

  test('label() resolves associated <label for="id">', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent(`
        <label for="email">Email Address</label>
        <input id="email" type="email" />
      `);

      const input = await page.find('#email');
      const label = await input.label();
      assert.strictEqual(label, 'Email Address');
    } finally {
      await b.close();
    }
  });

  test('fluent: find().label() chains', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<button aria-label="Submit form">Go</button>');

      const label = await page.find('button').label();
      assert.strictEqual(label, 'Submit form');
    } finally {
      await b.close();
    }
  });
});

// --- page.a11yTree() ---

describe('Page Accessibility: a11yTree()', () => {
  test('returns tree with WebArea root and document title', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const tree = await page.a11yTree();
      assert.strictEqual(tree.role, 'WebArea');
      assert.strictEqual(tree.name, 'Example Domain');
      assert.ok(Array.isArray(tree.children), 'tree should have children');
    } finally {
      await b.close();
    }
  });

  test('tree contains heading and link roles on example.com', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://example.com');

      const tree = await page.a11yTree();

      function findRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) {
            roles.push(...findRoles(child));
          }
        }
        return roles;
      }

      const roles = findRoles(tree);
      assert.ok(roles.includes('heading'), `tree should contain a heading role, got: ${roles.join(', ')}`);
      assert.ok(roles.includes('link'), `tree should contain a link role, got: ${roles.join(', ')}`);
    } finally {
      await b.close();
    }
  });

  test('interestingOnly: false includes generic nodes', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<div><span>hello</span></div>');

      const tree = await page.a11yTree({ interestingOnly: false });

      function findRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) {
            roles.push(...findRoles(child));
          }
        }
        return roles;
      }

      const roles = findRoles(tree);
      assert.ok(roles.includes('generic'), `interestingOnly:false should include generic roles, got: ${roles.join(', ')}`);
    } finally {
      await b.close();
    }
  });

  test('interestingOnly: true (default) filters generic nodes', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<div><span>hello</span></div>');

      const tree = await page.a11yTree();

      function findRoles(node) {
        const roles = [node.role];
        if (node.children) {
          for (const child of node.children) {
            roles.push(...findRoles(child));
          }
        }
        return roles;
      }

      const roles = findRoles(tree);
      assert.ok(!roles.includes('generic'), `interestingOnly:true should filter generic roles, got: ${roles.join(', ')}`);
    } finally {
      await b.close();
    }
  });

  test('captures checked state on checkbox', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent(`
        <input type="checkbox" id="cb" checked />
        <label for="cb">Accept</label>
      `);

      const tree = await page.a11yTree();

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
      await b.close();
    }
  });

  test('captures disabled state on button', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<button disabled>Submit</button>');

      const tree = await page.a11yTree();

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

      const btn = findByRole(tree, 'button');
      assert.ok(btn, 'tree should contain a button');
      assert.strictEqual(btn.disabled, true, 'button should be disabled');
    } finally {
      await b.close();
    }
  });

  test('captures heading levels', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent('<h1>One</h1><h2>Two</h2><h3>Three</h3>');

      const tree = await page.a11yTree();

      function findAll(node, role) {
        const found = [];
        if (node.role === role) found.push(node);
        if (node.children) {
          for (const child of node.children) {
            found.push(...findAll(child, role));
          }
        }
        return found;
      }

      const headings = findAll(tree, 'heading');
      assert.ok(headings.length >= 3, `should have at least 3 headings, got ${headings.length}`);

      const levels = headings.map(h => h.level);
      assert.ok(levels.includes(1), 'should have h1 level=1');
      assert.ok(levels.includes(2), 'should have h2 level=2');
      assert.ok(levels.includes(3), 'should have h3 level=3');
    } finally {
      await b.close();
    }
  });

  test('root option scopes tree to a subtree', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.setContent(`
        <div>
          <h1>Outside</h1>
          <nav id="sidebar">
            <a href="/a">Link A</a>
            <a href="/b">Link B</a>
          </nav>
        </div>
      `);

      const tree = await page.a11yTree({ root: '#sidebar' });

      function findAll(node, role) {
        const found = [];
        if (node.role === role) found.push(node);
        if (node.children) {
          for (const child of node.children) {
            found.push(...findAll(child, role));
          }
        }
        return found;
      }

      const links = findAll(tree, 'link');
      const headings = findAll(tree, 'heading');
      assert.ok(links.length >= 2, `scoped tree should have at least 2 links, got ${links.length}`);
      assert.strictEqual(headings.length, 0, 'scoped tree should not contain the heading outside the root');
    } finally {
      await b.close();
    }
  });
});
