#!/usr/bin/env node
// Find clicker binary from platform package and run it.
// Default: `vibium` → `clicker mcp` (MCP server mode)
// Known subcommands pass through directly.

const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');

// Subcommands that the binary handles directly
const KNOWN_SUBCOMMANDS = new Set([
  'mcp', 'navigate', 'click', 'type', 'find', 'find-all', 'screenshot',
  'text', 'html', 'url', 'title', 'eval', 'hover', 'scroll', 'select',
  'keys', 'wait', 'tabs', 'tab-new', 'tab-switch', 'tab-close', 'quit',
  'install', 'serve', 'version', 'paths', 'daemon', 'add-skill',
  'launch-test', 'ws-test', 'bidi-test', 'check-actionable',
  'help', 'completion',
]);

function getClickerPath() {
  const platform = os.platform();
  const arch = os.arch() === 'x64' ? 'x64' : 'arm64';
  const packageName = `@vibium/${platform}-${arch}`;
  const binaryName = platform === 'win32' ? 'clicker.exe' : 'clicker';

  try {
    const packagePath = require.resolve(`${packageName}/package.json`);
    return path.join(path.dirname(packagePath), 'bin', binaryName);
  } catch {
    console.error(`Could not find clicker binary for ${platform}-${arch}`);
    process.exit(1);
  }
}

const clickerPath = getClickerPath();
const userArgs = process.argv.slice(2);

// If no args or first arg is not a known subcommand, default to 'mcp'
const args = (userArgs.length === 0 || !KNOWN_SUBCOMMANDS.has(userArgs[0]))
  ? ['mcp', ...userArgs]
  : userArgs;

execFileSync(clickerPath, args, { stdio: 'inherit' });
