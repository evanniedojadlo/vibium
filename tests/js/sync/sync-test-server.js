/**
 * HTTP server for sync API tests.
 * Runs in a separate process because the sync API blocks the main thread
 * with Atomics.wait(), which would prevent an in-process server from responding.
 *
 * Usage: node sync-test-server.js
 * Prints the base URL to stdout, then serves until killed.
 */
const http = require('http');

const HOME_HTML = `<html><head><title>Test App</title></head><body>
  <h1 class="heading">Welcome to test-app</h1>
  <a href="/subpage">Go to subpage</a>
  <a href="/inputs">Inputs</a>
  <a href="/form">Form</a>
  <p id="info">Some info text</p>
</body></html>`;

const SUBPAGE_HTML = `<html><head><title>Subpage</title></head><body>
  <h3>Subpage Title</h3>
  <a href="/">Back home</a>
</body></html>`;

const INPUTS_HTML = `<html><head><title>Inputs</title></head><body>
  <input type="text" id="text-input" />
  <input type="number" id="num-input" />
  <textarea id="textarea"></textarea>
</body></html>`;

const FORM_HTML = `<html><head><title>Form</title></head><body>
  <form>
    <label for="name">Name</label>
    <input type="text" id="name" name="name" />

    <label for="email">Email</label>
    <input type="email" id="email" name="email" />

    <label for="agree"><input type="checkbox" id="agree" name="agree" /> I agree</label>

    <select id="color" name="color">
      <option value="red">Red</option>
      <option value="green">Green</option>
      <option value="blue">Blue</option>
    </select>

    <button type="submit">Submit</button>
  </form>
</body></html>`;

const LINKS_HTML = `<html><head><title>Links</title></head><body>
  <ul>
    <li><a href="/subpage" class="link">Link 1</a></li>
    <li><a href="/subpage" class="link">Link 2</a></li>
    <li><a href="/subpage" class="link">Link 3</a></li>
    <li><a href="/subpage" class="link special">Link 4</a></li>
  </ul>
  <div id="nested">
    <span class="inner">Nested span</span>
    <span class="inner">Another span</span>
  </div>
</body></html>`;

const EVAL_HTML = `<html><head><title>Eval</title></head><body>
  <div id="result"></div>
  <script>window.testVal = 42;</script>
</body></html>`;

const DIALOG_HTML = `<html><head><title>Dialog</title></head><body>
  <button id="alert-btn" onclick="alert('hello')">Alert</button>
  <button id="confirm-btn" onclick="document.getElementById('result').textContent = confirm('sure?')">Confirm</button>
  <div id="result"></div>
</body></html>`;

const CLOCK_HTML = `<html><head><title>Clock</title></head><body>
  <div id="time"></div>
</body></html>`;

const PROMPT_HTML = `<html><head><title>Prompt</title></head><body>
  <button id="prompt-btn" onclick="document.getElementById('result').textContent = prompt('Enter name:')">Prompt</button>
  <button id="confirm-btn" onclick="document.getElementById('result').textContent = confirm('sure?')">Confirm</button>
  <button id="alert-btn" onclick="alert('hello')">Alert</button>
  <div id="result"></div>
</body></html>`;

const FETCH_HTML = `<html><head><title>Fetch</title></head><body>
  <div id="result"></div>
  <script>
    async function doFetch() {
      const res = await fetch('/api/data');
      const json = await res.json();
      document.getElementById('result').textContent = JSON.stringify(json);
    }
  </script>
</body></html>`;

const routes = {
  '/': HOME_HTML,
  '/subpage': SUBPAGE_HTML,
  '/inputs': INPUTS_HTML,
  '/form': FORM_HTML,
  '/links': LINKS_HTML,
  '/eval': EVAL_HTML,
  '/dialog': DIALOG_HTML,
  '/clock': CLOCK_HTML,
  '/prompt': PROMPT_HTML,
  '/fetch': FETCH_HTML,
};

const server = http.createServer((req, res) => {
  if (req.url === '/api/data') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ message: 'real data', count: 42 }));
    return;
  }
  res.writeHead(200, { 'Content-Type': 'text/html' });
  res.end(routes[req.url] || HOME_HTML);
});

server.listen(0, '127.0.0.1', () => {
  const { port } = server.address();
  // Print URL to stdout so parent process can read it
  process.stdout.write(`http://127.0.0.1:${port}\n`);
});
