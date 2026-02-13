import { VibiumProcess } from './clicker';
import { BiDiClient } from './bidi';
import { Vibe } from './vibe';
import { debug, info } from './utils/debug';

export interface LaunchOptions {
  headless?: boolean;
  port?: number;
  executablePath?: string;
}

export const browser = {
  async launch(options: LaunchOptions = {}): Promise<Vibe> {
    const { headless = false, port, executablePath } = options;
    debug('launching browser', { headless, port, executablePath });

    // Start the vibium process
    const process = await VibiumProcess.start({
      headless,
      port,
      executablePath,
    });
    debug('vibium started', { port: process.port });

    // Connect to the proxy
    const client = await BiDiClient.connect(`ws://localhost:${process.port}`);
    info('browser launched', { port: process.port });

    return new Vibe(client, process);
  },
};
