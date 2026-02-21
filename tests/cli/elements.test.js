/**
 * CLI Tests: Element Finding, Click, and Type
 * Tests the vibium binary directly
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const { VIBIUM } = require('../helpers');

describe('CLI: Elements', () => {
  test('find command locates element and returns @ref', () => {
    const result = execSync(`${VIBIUM} find https://example.com "a"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /@e1/, 'Should return @e1 ref');
    assert.match(result, /\[a\]/, 'Should show [a] tag label');
    // Link text may be "More information..." or "Learn more" depending on page version
    assert.match(result, /(More information|Learn more)/i, 'Should show link text');
  });

  test('click command navigates via link', () => {
    const result = execSync(`${VIBIUM} click https://example.com "a"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /iana\.org/i, 'Should navigate to IANA after clicking link');
  });

  test('type command enters text into input', () => {
    const result = execSync(
      `${VIBIUM} type https://the-internet.herokuapp.com/inputs "input" "12345"`,
      {
        encoding: 'utf-8',
        timeout: 30000,
      }
    );
    assert.match(result, /12345/, 'Should show typed text in result');
  });
});
