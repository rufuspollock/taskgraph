# Landing Page Graph Animation Design

## Goal

Turn the landing-page graph from a synchronized pulsing diagram into a short, legible story that explains TaskGraph's central idea: choose at a high level, descend to test a branch, learn when it fails, return, and find a concrete next action elsewhere.

## Recommended Direction

Use a single CSS-only SVG sequence with a visible navigation cursor and a changing caption. The sequence starts from three peer projects rather than one generic `Project` node:

1. Orient across the project landscape.
2. Choose `Site launch`.
3. Descend until the branch reports `needs breakdown`.
4. Return to the project choice.
5. Choose `TaskGraph`.
6. Descend through `Graph model` to `review notes`.
7. Hold on the actionable leaf before resetting.

This is preferable to a simple successful descent because the failed branch and return motion are what distinguish TaskGraph from a generic hierarchy or dependency graph.

## Visual Design

- Preserve the current near-white background, monospace typography, thin rules, and editorial restraint.
- Replace the background grid with a quieter dot field so the graph feels like a working canvas without competing with labels.
- Keep all graph structure faintly visible for orientation.
- Use dark ink for the moving selection and traversed path.
- Use a restrained warm red only for the branch that needs breakdown.
- Use a restrained green only for the final actionable leaf.
- Use compact rectangular nodes with slightly rounded corners rather than large pills; this makes the hierarchy feel more like a working tool and less like a generic flowchart.
- Add a small status line beneath the SVG. Its verb changes in sync with the graph: `Orient`, `Choose`, `Inspect`, `Back up`, `Descend`, `Do next`.

## Animation Architecture

The SVG contains a static base graph plus story-specific overlays:

- Base edges remain faint throughout.
- Story edges draw on using `stroke-dasharray` and `stroke-dashoffset`.
- Node groups receive stage-specific border, fill, label, and opacity changes.
- A small cursor dot moves between nodes to make direction explicit.
- Status captions occupy the same layout slot and cross-fade at stage boundaries.
- One shared 14-second duration keeps the timeline coherent.

The animation uses only CSS keyframes and inline SVG. It does not require JavaScript or external assets, so it remains compatible with the current Flowershow Markdown page.

## Responsive and Accessible Behavior

- Use a compact `720 × 500` SVG coordinate system so labels remain readable when scaled down on phones.
- Keep the status caption outside the SVG so it retains normal text size at every viewport width.
- Give the SVG a full narrative `aria-label`; decorative moving overlays are hidden from assistive technology.
- Under `prefers-reduced-motion: reduce`, disable looping and show the successful final path and leaf as a stable diagram.
- Pause the sequence while the graph is hovered or focused so visitors can inspect a frame.

## Verification

Update `scripts/test-site-landing-page.sh` to target `site/index.md`, which replaced the removed `site/index.html`. The checks should verify the failed-branch state, return path, successful leaf, synchronized caption structure, and reduced-motion fallback. Render the page in a browser at desktop and mobile widths and inspect the initial, failure, return, and success stages before shipping.
