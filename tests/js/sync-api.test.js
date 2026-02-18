/**
 * JS Library Tests: Sync API
 * Tests browser.launch() and BrowserSync → PageSync → ElementSync.
 *
 * The HTTP server runs in a child process because the sync API blocks
 * the main thread with Atomics.wait(), which would deadlock an in-process server.
 */

const { test, describe, before, after } = require('node:test');
const assert = require('node:assert');
const { fork } = require('child_process');
const path = require('path');

const { browser } = require('../../clients/javascript/dist/sync');

// --- Server child process ---

let serverProcess;
let baseURL;

before(async () => {
  // Start the HTTP server in a child process
  serverProcess = fork(path.join(__dirname, 'sync-test-server.js'), [], { silent: true });

  // Read the base URL from the server's stdout
  baseURL = await new Promise((resolve, reject) => {
    let data = '';
    serverProcess.stdout.on('data', (chunk) => {
      data += chunk.toString();
      const line = data.trim();
      if (line.startsWith('http://')) resolve(line);
    });
    serverProcess.on('error', reject);
    setTimeout(() => reject(new Error('Server startup timeout')), 5000);
  });
});

after(() => {
  if (serverProcess) serverProcess.kill();
});

// --- Tests ---

describe('Sync API: Browser lifecycle', () => {
  test('browser.launch() and close()', () => {
    const bro = browser.launch({ headless: true });
    assert.ok(bro, 'Should return a BrowserSync instance');
    bro.close();
  });

  test('browser.page() returns default page', () => {
    const bro = browser.launch({ headless: true });
    try {
      const page = bro.page();
      assert.ok(page, 'Should return a PageSync');
    } finally {
      bro.close();
    }
  });
});

describe('Sync API: Multi-page', () => {
  test('newPage() creates a new tab', () => {
    const bro = browser.launch({ headless: true });
    try {
      const page1 = bro.page();
      const page2 = bro.newPage();
      assert.ok(page2, 'Should return a new PageSync');
      const allPages = bro.pages();
      assert.ok(allPages.length >= 2, 'Should have at least 2 pages');
    } finally {
      bro.close();
    }
  });
});

describe('Sync API: Navigation', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('go() navigates to URL', () => {
    const page = bro.page();
    page.go(baseURL);
    assert.ok(true, 'Navigation succeeded');
  });

  test('url() returns current URL', () => {
    const page = bro.page();
    page.go(baseURL);
    const url = page.url();
    assert.ok(url.includes('127.0.0.1'), 'Should contain host');
  });

  test('title() returns page title', () => {
    const page = bro.page();
    page.go(baseURL);
    const title = page.title();
    assert.strictEqual(title, 'Test App');
  });

  test('content() returns HTML', () => {
    const page = bro.page();
    page.go(baseURL);
    const html = page.content();
    assert.ok(html.includes('Welcome to test-app'), 'Should contain page content');
  });

  test('back() and forward()', () => {
    const page = bro.page();
    page.go(baseURL);
    page.go(`${baseURL}/subpage`);
    assert.strictEqual(page.title(), 'Subpage');

    page.back();
    assert.strictEqual(page.title(), 'Test App');

    page.forward();
    assert.strictEqual(page.title(), 'Subpage');
  });

  test('reload()', () => {
    const page = bro.page();
    page.go(baseURL);
    page.reload();
    assert.strictEqual(page.title(), 'Test App');
  });
});

describe('Sync API: Screenshots & PDF', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('screenshot() returns PNG buffer', () => {
    const page = bro.page();
    page.go(baseURL);
    const png = page.screenshot();

    assert.ok(Buffer.isBuffer(png), 'Should return a Buffer');
    assert.ok(png.length > 100, 'Should have reasonable size');
    assert.strictEqual(png[0], 0x89, 'PNG magic byte 1');
    assert.strictEqual(png[1], 0x50, 'PNG magic byte 2');
  });

  test('pdf() returns PDF buffer', () => {
    const page = bro.page();
    page.go(baseURL);
    const pdf = page.pdf();

    assert.ok(Buffer.isBuffer(pdf), 'Should return a Buffer');
    assert.ok(pdf.length > 100, 'Should have reasonable size');
    assert.strictEqual(pdf[0], 0x25, 'PDF magic byte');
    assert.strictEqual(pdf[1], 0x50, 'PDF magic byte');
  });
});

describe('Sync API: Evaluation', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('evaluate() executes JavaScript', () => {
    const page = bro.page();
    page.go(baseURL);
    const title = page.evaluate('return document.title');
    assert.strictEqual(title, 'Test App');
  });

  test('eval() evaluates expression', () => {
    const page = bro.page();
    page.go(`${baseURL}/eval`);
    const val = page.eval('window.testVal');
    assert.strictEqual(val, 42);
  });

  test('eval() returns computed value', () => {
    const page = bro.page();
    page.go(baseURL);
    const year = page.eval('new Date().getFullYear()');
    assert.strictEqual(typeof year, 'number');
    assert.ok(year >= 2025);
  });
});

describe('Sync API: Element finding', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('find() locates element by CSS selector', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1.heading');
    assert.ok(heading, 'Should return an ElementSync');
    assert.ok(heading.info, 'Should have info');
    assert.match(heading.info.tag, /^h1$/i);
  });

  test('find() with semantic selector', () => {
    const page = bro.page();
    page.go(baseURL);
    const link = page.find({ role: 'link', text: 'Go to subpage' });
    assert.ok(link, 'Should find link by role+text');
    assert.match(link.info.tag, /^a$/i);
  });

  test('findAll() returns ElementListSync', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const links = page.findAll('a.link');
    assert.ok(links, 'Should return an ElementListSync');
    assert.strictEqual(links.count(), 4);
  });

  test('waitFor() waits for element', () => {
    const page = bro.page();
    page.go(baseURL);
    const el = page.waitFor('h1.heading');
    assert.ok(el, 'Should return an ElementSync');
  });
});

describe('Sync API: ElementListSync', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('count(), first(), last(), nth()', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const items = page.findAll('a.link');

    assert.strictEqual(items.count(), 4);

    const first = items.first();
    assert.ok(first.text().includes('Link 1'));

    const last = items.last();
    assert.ok(last.text().includes('Link 4'));

    const second = items.nth(1);
    assert.ok(second.text().includes('Link 2'));
  });

  test('iteration with for...of', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const items = page.findAll('a.link');
    let count = 0;
    for (const item of items) {
      count++;
      assert.ok(item, 'Each item should be an ElementSync');
    }
    assert.strictEqual(count, 4);
  });

  test('filter()', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const all = page.findAll('a.link');
    const filtered = all.filter({ hasText: 'Link 4' });
    assert.strictEqual(filtered.count(), 1);
  });
});

describe('Sync API: Element interaction', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('click() navigates via link', () => {
    const page = bro.page();
    page.go(baseURL);
    const link = page.find('a[href="/subpage"]');
    link.click();
    page.waitFor('h3'); // wait for subpage to load
    assert.strictEqual(page.title(), 'Subpage');
  });

  test('fill() and value()', () => {
    const page = bro.page();
    page.go(`${baseURL}/inputs`);
    const input = page.find('#text-input');
    input.fill('hello world');
    assert.strictEqual(input.value(), 'hello world');
  });

  test('type() appends text', () => {
    const page = bro.page();
    page.go(`${baseURL}/inputs`);
    const input = page.find('#text-input');
    input.type('12345');
    const value = page.evaluate("return document.querySelector('#text-input').value");
    assert.strictEqual(value, '12345');
  });

  test('check() and uncheck()', () => {
    const page = bro.page();
    page.go(`${baseURL}/form`);
    const checkbox = page.find('#agree');
    checkbox.check();
    assert.strictEqual(checkbox.isChecked(), true);
    checkbox.uncheck();
    assert.strictEqual(checkbox.isChecked(), false);
  });

  test('selectOption()', () => {
    const page = bro.page();
    page.go(`${baseURL}/form`);
    const select = page.find('#color');
    select.selectOption('blue');
    assert.strictEqual(select.value(), 'blue');
  });

  test('hover()', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1');
    heading.hover();
    assert.ok(true, 'hover completed without error');
  });

  test('press() on element', () => {
    const page = bro.page();
    page.go(`${baseURL}/inputs`);
    const input = page.find('#text-input');
    input.click();
    input.press('a');
    const value = page.evaluate("return document.querySelector('#text-input').value");
    assert.ok(typeof value === 'string');
  });
});

describe('Sync API: Element state', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('text() returns textContent', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1.heading');
    const text = heading.text();
    assert.ok(text.includes('Welcome to test-app'));
  });

  test('innerText() returns rendered text', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1.heading');
    const text = heading.innerText();
    assert.ok(text.includes('Welcome to test-app'));
  });

  test('html() returns innerHTML', () => {
    const page = bro.page();
    page.go(baseURL);
    const info = page.find('#info');
    const html = info.html();
    assert.ok(html.includes('Some info text'));
  });

  test('attr() returns attribute value', () => {
    const page = bro.page();
    page.go(baseURL);
    const link = page.find('a[href="/subpage"]');
    assert.strictEqual(link.attr('href'), '/subpage');
  });

  test('bounds() returns bounding box', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1');
    const box = heading.bounds();
    assert.ok(typeof box.x === 'number');
    assert.ok(typeof box.y === 'number');
    assert.ok(box.width > 0);
    assert.ok(box.height > 0);
  });

  test('isVisible() returns true for visible elements', () => {
    const page = bro.page();
    page.go(baseURL);
    const heading = page.find('h1');
    assert.strictEqual(heading.isVisible(), true);
  });

  test('isEnabled() returns true for enabled elements', () => {
    const page = bro.page();
    page.go(`${baseURL}/form`);
    const input = page.find('#name');
    assert.strictEqual(input.isEnabled(), true);
  });

  test('isEditable() returns true for editable elements', () => {
    const page = bro.page();
    page.go(`${baseURL}/form`);
    const input = page.find('#name');
    assert.strictEqual(input.isEditable(), true);
  });
});

describe('Sync API: Scoped find', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('element.find() scoped to parent', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const nested = page.find('#nested');
    const span = nested.find('.inner');
    assert.ok(span.text().includes('span'));
  });

  test('element.findAll() scoped to parent', () => {
    const page = bro.page();
    page.go(`${baseURL}/links`);
    const nested = page.find('#nested');
    const spans = nested.findAll('.inner');
    assert.strictEqual(spans.count(), 2);
  });
});

describe('Sync API: Keyboard, Mouse, Touch', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('keyboard.type() types text', () => {
    const page = bro.page();
    page.go(`${baseURL}/inputs`);
    page.find('#text-input').click();
    page.keyboard.type('hello');
    const value = page.evaluate("return document.querySelector('#text-input').value");
    assert.strictEqual(value, 'hello');
  });

  test('keyboard.press() presses a key', () => {
    const page = bro.page();
    page.go(`${baseURL}/inputs`);
    page.find('#text-input').click();
    page.keyboard.press('a');
    assert.ok(true, 'press completed');
  });

  test('mouse.click() clicks at coordinates', () => {
    const page = bro.page();
    page.go(baseURL);
    page.mouse.click(100, 100);
    assert.ok(true, 'mouse click completed');
  });
});

describe('Sync API: Clock control', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('clock.install() and setFixedTime()', () => {
    const page = bro.page();
    page.go(`${baseURL}/clock`);
    page.clock.install({ time: new Date('2025-06-15T12:00:00Z') });
    page.clock.setFixedTime(new Date('2025-06-15T12:00:00Z'));
    const year = page.eval('new Date().getFullYear()');
    assert.strictEqual(year, 2025);
  });

  test('clock.fastForward()', () => {
    const page = bro.newPage();
    page.go(`${baseURL}/clock`);
    page.clock.install({ time: new Date('2025-06-15T12:00:00Z') });
    page.clock.fastForward(60000);
    const time = page.eval('Date.now()');
    assert.ok(typeof time === 'number');
  });
});

describe('Sync API: Viewport & emulation', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('setViewport() and viewport()', () => {
    const page = bro.page();
    page.setViewport({ width: 375, height: 812 });
    const vp = page.viewport();
    assert.strictEqual(vp.width, 375);
    assert.strictEqual(vp.height, 812);
  });

  test('setContent() replaces page HTML', () => {
    const page = bro.page();
    page.setContent('<html><body><h1>Custom</h1></body></html>');
    const heading = page.find('h1');
    assert.strictEqual(heading.text(), 'Custom');
  });

  test('a11yTree() returns accessibility tree', () => {
    const page = bro.page();
    page.go(baseURL);
    const tree = page.a11yTree();
    assert.ok(tree, 'Should return a tree');
    assert.ok(tree.role, 'Root should have a role');
  });
});

describe('Sync API: Context isolation', () => {
  test('newContext() creates isolated context', () => {
    const bro = browser.launch({ headless: true });
    try {
      const ctx = bro.newContext();
      const page = ctx.newPage();
      page.go(baseURL);
      assert.strictEqual(page.title(), 'Test App');
      ctx.close();
    } finally {
      bro.close();
    }
  });

  test('cookies in context', () => {
    const bro = browser.launch({ headless: true });
    try {
      const ctx = bro.newContext();
      const page = ctx.newPage();
      page.go(baseURL);
      ctx.setCookies([{ name: 'test', value: 'val', url: baseURL }]);
      const cookies = ctx.cookies();
      assert.ok(cookies.some(c => c.name === 'test'), 'Should have the test cookie');
      ctx.clearCookies();
      const cleared = ctx.cookies();
      assert.ok(!cleared.some(c => c.name === 'test'), 'Cookie should be cleared');
      ctx.close();
    } finally {
      bro.close();
    }
  });
});

describe('Sync API: Dialog auto-handling', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('onDialog("accept") auto-accepts alerts', () => {
    const page = bro.page();
    page.go(`${baseURL}/dialog`);
    page.onDialog('accept');
    page.find('#alert-btn').click();
    assert.ok(true, 'Dialog was auto-accepted');
  });
});

describe('Sync API: Full checkpoint', () => {
  test('Phase 8 checkpoint', () => {
    const bro = browser.launch({ headless: true });
    try {
      const page = bro.newPage();
      page.go(baseURL);
      assert.strictEqual(page.title(), 'Test App');
      assert.ok(page.url().includes('127.0.0.1'));

      const link = page.find('a[href="/subpage"]');
      link.click();
      page.waitFor('h3'); // wait for subpage to load
      assert.strictEqual(page.title(), 'Subpage');

      const png = page.screenshot();
      assert.ok(png.length > 100, 'Screenshot should have data');

      const year = page.eval('new Date().getFullYear()');
      assert.strictEqual(typeof year, 'number');
    } finally {
      bro.close();
    }
  });
});
