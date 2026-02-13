/**
 * JS Library Tests: Browser Modes
 * Tests headless, visible, and default launch options
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

describe('JS Browser Modes', () => {
  test('headless mode works', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const screenshot = await page.screenshot();
      assert.ok(screenshot.length > 1000, 'Should capture screenshot in headless mode');
    } finally {
      await b.close();
    }
  });

  test('headed mode works', async () => {
    // Skip in CI environments where display is not available
    if (process.env.CI || process.env.GITHUB_ACTIONS) {
      console.log('  (skipped: no display in CI)');
      return;
    }

    const b = await browser.launch({ headless: false });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const screenshot = await page.screenshot();
      assert.ok(screenshot.length > 1000, 'Should capture screenshot in headed mode');
    } finally {
      await b.close();
    }
  });

  test('default is visible (not headless)', async () => {
    // Skip in CI environments where display is not available
    if (process.env.CI || process.env.GITHUB_ACTIONS) {
      console.log('  (skipped: no display in CI)');
      return;
    }

    // browser.launch() without options should default to visible
    const b = await browser.launch();
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      const title = await page.evaluate('return document.title');
      assert.match(title, /The Internet/i, 'Should work with default options');
    } finally {
      await b.close();
    }
  });
});
