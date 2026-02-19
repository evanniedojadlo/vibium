/**
 * JS Library Tests: Tracing
 * Tests context.tracing.start/stop, screenshots, snapshots, chunks, and groups.
 *
 * Uses a local HTTP server — no external network dependencies.
 */

const { test, describe, before, after } = require('node:test');
const assert = require('node:assert');
const http = require('http');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const { browser } = require('../../../clients/javascript/dist');

// --- Local test server ---

let server;
let baseURL;

const HTML_PAGE = `
<html>
<head><title>Tracing Test</title></head>
<body>
  <h1 id="heading">Hello Tracing</h1>
  <button id="btn" onclick="document.getElementById('heading').textContent='Clicked'">Click Me</button>
  <a href="/page2">Go to page 2</a>
</body>
</html>
`;

const HTML_PAGE2 = `
<html>
<head><title>Page 2</title></head>
<body>
  <h1 id="heading">Page Two</h1>
</body>
</html>
`;

before(async () => {
  server = http.createServer((req, res) => {
    if (req.url === '/page2') {
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end(HTML_PAGE2);
    } else {
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end(HTML_PAGE);
    }
  });

  await new Promise((resolve) => {
    server.listen(0, '127.0.0.1', () => {
      const { port } = server.address();
      baseURL = `http://127.0.0.1:${port}`;
      resolve();
    });
  });
});

after(() => {
  if (server) server.close();
});

// --- Helper: unzip and inspect trace ---

function unzipTrace(zipBuffer) {
  // Use Node.js built-in zlib + manual zip parsing, or shell out to unzip
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'vibium-trace-test-'));
  const zipPath = path.join(tmpDir, 'trace.zip');
  fs.writeFileSync(zipPath, zipBuffer);
  execSync(`unzip -o "${zipPath}" -d "${tmpDir}/extracted"`, { stdio: 'pipe' });
  return { tmpDir, extractedDir: path.join(tmpDir, 'extracted') };
}

function cleanupDir(dir) {
  fs.rmSync(dir, { recursive: true, force: true });
}

function readTraceEvents(extractedDir) {
  const files = fs.readdirSync(extractedDir).filter(f => f.endsWith('.trace'));
  const events = [];
  for (const file of files) {
    const content = fs.readFileSync(path.join(extractedDir, file), 'utf-8');
    for (const line of content.split('\n')) {
      if (line.trim()) {
        events.push(JSON.parse(line));
      }
    }
  }
  return events;
}

function readNetworkEvents(extractedDir) {
  const files = fs.readdirSync(extractedDir).filter(f => f.endsWith('.network'));
  const events = [];
  for (const file of files) {
    const content = fs.readFileSync(path.join(extractedDir, file), 'utf-8');
    for (const line of content.split('\n')) {
      if (line.trim()) {
        events.push(JSON.parse(line));
      }
    }
  }
  return events;
}

// --- Tests ---

describe('Tracing: basic start/stop', () => {
  test('start and stop produces valid trace zip', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ name: 'basic-test' });
      await vibe.go(baseURL);
      await vibe.find('#btn').click();
      await vibe.wait(200);
      const zipBuffer = await ctx.tracing.stop();

      assert.ok(Buffer.isBuffer(zipBuffer), 'stop() should return a Buffer');
      assert.ok(zipBuffer.length > 0, 'zip should not be empty');

      // Verify zip structure
      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      const files = fs.readdirSync(extractedDir);
      assert.ok(files.some(f => f.endsWith('.trace')), 'zip should contain a .trace file');
      assert.ok(files.some(f => f.endsWith('.network')), 'zip should contain a .network file');

      // Verify first event is context-options
      const events = readTraceEvents(extractedDir);
      assert.ok(events.length > 0, 'should have trace events');
      assert.strictEqual(events[0].type, 'context-options');
      assert.strictEqual(events[0].browserName, 'chromium');

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });

  test('stop with path writes trace to file', async () => {
    const bro = await browser.launch({ headless: true });
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'vibium-trace-path-'));
    const tracePath = path.join(tmpDir, 'my-trace.zip');
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start();
      await vibe.go(baseURL);
      const zipBuffer = await ctx.tracing.stop({ path: tracePath });

      assert.ok(fs.existsSync(tracePath), 'trace file should exist at the given path');
      const fileSize = fs.statSync(tracePath).size;
      assert.ok(fileSize > 0, 'trace file should not be empty');

      await ctx.close();
    } finally {
      await bro.close();
      cleanupDir(tmpDir);
    }
  });
});

describe('Tracing: screenshots', () => {
  test('screenshots option captures PNG resources', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ screenshots: true });
      await vibe.go(baseURL);
      // Wait for some screenshots to be captured
      await vibe.wait(500);
      await vibe.find('#btn').click();
      await vibe.wait(500);
      const zipBuffer = await ctx.tracing.stop();

      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      // Check for PNG resources
      const resourcesDir = path.join(extractedDir, 'resources');
      assert.ok(fs.existsSync(resourcesDir), 'resources directory should exist');

      const resources = fs.readdirSync(resourcesDir);
      const pngs = resources.filter(f => f.endsWith('.png'));
      assert.ok(pngs.length > 0, `Should have PNG screenshots, got: ${resources.join(', ')}`);

      // Check for screencast-frame events in trace
      const events = readTraceEvents(extractedDir);
      const frames = events.filter(e => e.type === 'screencast-frame');
      assert.ok(frames.length > 0, 'should have screencast-frame events');
      assert.ok(frames[0].sha1, 'screencast-frame should have sha1');
      assert.ok(frames[0].width > 0, 'screencast-frame should have width');
      assert.ok(frames[0].height > 0, 'screencast-frame should have height');

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });
});

describe('Tracing: snapshots', () => {
  test('snapshots option captures HTML resources', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ snapshots: true });
      await vibe.go(baseURL);
      await vibe.wait(200);
      const zipBuffer = await ctx.tracing.stop();

      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      // Check for HTML resources
      const resourcesDir = path.join(extractedDir, 'resources');
      if (fs.existsSync(resourcesDir)) {
        const resources = fs.readdirSync(resourcesDir);
        const htmlFiles = resources.filter(f => f.endsWith('.html'));
        assert.ok(htmlFiles.length > 0, `Should have HTML snapshots, got: ${resources.join(', ')}`);
      }

      // Check for frame-snapshot events
      const events = readTraceEvents(extractedDir);
      const snapshots = events.filter(e => e.type === 'frame-snapshot');
      assert.ok(snapshots.length > 0, 'should have frame-snapshot events');

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });
});

describe('Tracing: chunks', () => {
  test('startChunk/stopChunk produces separate trace zips', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir1, tmpDir2;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ name: 'chunk-test' });
      await vibe.go(baseURL);
      await vibe.wait(200);

      // Stop first chunk
      const zip1 = await ctx.tracing.stopChunk();
      assert.ok(Buffer.isBuffer(zip1), 'first chunk should return a Buffer');

      // Start second chunk
      await ctx.tracing.startChunk({ name: 'chunk-2' });
      await vibe.go(baseURL + '/page2');
      await vibe.wait(200);

      // Stop second chunk
      const zip2 = await ctx.tracing.stopChunk();
      assert.ok(Buffer.isBuffer(zip2), 'second chunk should return a Buffer');

      // Verify both zips are valid
      const { tmpDir: td1, extractedDir: ed1 } = unzipTrace(zip1);
      tmpDir1 = td1;
      const events1 = readTraceEvents(ed1);
      assert.ok(events1.length > 0, 'first chunk should have events');

      const { tmpDir: td2, extractedDir: ed2 } = unzipTrace(zip2);
      tmpDir2 = td2;
      const events2 = readTraceEvents(ed2);
      assert.ok(events2.length > 0, 'second chunk should have events');

      // Stop tracing
      await ctx.tracing.stop();
      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir1) cleanupDir(tmpDir1);
      if (tmpDir2) cleanupDir(tmpDir2);
    }
  });
});

describe('Tracing: groups', () => {
  test('startGroup/stopGroup adds group markers to trace', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ name: 'group-test' });
      await vibe.go(baseURL);

      await ctx.tracing.startGroup('login flow');
      await vibe.find('#btn').click();
      await vibe.wait(200);
      await ctx.tracing.stopGroup();

      const zipBuffer = await ctx.tracing.stop();

      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      const events = readTraceEvents(extractedDir);

      // Look for before/after events from groups
      const beforeEvents = events.filter(e => e.type === 'before' && e.apiName === 'login flow');
      assert.ok(beforeEvents.length > 0, 'should have a before event for the group');

      const afterEvents = events.filter(e => e.type === 'after');
      assert.ok(afterEvents.length > 0, 'should have an after event for group end');

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });
});

describe('Tracing: network events', () => {
  test('trace captures network events from navigation', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ name: 'network-test' });
      await vibe.go(baseURL);
      await vibe.wait(500);
      const zipBuffer = await ctx.tracing.stop();

      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      const networkEvents = readNetworkEvents(extractedDir);
      assert.ok(networkEvents.length > 0, `should have network events, got ${networkEvents.length}`);

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });
});

describe('Tracing: zip structure', () => {
  test('trace zip has correct Playwright-compatible structure', async () => {
    const bro = await browser.launch({ headless: true });
    let tmpDir;
    try {
      const ctx = await bro.newContext();
      const vibe = await ctx.newPage();

      await ctx.tracing.start({ screenshots: true, snapshots: true });
      await vibe.go(baseURL);
      await vibe.wait(500);
      const zipBuffer = await ctx.tracing.stop();

      const { tmpDir: td, extractedDir } = unzipTrace(zipBuffer);
      tmpDir = td;

      const files = fs.readdirSync(extractedDir);

      // Must have trace file matching pattern <n>-trace.trace
      const traceFiles = files.filter(f => /^\d+-trace\.trace$/.test(f));
      assert.ok(traceFiles.length > 0, 'should have numbered trace file');

      // Must have network file matching pattern <n>-trace.network
      const networkFiles = files.filter(f => /^\d+-trace\.network$/.test(f));
      assert.ok(networkFiles.length > 0, 'should have numbered network file');

      // Parse trace and verify event types
      const events = readTraceEvents(extractedDir);
      const types = [...new Set(events.map(e => e.type))];
      assert.ok(types.includes('context-options'), `should include context-options, got: ${types.join(', ')}`);

      await ctx.close();
    } finally {
      await bro.close();
      if (tmpDir) cleanupDir(tmpDir);
    }
  });
});
