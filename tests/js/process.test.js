/**
 * JS Library Tests: Process Management
 * Tests that browser processes are cleaned up properly
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');
const { execSync } = require('node:child_process');

const { browser } = require('../../clients/javascript/dist');
const { browser: browserSync } = require('../../clients/javascript/dist/sync');

/**
 * Get PIDs of Chrome for Testing processes spawned by clicker
 * Returns a Set of PIDs
 */
function getClickerChromePids() {
  try {
    const platform = process.platform;
    let cmd;

    if (platform === 'darwin') {
      cmd = "pgrep -f 'Chrome for Testing.*--remote-debugging-port' 2>/dev/null || true";
    } else if (platform === 'linux') {
      cmd = "pgrep -f 'chrome.*--remote-debugging-port' 2>/dev/null || true";
    } else {
      return new Set();
    }

    const result = execSync(cmd, { encoding: 'utf-8', stdio: ['pipe', 'pipe', 'pipe'] });
    const pids = result.trim().split('\n').filter(Boolean).map(Number);
    return new Set(pids);
  } catch {
    return new Set();
  }
}

/**
 * Get new PIDs that appeared between two sets
 */
function getNewPids(before, after) {
  return [...after].filter(pid => !before.has(pid));
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

describe('JS Process Cleanup', () => {
  test('async API cleans up Chrome on close()', async () => {
    const pidsBefore = getClickerChromePids();

    const b = await browser.launch({ headless: true });
    const page = await b.page();
    await page.go('https://the-internet.herokuapp.com/');
    await b.close();

    await sleep(2000);

    const pidsAfter = getClickerChromePids();
    const newPids = getNewPids(pidsBefore, pidsAfter);

    assert.strictEqual(
      newPids.length,
      0,
      `Chrome processes should be cleaned up. New PIDs remaining: ${newPids.join(', ')}`
    );
  });

  test('sync API cleans up Chrome on close()', async () => {
    const pidsBefore = getClickerChromePids();

    const bro = browserSync.launch({ headless: true });
    const page = bro.page();
    page.go('https://the-internet.herokuapp.com/');
    bro.close();

    await sleep(2000);

    const pidsAfter = getClickerChromePids();
    const newPids = getNewPids(pidsBefore, pidsAfter);

    assert.strictEqual(
      newPids.length,
      0,
      `Chrome processes should be cleaned up. New PIDs remaining: ${newPids.join(', ')}`
    );
  });

  test('multiple sequential sessions clean up properly', async () => {
    const pidsBefore = getClickerChromePids();

    // Run 3 sessions sequentially
    for (let i = 0; i < 3; i++) {
      const b = await browser.launch({ headless: true });
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');
      await b.close();
    }

    await sleep(2000);

    const pidsAfter = getClickerChromePids();
    const newPids = getNewPids(pidsBefore, pidsAfter);

    assert.strictEqual(
      newPids.length,
      0,
      `All Chrome processes should be cleaned up after 3 sessions. New PIDs remaining: ${newPids.join(', ')}`
    );
  });
});
