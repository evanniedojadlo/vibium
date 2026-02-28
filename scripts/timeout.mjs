#!/usr/bin/env node
// Run a command with a timeout. Kills the command if it exceeds the limit.
// Usage: node scripts/timeout.mjs SECONDS command [args...]
import { spawn } from 'node:child_process';

const [seconds, ...cmdArgs] = process.argv.slice(2);
const limit = Number(seconds) * 1000;

const child = spawn(cmdArgs[0], cmdArgs.slice(1), {
  stdio: 'inherit',
  shell: process.platform === 'win32',
});

const timer = setTimeout(() => {
  process.stderr.write(`\n--- TIMEOUT: killed after ${seconds}s ---\n`);
  child.kill();
}, limit);

child.on('exit', (code) => {
  clearTimeout(timer);
  process.exit(code ?? 1);
});
