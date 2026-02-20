/**
 * JS Sync Tests: Downloads — onDownload
 * Tests page.onDownload() and removeAllListeners('download') for sync API.
 */

const { test, describe, before, after } = require('node:test');
const assert = require('node:assert');
const { fork } = require('child_process');
const path = require('path');

const { browser } = require('../../../clients/javascript/dist/sync');

// --- Server child process ---

let serverProcess;
let baseURL;

before(async () => {
  serverProcess = fork(path.join(__dirname, 'sync-test-server.js'), [], { silent: true });

  baseURL = await new Promise((resolve, reject) => {
    let data = '';
    serverProcess.stdout.on('data', (chunk) => {
      data += chunk.toString();
      const line = data.trim().split('\n')[0];
      if (line.startsWith('http://')) resolve(line);
    });
    serverProcess.on('error', reject);
    setTimeout(() => reject(new Error('Server startup timeout')), 5000);
  });
});

after(() => {
  if (serverProcess) serverProcess.kill();
});

// --- Tests ---

describe('Sync API: onDownload', () => {
  let bro;
  before(() => { bro = browser.launch({ headless: true }); });
  after(() => { bro.close(); });

  test('onDownload fires when download link clicked', () => {
    const vibe = bro.newPage();
    vibe.go(`${baseURL}/download`);

    const downloads = [];
    vibe.onDownload((dl) => {
      downloads.push(dl);
    });

    vibe.find('#download-link').click();
    vibe.wait(1000);

    assert.ok(downloads.length >= 1, `Expected at least 1 download, got ${downloads.length}`);
  });

  test('download has url and suggestedFilename', () => {
    const vibe = bro.newPage();
    vibe.go(`${baseURL}/download`);

    const downloads = [];
    vibe.onDownload((dl) => {
      downloads.push(dl);
    });

    vibe.find('#download-link').click();
    vibe.wait(1000);

    assert.ok(downloads.length >= 1);
    assert.ok(downloads[0].url.includes('/download-file'), `URL should contain /download-file, got: ${downloads[0].url}`);
    assert.strictEqual(downloads[0].suggestedFilename, 'test.txt');
  });

  test('removeAllListeners("download") stops onDownload callbacks', () => {
    const vibe = bro.newPage();
    vibe.go(`${baseURL}/download`);

    const downloads = [];
    vibe.onDownload((dl) => {
      downloads.push(dl);
    });

    vibe.removeAllListeners('download');

    vibe.find('#download-link').click();
    vibe.wait(1000);

    assert.strictEqual(downloads.length, 0, 'Should not capture downloads after removeAllListeners');
  });
});
