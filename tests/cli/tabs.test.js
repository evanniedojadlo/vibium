/**
 * CLI Tests: Tab Commands
 * Tab commands require daemon mode — tested via MCP server tests.
 * This file verifies error handling in oneshot mode.
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const path = require('node:path');

const CLICKER = path.join(__dirname, '../../clicker/bin/clicker');

describe('CLI: Tab Commands (oneshot errors)', () => {
  test('tabs command fails in oneshot mode', () => {
    assert.throws(() => {
      execSync(`${CLICKER} tabs`, {
        encoding: 'utf-8',
        timeout: 5000,
      });
    }, 'Should error in oneshot mode');
  });

  test('tab-new command fails in oneshot mode', () => {
    assert.throws(() => {
      execSync(`${CLICKER} tab-new`, {
        encoding: 'utf-8',
        timeout: 5000,
      });
    }, 'Should error in oneshot mode');
  });

  test('tab-switch command fails in oneshot mode', () => {
    assert.throws(() => {
      execSync(`${CLICKER} tab-switch 0`, {
        encoding: 'utf-8',
        timeout: 5000,
      });
    }, 'Should error in oneshot mode');
  });

  test('tab-close command fails in oneshot mode', () => {
    assert.throws(() => {
      execSync(`${CLICKER} tab-close`, {
        encoding: 'utf-8',
        timeout: 5000,
      });
    }, 'Should error in oneshot mode');
  });
});
