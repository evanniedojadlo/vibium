import { spawn, execFileSync, ChildProcess } from 'child_process';
import { getVibiumBinPath } from './binary';
import { TimeoutError, BrowserCrashedError } from '../utils/errors';

export interface VibiumProcessOptions {
  headless?: boolean;
  executablePath?: string;
}

export class VibiumProcess {
  private _process: ChildProcess;
  private _stopped: boolean = false;
  private _preReadyLines: string[] = [];

  private constructor(process: ChildProcess, preReadyLines: string[]) {
    this._process = process;
    this._preReadyLines = preReadyLines;
  }

  /** The child process stdin stream (for sending commands). */
  get stdin() { return this._process.stdin!; }

  /** The child process stdout stream (for receiving responses/events). */
  get stdout() { return this._process.stdout!; }

  /** Lines received before the vibium:ready signal (buffered events). */
  get preReadyLines(): string[] { return this._preReadyLines; }

  static async start(options: VibiumProcessOptions = {}): Promise<VibiumProcess> {
    const binaryPath = options.executablePath || getVibiumBinPath();

    const args = ['pipe'];
    if (options.headless === true) {
      args.push('--headless');
    }

    const proc = spawn(binaryPath, args, {
      stdio: ['pipe', 'pipe', 'pipe'],
    });

    // Read lines from stdout until we get the vibium:ready signal.
    // Events (e.g. browsingContext.contextCreated) may arrive first.
    const preReadyLines: string[] = [];
    await new Promise<void>((resolve, reject) => {
      let buffer = '';
      let resolved = false;

      const timeout = setTimeout(() => {
        if (!resolved) {
          resolved = true;
          reject(new TimeoutError('vibium', 30000, 'waiting for vibium ready signal'));
        }
      }, 30000);

      const handleData = (data: Buffer) => {
        buffer += data.toString();
        let newlineIdx: number;
        while ((newlineIdx = buffer.indexOf('\n')) !== -1) {
          const line = buffer.slice(0, newlineIdx).trim();
          buffer = buffer.slice(newlineIdx + 1);
          if (!line) continue;

          try {
            const msg = JSON.parse(line);
            if (msg.method === 'vibium:ready') {
              if (!resolved) {
                resolved = true;
                clearTimeout(timeout);
                // Stop listening for data — the BiDiClient will take over
                proc.stdout?.removeListener('data', handleData);
                resolve();
              }
              return;
            }
          } catch {
            // Not JSON, ignore
          }
          // Buffer pre-ready lines for replay
          preReadyLines.push(line);
        }
      };

      proc.stdout?.on('data', handleData);

      proc.on('error', (err) => {
        if (!resolved) {
          resolved = true;
          clearTimeout(timeout);
          reject(err);
        }
      });

      proc.on('exit', (code) => {
        if (!resolved) {
          resolved = true;
          clearTimeout(timeout);
          reject(new BrowserCrashedError(code ?? 1, buffer));
        }
      });
    });

    const vp = new VibiumProcess(proc, preReadyLines);

    // Clean up child process when Node exits unexpectedly
    const cleanup = () => vp.stop();
    process.on('exit', cleanup);
    process.on('SIGINT', cleanup);
    process.on('SIGTERM', cleanup);
    vp._cleanupListeners = cleanup;

    return vp;
  }

  private _cleanupListeners: (() => void) | null = null;

  async stop(): Promise<void> {
    if (this._stopped) {
      return;
    }
    this._stopped = true;

    // Remove process exit listeners to avoid leaks
    if (this._cleanupListeners) {
      process.removeListener('exit', this._cleanupListeners);
      process.removeListener('SIGINT', this._cleanupListeners);
      process.removeListener('SIGTERM', this._cleanupListeners);
      this._cleanupListeners = null;
    }

    return new Promise((resolve) => {
      this._process.on('exit', () => {
        resolve();
      });

      // Close stdin to signal the child to shut down
      try { this._process.stdin?.end(); } catch {}

      if (process.platform === 'win32') {
        try {
          execFileSync('taskkill', ['/T', '/F', '/PID', this._process.pid!.toString()], { stdio: 'ignore' });
        } catch {
          // Process may have already exited
        }
        resolve();
      } else {
        let exited = false;
        this._process.on('exit', () => { exited = true; });

        // Try graceful shutdown first (closing stdin should trigger exit)
        // Use SIGTERM as fallback
        setTimeout(() => {
          if (!exited) {
            try { this._process.kill('SIGTERM'); } catch {}
          }
        }, 1000);

        // Force kill after longer timeout
        setTimeout(() => {
          if (!exited) {
            try { this._process.kill('SIGKILL'); } catch {}
            setTimeout(() => resolve(), 500);
          } else {
            resolve();
          }
        }, 4000);
      }
    });
  }
}
