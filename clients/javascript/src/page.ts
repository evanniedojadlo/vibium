import { BiDiClient, ScreenshotResult } from './bidi';
import { Element, ElementInfo, SelectorOptions, FluentElement, fluent } from './element';
import { ElementList, ElementListInfo } from './element-list';
import { debug } from './utils/debug';

export interface FindOptions {
  /** Timeout in milliseconds to wait for element. Default: 30000 */
  timeout?: number;
}

interface VibiumFindResult {
  tag: string;
  text: string;
  box: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
}

interface VibiumFindAllResult {
  elements: ElementListInfo[];
  count: number;
}

export class Page {
  private client: BiDiClient;
  private contextId: string;

  constructor(client: BiDiClient, contextId: string) {
    this.client = client;
    this.contextId = contextId;
  }

  /** The browsing context ID for this page. */
  get id(): string {
    return this.contextId;
  }

  /** Navigate to a URL. */
  async go(url: string): Promise<void> {
    debug('page.go', { url, context: this.contextId });
    await this.client.send('vibium:page.navigate', {
      context: this.contextId,
      url,
    });
  }

  /** Navigate back in history. */
  async back(): Promise<void> {
    await this.client.send('vibium:page.back', { context: this.contextId });
  }

  /** Navigate forward in history. */
  async forward(): Promise<void> {
    await this.client.send('vibium:page.forward', { context: this.contextId });
  }

  /** Reload the page. */
  async reload(): Promise<void> {
    await this.client.send('vibium:page.reload', { context: this.contextId });
  }

  /** Get the current page URL. */
  async url(): Promise<string> {
    const result = await this.client.send<{ url: string }>('vibium:page.url', {
      context: this.contextId,
    });
    return result.url;
  }

  /** Get the current page title. */
  async title(): Promise<string> {
    const result = await this.client.send<{ title: string }>('vibium:page.title', {
      context: this.contextId,
    });
    return result.title;
  }

  /** Get the full HTML content of the page. */
  async content(): Promise<string> {
    const result = await this.client.send<{ content: string }>('vibium:page.content', {
      context: this.contextId,
    });
    return result.content;
  }

  /** Wait until the page URL matches a pattern. */
  async waitForURL(pattern: string, options?: { timeout?: number }): Promise<void> {
    await this.client.send('vibium:page.waitForURL', {
      context: this.contextId,
      pattern,
      timeout: options?.timeout,
    });
  }

  /** Wait until the page reaches a load state. */
  async waitForLoad(state?: string, options?: { timeout?: number }): Promise<void> {
    await this.client.send('vibium:page.waitForLoad', {
      context: this.contextId,
      state,
      timeout: options?.timeout,
    });
  }

  /** Bring this page/tab to the foreground. */
  async bringToFront(): Promise<void> {
    await this.client.send('vibium:page.activate', { context: this.contextId });
  }

  /** Close this page/tab. */
  async close(): Promise<void> {
    await this.client.send('vibium:page.close', { context: this.contextId });
  }

  /** Take a screenshot of the page. Returns a PNG buffer. */
  async screenshot(): Promise<Buffer> {
    const result = await this.client.send<ScreenshotResult>('browsingContext.captureScreenshot', {
      context: this.contextId,
    });
    return Buffer.from(result.data, 'base64');
  }

  /** Execute JavaScript in the page context. */
  async evaluate<T = unknown>(script: string): Promise<T> {
    const result = await this.client.send<{
      type: string;
      result: { type: string; value?: T };
    }>('script.callFunction', {
      functionDeclaration: `() => { ${script} }`,
      target: { context: this.contextId },
      arguments: [],
      awaitPromise: true,
      resultOwnership: 'root',
    });

    return result.result.value as T;
  }

  /** Find an element by CSS selector or semantic options. Waits for element to exist. */
  find(selector: string | SelectorOptions, options?: FindOptions): FluentElement {
    const promise = (async () => {
      const params: Record<string, unknown> = {
        context: this.contextId,
        timeout: options?.timeout,
      };

      if (typeof selector === 'string') {
        debug('page.find', { selector, timeout: options?.timeout });
        params.selector = selector;
      } else {
        debug('page.find', { ...selector, timeout: options?.timeout });
        Object.assign(params, selector);
        if (selector.timeout && !options?.timeout) params.timeout = selector.timeout;
      }

      const result = await this.client.send<VibiumFindResult>('vibium:find', params);

      const info: ElementInfo = {
        tag: result.tag,
        text: result.text,
        box: result.box,
      };

      const selectorStr = typeof selector === 'string' ? selector : '';
      const selectorParams = typeof selector === 'string' ? { selector } : { ...selector };
      return new Element(this.client, this.contextId, selectorStr, info, undefined, selectorParams);
    })();
    return fluent(promise);
  }

  /** Find all elements matching a CSS selector or semantic options. Waits for at least one. */
  async findAll(selector: string | SelectorOptions, options?: FindOptions): Promise<ElementList> {
    const params: Record<string, unknown> = {
      context: this.contextId,
      timeout: options?.timeout,
    };

    if (typeof selector === 'string') {
      debug('page.findAll', { selector, timeout: options?.timeout });
      params.selector = selector;
    } else {
      debug('page.findAll', { ...selector, timeout: options?.timeout });
      Object.assign(params, selector);
      if (selector.timeout && !options?.timeout) params.timeout = selector.timeout;
    }

    const result = await this.client.send<VibiumFindAllResult>('vibium:findAll', params);

    const selectorStr = typeof selector === 'string' ? selector : '';
    const selectorParams = typeof selector === 'string' ? { selector } : { ...selector };
    const elements = result.elements.map((el) => {
      const info: ElementInfo = { tag: el.tag, text: el.text, box: el.box };
      return new Element(this.client, this.contextId, selectorStr, info, el.index, selectorParams);
    });

    return new ElementList(this.client, this.contextId, selector, elements);
  }
}
