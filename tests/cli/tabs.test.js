/**
 * CLI Tests: Tab Commands
 * Tests tab management via the daemon.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const { VIBIUM } = require('../helpers');

describe('CLI: Tab Commands', () => {
  test('tabs lists open tabs', () => {
    const result = execSync(`${VIBIUM} tabs`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /\[0\]/, 'Should list tab 0');
  });

  test('tab-new creates a new tab', () => {
    const result = execSync(`${VIBIUM} tab-new`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /created|new tab/i, 'Should confirm new tab created');
  });

  test('tab-switch switches to a tab', () => {
    const result = execSync(`${VIBIUM} tab-switch 0`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /switched|tab 0/i, 'Should confirm tab switch');
  });

  test('tab-close closes a tab', () => {
    const result = execSync(`${VIBIUM} tab-close`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /closed/i, 'Should confirm tab closed');
  });
});
