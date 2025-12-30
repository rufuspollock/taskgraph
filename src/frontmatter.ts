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
