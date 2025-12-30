import { existsSync, readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { indexMarkdown } from "./indexer";
import { writeIndex, readIndex } from "./storage";
import { searchNodes } from "./search";
import { runInteractive } from "./interactive";

function printHelp() {
  console.log("taskgraph index <dir> [--out data/index.json]");
  console.log("taskgraph query <query> [--index data/index.json] [--limit 10] [--interactive]");
}

export function runCli(argv = process.argv.slice(2)) {
  if (argv.length === 0 || argv.includes("-h") || argv.includes("--help")) {
    printHelp();
    return 0;
  }
  const [command, ...rest] = argv;
  if (command === "index") {
    const source = rest[0];
    const outIdx = rest.indexOf("--out");
    const out = outIdx === -1 ? "data/index.json" : rest[outIdx + 1];
    const files = listMarkdownFiles(source);
    const nodes = files.flatMap((file) =>
      indexMarkdown(readFileSync(file, "utf-8"), file)
    );
    writeIndex(out, { nodes });
    console.log(`Indexed ${nodes.length} nodes to ${out}`);
    return 0;
  }
  if (command === "query") {
    const query = rest[0] ?? "";
    const idx = rest.indexOf("--index");
    const limitIdx = rest.indexOf("--limit");
    const interactive = rest.includes("--interactive");
    const indexPath = idx === -1 ? "data/index.json" : rest[idx + 1];
    const limit = limitIdx === -1 ? 10 : Number(rest[limitIdx + 1] ?? 10);
    if (!existsSync(indexPath)) {
      console.error(`Index not found at ${indexPath}`);
      console.error("Run: taskgraph index <dir> [--out data/index.json]");
      return 1;
    }
    if (interactive) return runInteractive(indexPath, limit);
    const data = readIndex(indexPath);
    const results = searchNodes(data.nodes ?? [], query, limit);
    for (const n of results) {
      console.log(`[${n.state}] ${n.context} (${n.path}:${n.line})`);
    }
    return 0;
  }
  console.error("Unknown command");
  printHelp();
  return 1;
}

function listMarkdownFiles(dir: string): string[] {
  const entries = readdirSync(dir, { withFileTypes: true });
  const out: string[] = [];
  for (const entry of entries) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...listMarkdownFiles(full));
    } else if (entry.isFile() && entry.name.endsWith(".md")) {
      out.push(full);
    }
  }
  return out;
}
