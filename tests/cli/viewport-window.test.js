/**
 * CLI Tests: Viewport & Window Commands
 * Tests viewport and window management via the daemon.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const { VIBIUM } = require('../helpers');

describe('CLI: Viewport & Window Commands', () => {
  test('viewport returns width and height', () => {
    const result = execSync(`${VIBIUM} viewport`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const vp = JSON.parse(result);
    assert.ok(typeof vp.width === 'number', 'Should have numeric width');
    assert.ok(typeof vp.height === 'number', 'Should have numeric height');
  });

  test('set-viewport then viewport verifies dimensions', () => {
    execSync(`${VIBIUM} set-viewport 800 600`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const result = execSync(`${VIBIUM} viewport`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const vp = JSON.parse(result);
    assert.strictEqual(vp.width, 800);
    assert.strictEqual(vp.height, 600);
  });

  test('window returns state and dimensions', () => {
    const result = execSync(`${VIBIUM} window`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const win = JSON.parse(result);
    assert.ok(typeof win.width === 'number', 'Should have numeric width');
    assert.ok(typeof win.height === 'number', 'Should have numeric height');
  });

  test('set-window then window verifies dimensions', () => {
    execSync(`${VIBIUM} set-window 900 700`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const result = execSync(`${VIBIUM} window`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    const win = JSON.parse(result);
    assert.strictEqual(win.width, 900);
    assert.strictEqual(win.height, 700);
  });
});
