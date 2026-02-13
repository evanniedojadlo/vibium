import { parentPort, workerData } from 'worker_threads';
import { browser, Browser } from '../browser';
import { Page, FindOptions } from '../page';
import { Element, ActionOptions } from '../element';

interface WorkerData {
  signal: Int32Array;
}

interface Command {
  id: number;
  method: string;
  args: unknown[];
}

const { signal } = workerData as WorkerData;

let browserInstance: Browser | null = null;
let page: Page | null = null;
let elements: Map<number, Element> = new Map();
let nextElementId = 1;

async function handleCommand(cmd: Command): Promise<unknown> {
  switch (cmd.method) {
    case 'launch': {
      const [options] = cmd.args as [{ headless?: boolean } | undefined];
      browserInstance = await browser.launch(options);
      page = await browserInstance.page();
      return { success: true };
    }

    case 'go': {
      if (!page) throw new Error('Browser not launched');
      const [url] = cmd.args as [string];
      await page.go(url);
      return { success: true };
    }

    case 'screenshot': {
      if (!page) throw new Error('Browser not launched');
      const buffer = await page.screenshot();
      return { data: buffer.toString('base64') };
    }

    case 'evaluate': {
      if (!page) throw new Error('Browser not launched');
      const [script] = cmd.args as [string];
      const result = await page.evaluate(script);
      return { result };
    }

    case 'find': {
      if (!page) throw new Error('Browser not launched');
      const [selector, options] = cmd.args as [string, FindOptions | undefined];
      const element = await page.find(selector, options);
      const elementId = nextElementId++;
      elements.set(elementId, element);
      return { elementId, info: element.info };
    }

    case 'element.click': {
      const [elementId, options] = cmd.args as [number, ActionOptions | undefined];
      const element = elements.get(elementId);
      if (!element) throw new Error(`Element ${elementId} not found`);
      await element.click(options);
      return { success: true };
    }

    case 'element.type': {
      const [elementId, text, options] = cmd.args as [number, string, ActionOptions | undefined];
      const element = elements.get(elementId);
      if (!element) throw new Error(`Element ${elementId} not found`);
      await element.type(text, options);
      return { success: true };
    }

    case 'element.text': {
      const [elementId] = cmd.args as [number];
      const element = elements.get(elementId);
      if (!element) throw new Error(`Element ${elementId} not found`);
      const text = await element.text();
      return { text };
    }

    case 'element.getAttribute': {
      const [elementId, name] = cmd.args as [number, string];
      const element = elements.get(elementId);
      if (!element) throw new Error(`Element ${elementId} not found`);
      const value = await element.getAttribute(name);
      return { value };
    }

    case 'element.boundingBox': {
      const [elementId] = cmd.args as [number];
      const element = elements.get(elementId);
      if (!element) throw new Error(`Element ${elementId} not found`);
      const box = await element.boundingBox();
      return { box };
    }

    case 'quit': {
      if (!browserInstance) throw new Error('Browser not launched');
      await browserInstance.close();
      browserInstance = null;
      page = null;
      elements.clear();
      return { success: true };
    }

    default:
      throw new Error(`Unknown method: ${cmd.method}`);
  }
}

parentPort!.on('message', async ({ cmd, port }: { cmd: Command; port: MessagePort }) => {
  let result: unknown;
  let error: string | null = null;

  try {
    result = await handleCommand(cmd);
  } catch (err) {
    error = err instanceof Error ? err.message : String(err);
  }

  // Send result through the port
  port.postMessage({ result, error });

  // Signal completion
  Atomics.store(signal, 0, 1);
  Atomics.notify(signal, 0);
});
