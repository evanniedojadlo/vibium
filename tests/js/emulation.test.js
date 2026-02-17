/**
 * JS Library Tests: Emulation
 * Tests page.setViewport, page.viewport, page.emulateMedia,
 * page.setContent, page.setGeolocation
 */

const { test, describe, after } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

describe('JS Emulation', () => {
  let b;

  test('setup', async () => {
    b = await browser.launch({ headless: true });
  });

  after(async () => {
    if (b) await b.close().catch(() => {});
  });

  // --- setViewport / viewport ---

  test('viewport() returns current size', async () => {
    const page = await b.page();
    const size = await page.viewport();
    assert.ok(typeof size.width === 'number' && size.width > 0, `width should be > 0, got ${size.width}`);
    assert.ok(typeof size.height === 'number' && size.height > 0, `height should be > 0, got ${size.height}`);
  });

  test('setViewport() changes viewport size', async () => {
    const page = await b.page();
    await page.setViewport({ width: 800, height: 600 });
    const size = await page.viewport();
    assert.strictEqual(size.width, 800, `width should be 800, got ${size.width}`);
    assert.strictEqual(size.height, 600, `height should be 600, got ${size.height}`);
  });

  // --- setContent ---

  test('setContent() replaces page HTML', async () => {
    const page = await b.page();
    await page.setContent('<html><body><h1>Hello Vibium</h1></body></html>');
    const el = await page.find('h1');
    const text = await el.text();
    assert.strictEqual(text, 'Hello Vibium', `h1 text should be "Hello Vibium", got "${text}"`);
  });

  test('setContent() with full document including title', async () => {
    const page = await b.page();
    await page.setContent('<!DOCTYPE html><html><head><title>Custom Title</title></head><body><p>content</p></body></html>');
    const title = await page.title();
    assert.strictEqual(title, 'Custom Title', `title should be "Custom Title", got "${title}"`);
  });

  // --- emulateMedia ---

  test('emulateMedia({ colorScheme: "dark" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ colorScheme: 'dark' });
    const matches = await page.eval('window.matchMedia("(prefers-color-scheme: dark)").matches');
    assert.strictEqual(matches, true, 'prefers-color-scheme: dark should match');
  });

  test('emulateMedia({ colorScheme: "light" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ colorScheme: 'light' });
    const matches = await page.eval('window.matchMedia("(prefers-color-scheme: light)").matches');
    assert.strictEqual(matches, true, 'prefers-color-scheme: light should match');
    const darkMatches = await page.eval('window.matchMedia("(prefers-color-scheme: dark)").matches');
    assert.strictEqual(darkMatches, false, 'prefers-color-scheme: dark should NOT match');
  });

  test('emulateMedia({ media: "print" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ media: 'print' });
    const matches = await page.eval('window.matchMedia("print").matches');
    assert.strictEqual(matches, true, 'print media should match');
  });

  test('emulateMedia({ reducedMotion: "reduce" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ reducedMotion: 'reduce' });
    const matches = await page.eval('window.matchMedia("(prefers-reduced-motion: reduce)").matches');
    assert.strictEqual(matches, true, 'prefers-reduced-motion: reduce should match');
  });

  test('emulateMedia({ forcedColors: "active" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ forcedColors: 'active' });
    const matches = await page.eval('window.matchMedia("(forced-colors: active)").matches');
    assert.strictEqual(matches, true, 'forced-colors: active should match');
  });

  test('emulateMedia({ contrast: "more" })', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.emulateMedia({ contrast: 'more' });
    const matches = await page.eval('window.matchMedia("(prefers-contrast: more)").matches');
    assert.strictEqual(matches, true, 'prefers-contrast: more should match');
  });

  test('emulateMedia(null) resets overrides', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');

    // Set override
    await page.emulateMedia({ colorScheme: 'dark' });
    let matches = await page.eval('window.matchMedia("(prefers-color-scheme: dark)").matches');
    assert.strictEqual(matches, true, 'dark should match after setting');

    // Reset
    await page.emulateMedia({ colorScheme: null });
    // After reset, the query should use browser default (which may or may not be dark)
    // The key test is that the override was removed — query passthrough to native matchMedia
    const result = await page.eval('typeof window.__vibiumMediaOverrides.colorScheme');
    assert.strictEqual(result, 'undefined', 'colorScheme override should be removed');
  });

  // --- setGeolocation ---

  test('setGeolocation() overrides position', async () => {
    const page = await b.page();
    await page.setContent('<html><body></body></html>');
    await page.setGeolocation({ latitude: 51.5074, longitude: -0.1278 });

    const coords = await page.eval(`
      new Promise((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(
          pos => resolve({ lat: pos.coords.latitude, lng: pos.coords.longitude }),
          err => reject(err),
          { timeout: 5000 }
        );
      })
    `);

    assert.ok(Math.abs(coords.lat - 51.5074) < 0.001, `latitude should be ~51.5074, got ${coords.lat}`);
    assert.ok(Math.abs(coords.lng - (-0.1278)) < 0.001, `longitude should be ~-0.1278, got ${coords.lng}`);
  });
});
