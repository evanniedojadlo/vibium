import { BiDiClient } from './bidi';

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

export class Element {
  private client: BiDiClient;
  private context: string;
  private selector: string;
  private info: ElementInfo;

  constructor(
    client: BiDiClient,
    context: string,
    selector: string,
    info: ElementInfo
  ) {
    this.client = client;
    this.context = context;
    this.selector = selector;
    this.info = info;
  }

  async click(): Promise<void> {
    const { x, y } = this.getCenter();

    const actions = [
      {
        type: 'pointer',
        id: 'mouse',
        parameters: { pointerType: 'mouse' },
        actions: [
          { type: 'pointerMove', x: Math.round(x), y: Math.round(y), duration: 0 },
          { type: 'pointerDown', button: 0 },
          { type: 'pointerUp', button: 0 },
        ],
      },
    ];

    await this.client.send('input.performActions', {
      context: this.context,
      actions,
    });
  }

  async type(text: string): Promise<void> {
    // Click to focus first
    await this.click();

    // Build key actions for each character
    const keyActions: Array<{ type: string; value: string }> = [];
    for (const char of text) {
      keyActions.push(
        { type: 'keyDown', value: char },
        { type: 'keyUp', value: char }
      );
    }

    const actions = [
      {
        type: 'key',
        id: 'keyboard',
        actions: keyActions,
      },
    ];

    await this.client.send('input.performActions', {
      context: this.context,
      actions,
    });
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
      throw new Error(`Element not found: ${this.selector}`);
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
      throw new Error(`Element not found: ${this.selector}`);
    }

    return JSON.parse(result.result.value as string) as BoundingBox;
  }

  private getCenter(): { x: number; y: number } {
    return {
      x: this.info.box.x + this.info.box.width / 2,
      y: this.info.box.y + this.info.box.height / 2,
    };
  }
}
