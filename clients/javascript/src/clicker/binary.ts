import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { createRequire } from 'module';
import { getPlatform, getArch } from './platform';

/**
 * Resolve the path to the clicker binary.
 *
 * Search order:
 * 1. CLICKER_PATH environment variable
 * 2. Platform-specific npm package (@vibium/{platform}-{arch})
 * 3. Local development paths (relative to cwd)
 */
export function getClickerPath(): string {
  // 1. Check environment variable
  const envPath = process.env.CLICKER_PATH;
  if (envPath && fs.existsSync(envPath)) {
    return envPath;
  }

  const platform = getPlatform();
  const arch = getArch();
  const packageName = `@vibium/${platform}-${arch}`;
  const binaryName = platform === 'win32' ? 'clicker.exe' : 'clicker';

  // 2. Check platform-specific npm package
  try {
    // Use createRequire to make this work in both CJS and ESM
    // In ESM, use import.meta.url; in CJS, construct from __filename
    let currentModulePath: string;
    if (typeof __filename !== 'undefined') {
      // CommonJS
      currentModulePath = __filename;
    } else {
      // ESM - TypeScript won't compile this, but it works at runtime in ESM
      // @ts-ignore - import.meta.url is available in ESM at runtime
      currentModulePath = fileURLToPath(import.meta.url);
    }
    const require = createRequire(currentModulePath);
    const packagePath = require.resolve(`${packageName}/package.json`);
    const packageDir = path.dirname(packagePath);
    const binaryPath = path.join(packageDir, 'bin', binaryName);

    if (fs.existsSync(binaryPath)) {
      return binaryPath;
    }
  } catch {
    // Package not installed, continue to fallback
  }

  // 3. Check local development paths (relative to cwd)
  const localPaths = [
    // From vibium/ root
    path.resolve(process.cwd(), 'clicker', 'bin', binaryName),
    // From clients/javascript/
    path.resolve(process.cwd(), '..', '..', 'clicker', 'bin', binaryName),
  ];

  for (const localPath of localPaths) {
    if (fs.existsSync(localPath)) {
      return localPath;
    }
  }

  throw new Error(
    `Could not find clicker binary. ` +
    `Set CLICKER_PATH environment variable or install ${packageName}`
  );
}
