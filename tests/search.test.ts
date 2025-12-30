import { describe, it, expect } from "vitest";
import { scoreMatch } from "../src/search";

describe("scoreMatch", () => {
  it("scores higher for repeated tokens", () => {
    expect(scoreMatch("alpha beta beta", "beta")).toBeGreaterThan(
      scoreMatch("alpha beta", "alpha")
    );
  });
});
