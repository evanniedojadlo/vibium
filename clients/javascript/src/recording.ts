import { BiDiClient } from './bidi';

export interface RecordingStartOptions {
  name?: string;
  screenshots?: boolean;
  snapshots?: boolean;
  sources?: boolean;
  title?: string;
  bidi?: boolean;
  /** Screenshot format: 'jpeg' (default, faster/smaller) or 'png' (lossless). */
  format?: 'jpeg' | 'png';
  /** JPEG quality 0.0-1.0 (default 0.5). Ignored for PNG. */
  quality?: number;
}

export interface RecordingStopOptions {
  path?: string;
}

export class Recording {
  private client: BiDiClient;
  private userContextId: string;

  constructor(client: BiDiClient, userContextId: string) {
    this.client = client;
    this.userContextId = userContextId;
  }

  /** Start recording. */
  async start(options: RecordingStartOptions = {}): Promise<void> {
    await this.client.send('vibium:recording.start', {
      userContext: this.userContextId,
      ...options,
    });
  }

  /** Stop recording and return the recording zip as a Buffer. */
  async stop(options: RecordingStopOptions = {}): Promise<Buffer> {
    const result = await this.client.send<{ path?: string; data?: string }>('vibium:recording.stop', {
      userContext: this.userContextId,
      ...options,
    });

    if (options.path) {
      // File was written by the engine; read it back
      const fs = await import('fs');
      return fs.readFileSync(result.path!);
    }

    // Base64-encoded zip returned inline
    return Buffer.from(result.data!, 'base64');
  }

  /** Start a new recording chunk (resets event buffer, keeps resources). */
  async startChunk(options: { name?: string; title?: string } = {}): Promise<void> {
    await this.client.send('vibium:recording.startChunk', {
      userContext: this.userContextId,
      ...options,
    });
  }

  /** Stop the current chunk and return the recording zip as a Buffer. */
  async stopChunk(options: RecordingStopOptions = {}): Promise<Buffer> {
    const result = await this.client.send<{ path?: string; data?: string }>('vibium:recording.stopChunk', {
      userContext: this.userContextId,
      ...options,
    });

    if (options.path) {
      const fs = await import('fs');
      return fs.readFileSync(result.path!);
    }

    return Buffer.from(result.data!, 'base64');
  }

  /** Start a named group of actions in the recording. */
  async startGroup(name: string, options: { location?: { file: string; line?: number; column?: number } } = {}): Promise<void> {
    await this.client.send('vibium:recording.startGroup', {
      userContext: this.userContextId,
      name,
      ...options,
    });
  }

  /** End the current group. */
  async stopGroup(): Promise<void> {
    await this.client.send('vibium:recording.stopGroup', {
      userContext: this.userContextId,
    });
  }
}
