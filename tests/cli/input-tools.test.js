/**
 * CLI Tests: Input Tools
 * Tests hover command in oneshot mode
 * Note: scroll, keys, select require daemon mode and are tested via MCP
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const path = require('node:path');

const CLICKER = path.join(__dirname, '../../clicker/bin/clicker');

describe('CLI: Input Tools', () => {
  test('hover command hovers over element', () => {
    const result = execSync(`${CLICKER} hover https://example.com "a"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /Hovered/, 'Should confirm hover');
  });

  test('skill command outputs markdown', () => {
    const result = execSync(`${CLICKER} skill`, {
      encoding: 'utf-8',
      timeout: 5000,
    });
    assert.match(result, /# Vibium Clicker/, 'Should have title');
    assert.match(result, /navigate/, 'Should list navigate');
    assert.match(result, /click/, 'Should list click');
    assert.match(result, /screenshot/, 'Should list screenshot');
    assert.match(result, /tab-new/, 'Should list tab-new');
    assert.match(result, /scroll/, 'Should list scroll');
    assert.match(result, /keys/, 'Should list keys');
  });
});
