/**
 * Parses annotated markdown tutorials and creates tests from code blocks.
 *
 * Annotations:
 *   <!-- helpers -->              — next code block defines shared helper functions
 *   <!-- test: async "name" -->   — next code block is an async test
 *   <!-- test: sync "name" -->    — next code block is a sync test
 */

const { test } = require('node:test');
const assert = require('node:assert');
const { readFileSync } = require('fs');
const { resolve } = require('path');

const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;

function extractBlocks(mdPath) {
  const fullPath = resolve(__dirname, '../../..', mdPath);
  const content = readFileSync(fullPath, 'utf8');
  const lines = content.split('\n');
  const blocks = [];

  let pending = null;
  let inCodeBlock = false;
  let isAnnotated = false;
  let codeLines = [];

  for (const line of lines) {
    if (!inCodeBlock) {
      const helpersMatch = line.match(/<!--\s*helpers\s*-->/);
      if (helpersMatch) {
        pending = { type: 'helpers' };
        continue;
      }

      const testMatch = line.match(/<!--\s*test:\s*(async|sync)\s+"([^"]+)"\s*-->/);
      if (testMatch) {
        pending = { type: 'test', mode: testMatch[1], name: testMatch[2] };
        continue;
      }

      if (line.match(/^```javascript\s*$/)) {
        inCodeBlock = true;
        if (pending) {
          isAnnotated = true;
          codeLines = [];
        } else {
          isAnnotated = false;
        }
        continue;
      }
    } else {
      if (line.match(/^```\s*$/)) {
        inCodeBlock = false;
        if (isAnnotated && pending) {
          blocks.push({ ...pending, code: codeLines.join('\n') });
          pending = null;
        }
        continue;
      }
      if (isAnnotated) {
        codeLines.push(line);
      }
    }
  }

  return blocks;
}

function runTutorial(mdPath, { browser, mode }) {
  const blocks = extractBlocks(mdPath);
  let helpers = '';

  for (const block of blocks) {
    if (block.type === 'helpers') {
      helpers += block.code + '\n';
      continue;
    }
    if (block.mode !== mode) continue;

    if (mode === 'async') {
      test(block.name, async () => {
        const bro = await browser.launch({ headless: true });
        try {
          const vibe = await bro.page();
          const fn = new AsyncFunction('vibe', 'assert', helpers + block.code);
          await fn(vibe, assert);
        } finally {
          await bro.close();
        }
      });
    } else {
      test(block.name, () => {
        const bro = browser.launch({ headless: true });
        try {
          const vibe = bro.page();
          const fn = new Function('vibe', 'assert', helpers + block.code);
          fn(vibe, assert);
        } finally {
          bro.close();
        }
      });
    }
  }
}

module.exports = { runTutorial, extractBlocks };
