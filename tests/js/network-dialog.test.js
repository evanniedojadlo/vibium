/**
 * JS Library Tests: Network Interception & Dialogs
 * Tests page.route, route.fulfill/continue/abort, page.onRequest/onResponse,
 * page.waitForRequest/waitForResponse, page.onDialog, dialog.accept/dismiss.
 *
 * Uses a local HTTP server — no external network dependencies.
 */

const { test, describe, before, after } = require('node:test');
const assert = require('node:assert');
const http = require('http');

const { browser } = require('../../clients/javascript/dist');

// --- Local test server ---

let server;
let baseURL;

before(async () => {
  server = http.createServer((req, res) => {
    if (req.url === '/json') {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ name: 'vibium', version: 1 }));
    } else if (req.url === '/text') {
      res.writeHead(200, { 'Content-Type': 'text/plain' });
      res.end('hello world');
    } else {
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end('<html><head><title>Test Page</title></head><body>Test Content</body></html>');
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

// --- Network Interception ---

describe('Network Interception: page.route', () => {
  test('route.abort() blocks a request', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      // Block all .png requests
      await page.route('**/*.png', (route) => {
        route.abort();
      });

      await page.go(baseURL);

      // Verify the page loaded (route didn't break navigation)
      const title = await page.title();
      assert.strictEqual(title, 'Test Page');
    } finally {
      await b.close();
    }
  });

  test('route.fulfill() returns a mock response', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      await page.route('**/json', (route) => {
        route.fulfill({
          status: 200,
          body: JSON.stringify({ mocked: true }),
          contentType: 'application/json',
        });
      });

      const result = await page.eval(`
        fetch('${baseURL}/json')
          .then(r => r.json())
      `);

      assert.deepStrictEqual(result, { mocked: true });
    } finally {
      await b.close();
    }
  });

  test('route.fulfill() with custom headers', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      await page.route('**/text', (route) => {
        route.fulfill({
          status: 201,
          headers: { 'X-Custom': 'test-value', 'Content-Type': 'text/plain' },
          body: 'custom body',
        });
      });

      const result = await page.eval(`
        fetch('${baseURL}/text')
          .then(r => r.text().then(body => ({ status: r.status, body, custom: r.headers.get('X-Custom') })))
      `);

      assert.strictEqual(result.status, 201);
      assert.strictEqual(result.body, 'custom body');
      assert.strictEqual(result.custom, 'test-value');
    } finally {
      await b.close();
    }
  });

  test('route.continue() lets request through', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      let intercepted = false;
      await page.route('**', (route) => {
        intercepted = true;
        route.continue();
      });

      // Fetch triggers the intercept
      await page.eval(`fetch('${baseURL}/text')`);
      await page.wait(200);

      assert.ok(intercepted, 'Route handler should have been called');
    } finally {
      await b.close();
    }
  });

  test('unroute() removes a route', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      let callCount = 0;
      await page.route('**/text', (route) => {
        callCount++;
        route.continue();
      });

      // First fetch — should be intercepted
      await page.eval(`fetch('${baseURL}/text')`);
      await page.wait(200);
      assert.ok(callCount > 0, 'Route handler should have been called');

      const countBefore = callCount;
      await page.unroute('**/text');

      // Second fetch — should NOT be intercepted
      await page.eval(`fetch('${baseURL}/text')`);
      await page.wait(200);
      assert.strictEqual(callCount, countBefore, 'Route should not fire after unroute');
    } finally {
      await b.close();
    }
  });
});

// --- Network Events & Waiters ---

describe('Network Events: onRequest/onResponse', () => {
  test('onRequest() fires for page navigation', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      const urls = [];
      page.onRequest((req) => {
        urls.push(req.url());
      });

      await page.go(baseURL);
      await page.wait(200);

      assert.ok(urls.length > 0, 'Should have captured at least one request');
      assert.ok(
        urls.some(u => u.includes('127.0.0.1')),
        `Should have a request to local server, got: ${urls.join(', ')}`
      );
    } finally {
      await b.close();
    }
  });

  test('onResponse() fires for page navigation', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      const statuses = [];
      page.onResponse((resp) => {
        statuses.push(resp.status());
      });

      await page.go(baseURL);
      await page.wait(200);

      assert.ok(statuses.length > 0, 'Should have captured at least one response');
      assert.ok(statuses.includes(200), `Should have a 200 response, got: ${statuses.join(', ')}`);
    } finally {
      await b.close();
    }
  });

  test('request.method() and request.headers() work', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      let capturedMethod = '';
      let capturedHeaders = {};
      page.onRequest((req) => {
        if (req.url().includes('127.0.0.1') && !capturedMethod) {
          capturedMethod = req.method();
          capturedHeaders = req.headers();
        }
      });

      await page.go(baseURL);
      await page.wait(200);

      assert.strictEqual(capturedMethod, 'GET');
      assert.ok(typeof capturedHeaders === 'object');
    } finally {
      await b.close();
    }
  });

  test('response.url() and response.status() work', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      const responsePromise = page.waitForResponse('**/json');
      await page.eval(`fetch('${baseURL}/json')`);
      const resp = await responsePromise;

      assert.ok(resp.url().includes('/json'));
      assert.strictEqual(resp.status(), 200);
      assert.ok(typeof resp.headers() === 'object');
    } finally {
      await b.close();
    }
  });
});

describe('Network Waiters: waitForRequest/waitForResponse', () => {
  test('waitForResponse() resolves on matching response', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      const responsePromise = page.waitForResponse('**/json');
      await page.eval(`fetch('${baseURL}/json')`);

      const resp = await responsePromise;
      assert.ok(resp.url().includes('/json'), `Response URL should include /json, got: ${resp.url()}`);
      assert.strictEqual(resp.status(), 200);
    } finally {
      await b.close();
    }
  });

  test('waitForRequest() resolves on matching request', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      const requestPromise = page.waitForRequest('**/text');
      await page.eval(`fetch('${baseURL}/text')`);

      const req = await requestPromise;
      assert.ok(req.url().includes('/text'), `Request URL should include /text, got: ${req.url()}`);
      assert.strictEqual(req.method(), 'GET');
    } finally {
      await b.close();
    }
  });
});

// --- Response Body ---

describe('Response Body: response.body() and response.json()', () => {
  test('response.body() returns text content via onResponse', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      let captured = null;
      page.onResponse((resp) => {
        if (resp.url().includes('/text')) {
          captured = resp;
        }
      });

      await page.eval(`fetch('${baseURL}/text')`);
      await page.wait(500);

      assert.ok(captured, 'Should have captured the /text response');
      const body = await captured.body();
      assert.strictEqual(body, 'hello world');
    } finally {
      await b.close();
    }
  });

  test('response.json() parses JSON content via onResponse', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      let captured = null;
      page.onResponse((resp) => {
        if (resp.url().includes('/json')) {
          captured = resp;
        }
      });

      await page.eval(`fetch('${baseURL}/json')`);
      await page.wait(500);

      assert.ok(captured, 'Should have captured the /json response');
      const data = await captured.json();
      assert.deepStrictEqual(data, { name: 'vibium', version: 1 });
    } finally {
      await b.close();
    }
  });

  test('response.body() works with waitForResponse', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      const responsePromise = page.waitForResponse('**/text');
      await page.eval(`fetch('${baseURL}/text')`);
      const resp = await responsePromise;

      const body = await resp.body();
      assert.strictEqual(body, 'hello world');
    } finally {
      await b.close();
    }
  });

  test('response.json() works with waitForResponse', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      const responsePromise = page.waitForResponse('**/json');
      await page.eval(`fetch('${baseURL}/json')`);
      const resp = await responsePromise;

      const data = await resp.json();
      assert.deepStrictEqual(data, { name: 'vibium', version: 1 });
    } finally {
      await b.close();
    }
  });
});

// --- Dialogs ---

describe('Dialogs: page.onDialog', () => {
  test('onDialog() handles alert', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      let dialogMessage = '';
      let dialogType = '';
      page.onDialog((dialog) => {
        dialogMessage = dialog.message();
        dialogType = dialog.type();
        dialog.accept();
      });

      await page.eval('alert("Hello from test")');

      assert.strictEqual(dialogMessage, 'Hello from test');
      assert.strictEqual(dialogType, 'alert');
    } finally {
      await b.close();
    }
  });

  test('onDialog() handles confirm with accept', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      page.onDialog((dialog) => {
        dialog.accept();
      });

      const result = await page.eval('confirm("Are you sure?")');
      assert.strictEqual(result, true);
    } finally {
      await b.close();
    }
  });

  test('onDialog() handles confirm with dismiss', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      page.onDialog((dialog) => {
        dialog.dismiss();
      });

      const result = await page.eval('confirm("Are you sure?")');
      assert.strictEqual(result, false);
    } finally {
      await b.close();
    }
  });

  test('onDialog() handles prompt with text', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      page.onDialog((dialog) => {
        assert.strictEqual(dialog.type(), 'prompt');
        dialog.accept('my answer');
      });

      const result = await page.eval('prompt("Enter name:")');
      assert.strictEqual(result, 'my answer');
    } finally {
      await b.close();
    }
  });

  test('dialogs are auto-dismissed when no handler registered', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go(baseURL);

      // No onDialog handler — should auto-dismiss
      const result = await page.eval('confirm("Auto dismiss?")');
      assert.strictEqual(result, false);
    } finally {
      await b.close();
    }
  });
});

// --- WebSocket Stubs ---

describe('Stubs: WebSocket methods', () => {
  test('routeWebSocket() throws not implemented', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      assert.throws(
        () => page.routeWebSocket('**', () => {}),
        /Not implemented/
      );
    } finally {
      await b.close();
    }
  });

  test('onWebSocket() throws not implemented', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      assert.throws(
        () => page.onWebSocket(() => {}),
        /Not implemented/
      );
    } finally {
      await b.close();
    }
  });
});

// --- Checkpoint ---

describe('Network & Dialog Checkpoint', () => {
  test('route.continue, onResponse, and onDialog work together', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();

      // Set up route that intercepts and continues
      let intercepted = false;
      await page.route('**', (route) => {
        intercepted = true;
        route.continue();
      });

      // Track responses
      const responseUrls = [];
      page.onResponse((resp) => {
        responseUrls.push(resp.url());
      });

      await page.go(baseURL);
      await page.wait(200);

      assert.ok(intercepted, 'Route should have intercepted');
      assert.ok(responseUrls.length > 0, 'Should have captured responses');

      // Set up dialog handler and trigger a dialog
      let dialogHandled = false;
      page.onDialog((dialog) => {
        dialogHandled = true;
        dialog.accept();
      });

      await page.eval('alert("checkpoint")');
      assert.ok(dialogHandled);
    } finally {
      await b.close();
    }
  });
});
