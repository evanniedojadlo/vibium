/**
 * JS Library Tests: Object Model
 * Verifies the Browser → Page → BrowserContext object model works end-to-end.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser, Browser, Page, BrowserContext } = require('../../clients/javascript/dist');

describe('JS Object Model', () => {
  test('browser.launch() returns Browser instance', async () => {
    const b = await browser.launch({ headless: true });
    try {
      assert.ok(b instanceof Browser, 'Should return a Browser instance');
    } finally {
      await b.close();
    }
  });

  test('browser.page() returns Page for default tab', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      assert.ok(page instanceof Page, 'Should return a Page instance');
      assert.ok(page.id, 'Page should have an id');
    } finally {
      await b.close();
    }
  });

  test('browser.newPage() creates a new tab', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page1 = await b.page();
      const page2 = await b.newPage();

      assert.ok(page2 instanceof Page, 'Should return a Page instance');
      assert.notStrictEqual(page1.id, page2.id, 'New page should have different context ID');
    } finally {
      await b.close();
    }
  });

  test('browser.pages() returns all open pages', async () => {
    const b = await browser.launch({ headless: true });
    try {
      await b.newPage();
      const pages = await b.pages();

      // At least 2 pages: initial tab + newly created
      assert.ok(pages.length >= 2, `Should have at least 2 pages, got ${pages.length}`);
      for (const page of pages) {
        assert.ok(page instanceof Page, 'Each page should be a Page instance');
      }
    } finally {
      await b.close();
    }
  });

  test('browser.newContext() creates isolated context', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const ctx = await b.newContext();
      assert.ok(ctx instanceof BrowserContext, 'Should return a BrowserContext instance');
      assert.ok(ctx.id, 'Context should have an id');

      const page = await ctx.newPage();
      assert.ok(page instanceof Page, 'context.newPage() should return a Page');

      await ctx.close();
    } finally {
      await b.close();
    }
  });

  test('page.go() + page.url() round-trip', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const url = await page.url();
      assert.ok(url.includes('the-internet.herokuapp.com'), `URL should contain domain, got: ${url}`);
    } finally {
      await b.close();
    }
  });

  test('page.title() returns page title', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const title = await page.title();
      assert.match(title, /The Internet/i, 'Should return page title');
    } finally {
      await b.close();
    }
  });

  test('page.close() closes a tab', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page2 = await b.newPage();
      const pagesBefore = await b.pages();

      await page2.close();

      const pagesAfter = await b.pages();
      assert.ok(pagesAfter.length < pagesBefore.length, 'Should have fewer pages after closing');
    } finally {
      await b.close();
    }
  });
});
