/**
 * CLI Tests: Page Reading Tools
 * Tests text, html, find-all commands in oneshot mode
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');
const { CLICKER } = require('../helpers');

describe('CLI: Page Reading', () => {
  test('text command returns page text', () => {
    const result = execSync(`${CLICKER} text https://example.com`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /Example Domain/, 'Should contain page text');
  });

  test('text command with selector returns element text', () => {
    const result = execSync(`${CLICKER} text https://example.com "h1"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /Example Domain/, 'Should contain h1 text');
  });

  test('html command returns page HTML', () => {
    const result = execSync(`${CLICKER} html https://example.com "h1"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /Example Domain/, 'Should contain HTML');
  });

  test('html command with --outer returns outer HTML', () => {
    const result = execSync(`${CLICKER} html https://example.com "h1" --outer`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /<h1>/, 'Should contain h1 tag');
    assert.match(result, /Example Domain/, 'Should contain text');
  });

  test('find-all command returns multiple elements', () => {
    const result = execSync(`${CLICKER} find-all https://example.com "p"`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /\[0\]/, 'Should contain indexed results');
    assert.match(result, /tag=p/, 'Should contain tag info');
  });

  test('find-all command with --limit', () => {
    const result = execSync(`${CLICKER} find-all https://example.com "p" --limit 1`, {
      encoding: 'utf-8',
      timeout: 30000,
    });
    assert.match(result, /\[0\]/, 'Should contain first result');
    assert.ok(!result.includes('[1]'), 'Should not contain second result');
  });
});
