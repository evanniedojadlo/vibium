import { BiDiClient, BrowsingContextTree, NavigationResult, ScreenshotResult } from './bidi';
import { ClickerProcess } from './clicker';
import { Element, ElementInfo, ScriptResult } from './element';

export class Vibe {
  private client: BiDiClient;
  private process: ClickerProcess | null;
  private context: string | null = null;

  constructor(client: BiDiClient, process: ClickerProcess | null) {
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
    const context = await this.getContext();
    await this.client.send<NavigationResult>('browsingContext.navigate', {
      context,
      url,
      wait: 'complete',
    });
  }

  async screenshot(): Promise<Buffer> {
    const context = await this.getContext();
    const result = await this.client.send<ScreenshotResult>('browsingContext.captureScreenshot', {
      context,
    });
    return Buffer.from(result.data, 'base64');
  }

  async find(selector: string): Promise<Element> {
    const context = await this.getContext();

    const result = await this.client.send<ScriptResult>('script.callFunction', {
      functionDeclaration: `(selector) => {
        const el = document.querySelector(selector);
        if (!el) return null;
        const rect = el.getBoundingClientRect();
        return JSON.stringify({
          tag: el.tagName,
          text: (el.textContent || '').trim().substring(0, 100),
          box: {
            x: rect.x,
            y: rect.y,
            width: rect.width,
            height: rect.height
          }
        });
      }`,
      target: { context },
      arguments: [{ type: 'string', value: selector }],
      awaitPromise: false,
      resultOwnership: 'root',
    });

    if (result.result.type === 'null') {
      throw new Error(`Element not found: ${selector}`);
    }

    const info: ElementInfo = JSON.parse(result.result.value as string);
    return new Element(this.client, context, selector, info);
  }

  async quit(): Promise<void> {
    await this.client.close();
    if (this.process) {
      await this.process.stop();
    }
  }
}
