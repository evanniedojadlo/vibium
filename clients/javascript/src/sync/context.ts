import { SyncBridge } from './bridge';
import { PageSync } from './page';
import { RecordingSync } from './recording';
import { Cookie, SetCookieParam, StorageState } from '../context';

export class BrowserContextSync {
  private bridge: SyncBridge;
  private contextId: number;
  readonly recording: RecordingSync;

  constructor(bridge: SyncBridge, contextId: number) {
    this.bridge = bridge;
    this.contextId = contextId;
    this.recording = new RecordingSync(bridge, contextId);
  }

  newPage(): PageSync {
    const result = this.bridge.call<{ pageId: number }>('context.newPage', [this.contextId]);
    return new PageSync(this.bridge, result.pageId);
  }

  close(): void {
    this.bridge.call('context.close', [this.contextId]);
  }

  cookies(urls?: string[]): Cookie[] {
    const result = this.bridge.call<{ cookies: Cookie[] }>('context.cookies', [this.contextId, urls]);
    return result.cookies;
  }

  setCookies(cookies: SetCookieParam[]): void {
    this.bridge.call('context.setCookies', [this.contextId, cookies]);
  }

  clearCookies(): void {
    this.bridge.call('context.clearCookies', [this.contextId]);
  }

  storageState(): StorageState {
    return this.bridge.call<StorageState>('context.storageState', [this.contextId]);
  }

  addInitScript(script: string): string {
    const result = this.bridge.call<{ script: string }>('context.addInitScript', [this.contextId, script]);
    return result.script;
  }
}
