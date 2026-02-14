import { BiDiClient, BrowsingContextTree, NavigationResult, ScreenshotResult } from './bidi';
import { VibiumProcess } from './clicker';
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

export class Vibe {
  private client: BiDiClient;
  private process: VibiumProcess | null;
  private context: string | null = null;

  constructor(client: BiDiClient, process: VibiumProcess | null) {
    this.client = client;
    this.process = process;
  }

  private async getContext(): Promise<string> {
    if (this.context) {
      return this.context;
    }

    const tree = await this.client.send<BrowsingContextTree>('browsingContext.getTree', {});
    if (!tree.contexts || tree.contexts.length === 0) {
      throw new Error('No browsing context available');
    }

    this.context = tree.contexts[0].context;
    return this.context;
  }

  async go(url: string): Promise<void> {
    debug('navigating', { url });
    const context = await this.getContext();
    await this.client.send<NavigationResult>('browsingContext.navigate', {
      context,
      url,
      wait: 'complete',
    });
    debug('navigation complete', { url });
  }

  async screenshot(): Promise<Buffer> {
    const context = await this.getContext();
    const result = await this.client.send<ScreenshotResult>('browsingContext.captureScreenshot', {
      context,
    });
    return Buffer.from(result.data, 'base64');
  }

  /**
   * Execute JavaScript in the page context.
   */
  async evaluate<T = unknown>(script: string): Promise<T> {
    const context = await this.getContext();
    const result = await this.client.send<{
      type: string;
      result: { type: string; value?: T };
    }>('script.callFunction', {
      functionDeclaration: `() => { ${script} }`,
      target: { context },
      arguments: [],
      awaitPromise: true,
      resultOwnership: 'root',
    });

    return result.result.value as T;
  }

  /**
   * Find an element by CSS selector or semantic options.
   * Waits for element to exist before returning.
   */
  find(selector: string | SelectorOptions, options?: FindOptions): FluentElement {
    const promise = (async () => {
      const context = await this.getContext();
      const params: Record<string, unknown> = {
        context,
        timeout: options?.timeout,
      };

      if (typeof selector === 'string') {
        debug('finding element', { selector, timeout: options?.timeout });
        params.selector = selector;
      } else {
        debug('finding element', { ...selector, timeout: options?.timeout });
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
      debug('element found', { selector: selectorStr, tag: result.tag });
      return new Element(this.client, context, selectorStr, info);
    })();
    return fluent(promise);
  }

  /**
   * Find all elements matching a CSS selector or semantic options.
   * Waits for at least one element to exist before returning.
   */
  async findAll(selector: string | SelectorOptions, options?: FindOptions): Promise<ElementList> {
    const context = await this.getContext();
    const params: Record<string, unknown> = {
      context,
      timeout: options?.timeout,
    };

    if (typeof selector === 'string') {
      debug('finding all elements', { selector, timeout: options?.timeout });
      params.selector = selector;
    } else {
      debug('finding all elements', { ...selector, timeout: options?.timeout });
      Object.assign(params, selector);
      if (selector.timeout && !options?.timeout) params.timeout = selector.timeout;
    }

    const result = await this.client.send<VibiumFindAllResult>('vibium:findAll', params);

    const selectorStr = typeof selector === 'string' ? selector : '';
    const elements = result.elements.map((el) => {
      const info: ElementInfo = { tag: el.tag, text: el.text, box: el.box };
      return new Element(this.client, context, selectorStr, info, el.index);
    });

    return new ElementList(this.client, context, selector, elements);
  }

  async quit(): Promise<void> {
    await this.client.close();
    if (this.process) {
      await this.process.stop();
    }
  }
}
