/**
 * JS Library Tests: Lifecycle
 * Tests browser.page(), newPage(), newContext(), pages(), close(),
 * context.newPage(), context.close(), page.activate(), page.close()
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

describe('JS Lifecycle', () => {
  test('browser.page() returns default page', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      assert.ok(page, 'Should return a page');
      assert.ok(page.id, 'Page should have an id');
    } finally {
      await b.close();
    }
  });

  test('browser.newPage() creates new tab with unique ID', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page1 = await b.page();
      const page2 = await b.newPage();
      assert.notStrictEqual(page1.id, page2.id, 'Pages should have different IDs');
    } finally {
      await b.close();
    }
  });

  test('browser.pages() lists all tabs', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const pagesBefore = await b.pages();
      await b.newPage();
      await b.newPage();
      const pagesAfter = await b.pages();

      assert.ok(
        pagesAfter.length >= pagesBefore.length + 2,
        `Should have at least 2 more pages. Before: ${pagesBefore.length}, After: ${pagesAfter.length}`
      );
    } finally {
      await b.close();
    }
  });

  test('page.close() removes a tab', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const newPage = await b.newPage();
      const pagesBefore = await b.pages();

      await newPage.close();

      const pagesAfter = await b.pages();
      assert.strictEqual(
        pagesAfter.length,
        pagesBefore.length - 1,
        'Should have one fewer page'
      );
    } finally {
      await b.close();
    }
  });

  test('page.bringToFront() activates a tab', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page1 = await b.page();
      const page2 = await b.newPage();

      // Activate page1 (should not throw)
      await page1.bringToFront();
      assert.ok(true, 'bringToFront should succeed');
    } finally {
      await b.close();
    }
  });

  test('browser.newContext() creates isolated context', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const ctx = await b.newContext();
      assert.ok(ctx.id, 'Context should have an id');

      const page = await ctx.newPage();
      assert.ok(page.id, 'Page in new context should have an id');

      // Navigate in the new context
      await page.go('https://the-internet.herokuapp.com/');
      const title = await page.title();
      assert.match(title, /The Internet/i, 'Should navigate in new context');

      await ctx.close();
    } finally {
      await b.close();
    }
  });

  test('context.close() removes all pages in context', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const ctx = await b.newContext();
      await ctx.newPage();
      await ctx.newPage();

      const pagesBefore = await b.pages();
      await ctx.close();
      const pagesAfter = await b.pages();

      assert.ok(
        pagesAfter.length < pagesBefore.length,
        'Closing context should remove its pages'
      );
    } finally {
      await b.close();
    }
  });

  test('multiple pages can navigate independently', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page1 = await b.page();
      const page2 = await b.newPage();

      await page1.go('https://the-internet.herokuapp.com/');
      await page2.go('https://the-internet.herokuapp.com/login');

      const url1 = await page1.url();
      const url2 = await page2.url();

      assert.ok(!url1.includes('/login'), 'Page 1 should not be on login');
      assert.ok(url2.includes('/login'), 'Page 2 should be on login');
    } finally {
      await b.close();
    }
  });

  test('browser.close() shuts down cleanly', async () => {
    const b = await browser.launch({ headless: true });
    const page = await b.page();
    await page.go('https://the-internet.herokuapp.com/');

    // close() should not throw
    await b.close();
    assert.ok(true, 'browser.close() should complete without error');
  });
});
