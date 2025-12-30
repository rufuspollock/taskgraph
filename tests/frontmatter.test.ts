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
