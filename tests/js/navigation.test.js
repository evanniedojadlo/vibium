/**
 * JS Library Tests: Navigation
 * Tests page.go(), back(), forward(), reload(), url(), title(), content(),
 * waitForURL(), waitForLoad()
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

describe('JS Navigation', () => {
  test('page.go() navigates to URL', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const url = await page.url();
      assert.ok(url.includes('the-internet.herokuapp.com'), 'Should have navigated');
    } finally {
      await b.close();
    }
  });

  test('page.back() and page.forward()', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      await page.go('https://the-internet.herokuapp.com/login');

      const urlAfterNav = await page.url();
      assert.ok(urlAfterNav.includes('/login'), 'Should be on login page');

      await page.back();
      const urlAfterBack = await page.url();
      assert.ok(!urlAfterBack.includes('/login'), 'Should have gone back');

      await page.forward();
      const urlAfterForward = await page.url();
      assert.ok(urlAfterForward.includes('/login'), 'Should have gone forward');
    } finally {
      await b.close();
    }
  });

  test('page.reload() reloads the page', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      await page.reload();
      const url = await page.url();
      assert.ok(url.includes('the-internet.herokuapp.com'), 'Should still be on same page after reload');
    } finally {
      await b.close();
    }
  });

  test('page.url() returns current URL', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/login');
      const url = await page.url();
      assert.ok(url.includes('/login'), 'URL should contain /login');
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

  test('page.content() returns full HTML', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const content = await page.content();
      assert.ok(content.includes('<html'), 'Should contain <html tag');
      assert.ok(content.includes('Welcome to the-internet'), 'Should contain page content');
    } finally {
      await b.close();
    }
  });

  test('page.waitForURL() waits for matching URL', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/login');

      // URL should already match — waitForURL should return immediately
      await page.waitForURL('/login', { timeout: 5000 });

      const url = await page.url();
      assert.ok(url.includes('/login'), 'Should have matched login URL');
    } finally {
      await b.close();
    }
  });

  test('page.waitForLoad() waits for load state', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      await page.waitForLoad('complete', { timeout: 10000 });
      // If we get here, it passed
      assert.ok(true);
    } finally {
      await b.close();
    }
  });

  test('page.waitForURL() times out on mismatch', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');

      await assert.rejects(
        () => page.waitForURL('**/nonexistent-page-xyz', { timeout: 1000 }),
        /timeout/i,
        'Should timeout when URL does not match'
      );
    } finally {
      await b.close();
    }
  });
});
