import WebSocket from 'ws';
import { BiDiCommand, BiDiResponse, BiDiEvent, BiDiMessage, isResponse, isEvent } from './types';
import { ConnectionError } from '../utils/errors';

export type EventHandler = (event: BiDiEvent) => void;

const DEFAULT_CONNECT_TIMEOUT = 30_000;
const DEFAULT_COMMAND_TIMEOUT = 60_000;

export class BiDiClient {
  private ws: WebSocket;
  private nextId: number = 1;
  private pendingCommands: Map<number, {
    resolve: (result: unknown) => void;
    reject: (error: Error) => void;
    timer: ReturnType<typeof setTimeout>;
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
          clearTimeout(pending.timer);
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

  static connect(url: string, timeout: number = DEFAULT_CONNECT_TIMEOUT): Promise<BiDiClient> {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(url);
      let settled = false;

      const timer = setTimeout(() => {
        if (!settled) {
          settled = true;
          ws.close();
          reject(new ConnectionError(url, new Error(`Timed out after ${timeout}ms`)));
        }
      }, timeout);

      ws.on('open', () => {
        if (!settled) {
          settled = true;
          clearTimeout(timer);
          resolve(new BiDiClient(ws));
        }
      });

      ws.on('error', (err) => {
        if (!settled) {
          settled = true;
          clearTimeout(timer);
          reject(new ConnectionError(url, err as Error));
        }
      });
    });
  }

  private handleResponse(response: BiDiResponse): void {
    const pending = this.pendingCommands.get(response.id);
    if (!pending) {
      console.warn('Received response for unknown command:', response.id);
      return;
    }

    clearTimeout(pending.timer);
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

  send<T = unknown>(method: string, params: Record<string, unknown> = {}, timeout: number = DEFAULT_COMMAND_TIMEOUT): Promise<T> {
    return new Promise((resolve, reject) => {
      const id = this.nextId++;
      const command: BiDiCommand = { id, method, params };

      const timer = setTimeout(() => {
        this.pendingCommands.delete(id);
        reject(new Error(`Command '${method}' timed out after ${timeout}ms`));
      }, timeout);

      this.pendingCommands.set(id, {
        resolve: resolve as (result: unknown) => void,
        reject,
        timer,
      });

      try {
        if (this._closed) {
          throw new Error('Connection closed');
        }
        this.ws.send(JSON.stringify(command));
      } catch (err) {
        clearTimeout(timer);
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
      clearTimeout(pending.timer);
      pending.reject(new Error('Connection closed'));
      this.pendingCommands.delete(id);
    }

    this.ws.close();
    await this.closePromise;
  }
}
