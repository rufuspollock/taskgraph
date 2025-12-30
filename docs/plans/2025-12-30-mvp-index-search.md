# TaskGraph MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a minimal CLI that indexes Markdown tasks into a JSON database and supports live, interactive search.

**Architecture:** A small Python CLI parses Markdown files into flat task nodes with parent pointers, writes `data/index.json`, and provides a query command that scores matches against a precomputed `search_text` field. Interactive search runs in a simple raw-terminal loop that re-renders results on each keystroke.

**Tech Stack:** Python 3 (stdlib only: argparse, json, re, hashlib, pathlib, termios, tty, textwrap)

### Task 1: Define node schema + ID helpers

**Files:**
- Create: `taskgraph/models.py`
- Test: `tests/test_models.py`

**Step 1: Write the failing test**

```python
from taskgraph.models import build_node_id, TaskNode

def test_build_node_id_stable():
    id1 = build_node_id("fixtures/a.md", ["Section"], 12)
    id2 = build_node_id("fixtures/a.md", ["Section"], 12)
    assert id1 == id2
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_models.py -q`
Expected: FAIL with import error or missing function

**Step 3: Write minimal implementation**

```python
import hashlib
from dataclasses import dataclass
from typing import Dict, Optional, List

@dataclass
class TaskNode:
    id: str
    kind: str
    title: str
    state: str
    path: str
    line: int
    parent_id: Optional[str]
    context: str
    search_text: str
    frontmatter: Dict[str, str]


def build_node_id(path: str, heading_path: List[str], line: int) -> str:
    raw = "::".join([path] + heading_path + [str(line)])
    return hashlib.sha1(raw.encode("utf-8")).hexdigest()
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_models.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/models.py tests/test_models.py
git commit -m "feat: add task node model and id helper"
```

### Task 2: Parse front matter (simple key/value)

**Files:**
- Create: `taskgraph/frontmatter.py`
- Test: `tests/test_frontmatter.py`

**Step 1: Write the failing test**

```python
from taskgraph.frontmatter import parse_frontmatter

def test_parse_frontmatter_simple():
    text = """---\ncreated: 2024-08-13\ncompleted: \nkind: product\n---\n\nBody"""
    fm, body = parse_frontmatter(text)
    assert fm["created"] == "2024-08-13"
    assert fm["completed"] == ""
    assert fm["kind"] == "product"
    assert body.startswith("Body")
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_frontmatter.py -q`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```python
def parse_frontmatter(text: str):
    if not text.startswith("---\n"):
        return {}, text
    parts = text.split("\n---\n", 1)
    if len(parts) != 2:
        return {}, text
    fm_block, body = parts
    lines = fm_block.splitlines()[1:]
    data = {}
    for line in lines:
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        data[key.strip()] = value.strip()
    return data, body.lstrip("\n")
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_frontmatter.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/frontmatter.py tests/test_frontmatter.py
git commit -m "feat: add simple front matter parser"
```

### Task 3: Markdown scanner for headings + checklists

**Files:**
- Create: `taskgraph/indexer.py`
- Test: `tests/test_indexer.py`

**Step 1: Write the failing test**

```python
from taskgraph.indexer import index_markdown_text

def test_index_builds_hierarchy():
    text = """---\nkind: project\n---\n\n# Alpha\n\n## Build\n- [ ] Task one\n- [x] Task two\n"""
    nodes = index_markdown_text(text, "fixtures/alpha.md")
    kinds = [n.kind for n in nodes]
    assert "file" in kinds
    assert "heading" in kinds
    assert "checklist" in kinds
    task_one = [n for n in nodes if n.title == "Task one"][0]
    assert task_one.state == "open"
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_indexer.py -q`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```python
import re
from typing import List
from taskgraph.frontmatter import parse_frontmatter
from taskgraph.models import TaskNode, build_node_id

CHECKBOX_RE = re.compile(r"^\s*-\s*\[( |x|X)\]\s+(.*)$")
HEADING_RE = re.compile(r"^(#{1,6})\s+(.*)$")


def index_markdown_text(text: str, path: str) -> List[TaskNode]:
    frontmatter, body = parse_frontmatter(text)
    file_title = path.rsplit("/", 1)[-1].rsplit(".", 1)[0].replace("-", " ")
    nodes: List[TaskNode] = []
    heading_stack: List[str] = []
    heading_ids: List[str] = []

    file_id = build_node_id(path, [], 0)
    file_context = file_title
    file_search = build_search_text(file_context, frontmatter)
    nodes.append(TaskNode(
        id=file_id,
        kind="file",
        title=file_title,
        state="unknown",
        path=path,
        line=0,
        parent_id=None,
        context=file_context,
        search_text=file_search,
        frontmatter=frontmatter,
    ))

    for idx, line in enumerate(body.splitlines(), start=1):
        heading_match = HEADING_RE.match(line)
        if heading_match:
            level = len(heading_match.group(1))
            title = heading_match.group(2).strip()
            while len(heading_stack) >= level:
                heading_stack.pop()
                heading_ids.pop()
            heading_stack.append(title)
            node_id = build_node_id(path, heading_stack, idx)
            parent_id = heading_ids[-1] if heading_ids else file_id
            heading_ids.append(node_id)
            context = " > ".join([file_title] + heading_stack)
            search_text = build_search_text(context, frontmatter)
            nodes.append(TaskNode(
                id=node_id,
                kind="heading",
                title=title,
                state="unknown",
                path=path,
                line=idx,
                parent_id=parent_id,
                context=context,
                search_text=search_text,
                frontmatter=frontmatter,
            ))
            continue

        checkbox_match = CHECKBOX_RE.match(line)
        if checkbox_match:
            checked = checkbox_match.group(1).lower() == "x"
            title = checkbox_match.group(2).strip()
            state = "closed" if checked else "open"
            context = " > ".join([file_title] + heading_stack + [title])
            search_text = build_search_text(context, frontmatter)
            parent_id = heading_ids[-1] if heading_ids else file_id
            node_id = build_node_id(path, heading_stack + [title], idx)
            nodes.append(TaskNode(
                id=node_id,
                kind="checklist",
                title=title,
                state=state,
                path=path,
                line=idx,
                parent_id=parent_id,
                context=context,
                search_text=search_text,
                frontmatter=frontmatter,
            ))

    return nodes


def build_search_text(context: str, frontmatter: dict) -> str:
    parts = [context]
    for key, value in frontmatter.items():
        parts.append(f"{key}:{value}")
    return " ".join(parts).strip()
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_indexer.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/indexer.py tests/test_indexer.py
git commit -m "feat: parse headings and checklists into nodes"
```

### Task 4: Build JSON index from directory

**Files:**
- Modify: `taskgraph/indexer.py`
- Create: `taskgraph/storage.py`
- Test: `tests/test_storage.py`

**Step 1: Write the failing test**

```python
import json
from pathlib import Path
from taskgraph.storage import write_index


def test_write_index(tmp_path: Path):
    nodes = []
    out_path = tmp_path / "index.json"
    write_index(out_path, nodes)
    data = json.loads(out_path.read_text())
    assert "nodes" in data
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_storage.py -q`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```python
import json
from pathlib import Path
from typing import List
from taskgraph.models import TaskNode


def write_index(path: Path, nodes: List[TaskNode]) -> None:
    payload = {"nodes": [n.__dict__ for n in nodes]}
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2), encoding="utf-8")


def read_index(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_storage.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/storage.py tests/test_storage.py
git commit -m "feat: write json index"
```

### Task 5: Search scoring + non-interactive query

**Files:**
- Create: `taskgraph/search.py`
- Test: `tests/test_search.py`

**Step 1: Write the failing test**

```python
from taskgraph.search import score_match


def test_score_match_counts_tokens():
    text = "alpha beta beta"
    assert score_match(text, "beta") > score_match(text, "alpha")
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_search.py -q`
Expected: FAIL (function missing)

**Step 3: Write minimal implementation**

```python
import re
from typing import List
from taskgraph.models import TaskNode


def score_match(search_text: str, query: str) -> int:
    tokens = [t for t in re.split(r"\s+", query.lower()) if t]
    hay = search_text.lower()
    score = 0
    for token in tokens:
        score += hay.count(token)
    return score


def search_nodes(nodes: List[TaskNode], query: str, limit: int = 10) -> List[TaskNode]:
    scored = [(score_match(n.search_text, query), n) for n in nodes]
    scored = [s for s in scored if s[0] > 0]
    scored.sort(key=lambda item: item[0], reverse=True)
    return [n for _, n in scored[:limit]]
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_search.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/search.py tests/test_search.py
git commit -m "feat: basic search scoring"
```

### Task 6: CLI wiring + interactive query

**Files:**
- Create: `taskgraph/cli.py`
- Create: `taskgraph/interactive.py`
- Create: `taskgraph/__main__.py`
- Test: `tests/test_cli_smoke.py`

**Step 1: Write the failing test**

```python
import subprocess
import sys


def test_cli_help():
    result = subprocess.run([sys.executable, "-m", "taskgraph", "-h"], capture_output=True, text=True)
    assert result.returncode == 0
    assert "index" in result.stdout
```

**Step 2: Run test to verify it fails**

Run: `python -m pytest tests/test_cli_smoke.py -q`
Expected: FAIL (module missing)

**Step 3: Write minimal implementation**

```python
# taskgraph/cli.py
import argparse
from pathlib import Path
from taskgraph.indexer import index_markdown_text
from taskgraph.storage import write_index
from taskgraph.search import search_nodes
from taskgraph.interactive import run_interactive


def load_nodes(index_path: Path):
    import json
    data = json.loads(index_path.read_text())
    return data["nodes"]


def main(argv=None):
    parser = argparse.ArgumentParser(prog="taskgraph")
    sub = parser.add_subparsers(dest="command", required=True)

    p_index = sub.add_parser("index")
    p_index.add_argument("source")
    p_index.add_argument("--out", default="data/index.json")

    p_query = sub.add_parser("query")
    p_query.add_argument("query")
    p_query.add_argument("--index", default="data/index.json")
    p_query.add_argument("--limit", type=int, default=10)
    p_query.add_argument("--interactive", action="store_true")

    args = parser.parse_args(argv)

    if args.command == "index":
        from pathlib import Path
        source = Path(args.source)
        files = sorted(source.rglob("*.md"))
        nodes = []
        for file_path in files:
            nodes.extend(index_markdown_text(file_path.read_text(encoding="utf-8"), str(file_path)))
        write_index(Path(args.out), nodes)
        print(f"Indexed {len(nodes)} nodes to {args.out}")
        return 0

    if args.command == "query":
        from taskgraph.models import TaskNode
        from taskgraph.storage import read_index
        if args.interactive:
            return run_interactive(Path(args.index), args.limit)
        data = read_index(Path(args.index))
        nodes = [TaskNode(**n) for n in data["nodes"]]
        results = search_nodes(nodes, args.query, args.limit)
        for n in results:
            print(f"[{n.state}] {n.context} ({n.path}:{n.line})")
        return 0


# taskgraph/interactive.py
import json
import sys
import termios
import tty
from pathlib import Path
from taskgraph.models import TaskNode
from taskgraph.search import search_nodes


def run_interactive(index_path: Path, limit: int) -> int:
    data = json.loads(index_path.read_text())
    nodes = [TaskNode(**n) for n in data["nodes"]]
    buf = ""
    fd = sys.stdin.fileno()
    old = termios.tcgetattr(fd)
    try:
        tty.setcbreak(fd)
        while True:
            print("\x1b[2J\x1b[H", end="")
            print(f"Query: {buf}")
            if buf:
                results = search_nodes(nodes, buf, limit)
                for n in results:
                    print(f"[{n.state}] {n.context} ({n.path}:{n.line})")
            ch = sys.stdin.read(1)
            if ch in ("\n", "\r"):
                return 0
            if ch == "\x03":
                return 1
            if ch == "\x7f":
                buf = buf[:-1]
            else:
                buf += ch
    finally:
        termios.tcsetattr(fd, termios.TCSADRAIN, old)


# taskgraph/__main__.py
from taskgraph.cli import main

if __name__ == "__main__":
    raise SystemExit(main())
```

**Step 4: Run test to verify it passes**

Run: `python -m pytest tests/test_cli_smoke.py -q`
Expected: PASS

**Step 5: Commit**

```bash
git add taskgraph/cli.py taskgraph/interactive.py taskgraph/__main__.py tests/test_cli_smoke.py
git commit -m "feat: add cli with interactive query"
```

### Task 7: README usage + workflow notes

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

```python
# No tests for docs
```

**Step 2: Run test to verify it fails**

Run: `true`
Expected: PASS

**Step 3: Write minimal implementation**

Add usage examples:
- `python -m taskgraph index fixtures/rufus-projects`
- `python -m taskgraph query "meeting"`
- `python -m taskgraph query "meeting" --interactive`

**Step 4: Run test to verify it passes**

Run: `true`
Expected: PASS

**Step 5: Commit**

```bash
git add README.md
git commit -m "docs: add mvp usage"
```
