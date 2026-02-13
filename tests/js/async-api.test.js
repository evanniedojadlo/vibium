/**
 * JS Library Tests: Async API
 * Tests browser.launch() and Browser → Page methods
 */

const { test, describe, after } = require('node:test');
const assert = require('node:assert');
const fs = require('node:fs');
const path = require('node:path');

// Import from built library
const { browser } = require('../../clients/javascript/dist');

describe('JS Async API', () => {
  test('browser.launch() and browser.close() work', async () => {
    const b = await browser.launch({ headless: true });
    assert.ok(b, 'Should return a Browser instance');
    await b.close();
  });

  test('page.go() navigates to URL', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      // If we get here without error, navigation worked
      assert.ok(true);
    } finally {
      await b.close();
    }
  });

  test('page.screenshot() returns PNG buffer', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const screenshot = await page.screenshot();

      assert.ok(Buffer.isBuffer(screenshot), 'Should return a Buffer');
      assert.ok(screenshot.length > 1000, 'Screenshot should have reasonable size');

      // Check PNG magic bytes
      assert.strictEqual(screenshot[0], 0x89, 'Should be valid PNG');
      assert.strictEqual(screenshot[1], 0x50, 'Should be valid PNG');
      assert.strictEqual(screenshot[2], 0x4E, 'Should be valid PNG');
      assert.strictEqual(screenshot[3], 0x47, 'Should be valid PNG');
    } finally {
      await b.close();
    }
  });

  test('page.evaluate() executes JavaScript', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const title = await page.evaluate('return document.title');
      assert.match(title, /The Internet/i, 'Should return page title');
    } finally {
      await b.close();
    }
  });

  test('page.find() locates element', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const heading = await page.find('h1.heading');

      assert.ok(heading, 'Should return an Element');
      assert.ok(heading.info, 'Element should have info');
      assert.match(heading.info.tag, /^h1$/i, 'Should be an h1 tag');
      assert.match(heading.info.text, /Welcome to the-internet/i, 'Should have heading text');
    } finally {
      await b.close();
    }
  });

  test('element.click() works', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const link = await page.find('a[href="/add_remove_elements/"]');
      await link.click();

      // After click, we should have navigated
      await new Promise(resolve => setTimeout(resolve, 1000));

      const heading = await page.find('h3');
      assert.match(heading.info.text, /Add\/Remove Elements/i, 'Should have navigated to new page');
    } finally {
      await b.close();
    }
  });

  test('element.type() enters text', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/inputs');
      const input = await page.find('input');
      await input.type('12345');

      // Verify the value was entered
      const value = await page.evaluate(`
        return document.querySelector('input').value;
      `);
      assert.strictEqual(value, '12345', 'Input should have typed value');
    } finally {
      await b.close();
    }
  });

  test('element.text() returns element text', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const heading = await page.find('h1.heading');
      const text = await heading.text();
      assert.match(text, /Welcome to the-internet/i, 'Should return heading text');
    } finally {
      await b.close();
    }
  });
});
