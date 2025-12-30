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
