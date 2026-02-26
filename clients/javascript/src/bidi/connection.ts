import WebSocket from 'ws';
import { BiDiMessage } from './types';
import { ConnectionError } from '../utils/errors';

export type MessageHandler = (msg: BiDiMessage) => void;

export class BiDiConnection {
  private ws: WebSocket;
  private messageHandler: MessageHandler | null = null;
  private closeHandler: (() => void) | null = null;
  private closePromise: Promise<void>;
  private _closed: boolean = false;

  private constructor(ws: WebSocket) {
    this.ws = ws;

    this.closePromise = new Promise((resolve) => {
      ws.on('close', () => {
        this._closed = true;
        if (this.closeHandler) this.closeHandler();
        resolve();
      });
    });

    ws.on('message', (data: WebSocket.Data) => {
      if (this.messageHandler) {
        try {
          const msg = JSON.parse(data.toString()) as BiDiMessage;
          this.messageHandler(msg);
        } catch (err) {
          console.error('Failed to parse BiDi message:', err);
        }
      }
    });
  }

  static connect(url: string): Promise<BiDiConnection> {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(url);

      ws.on('open', () => {
        resolve(new BiDiConnection(ws));
      });

      ws.on('error', (err) => {
        reject(new ConnectionError(url, err as Error));
      });
    });
  }

  get closed(): boolean {
    return this._closed;
  }

  onMessage(handler: MessageHandler): void {
    this.messageHandler = handler;
  }

  /** Register a handler called when the connection closes unexpectedly. */
  onClose(handler: () => void): void {
    this.closeHandler = handler;
  }

  send(message: string): void {
    if (this._closed) {
      throw new Error('Connection closed');
    }
    this.ws.send(message);
  }

  async close(): Promise<void> {
    if (this._closed) {
      return;
    }
    this.ws.close();
    await this.closePromise;
  }
}
