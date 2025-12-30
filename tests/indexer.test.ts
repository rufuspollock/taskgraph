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
