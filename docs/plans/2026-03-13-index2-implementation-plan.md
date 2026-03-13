# Index2 Landing Page Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a second single-page static landing page at `site/index2.html` with a cleaner ScreenshotIt-inspired design, animated graph hero, install section, quickstart, and vision section.

**Architecture:** Build one self-contained HTML file using Tailwind CDN plus small inline CSS animations. Keep the layout editorial and minimal, with a CSS-only hero animation that demonstrates moving from a high-level branch to a concrete leaf task.

**Tech Stack:** HTML, Tailwind CSS CDN, inline CSS, POSIX shell verification script

---

### Task 1: Add a minimal verification script for `index2.html`

**Files:**
- Modify: `scripts/test-site-landing-page.sh`

**Step 1: Write the failing test**

Extend the script to check:
- `site/index2.html` exists
- the file contains `TaskGraph`
- the file contains `id="quickstart"`
- the file contains `id="vision"`
- the file contains the install command `curl -fsSL`

**Step 2: Run test to verify it fails**

Run: `bash scripts/test-site-landing-page.sh`
Expected: FAIL because `site/index2.html` does not exist yet.

**Step 3: Write minimal implementation**

Create `site/index2.html` with those required sections and strings.

**Step 4: Run test to verify it passes**

Run: `bash scripts/test-site-landing-page.sh`
Expected: PASS with checks for both landing pages.

### Task 2: Build the simplified long-form landing page

**Files:**
- Create: `site/index2.html`

**Step 1: Implement hero**

Include:
- minimal centered brand and headline
- short explanatory subhead
- install and quickstart CTAs
- CSS-only animated task graph with highlighted branch-to-leaf movement

**Step 2: Implement lower sections**

Include:
- install
- quickstart
- vision

Use distilled content from:
- `README.md`
- `docs/DESIGN.md`
- `docs/task-graph-vision-essay-2026-03-02.md`

**Step 3: Keep styling minimal**

Use thin dividers, mono typography, flat backgrounds, and restrained accent color. Avoid the heavier card-based look from `site/index.html`.

### Task 3: Verify completion

**Files:**
- Verify: `site/index2.html`
- Verify: `scripts/test-site-landing-page.sh`

**Step 1: Run verification**

Run: `bash scripts/test-site-landing-page.sh`
Expected: PASS

**Step 2: Run HTML parse sanity check**

Run: `python3 - <<'PY'`
from html.parser import HTMLParser
HTMLParser().feed(open("site/index2.html", "r", encoding="utf-8").read())
print("html-parse-ok")
PY

Expected: `html-parse-ok`
