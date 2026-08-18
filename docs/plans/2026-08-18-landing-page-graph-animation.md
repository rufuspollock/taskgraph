# Landing Page Graph Animation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the synchronized hero pulse with a responsive CSS-only narrative showing a failed descent, backtracking, and a successful descent to an actionable leaf.

**Architecture:** Keep the landing page as a self-contained Flowershow Markdown document. Build the hero from a static inline SVG base graph, animated overlay paths and node states, a cursor, and synchronized HTML captions; provide a static reduced-motion end state.

**Tech Stack:** Flowershow Markdown, HTML, inline SVG, CSS keyframes, POSIX shell verification

---

### Task 1: Correct and extend landing-page verification

**Files:**
- Modify: `scripts/test-site-landing-page.sh`

**Step 1: Write the failing checks**

Point the primary-page checks at `site/index.md`. Require markers for `story-path`, `story-cursor`, `story-caption`, `needs breakdown`, `review notes`, and `prefers-reduced-motion`. Remove checks that refer to the deleted `site/index.html`.

**Step 2: Run the test to verify it fails**

Run: `bash scripts/test-site-landing-page.sh`

Expected: FAIL because the existing graph lacks the narrative story markers.

**Step 3: Commit the red test**

Run: `git add scripts/test-site-landing-page.sh && git commit -m "test(site): specify narrative graph animation"`

### Task 2: Implement the narrative graph

**Files:**
- Modify: `site/index.md`

**Step 1: Replace the graph animation CSS**

Add a shared story duration, base graph styling, stage-specific node and edge keyframes, cursor motion, synchronized caption transitions, hover/focus pause behavior, and a `prefers-reduced-motion` fallback.

**Step 2: Replace the SVG and caption markup**

Create three project choices, a failed `Site launch → Launch plan → needs breakdown` branch, and a successful `TaskGraph → Graph model → review notes` branch. Add animated overlay paths, a cursor, and six status captions.

**Step 3: Run the focused verification**

Run: `bash scripts/test-site-landing-page.sh`

Expected: PASS for both landing-page variants.

**Step 4: Commit the implementation**

Run: `git add site/index.md && git commit -m "feat(site): animate graph navigation story"`

### Task 3: Render and refine

**Files:**
- Modify if needed: `site/index.md`

**Step 1: Render at desktop width**

Open the Flowershow page or an equivalent local render at approximately 1440px wide. Inspect the overview, failed branch, backtrack, and successful leaf frames.

**Step 2: Render at mobile width**

Inspect at approximately 390px wide. Verify node labels and the caption remain readable without horizontal scrolling.

**Step 3: Inspect reduced motion**

Emulate `prefers-reduced-motion: reduce` and verify the successful path appears as a stable final state.

**Step 4: Refine and rerun tests**

Adjust timing, spacing, contrast, or copy based on the renders. Run `bash scripts/test-site-landing-page.sh` after each refinement.

### Task 4: Final verification

**Files:**
- Verify: `site/index.md`
- Verify: `scripts/test-site-landing-page.sh`

**Step 1: Run site verification**

Run: `bash scripts/test-site-landing-page.sh`

Expected: both landing-page checks pass.

**Step 2: Run repository tests**

Run: `go test ./...`

Expected: all packages pass.

**Step 3: Check formatting and diff**

Run: `git diff --check && git status --short`

Expected: no whitespace errors and only intentional files changed.
