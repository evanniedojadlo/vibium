import WebSocket from 'ws';
import { BiDiCommand, BiDiResponse, BiDiEvent, BiDiMessage, isResponse, isEvent } from './types';
import { ConnectionError } from '../utils/errors';

export type EventHandler = (event: BiDiEvent) => void;

export class BiDiClient {
  private ws: WebSocket;
  private nextId: number = 1;
  private pendingCommands: Map<number, {
    resolve: (result: unknown) => void;
    reject: (error: Error) => void;
  }> = new Map();
  private eventHandlers: EventHandler[] = [];
  private _closed: boolean = false;
  private closePromise: Promise<void>;

  private constructor(ws: WebSocket) {
    this.ws = ws;

    this.closePromise = new Promise((resolve) => {
      ws.on('close', () => {
        this._closed = true;
        for (const [id, pending] of this.pendingCommands) {
          pending.reject(new Error('Connection closed unexpectedly'));
          this.pendingCommands.delete(id);
        }
        resolve();
      });
    });

    ws.on('message', (data: WebSocket.Data) => {
      try {
        const msg = JSON.parse(data.toString()) as BiDiMessage;
        if (isResponse(msg)) {
          this.handleResponse(msg);
        } else if (isEvent(msg)) {
          this.handleEvent(msg);
        }
      } catch (err) {
        console.error('Failed to parse BiDi message:', err);
      }
    });
  }

  static connect(url: string): Promise<BiDiClient> {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(url);

      ws.on('open', () => {
        resolve(new BiDiClient(ws));
      });

      ws.on('error', (err) => {
        reject(new ConnectionError(url, err as Error));
      });
    });
  }

  private handleResponse(response: BiDiResponse): void {
    const pending = this.pendingCommands.get(response.id);
    if (!pending) {
      console.warn('Received response for unknown command:', response.id);
      return;
    }

    this.pendingCommands.delete(response.id);

    if (response.type === 'error' && response.error) {
      pending.reject(new Error(`${response.error}: ${response.message}`));
    } else {
      pending.resolve(response.result);
    }
  }

  private handleEvent(event: BiDiEvent): void {
    for (const handler of this.eventHandlers) {
      handler(event);
    }
  }

  onEvent(handler: EventHandler): void {
    this.eventHandlers.push(handler);
  }

  offEvent(handler: EventHandler): void {
    this.eventHandlers = this.eventHandlers.filter(h => h !== handler);
  }

  send<T = unknown>(method: string, params: Record<string, unknown> = {}): Promise<T> {
    return new Promise((resolve, reject) => {
      const id = this.nextId++;
      const command: BiDiCommand = { id, method, params };

      this.pendingCommands.set(id, {
        resolve: resolve as (result: unknown) => void,
        reject,
      });

      try {
        if (this._closed) {
          throw new Error('Connection closed');
        }
        this.ws.send(JSON.stringify(command));
      } catch (err) {
        this.pendingCommands.delete(id);
        reject(err);
      }
    });
  }

  async close(): Promise<void> {
    if (this._closed) {
      return;
    }
    // Reject all pending commands
    for (const [id, pending] of this.pendingCommands) {
      pending.reject(new Error('Connection closed'));
      this.pendingCommands.delete(id);
    }

    this.ws.close();
    await this.closePromise;
  }
}
