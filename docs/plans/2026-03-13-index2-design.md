# Index2 Landing Page Design

## Goal

Create a second landing page variant at `site/index2.html` that is cleaner and more stripped down than `site/index.html`, using the visual language of `../screenshotit.app/src/homepage.ts` as the primary reference.

## Visual Direction

- Flat, near-white page background
- Monospace typography
- Thin borders and dividers
- Minimal color usage
- Centered hero with very limited chrome
- Editorial long-form layout rather than marketing-card layout

## Content Structure

1. Hero
   - State the core idea of TaskGraph in one sentence
   - Include a subtle animated task-graph visualization that shows selective descent from higher-level nodes to a leaf task
   - Primary CTA should go to install
   - Secondary CTA should go to quickstart
2. Install
   - Show the release install command
   - Keep explanation brief
3. Quickstart
   - Show compact CLI sequence grounded in README
   - Explain markdown source of truth and derived SQLite index
4. Vision
   - Use distilled language from `docs/task-graph-vision-essay-2026-03-02.md`
   - Emphasize movement between levels, selective descent, leaf tasks, and learning when a branch is blocked

## Interaction

- Lightweight CSS animation only
- Hero animation should feel explanatory, not decorative
- No heavy gradients, no large image assets, no JS framework

## Constraints

- Single static HTML file
- Tailwind via CDN is acceptable
- Keep claims aligned with README and design docs
