import { readIndex } from "./storage";
import { searchNodes } from "./search";

export function runInteractive(indexPath: string, limit: number) {
  const data = readIndex(indexPath);
  let buf = "";
  const stdin = process.stdin;
  if (stdin.setRawMode) {
    stdin.setRawMode(true);
  }
  stdin.resume();
  stdin.setEncoding("utf-8");
  stdin.on("data", (chunk) => {
    const ch = String(chunk);
    if (ch === "\u0003") process.exit(1);
    if (ch === "\r" || ch === "\n") process.exit(0);
    if (ch === "\u007f") buf = buf.slice(0, -1);
    else buf += ch;
    process.stdout.write("\x1b[2J\x1b[H");
    process.stdout.write(`Query: ${buf}\n`);
    const results = buf ? searchNodes(data.nodes ?? [], buf, limit) : [];
    for (const n of results) {
      process.stdout.write(`[${n.state}] ${n.context} (${n.path}:${n.line})\n`);
    }
  });
  return 0;
}
