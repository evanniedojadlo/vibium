import { SyncBridge } from './bridge';
import { RecordingStartOptions, RecordingStopOptions } from '../recording';

export class RecordingSync {
  private bridge: SyncBridge;
  private contextId: number;

  constructor(bridge: SyncBridge, contextId: number) {
    this.bridge = bridge;
    this.contextId = contextId;
  }

  start(options: RecordingStartOptions = {}): void {
    this.bridge.call('recording.start', [this.contextId, options]);
  }

  stop(options: RecordingStopOptions = {}): Buffer {
    const result = this.bridge.call<{ data: string }>('recording.stop', [this.contextId, options]);
    return Buffer.from(result.data, 'base64');
  }

  startChunk(options: { name?: string; title?: string } = {}): void {
    this.bridge.call('recording.startChunk', [this.contextId, options]);
  }

  stopChunk(options: RecordingStopOptions = {}): Buffer {
    const result = this.bridge.call<{ data: string }>('recording.stopChunk', [this.contextId, options]);
    return Buffer.from(result.data, 'base64');
  }

  startGroup(name: string, options: { location?: { file: string; line?: number; column?: number } } = {}): void {
    this.bridge.call('recording.startGroup', [this.contextId, name, options]);
  }

  stopGroup(): void {
    this.bridge.call('recording.stopGroup', [this.contextId]);
  }
}
