# TaskGraph Landing Page Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a single-page static landing page in `site/` that presents TaskGraph with an Inbox-inspired visual style and clear install/demo calls to action.

**Architecture:** Use one self-contained `site/index.html` file with Tailwind CDN plus inline CSS variables and small animation helpers. Keep content grounded in `README.md` and `docs/DESIGN.md`, with no build step and no external assets required beyond Tailwind CDN.

**Tech Stack:** HTML, Tailwind CSS CDN, inline CSS, POSIX shell verification script

---

### Task 1: Add a minimal verification script

**Files:**
- Create: `scripts/test-site-landing-page.sh`

**Step 1: Write the failing test**

Create a shell script that checks:
- `site/index.html` exists
- the file contains `TaskGraph`
- the file contains the install command `curl -fsSL`
- the file contains a demo anchor such as `id="demo"`

**Step 2: Run test to verify it fails**

Run: `bash scripts/test-site-landing-page.sh`
Expected: FAIL because `site/index.html` does not exist yet.

**Step 3: Write minimal implementation**

Create `site/index.html` with those required sections and strings.

**Step 4: Run test to verify it passes**

Run: `bash scripts/test-site-landing-page.sh`
Expected: PASS with a short success message.

### Task 2: Build the single-page landing page

**Files:**
- Create: `site/index.html`

**Step 1: Implement page structure**

Include:
- metadata and title
- warm paper-style background and mono typography inspired by Inbox
- hero with install primary CTA and on-page demo secondary CTA
- product value sections grounded in README and design doc
- CLI demo section with example commands
- footer with docs and GitHub links

**Step 2: Keep implementation self-contained**

Use inline SVG or pure HTML/CSS for decorative elements. Avoid requiring image assets or a JS build.

**Step 3: Verify output quality**

Run HTML sanity checks and inspect that the page copy matches current product capabilities.

### Task 3: Verify completion

**Files:**
- Verify: `site/index.html`
- Verify: `scripts/test-site-landing-page.sh`

**Step 1: Run page verification**

Run: `bash scripts/test-site-landing-page.sh`
Expected: PASS

**Step 2: Run HTML parse check**

Run: `python3 - <<'PY'`
from html.parser import HTMLParser
HTMLParser().feed(open("site/index.html", "r", encoding="utf-8").read())
print("html-parse-ok")
PY

Expected: `html-parse-ok`
