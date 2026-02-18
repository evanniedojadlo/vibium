import { SyncBridge } from './bridge';
import { ElementSync } from './element';
import { ElementListSync } from './element-list';
import { KeyboardSync, MouseSync, TouchSync } from './keyboard';
import { ClockSync } from './clock';
import { ElementInfo, SelectorOptions } from '../element';
import { A11yNode, ScreenshotOptions, FindOptions } from '../page';

export class PageSync {
  /** @internal */
  readonly _bridge: SyncBridge;
  /** @internal */
  readonly _pageId: number;

  readonly keyboard: KeyboardSync;
  readonly mouse: MouseSync;
  readonly touch: TouchSync;
  readonly clock: ClockSync;

  constructor(bridge: SyncBridge, pageId: number) {
    this._bridge = bridge;
    this._pageId = pageId;
    this.keyboard = new KeyboardSync(bridge, pageId);
    this.mouse = new MouseSync(bridge, pageId);
    this.touch = new TouchSync(bridge, pageId);
    this.clock = new ClockSync(bridge, pageId);
  }

  // --- Navigation ---

  go(url: string): void {
    this._bridge.call('page.go', [this._pageId, url]);
  }

  back(): void {
    this._bridge.call('page.back', [this._pageId]);
  }

  forward(): void {
    this._bridge.call('page.forward', [this._pageId]);
  }

  reload(): void {
    this._bridge.call('page.reload', [this._pageId]);
  }

  // --- Info ---

  url(): string {
    const result = this._bridge.call<{ url: string }>('page.url', [this._pageId]);
    return result.url;
  }

  title(): string {
    const result = this._bridge.call<{ title: string }>('page.title', [this._pageId]);
    return result.title;
  }

  content(): string {
    const result = this._bridge.call<{ content: string }>('page.content', [this._pageId]);
    return result.content;
  }

  // --- Finding ---

  find(selector: string | SelectorOptions, options?: FindOptions): ElementSync {
    const result = this._bridge.call<{ elementId: number; info: ElementInfo }>('page.find', [this._pageId, selector, options]);
    return new ElementSync(this._bridge, result.elementId, result.info);
  }

  findAll(selector: string | SelectorOptions, options?: FindOptions): ElementListSync {
    const result = this._bridge.call<{ listId: number; elementIds: number[]; count: number }>('page.findAll', [this._pageId, selector, options]);
    return new ElementListSync(this._bridge, result.listId);
  }

  waitFor(selector: string | SelectorOptions, options?: FindOptions): ElementSync {
    const result = this._bridge.call<{ elementId: number; info: ElementInfo }>('page.waitFor', [this._pageId, selector, options]);
    return new ElementSync(this._bridge, result.elementId, result.info);
  }

  // --- Waiting ---

  wait(ms: number): void {
    this._bridge.call('page.wait', [this._pageId, ms]);
  }

  waitForURL(pattern: string, options?: { timeout?: number }): void {
    this._bridge.call('page.waitForURL', [this._pageId, pattern, options]);
  }

  waitForLoad(state?: string, options?: { timeout?: number }): void {
    this._bridge.call('page.waitForLoad', [this._pageId, state, options]);
  }

  waitForFunction<T = unknown>(fn: string, options?: { timeout?: number }): T {
    const result = this._bridge.call<{ value: T }>('page.waitForFunction', [this._pageId, fn, options]);
    return result.value;
  }

  // --- Screenshots & PDF ---

  screenshot(options?: ScreenshotOptions): Buffer {
    const result = this._bridge.call<{ data: string }>('page.screenshot', [this._pageId, options]);
    return Buffer.from(result.data, 'base64');
  }

  pdf(): Buffer {
    const result = this._bridge.call<{ data: string }>('page.pdf', [this._pageId]);
    return Buffer.from(result.data, 'base64');
  }

  // --- Evaluation ---

  evaluate<T = unknown>(script: string): T {
    const result = this._bridge.call<{ result: T }>('page.evaluate', [this._pageId, script]);
    return result.result;
  }

  eval<T = unknown>(expression: string): T {
    const result = this._bridge.call<{ value: T }>('page.eval', [this._pageId, expression]);
    return result.value;
  }

  evalHandle(expression: string): string {
    const result = this._bridge.call<{ handle: string }>('page.evalHandle', [this._pageId, expression]);
    return result.handle;
  }

  addScript(source: string): void {
    this._bridge.call('page.addScript', [this._pageId, source]);
  }

  addStyle(source: string): void {
    this._bridge.call('page.addStyle', [this._pageId, source]);
  }

  // --- Lifecycle ---

  bringToFront(): void {
    this._bridge.call('page.bringToFront', [this._pageId]);
  }

  close(): void {
    this._bridge.call('page.close', [this._pageId]);
  }

  // --- Emulation ---

  setViewport(size: { width: number; height: number }): void {
    this._bridge.call('page.setViewport', [this._pageId, size]);
  }

  viewport(): { width: number; height: number } {
    return this._bridge.call<{ width: number; height: number }>('page.viewport', [this._pageId]);
  }

  emulateMedia(opts: {
    media?: 'screen' | 'print' | null;
    colorScheme?: 'light' | 'dark' | 'no-preference' | null;
    reducedMotion?: 'reduce' | 'no-preference' | null;
    forcedColors?: 'active' | 'none' | null;
    contrast?: 'more' | 'no-preference' | null;
  }): void {
    this._bridge.call('page.emulateMedia', [this._pageId, opts]);
  }

  setContent(html: string): void {
    this._bridge.call('page.setContent', [this._pageId, html]);
  }

  setGeolocation(coords: { latitude: number; longitude: number; accuracy?: number }): void {
    this._bridge.call('page.setGeolocation', [this._pageId, coords]);
  }

  setWindow(options: {
    width?: number;
    height?: number;
    x?: number;
    y?: number;
    state?: 'normal' | 'maximized' | 'minimized' | 'fullscreen';
  }): void {
    this._bridge.call('page.setWindow', [this._pageId, options]);
  }

  window(): { state: string; width: number; height: number; x: number; y: number } {
    return this._bridge.call('page.window', [this._pageId]);
  }

  // --- Frames ---

  frames(): PageSync[] {
    const result = this._bridge.call<{ frameIds: number[] }>('page.frames', [this._pageId]);
    return result.frameIds.map(id => new PageSync(this._bridge, id));
  }

  frame(nameOrUrl: string): PageSync | null {
    const result = this._bridge.call<{ frameId: number | null }>('page.frame', [this._pageId, nameOrUrl]);
    if (result.frameId === null) return null;
    return new PageSync(this._bridge, result.frameId);
  }

  mainFrame(): PageSync {
    return this;
  }

  // --- Accessibility ---

  a11yTree(options?: { interestingOnly?: boolean; root?: string }): A11yNode {
    const result = this._bridge.call<{ tree: A11yNode }>('page.a11yTree', [this._pageId, options]);
    return result.tree;
  }

  // --- Network ---

  route(pattern: string, action: 'continue' | 'abort' | { status?: number; body?: string; headers?: Record<string, string> }): void {
    this._bridge.call('page.route', [this._pageId, pattern, action]);
  }

  unroute(pattern: string): void {
    this._bridge.call('page.unroute', [this._pageId, pattern]);
  }

  setHeaders(headers: Record<string, string>): void {
    this._bridge.call('page.setHeaders', [this._pageId, headers]);
  }

  waitForRequest(pattern: string, options?: { timeout?: number }): { url: string; method: string; headers: Record<string, string>; postData: string | null } {
    return this._bridge.call('page.waitForRequest', [this._pageId, pattern, options]);
  }

  waitForResponse(pattern: string, options?: { timeout?: number }): { url: string; status: number; headers: Record<string, string> } {
    return this._bridge.call('page.waitForResponse', [this._pageId, pattern, options]);
  }

  // --- Events (simplified) ---

  onDialog(action: 'accept' | 'dismiss'): void {
    this._bridge.call('page.onDialog', [this._pageId, action]);
  }

  onConsole(mode: 'collect'): void {
    this._bridge.call('page.onConsole', [this._pageId, mode]);
  }

  consoleMessages(): { type: string; text: string }[] {
    const result = this._bridge.call<{ messages: { type: string; text: string }[] }>('page.consoleMessages', [this._pageId]);
    return result.messages;
  }

  onError(mode: 'collect'): void {
    this._bridge.call('page.onError', [this._pageId, mode]);
  }

  errors(): { message: string }[] {
    const result = this._bridge.call<{ errors: { message: string }[] }>('page.errors', [this._pageId]);
    return result.errors;
  }

  removeAllListeners(event?: 'request' | 'response' | 'dialog' | 'console' | 'error'): void {
    this._bridge.call('page.removeAllListeners', [this._pageId, event]);
  }
}
