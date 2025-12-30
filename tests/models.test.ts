import { describe, it, expect } from "vitest";
import { buildNodeId } from "../src/models";

describe("buildNodeId", () => {
  it("is stable", () => {
    const id1 = buildNodeId("fixtures/a.md", ["Section"], 12);
    const id2 = buildNodeId("fixtures/a.md", ["Section"], 12);
    expect(id1).toBe(id2);
  });
});
