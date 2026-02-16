import { BiDiClient } from './bidi';
import { Page } from './page';

export class BrowserContext {
  private client: BiDiClient;
  private userContextId: string;

  constructor(client: BiDiClient, userContextId: string) {
    this.client = client;
    this.userContextId = userContextId;
  }

  /** The user context ID for this browser context. */
  get id(): string {
    return this.userContextId;
  }

  /** Create a new page (tab) in this context. */
  async newPage(): Promise<Page> {
    const result = await this.client.send<{ context: string }>('vibium:context.newPage', {
      userContext: this.userContextId,
    });
    return new Page(this.client, result.context);
  }

  /** Close this context and all its pages. */
  async close(): Promise<void> {
    await this.client.send('browser.removeUserContext', {
      userContext: this.userContextId,
    });
  }
}
