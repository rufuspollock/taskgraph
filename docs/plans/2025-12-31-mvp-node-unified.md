# TaskGraph MVP Node/TS Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a TypeScript CLI that indexes Markdown tasks into a JSON database and supports interactive search updating on each keystroke.

**Architecture:** A small Node/TS CLI uses `unified` + `remark-parse` to parse Markdown, walks the AST to emit file/heading/checklist task nodes, writes `data/index.json`, and provides a query command that scores matches against a precomputed `searchText` field. Interactive query mode uses raw TTY input to redraw results on each keystroke.

**Tech Stack:** Node.js, TypeScript, `tsx`, `vitest`, `unified`, `remark-parse`, `remark-frontmatter`, `mdast-util-frontmatter`

### Task 1: Initialize Node/TS tooling

**Files:**
- Create: `package.json`
- Create: `tsconfig.json`
- Create: `vitest.config.ts`
- Create: `src/`
- Create: `tests/`

**Step 1: Write the failing test**

```ts
// tests/smoke.test.ts
import { describe, it, expect } from "vitest";

describe("smoke", () => {
  it("runs", () => {
    expect(true).toBe(true);
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test`
Expected: FAIL (package.json/scripts missing)

**Step 3: Write minimal implementation**

Create `package.json`:

```json
{
  "name": "taskgraph",
  "private": true,
  "type": "module",
  "scripts": {
    "test": "vitest run",
    "dev": "tsx src/cli.ts"
  },
  "devDependencies": {
    "tsx": "^4.19.2",
    "typescript": "^5.7.3",
    "vitest": "^2.1.3"
  },
  "dependencies": {
    "unified": "^11.0.5",
    "remark-parse": "^11.0.0",
    "remark-frontmatter": "^5.0.0",
    "mdast-util-frontmatter": "^2.0.0",
    "unist-util-visit": "^5.0.0",
    "remark-gfm": "^4.0.1"
  }
}
```

Create `tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ES2022",
    "moduleResolution": "bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "outDir": "dist"
  },
  "include": ["src", "tests"]
}
```

Create `vitest.config.ts`:

```ts
import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "node",
  },
});
```

**Step 4: Run test to verify it passes**

Run: `pnpm install && pnpm test`
Expected: PASS

**Step 5: Commit**

```bash
git add package.json tsconfig.json vitest.config.ts tests/smoke.test.ts pnpm-lock.yaml
git commit -m "chore: set up tsx and vitest"
```

### Task 2: Define node schema + ID helper

**Files:**
- Create: `src/models.ts`
- Test: `tests/models.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { buildNodeId } from "../src/models";

describe("buildNodeId", () => {
  it("is stable", () => {
    const id1 = buildNodeId("fixtures/a.md", ["Section"], 12);
    const id2 = buildNodeId("fixtures/a.md", ["Section"], 12);
    expect(id1).toBe(id2);
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/models.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/models.ts
import crypto from "node:crypto";

export type TaskKind = "file" | "heading" | "checklist";
export type TaskState = "open" | "closed" | "unknown";

export interface TaskNode {
  id: string;
  kind: TaskKind;
  title: string;
  state: TaskState;
  path: string;
  line: number;
  parentId: string | null;
  context: string;
  searchText: string;
  frontmatter: Record<string, string>;
}

export function buildNodeId(path: string, headingPath: string[], line: number): string {
  const raw = [path, ...headingPath, String(line)].join("::");
  return crypto.createHash("sha1").update(raw).digest("hex");
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/models.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/models.ts tests/models.test.ts
git commit -m "feat: add task node model"
```

### Task 3: Front matter parsing (key/value)

**Files:**
- Create: `src/frontmatter.ts`
- Test: `tests/frontmatter.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { parseFrontmatter } from "../src/frontmatter";

describe("parseFrontmatter", () => {
  it("extracts simple key values", () => {
    const text = "---\ncreated: 2024-08-13\ncompleted: \nkind: product\n---\n\nBody";
    const { frontmatter, body } = parseFrontmatter(text);
    expect(frontmatter.created).toBe("2024-08-13");
    expect(frontmatter.completed).toBe("");
    expect(frontmatter.kind).toBe("product");
    expect(body.startsWith("Body")).toBe(true);
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/frontmatter.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/frontmatter.ts
export function parseFrontmatter(text: string) {
  if (!text.startsWith("---\n")) {
    return { frontmatter: {}, body: text };
  }
  const parts = text.split("\n---\n", 2);
  if (parts.length !== 2) {
    return { frontmatter: {}, body: text };
  }
  const [raw, body] = parts;
  const lines = raw.split("\n").slice(1);
  const frontmatter: Record<string, string> = {};
  for (const line of lines) {
    const idx = line.indexOf(":");
    if (idx === -1) continue;
    const key = line.slice(0, idx).trim();
    const value = line.slice(idx + 1).trim();
    frontmatter[key] = value;
  }
  return { frontmatter, body: body.replace(/^\n/, "") };
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/frontmatter.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/frontmatter.ts tests/frontmatter.test.ts
git commit -m "feat: parse frontmatter"
```

### Task 4: Index markdown AST into nodes

**Files:**
- Create: `src/indexer.ts`
- Test: `tests/indexer.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { indexMarkdown } from "../src/indexer";

describe("indexMarkdown", () => {
  it("builds hierarchy and state", () => {
    const text = "---\nkind: project\n---\n\n# Alpha\n\n## Build\n- [ ] Task one\n- [x] Task two\n";
    const nodes = indexMarkdown(text, "fixtures/alpha.md");
    const kinds = nodes.map((n) => n.kind);
    expect(kinds).toContain("file");
    expect(kinds).toContain("heading");
    expect(kinds).toContain("checklist");
    const taskOne = nodes.find((n) => n.title === "Task one");
    expect(taskOne?.state).toBe("open");
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/indexer.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/indexer.ts
import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import { visit } from "unist-util-visit";
import { parseFrontmatter } from "./frontmatter";
import { buildNodeId, TaskNode, TaskState } from "./models";

function fileTitleFromPath(path: string) {
  return path.split("/").pop()?.replace(/\.md$/i, "").replace(/-/g, " ") ?? path;
}

function buildSearchText(context: string, frontmatter: Record<string, string>) {
  const parts = [context];
  for (const [key, value] of Object.entries(frontmatter)) {
    parts.push(`${key}:${value}`);
  }
  return parts.join(" ").trim();
}

function collectText(node: any): string {
  if (!node) return "";
  if (node.type === "text") return node.value ?? "";
  if (Array.isArray(node.children)) {
    return node.children.map(collectText).join("");
  }
  return "";
}

function checkboxState(node: any): TaskState | null {
  if (node?.checked === true) return "closed";
  if (node?.checked === false) return "open";
  return null;
}

export function indexMarkdown(text: string, path: string): TaskNode[] {
  const { frontmatter, body } = parseFrontmatter(text);
  const fileTitle = fileTitleFromPath(path);
  const nodes: TaskNode[] = [];
  const headingStack: string[] = [];
  const headingIds: string[] = [];

  const fileId = buildNodeId(path, [], 0);
  nodes.push({
    id: fileId,
    kind: "file",
    title: fileTitle,
    state: "unknown",
    path,
    line: 0,
    parentId: null,
    context: fileTitle,
    searchText: buildSearchText(fileTitle, frontmatter),
    frontmatter,
  });

  const tree = unified()
    .use(remarkParse)
    .use(remarkFrontmatter, ["yaml"])
    .use(remarkGfm)
    .parse(body);

  visit(tree, (node: any) => {
    if (node.type === "heading") {
      const title = collectText(node).trim();
      const level = node.depth ?? 1;
      while (headingStack.length >= level) {
        headingStack.pop();
        headingIds.pop();
      }
      headingStack.push(title);
      const id = buildNodeId(path, headingStack, node.position?.start?.line ?? 0);
      const parentId = headingIds.length ? headingIds[headingIds.length - 1] : fileId;
      headingIds.push(id);
      const context = [fileTitle, ...headingStack].join(" > ");
      nodes.push({
        id,
        kind: "heading",
        title,
        state: "unknown",
        path,
        line: node.position?.start?.line ?? 0,
        parentId,
        context,
        searchText: buildSearchText(context, frontmatter),
        frontmatter,
      });
      return;
    }

    if (node.type === "listItem") {
      const state = checkboxState(node);
      if (!state) return;
      const title = collectText(node).trim();
      const context = [fileTitle, ...headingStack, title].join(" > ");
      const parentId = headingIds.length ? headingIds[headingIds.length - 1] : fileId;
      const id = buildNodeId(path, [...headingStack, title], node.position?.start?.line ?? 0);
      nodes.push({
        id,
        kind: "checklist",
        title,
        state,
        path,
        line: node.position?.start?.line ?? 0,
        parentId,
        context,
        searchText: buildSearchText(context, frontmatter),
        frontmatter,
      });
    }
  });

  return nodes;
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/indexer.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/indexer.ts tests/indexer.test.ts
git commit -m "feat: index markdown tasks"
```

### Task 5: JSON storage

**Files:**
- Create: `src/storage.ts`
- Test: `tests/storage.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { writeIndex, readIndex } from "../src/storage";
import { mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const tmp = mkdtempSync(join(tmpdir(), "taskgraph-"));

describe("storage", () => {
  it("writes and reads index", () => {
    const path = join(tmp, "index.json");
    writeIndex(path, { nodes: [] });
    const data = readIndex(path);
    expect(data.nodes).toBeTruthy();
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/storage.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/storage.ts
import { readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { dirname } from "node:path";

export function writeIndex(path: string, data: unknown) {
  mkdirSync(dirname(path), { recursive: true });
  writeFileSync(path, JSON.stringify(data, null, 2), "utf-8");
}

export function readIndex(path: string) {
  return JSON.parse(readFileSync(path, "utf-8"));
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/storage.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/storage.ts tests/storage.test.ts
git commit -m "feat: add json storage"
```

### Task 6: Search scoring

**Files:**
- Create: `src/search.ts`
- Test: `tests/search.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { scoreMatch } from "../src/search";

describe("scoreMatch", () => {
  it("scores higher for repeated tokens", () => {
    expect(scoreMatch("alpha beta beta", "beta")).toBeGreaterThan(scoreMatch("alpha beta", "alpha"));
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/search.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/search.ts
export function scoreMatch(text: string, query: string) {
  const tokens = query.toLowerCase().split(/\s+/).filter(Boolean);
  const hay = text.toLowerCase();
  let score = 0;
  for (const token of tokens) {
    const matches = hay.split(token).length - 1;
    score += matches;
  }
  return score;
}

export function searchNodes(nodes: any[], query: string, limit = 10) {
  const scored = nodes
    .map((n) => ({ score: scoreMatch(n.searchText ?? "", query), node: n }))
    .filter((s) => s.score > 0)
    .sort((a, b) => b.score - a.score);
  return scored.slice(0, limit).map((s) => s.node);
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/search.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/search.ts tests/search.test.ts
git commit -m "feat: add search scoring"
```

### Task 7: CLI index + query + interactive mode

**Files:**
- Create: `src/cli.ts`
- Create: `src/interactive.ts`
- Create: `src/types.ts`
- Test: `tests/cli.test.ts`

**Step 1: Write the failing test**

```ts
import { describe, it, expect } from "vitest";
import { runCli } from "../src/cli";

describe("cli", () => {
  it("shows help", () => {
    const code = runCli(["-h"]);
    expect(code).toBe(0);
  });
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm test -- tests/cli.test.ts`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```ts
// src/cli.ts
import { readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { indexMarkdown } from "./indexer";
import { writeIndex, readIndex } from "./storage";
import { searchNodes } from "./search";
import { runInteractive } from "./interactive";

export function runCli(argv = process.argv.slice(2)) {
  if (argv.includes("-h") || argv.includes("--help")) {
    console.log("taskgraph index <dir> [--out data/index.json]");
    console.log("taskgraph query <query> [--index data/index.json] [--limit 10] [--interactive]");
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
    if (interactive) return runInteractive(indexPath, limit);
    const data = readIndex(indexPath);
    const results = searchNodes(data.nodes ?? [], query, limit);
    for (const n of results) {
      console.log(`[${n.state}] ${n.context} (${n.path}:${n.line})`);
    }
    return 0;
  }
  console.error("Unknown command");
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

// src/interactive.ts
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

// src/types.ts
export interface IndexFile {
  nodes: any[];
}
```

**Step 4: Run test to verify it passes**

Run: `pnpm test -- tests/cli.test.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add src/cli.ts src/interactive.ts src/types.ts tests/cli.test.ts
git commit -m "feat: add cli commands"
```

### Task 8: Add bin entry + README usage

**Files:**
- Modify: `package.json`
- Modify: `README.md`

**Step 1: Write the failing test**

```ts
// no tests for docs
```

**Step 2: Run test to verify it fails**

Run: `true`
Expected: PASS

**Step 3: Write minimal implementation**

Add to `package.json`:

```json
"bin": {
  "taskgraph": "src/cli.ts"
}
```

Add usage to `README.md`:
- `pnpm install`
- `pnpm dev -- index fixtures/rufus-projects`
- `pnpm dev -- query "meeting"`
- `pnpm dev -- query "meeting" --interactive`

**Step 4: Run test to verify it passes**

Run: `true`
Expected: PASS

**Step 5: Commit**

```bash
git add package.json README.md
git commit -m "docs: add cli usage"
```
