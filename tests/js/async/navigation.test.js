/**
 * JS Library Tests: Navigation
 * Tests page.go(), back(), forward(), reload(), url(), title(), content(),
 * waitForURL(), waitForLoad()
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../../clients/javascript/dist');

describe('JS Navigation', () => {
  test('page.go() navigates to URL', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      const url = await vibe.url();
      assert.ok(url.includes('the-internet.herokuapp.com'), 'Should have navigated');
    } finally {
      await bro.close();
    }
  });

  test('page.back() and page.forward()', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      await vibe.go('https://the-internet.herokuapp.com/login');

      const urlAfterNav = await vibe.url();
      assert.ok(urlAfterNav.includes('/login'), 'Should be on login page');

      await vibe.back();
      const urlAfterBack = await vibe.url();
      assert.ok(!urlAfterBack.includes('/login'), 'Should have gone back');

      await vibe.forward();
      const urlAfterForward = await vibe.url();
      assert.ok(urlAfterForward.includes('/login'), 'Should have gone forward');
    } finally {
      await bro.close();
    }
  });

  test('page.reload() reloads the page', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      await vibe.reload();
      const url = await vibe.url();
      assert.ok(url.includes('the-internet.herokuapp.com'), 'Should still be on same page after reload');
    } finally {
      await bro.close();
    }
  });

  test('page.url() returns current URL', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/login');
      const url = await vibe.url();
      assert.ok(url.includes('/login'), 'URL should contain /login');
    } finally {
      await bro.close();
    }
  });

  test('page.title() returns page title', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      const title = await vibe.title();
      assert.match(title, /The Internet/i, 'Should return page title');
    } finally {
      await bro.close();
    }
  });

  test('page.content() returns full HTML', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      const content = await vibe.content();
      assert.ok(content.includes('<html'), 'Should contain <html tag');
      assert.ok(content.includes('Welcome to the-internet'), 'Should contain page content');
    } finally {
      await bro.close();
    }
  });

  test('page.waitForURL() waits for matching URL', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/login');

      // URL should already match — waitForURL should return immediately
      await vibe.waitForURL('/login', { timeout: 5000 });

      const url = await vibe.url();
      assert.ok(url.includes('/login'), 'Should have matched login URL');
    } finally {
      await bro.close();
    }
  });

  test('page.waitForLoad() waits for load state', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');
      await vibe.waitForLoad('complete', { timeout: 10000 });
      // If we get here, it passed
      assert.ok(true);
    } finally {
      await bro.close();
    }
  });

  test('page.waitForURL() times out on mismatch', async () => {
    const bro = await browser.launch({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go('https://the-internet.herokuapp.com/');

      await assert.rejects(
        () => vibe.waitForURL('**/nonexistent-page-xyz', { timeout: 1000 }),
        /timeout/i,
        'Should timeout when URL does not match'
      );
    } finally {
      await bro.close();
    }
  });
});
