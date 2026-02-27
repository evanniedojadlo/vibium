import { createInterface, Interface as ReadlineInterface } from 'readline';
import { Writable, Readable } from 'stream';
import { BiDiCommand, BiDiResponse, BiDiEvent, BiDiMessage, isResponse, isEvent } from './types';

export type EventHandler = (event: BiDiEvent) => void;

const DEFAULT_COMMAND_TIMEOUT = 60_000;

export class BiDiClient {
  private stdin: Writable;
  private rl: ReadlineInterface;
  private nextId: number = 1;
  private pendingCommands: Map<number, {
    resolve: (result: unknown) => void;
    reject: (error: Error) => void;
    timer: ReturnType<typeof setTimeout>;
  }> = new Map();
  private eventHandlers: EventHandler[] = [];
  private _closed: boolean = false;

  private constructor(stdin: Writable, stdout: Readable) {
    this.stdin = stdin;

    // Read responses/events line by line from stdout
    this.rl = createInterface({ input: stdout, crlfDelay: Infinity });

    this.rl.on('line', (line: string) => {
      const trimmed = line.trim();
      if (!trimmed) return;
      try {
        const msg = JSON.parse(trimmed) as BiDiMessage;
        if (isResponse(msg)) {
          this.handleResponse(msg);
        } else if (isEvent(msg)) {
          this.handleEvent(msg);
        }
      } catch (err) {
        // Ignore unparseable lines
      }
    });

    this.rl.on('close', () => {
      this._closed = true;
      for (const [id, pending] of this.pendingCommands) {
        clearTimeout(pending.timer);
        pending.reject(new Error('Connection closed unexpectedly'));
        this.pendingCommands.delete(id);
      }
    });
  }

  /**
   * Create a BiDiClient from stdin/stdout streams.
   * Optionally replay pre-ready lines (events received before the ready signal).
   */
  static fromStreams(stdin: Writable, stdout: Readable, preReadyLines: string[] = []): BiDiClient {
    const client = new BiDiClient(stdin, stdout);
    // Replay buffered events
    for (const line of preReadyLines) {
      try {
        const msg = JSON.parse(line) as BiDiMessage;
        if (isEvent(msg)) {
          client.handleEvent(msg);
        }
      } catch {
        // Ignore
      }
    }
    return client;
  }

  private handleResponse(response: BiDiResponse): void {
    const pending = this.pendingCommands.get(response.id);
    if (!pending) {
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
        this.stdin.write(JSON.stringify(command) + '\n');
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
    this._closed = true;

    // Reject all pending commands
    for (const [id, pending] of this.pendingCommands) {
      clearTimeout(pending.timer);
      pending.reject(new Error('Connection closed'));
      this.pendingCommands.delete(id);
    }

    this.rl.close();
    try { this.stdin.end(); } catch {}
  }
}
