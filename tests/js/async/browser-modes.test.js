/**
 * JS Library Tests: Browser Modes
 * Tests headless, visible, and default launch options
 */

const { test, describe, before, after } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../../clients/javascript/dist');
const { createTestServer } = require('../../helpers/test-server');

let server, baseURL;

before(async () => {
  ({ server, baseURL } = await createTestServer());
});

after(() => {
  if (server) server.close();
});

describe('JS Browser Modes', () => {
  test('headless mode works', async () => {
    const bro = await browser.start({ headless: true });
    try {
      const vibe = await bro.page();
      await vibe.go(baseURL);
      const screenshot = await vibe.screenshot();
      assert.ok(screenshot.length > 1000, 'Should capture screenshot in headless mode');
    } finally {
      await bro.stop();
    }
  });

  test('headed mode works', async () => {
    // Skip in CI environments where display is not available
    if (process.env.CI || process.env.GITHUB_ACTIONS) {
      console.log('  (skipped: no display in CI)');
      return;
    }

    const bro = await browser.start({ headless: false });
    try {
      const vibe = await bro.page();
      await vibe.go(baseURL);
      const screenshot = await vibe.screenshot();
      assert.ok(screenshot.length > 1000, 'Should capture screenshot in headed mode');
    } finally {
      await bro.stop();
    }
  });

  test('default is visible (not headless)', async () => {
    // Skip in CI environments where display is not available
    if (process.env.CI || process.env.GITHUB_ACTIONS) {
      console.log('  (skipped: no display in CI)');
      return;
    }

    // browser.start() without options should default to visible
    const bro = await browser.start();
    try {
      const vibe = await bro.page();
      await vibe.go(baseURL);
      const title = await vibe.evaluate('document.title');
      assert.match(title, /The Internet/i, 'Should work with default options');
    } finally {
      await bro.stop();
    }
  });
});
