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
