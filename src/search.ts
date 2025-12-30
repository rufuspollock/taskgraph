export function scoreMatch(text: string, query: string) {
  const tokens = query.toLowerCase().split(/\s+/).filter(Boolean);
  const hay = text.toLowerCase();
  let score = 0;
  for (const token of tokens) {
    const matches = hay.split(token).length - 1;
    score += matches;
  }
  return score;
}

export function searchNodes(nodes: any[], query: string, limit = 10) {
  const scored = nodes
    .map((n) => ({ score: scoreMatch(n.searchText ?? "", query), node: n }))
    .filter((s) => s.score > 0)
    .sort((a, b) => b.score - a.score);
  return scored.slice(0, limit).map((s) => s.node);
}
