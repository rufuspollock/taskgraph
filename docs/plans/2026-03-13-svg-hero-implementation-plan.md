# SVG Hero Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the fragile HTML/CSS hero graph in `site/index.html` with an inline SVG graph whose connectors align exactly with the labeled nodes.

**Architecture:** Keep the rest of the page unchanged. Rebuild the hero visualization as a single SVG with fixed coordinates, lines and pills in one coordinate space, and CSS animation on SVG strokes and node fills.

**Tech Stack:** HTML, inline SVG, inline CSS, POSIX shell verification script

---

### Task 1: Extend verification for SVG hero

**Files:**
- Modify: `scripts/test-site-landing-page.sh`

**Step 1: Write the failing test**

Add checks for:
- `<svg`
- `graph-svg`
- absence of the old `class="edge active"` HTML-line hero implementation

**Step 2: Run test to verify it fails**

Run: `bash scripts/test-site-landing-page.sh`
Expected: FAIL because `site/index.html` still uses the old hero structure.

### Task 2: Replace the hero visualization

**Files:**
- Modify: `site/index.html`

**Step 1: Remove fragile line-box implementation**

Delete the old absolutely positioned node and edge layout.

**Step 2: Add inline SVG hero**

Implement:
- one root project
- labeled pill nodes
- exact connector lines
- subtle line and node highlight animation

### Task 3: Verify completion

**Files:**
- Verify: `site/index.html`
- Verify: `scripts/test-site-landing-page.sh`

**Step 1: Run verification**

Run: `bash scripts/test-site-landing-page.sh`
Expected: PASS

**Step 2: Run HTML parse check**

Run: `python3 - <<'PY'`
from html.parser import HTMLParser
HTMLParser().feed(open("site/index.html", "r", encoding="utf-8").read())
print("html-parse-ok")
PY

Expected: `html-parse-ok`
