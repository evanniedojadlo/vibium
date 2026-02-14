import { BiDiClient } from './bidi';
import { ElementList, ElementListInfo } from './element-list';
import { ElementNotFoundError } from './utils/errors';

export interface BoundingBox {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface ElementInfo {
  tag: string;
  text: string;
  box: BoundingBox;
}

export interface ScriptResult {
  type: string;
  result: {
    type: string;
    value?: unknown;
  };
}

export interface ActionOptions {
  /** Timeout in milliseconds for actionability checks. Default: 30000 */
  timeout?: number;
}

export interface SelectorOptions {
  role?: string;
  text?: string;
  label?: string;
  placeholder?: string;
  alt?: string;
  title?: string;
  testid?: string;
  xpath?: string;
  near?: string;
  timeout?: number;
}

export class Element {
  private client: BiDiClient;
  private context: string;
  private selector: string;
  private _index?: number;
  private _params: Record<string, unknown>;
  readonly info: ElementInfo;

  constructor(
    client: BiDiClient,
    context: string,
    selector: string,
    info: ElementInfo,
    index?: number,
    params?: Record<string, unknown>
  ) {
    this.client = client;
    this.context = context;
    this.selector = selector;
    this.info = info;
    this._index = index;
    this._params = params || {};
  }

  /** Build the common params sent to vibium: commands for element resolution. */
  private commandParams(extra?: Record<string, unknown>): Record<string, unknown> {
    return {
      ...this._params,
      context: this.context,
      selector: this.selector,
      index: this._index,
      ...extra,
    };
  }

  /** Return params that can identify this element for use as a target (e.g. dragTo). */
  toParams(): Record<string, unknown> {
    return {
      ...this._params,
      selector: this.selector,
      index: this._index,
    };
  }

  /**
   * Click the element.
   * Waits for element to be visible, stable, receive events, and enabled.
   */
  async click(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:click', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Double-click the element. */
  async dblclick(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:dblclick', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /**
   * Fill the element with text (clears existing content first).
   * For inputs and textareas.
   */
  async fill(value: string, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:fill', this.commandParams({
      value,
      timeout: options?.timeout,
    }));
  }

  /**
   * Type text into the element (appends to existing content).
   * Waits for element to be visible, stable, receive events, enabled, and editable.
   */
  async type(text: string, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:type', this.commandParams({
      text,
      timeout: options?.timeout,
    }));
  }

  /**
   * Press a key while the element is focused.
   * Supports key names ("Enter", "Tab") and combos ("Control+a").
   */
  async press(key: string, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:press', this.commandParams({
      key,
      timeout: options?.timeout,
    }));
  }

  /** Clear the element's content (select all + delete). */
  async clear(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:clear', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Check a checkbox (no-op if already checked). */
  async check(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:check', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Uncheck a checkbox (no-op if already unchecked). */
  async uncheck(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:uncheck', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Select an option in a <select> element by value. */
  async selectOption(value: string, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:selectOption', this.commandParams({
      value,
      timeout: options?.timeout,
    }));
  }

  /** Hover over the element (move mouse to center, no click). */
  async hover(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:hover', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Focus the element. */
  async focus(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:focus', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Drag this element to a target element. */
  async dragTo(target: Element, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:dragTo', this.commandParams({
      target: target.toParams(),
      timeout: options?.timeout,
    }));
  }

  /** Tap the element (touch action). */
  async tap(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:tap', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Scroll the element into view. */
  async scrollIntoView(options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:scrollIntoView', this.commandParams({
      timeout: options?.timeout,
    }));
  }

  /** Dispatch a DOM event on the element. */
  async dispatchEvent(eventType: string, eventInit?: Record<string, unknown>, options?: ActionOptions): Promise<void> {
    await this.client.send('vibium:dispatchEvent', this.commandParams({
      eventType,
      eventInit,
      timeout: options?.timeout,
    }));
  }

  async text(): Promise<string> {
    const result = await this.client.send<ScriptResult>('script.callFunction', {
      functionDeclaration: `(selector) => {
        const el = document.querySelector(selector);
        return el ? (el.textContent || '').trim() : null;
      }`,
      target: { context: this.context },
      arguments: [{ type: 'string', value: this.selector }],
      awaitPromise: false,
      resultOwnership: 'root',
    });

    if (result.result.type === 'null') {
      throw new ElementNotFoundError(this.selector);
    }

    return result.result.value as string;
  }

  async getAttribute(name: string): Promise<string | null> {
    const result = await this.client.send<ScriptResult>('script.callFunction', {
      functionDeclaration: `(selector, attrName) => {
        const el = document.querySelector(selector);
        return el ? el.getAttribute(attrName) : null;
      }`,
      target: { context: this.context },
      arguments: [
        { type: 'string', value: this.selector },
        { type: 'string', value: name },
      ],
      awaitPromise: false,
      resultOwnership: 'root',
    });

    if (result.result.type === 'null') {
      return null;
    }

    return result.result.value as string;
  }

  async boundingBox(): Promise<BoundingBox> {
    const result = await this.client.send<ScriptResult>('script.callFunction', {
      functionDeclaration: `(selector) => {
        const el = document.querySelector(selector);
        if (!el) return null;
        const rect = el.getBoundingClientRect();
        return JSON.stringify({
          x: rect.x,
          y: rect.y,
          width: rect.width,
          height: rect.height
        });
      }`,
      target: { context: this.context },
      arguments: [{ type: 'string', value: this.selector }],
      awaitPromise: false,
      resultOwnership: 'root',
    });

    if (result.result.type === 'null') {
      throw new ElementNotFoundError(this.selector);
    }

    return JSON.parse(result.result.value as string) as BoundingBox;
  }

  /** Find a child element by CSS selector or semantic options. Scoped to this element. */
  find(selector: string | SelectorOptions, options?: { timeout?: number }): FluentElement {
    const promise = (async () => {
      const params: Record<string, unknown> = {
        context: this.context,
        scope: this.selector,
        timeout: options?.timeout,
      };

      if (typeof selector === 'string') {
        params.selector = selector;
      } else {
        Object.assign(params, selector);
        if (selector.timeout) params.timeout = selector.timeout;
      }

      const result = await this.client.send<{
        tag: string;
        text: string;
        box: BoundingBox;
      }>('vibium:find', params);

      const info: ElementInfo = { tag: result.tag, text: result.text, box: result.box };
      const childSelector = typeof selector === 'string' ? selector : '';
      const childParams = typeof selector === 'string' ? { selector } : { ...selector };
      return new Element(this.client, this.context, childSelector, info, undefined, childParams);
    })();
    return fluent(promise);
  }

  /** Find all child elements by CSS selector or semantic options. Scoped to this element. */
  async findAll(selector: string | SelectorOptions, options?: { timeout?: number }): Promise<ElementList> {
    const params: Record<string, unknown> = {
      context: this.context,
      scope: this.selector,
      timeout: options?.timeout,
    };

    if (typeof selector === 'string') {
      params.selector = selector;
    } else {
      Object.assign(params, selector);
      if (selector.timeout) params.timeout = selector.timeout;
    }

    const result = await this.client.send<{
      elements: Array<{ tag: string; text: string; box: BoundingBox; index: number }>;
      count: number;
    }>('vibium:findAll', params);

    const selectorStr = typeof selector === 'string' ? selector : '';
    const selectorParams = typeof selector === 'string' ? { selector } : { ...selector };
    const elements = result.elements.map((el) => {
      const info: ElementInfo = { tag: el.tag, text: el.text, box: el.box };
      return new Element(this.client, this.context, selectorStr, info, el.index, selectorParams);
    });

    return new ElementList(this.client, this.context, selector, elements);
  }

  private getCenter(): { x: number; y: number } {
    return {
      x: this.info.box.x + this.info.box.width / 2,
      y: this.info.box.y + this.info.box.height / 2,
    };
  }
}

/** A Promise<Element> that also exposes Element methods for chaining. */
export type FluentElement = Promise<Element> & {
  click(options?: ActionOptions): Promise<void>;
  dblclick(options?: ActionOptions): Promise<void>;
  fill(value: string, options?: ActionOptions): Promise<void>;
  type(text: string, options?: ActionOptions): Promise<void>;
  press(key: string, options?: ActionOptions): Promise<void>;
  clear(options?: ActionOptions): Promise<void>;
  check(options?: ActionOptions): Promise<void>;
  uncheck(options?: ActionOptions): Promise<void>;
  selectOption(value: string, options?: ActionOptions): Promise<void>;
  hover(options?: ActionOptions): Promise<void>;
  focus(options?: ActionOptions): Promise<void>;
  dragTo(target: Element, options?: ActionOptions): Promise<void>;
  tap(options?: ActionOptions): Promise<void>;
  scrollIntoView(options?: ActionOptions): Promise<void>;
  dispatchEvent(eventType: string, eventInit?: Record<string, unknown>, options?: ActionOptions): Promise<void>;
  text(): Promise<string>;
  getAttribute(name: string): Promise<string | null>;
  boundingBox(): Promise<BoundingBox>;
  find(selector: string | SelectorOptions, options?: { timeout?: number }): FluentElement;
  findAll(selector: string | SelectorOptions, options?: { timeout?: number }): Promise<ElementList>;
};

export function fluent(promise: Promise<Element>): FluentElement {
  const p = promise as FluentElement;
  p.click = (opts?) => promise.then(el => el.click(opts));
  p.dblclick = (opts?) => promise.then(el => el.dblclick(opts));
  p.fill = (value, opts?) => promise.then(el => el.fill(value, opts));
  p.type = (text, opts?) => promise.then(el => el.type(text, opts));
  p.press = (key, opts?) => promise.then(el => el.press(key, opts));
  p.clear = (opts?) => promise.then(el => el.clear(opts));
  p.check = (opts?) => promise.then(el => el.check(opts));
  p.uncheck = (opts?) => promise.then(el => el.uncheck(opts));
  p.selectOption = (value, opts?) => promise.then(el => el.selectOption(value, opts));
  p.hover = (opts?) => promise.then(el => el.hover(opts));
  p.focus = (opts?) => promise.then(el => el.focus(opts));
  p.dragTo = (target, opts?) => promise.then(el => el.dragTo(target, opts));
  p.tap = (opts?) => promise.then(el => el.tap(opts));
  p.scrollIntoView = (opts?) => promise.then(el => el.scrollIntoView(opts));
  p.dispatchEvent = (type, init?, opts?) => promise.then(el => el.dispatchEvent(type, init, opts));
  p.text = () => promise.then(el => el.text());
  p.getAttribute = (name) => promise.then(el => el.getAttribute(name));
  p.boundingBox = () => promise.then(el => el.boundingBox());
  p.find = (sel, opts?) => fluent(promise.then(el => el.find(sel, opts)));
  p.findAll = (sel, opts?) => promise.then(el => el.findAll(sel, opts));
  return p;
}

export { ElementList } from './element-list';
