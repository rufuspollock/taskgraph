import { describe, it, expect, vi } from "vitest";
import { runCli } from "../src/cli";

describe("cli", () => {
  it("shows help", async () => {
    const code = await runCli(["-h"]);
    expect(code).toBe(0);
  });

  it("warns when index is missing", async () => {
    const err = vi.spyOn(console, "error").mockImplementation(() => {});
    const code = await runCli(["query", "meeting", "--index", "does-not-exist.json"]);
    expect(code).toBe(1);
    expect(err).toHaveBeenCalled();
    err.mockRestore();
  });
});
