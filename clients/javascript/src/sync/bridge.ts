import { Worker, MessageChannel, receiveMessageOnPort } from 'worker_threads';
import * as path from 'path';

interface CommandResult {
  result?: unknown;
  error?: string;
}

// Track active bridges for cleanup on process exit
const activeBridges: Set<SyncBridge> = new Set();

function cleanup() {
  for (const bridge of activeBridges) {
    try {
      // Try graceful quit first (with short timeout)
      bridge.tryQuit();
    } catch {
      // Ignore errors during cleanup
    }
  }
  activeBridges.clear();
}

// Register cleanup handlers once
let handlersRegistered = false;
function registerCleanupHandlers() {
  if (handlersRegistered) return;
  handlersRegistered = true;

  process.on('exit', cleanup);
  process.on('SIGINT', () => {
    cleanup();
    process.exit(130);
  });
  process.on('SIGTERM', () => {
    cleanup();
    process.exit(143);
  });
}

export class SyncBridge {
  private worker: Worker;
  private signal: Int32Array;
  private commandId = 0;
  private terminated = false;

  private constructor(worker: Worker, signal: Int32Array) {
    this.worker = worker;
    this.signal = signal;
  }

  static create(): SyncBridge {
    registerCleanupHandlers();

    const signal = new Int32Array(new SharedArrayBuffer(4));

    // Resolve worker path - works for both dev and built scenarios
    const workerPath = path.join(__dirname, 'worker.js');

    const worker = new Worker(workerPath, {
      workerData: { signal },
    });

    const bridge = new SyncBridge(worker, signal);
    activeBridges.add(bridge);

    return bridge;
  }

  call<T = unknown>(method: string, args: unknown[] = []): T {
    const cmd = { id: this.commandId++, method, args };

    // Create a channel for this call's result
    const { port1, port2 } = new MessageChannel();

    // Reset signal
    Atomics.store(this.signal, 0, 0);

    // Send command with the port
    this.worker.postMessage({ cmd, port: port2 }, [port2]);

    // Block until worker signals completion
    Atomics.wait(this.signal, 0, 0);

    // Synchronously receive the result
    const message = receiveMessageOnPort(port1);
    port1.close();

    if (!message) {
      throw new Error('No response from worker');
    }

    const response = message.message as CommandResult;

    if (response.error) {
      throw new Error(response.error);
    }

    return response.result as T;
  }

  tryQuit(): void {
    if (this.terminated) return;

    try {
      // Try to send quit command with timeout for graceful shutdown
      const cmd = { id: this.commandId++, method: 'quit', args: [] };
      const { port1, port2 } = new MessageChannel();

      Atomics.store(this.signal, 0, 0);
      this.worker.postMessage({ cmd, port: port2 }, [port2]);

      // Wait with timeout (5s for cleanup - clicker needs time to kill Chrome)
      const result = Atomics.wait(this.signal, 0, 0, 5000);
      port1.close();

      // Only force terminate if timed out
      if (result === 'timed-out') {
        this.terminate();
      } else {
        // Quit succeeded, terminate worker thread so process can exit
        this.terminated = true;
        activeBridges.delete(this);
        this.worker.terminate();
      }
    } catch {
      // If anything fails, force terminate
      this.terminate();
    }
  }

  terminate(): void {
    if (this.terminated) return;
    this.terminated = true;
    activeBridges.delete(this);
    this.worker.terminate();
  }
}
