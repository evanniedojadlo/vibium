import { SyncBridge } from './bridge';
import { PageSync } from './page';
import { BrowserContextSync } from './context';

export interface LaunchOptions {
  headless?: boolean;
}

export class BrowserSync {
  /** @internal */
  readonly _bridge: SyncBridge;

  constructor(bridge: SyncBridge) {
    this._bridge = bridge;
  }

  page(): PageSync {
    const result = this._bridge.call<{ pageId: number }>('browser.page');
    return new PageSync(this._bridge, result.pageId);
  }

  newPage(): PageSync {
    const result = this._bridge.call<{ pageId: number }>('browser.newPage');
    return new PageSync(this._bridge, result.pageId);
  }

  pages(): PageSync[] {
    const result = this._bridge.call<{ pageIds: number[] }>('browser.pages');
    return result.pageIds.map(id => new PageSync(this._bridge, id));
  }

  newContext(): BrowserContextSync {
    const result = this._bridge.call<{ contextId: number }>('browser.newContext');
    return new BrowserContextSync(this._bridge, result.contextId);
  }

  waitForPage(options?: { timeout?: number }): PageSync {
    const result = this._bridge.call<{ pageId: number }>('browser.waitForPage', [options]);
    return new PageSync(this._bridge, result.pageId);
  }

  waitForPopup(options?: { timeout?: number }): PageSync {
    const result = this._bridge.call<{ pageId: number }>('browser.waitForPopup', [options]);
    return new PageSync(this._bridge, result.pageId);
  }

  close(): void {
    this._bridge.tryQuit();
  }
}

export const browser = {
  launch(options: LaunchOptions = {}): BrowserSync {
    const bridge = SyncBridge.create();
    bridge.call('browser.launch', [options]);
    return new BrowserSync(bridge);
  },
};
