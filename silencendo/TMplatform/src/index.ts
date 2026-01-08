#!/usr/bin/env node

import 'dotenv/config';
import { CLIInterface } from './cli/repl';

async function main() {
  const cli = new CLIInterface();
  await cli.start();
}

// Handle uncaught exceptions
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
  process.exit(1);
});

process.on('unhandledRejection', (reason) => {
  console.error('Unhandled Rejection:', reason);
  process.exit(1);
});

// Run the application
main().catch(console.error);
