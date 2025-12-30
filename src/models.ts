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
