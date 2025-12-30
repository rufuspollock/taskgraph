import { describe, it, expect } from "vitest";
import { runCli } from "../src/cli";

describe("cli", () => {
  it("shows help", () => {
    const code = runCli(["-h"]);
    expect(code).toBe(0);
  });
});
