#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const tsxPath = resolve(root, "node_modules", ".bin", "tsx");
const cliPath = resolve(root, "src", "cli.ts");
const args = process.argv.slice(2);

const result = spawnSync(tsxPath, [cliPath, ...args], { stdio: "inherit" });
process.exit(result.status ?? 1);
