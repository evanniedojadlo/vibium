/**
 * JS Library Tests: Auto-Wait Behavior
 * Tests that actions wait for elements to be actionable
 */

const { test, describe } = require('node:test');
const assert = require('node:assert');

const { browser } = require('../../clients/javascript/dist');

describe('JS Auto-Wait', () => {
  test('find() waits for element to appear', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/dynamic_loading/1');

      // Click the start button to trigger dynamic loading
      const startBtn = await page.find('#start button', { timeout: 5000 });
      await startBtn.click();

      // find() should wait for the dynamically loaded element
      const result = await page.find('#finish h4', { timeout: 10000 });
      assert.ok(result, 'Should find the dynamically loaded element');
      assert.strictEqual(result.info.text, 'Hello World!', 'Should have correct text');
    } finally {
      await b.close();
    }
  });

  test('click() waits for element to be actionable', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/add_remove_elements/');

      // Click the "Add Element" button
      const addBtn = await page.find('button[onclick="addElement()"]', { timeout: 5000 });
      await addBtn.click({ timeout: 5000 });

      // Verify the delete button appeared
      const deleteBtn = await page.find('.added-manually', { timeout: 5000 });
      assert.ok(deleteBtn, 'Delete button should have appeared after click');
    } finally {
      await b.close();
    }
  });

  test('find() times out for non-existent element', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');

      await assert.rejects(
        async () => {
          await page.find('#does-not-exist', { timeout: 1000 });
        },
        /timeout/i,
        'Should throw timeout error'
      );
    } finally {
      await b.close();
    }
  });

  test('timeout error message is clear', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await page.go('https://the-internet.herokuapp.com/');

      try {
        await page.find('#nonexistent-element-xyz', { timeout: 1000 });
        assert.fail('Should have thrown');
      } catch (err) {
        // Error should mention the selector or timeout
        assert.ok(
          err.message.includes('timeout') || err.message.includes('#nonexistent-element-xyz'),
          `Error message should be clear: ${err.message}`
        );
      }
    } finally {
      await b.close();
    }
  });

  test('navigation error message is clear', async () => {
    const b = await browser.launch({ headless: true });
    try {
      const page = await b.page();
      await assert.rejects(
        async () => {
          await page.go('https://test.invalid');
        },
        /error/i,
        'Should throw error for invalid domain'
      );
    } finally {
      await b.close();
    }
  });
});
