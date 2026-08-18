---
title: TaskGraph - Decide what to do next
description: "TaskGraph is a local-first CLI for moving from a high-level area of work to a concrete next task, with markdown as source of truth."
layout: plain
---

<style>
  @import url("https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600&display=swap");
  :root {
    --bg: #fafafa;
    --ink: #111111;
    --muted: #5f5f5f;
    --line: #e4e4e4;
    --soft: #f3f3f3;
    --accent: #111111;
  }
  html {
    scroll-behavior: smooth;
  }
  body {
    margin: 0;
    background: var(--bg);
    color: var(--ink);
    font-family: "IBM Plex Mono", "SF Mono", "Menlo", "Consolas", monospace;
    -webkit-font-smoothing: antialiased;
    text-rendering: optimizeLegibility;
  }
  ::selection {
    background: #e9e9e9;
  }
  .section-rule {
    border-top: 1px solid var(--line);
  }
  .quickstart-grid > * {
    min-width: 0;
  }
  .fade-in {
    animation: fade-in 700ms ease both;
  }
  .graph-stage {
    --story-duration: 14s;
    --blocked: #a94f3a;
    --complete: #16704a;
    position: relative;
    border: 1px solid var(--line);
    background: white;
    overflow: hidden;
  }
  .graph-svg {
    display: block;
    width: 100%;
    height: auto;
  }
  .graph-grid {
    fill: rgba(17, 17, 17, 0.065);
  }
  .graph-edge {
    fill: none;
    stroke: #e2e2e2;
    stroke-width: 1.5;
    stroke-linecap: round;
  }
  .story-path {
    fill: none;
    stroke: var(--ink);
    stroke-width: 2.5;
    stroke-linecap: round;
    stroke-dasharray: 1;
    stroke-dashoffset: 1;
    opacity: 0;
  }
  .story-path-fail-1 {
    animation: draw-fail-1 var(--story-duration) ease-in-out infinite;
  }
  .story-path-fail-2 {
    animation: draw-fail-2 var(--story-duration) ease-in-out infinite;
  }
  .story-path-fail-3 {
    stroke: var(--blocked);
    animation: draw-fail-3 var(--story-duration) ease-in-out infinite;
  }
  .story-path-return {
    stroke-dasharray: 0.04 0.045;
    animation: draw-return var(--story-duration) ease-in-out infinite;
  }
  .story-path-success-1 {
    animation: draw-success-1 var(--story-duration) ease-in-out infinite;
  }
  .story-path-success-2 {
    animation: draw-success-2 var(--story-duration) ease-in-out infinite;
  }
  .story-path-success-3 {
    stroke: var(--complete);
    animation: draw-success-3 var(--story-duration) ease-in-out infinite;
  }
  .story-node {
    color: #999999;
    opacity: 0.68;
  }
  .story-node .node-box {
    fill: #ffffff;
    stroke: currentColor;
    stroke-width: 1.5;
    vector-effect: non-scaling-stroke;
  }
  .story-node .node-label {
    fill: currentColor;
    font-size: 15px;
    font-weight: 500;
    letter-spacing: -0.01em;
  }
  .story-node .node-kicker {
    fill: currentColor;
    font-size: 10px;
    letter-spacing: 0.09em;
    text-transform: uppercase;
  }
  .node-site {
    animation: node-site-state var(--story-duration) ease-in-out infinite;
  }
  .node-launch-plan {
    animation: node-launch-state var(--story-duration) ease-in-out infinite;
  }
  .node-deploy-checklist {
    animation: node-deploy-state var(--story-duration) ease-in-out infinite;
  }
  .node-blocked {
    animation: node-blocked-state var(--story-duration) ease-in-out infinite;
  }
  .node-taskgraph {
    animation: node-taskgraph-state var(--story-duration) ease-in-out infinite;
  }
  .node-graph-model {
    animation: node-graph-state var(--story-duration) ease-in-out infinite;
  }
  .node-review-work {
    animation: node-review-state var(--story-duration) ease-in-out infinite;
  }
  .node-leaf {
    animation: node-leaf-state var(--story-duration) ease-in-out infinite;
  }
  .story-cursor {
    fill: var(--ink);
    stroke: white;
    stroke-width: 3;
    vector-effect: non-scaling-stroke;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.25));
    animation: move-cursor var(--story-duration) cubic-bezier(0.55, 0, 0.25, 1) infinite;
  }
  .story-status {
    display: grid;
    min-height: 56px;
    grid-template-columns: 82px minmax(0, 1fr);
    align-items: center;
    border-top: 1px solid var(--line);
    padding: 0 22px;
    font-size: 12px;
    line-height: 1.5;
  }
  .story-status-label {
    color: var(--muted);
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }
  .story-caption-stack {
    display: grid;
    min-width: 0;
  }
  .story-caption {
    grid-area: 1 / 1;
    color: var(--ink);
    opacity: 0;
    transform: translateY(4px);
  }
  .story-caption-1 { animation: caption-1 var(--story-duration) ease infinite; }
  .story-caption-2 { animation: caption-2 var(--story-duration) ease infinite; }
  .story-caption-3 { color: var(--blocked); animation: caption-3 var(--story-duration) ease infinite; }
  .story-caption-4 { animation: caption-4 var(--story-duration) ease infinite; }
  .story-caption-5 { animation: caption-5 var(--story-duration) ease infinite; }
  .story-caption-6 { color: var(--complete); animation: caption-6 var(--story-duration) ease infinite; }
  .graph-stage:hover *,
  .graph-stage:focus-within * {
    animation-play-state: paused;
  }
  @keyframes fade-in {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  @keyframes draw-fail-1 {
    0%, 11% { opacity: 0; stroke-dashoffset: 1; }
    18%, 39% { opacity: 1; stroke-dashoffset: 0; }
    47%, 100% { opacity: 0; stroke-dashoffset: 0; }
  }
  @keyframes draw-fail-2 {
    0%, 18% { opacity: 0; stroke-dashoffset: 1; }
    25%, 39% { opacity: 1; stroke-dashoffset: 0; }
    47%, 100% { opacity: 0; stroke-dashoffset: 0; }
  }
  @keyframes draw-fail-3 {
    0%, 25% { opacity: 0; stroke-dashoffset: 1; }
    32%, 39% { opacity: 1; stroke-dashoffset: 0; }
    47%, 100% { opacity: 0; stroke-dashoffset: 0; }
  }
  @keyframes draw-return {
    0%, 38% { opacity: 0; stroke-dashoffset: 1; }
    49%, 55% { opacity: 0.72; stroke-dashoffset: 0; }
    62%, 100% { opacity: 0; stroke-dashoffset: 0; }
  }
  @keyframes draw-success-1 {
    0%, 55% { opacity: 0; stroke-dashoffset: 1; }
    64%, 100% { opacity: 1; stroke-dashoffset: 0; }
  }
  @keyframes draw-success-2 {
    0%, 63% { opacity: 0; stroke-dashoffset: 1; }
    73%, 100% { opacity: 1; stroke-dashoffset: 0; }
  }
  @keyframes draw-success-3 {
    0%, 72% { opacity: 0; stroke-dashoffset: 1; }
    82%, 100% { opacity: 1; stroke-dashoffset: 0; }
  }
  @keyframes node-site-state {
    0%, 8%, 49%, 100% { color: #999999; opacity: 0.68; }
    13%, 25% { color: var(--ink); opacity: 1; }
    31%, 40% { color: var(--blocked); opacity: 1; }
  }
  @keyframes node-launch-state {
    0%, 15%, 47%, 100% { color: #aaaaaa; opacity: 0.5; }
    21%, 40% { color: var(--ink); opacity: 1; }
  }
  @keyframes node-deploy-state {
    0%, 22%, 47%, 100% { color: #aaaaaa; opacity: 0.5; }
    28%, 40% { color: var(--ink); opacity: 1; }
  }
  @keyframes node-blocked-state {
    0%, 29%, 47%, 100% { color: #aaaaaa; opacity: 0.5; }
    34%, 40% { color: var(--blocked); opacity: 1; }
  }
  @keyframes node-taskgraph-state {
    0%, 52%, 100% { color: #999999; opacity: 0.68; }
    58%, 94% { color: var(--ink); opacity: 1; }
  }
  @keyframes node-graph-state {
    0%, 60%, 100% { color: #aaaaaa; opacity: 0.5; }
    67%, 94% { color: var(--ink); opacity: 1; }
  }
  @keyframes node-review-state {
    0%, 69%, 100% { color: #aaaaaa; opacity: 0.5; }
    76%, 94% { color: var(--ink); opacity: 1; }
  }
  @keyframes node-leaf-state {
    0%, 78%, 100% { color: #aaaaaa; opacity: 0.5; }
    85%, 94% { color: var(--complete); opacity: 1; }
  }
  @keyframes move-cursor {
    0%, 7% { opacity: 0; transform: translate(0, 0); }
    10% { opacity: 1; transform: translate(230px, 0); }
    18% { transform: translate(230px, 100px); }
    25% { transform: translate(230px, 210px); }
    32%, 39% { opacity: 1; transform: translate(230px, 320px); }
    47% { transform: translate(230px, 0); }
    55% { transform: translate(0, 0); }
    63% { transform: translate(0, 100px); }
    72% { transform: translate(0, 210px); }
    82%, 94% { opacity: 1; transform: translate(0, 320px); }
    100% { opacity: 0; transform: translate(0, 320px); }
  }
  @keyframes caption-1 {
    0%, 2% { opacity: 0; transform: translateY(4px); }
    5%, 9% { opacity: 1; transform: translateY(0); }
    12%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @keyframes caption-2 {
    0%, 9% { opacity: 0; transform: translateY(4px); }
    12%, 27% { opacity: 1; transform: translateY(0); }
    30%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @keyframes caption-3 {
    0%, 28% { opacity: 0; transform: translateY(4px); }
    31%, 40% { opacity: 1; transform: translateY(0); }
    43%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @keyframes caption-4 {
    0%, 41% { opacity: 0; transform: translateY(4px); }
    44%, 55% { opacity: 1; transform: translateY(0); }
    58%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @keyframes caption-5 {
    0%, 56% { opacity: 0; transform: translateY(4px); }
    59%, 79% { opacity: 1; transform: translateY(0); }
    82%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @keyframes caption-6 {
    0%, 80% { opacity: 0; transform: translateY(4px); }
    83%, 96% { opacity: 1; transform: translateY(0); }
    99%, 100% { opacity: 0; transform: translateY(-4px); }
  }
  @media (max-width: 640px) {
    .graph-svg {
      width: 110%;
      max-width: none;
      margin-left: -5%;
    }
    .story-node .node-label {
      font-size: 17px;
    }
    .story-node .node-kicker {
      font-size: 12px;
      letter-spacing: 0.06em;
    }
    .story-status {
      min-height: 64px;
      grid-template-columns: 70px minmax(0, 1fr);
      padding: 0 14px;
      font-size: 11px;
    }
  }
  @media (prefers-reduced-motion: reduce) {
    .fade-in,
    .graph-stage * {
      animation: none;
    }
    .story-cursor,
    .story-path-fail-1,
    .story-path-fail-2,
    .story-path-fail-3,
    .story-path-return,
    .story-caption {
      display: none;
    }
    .story-path-success-1,
    .story-path-success-2,
    .story-path-success-3 {
      opacity: 1;
      stroke-dashoffset: 0;
    }
    .node-taskgraph,
    .node-graph-model,
    .node-review-work {
      color: var(--ink);
      opacity: 1;
    }
    .node-leaf {
      color: var(--complete);
      opacity: 1;
    }
    .story-caption-6 {
      display: block;
      opacity: 1;
      transform: none;
    }
  }
</style>

<main class="mx-auto max-w-5xl px-6 py-12 sm:px-8 sm:py-16">
  <section class="mx-auto max-w-4xl text-center fade-in">
    <p class="text-[15px] text-[color:var(--ink)]">TaskGraph</p>
    <h1 class="mx-auto mt-4 max-w-3xl text-4xl font-medium leading-tight sm:text-6xl">
      Decide what to do next.
    </h1>
    <p class="mx-auto mt-5 max-w-3xl text-sm leading-7 text-[color:var(--muted)] sm:text-[15px]">
      TaskGraph is a local-first CLI for quick task capture and graph-native planning. It keeps markdown
      authoritative, uses a derived index for fast queries, and helps you move from a high-level area of
      work to a concrete leaf task.
    </p>
    <div class="mt-8 flex flex-wrap items-center justify-center gap-3 text-sm">
      <a class="border border-[color:var(--ink)] bg-[color:var(--ink)] px-4 py-2 text-white" href="#install">
        Install
      </a>
      <a class="border border-[color:var(--line)] px-4 py-2 text-[color:var(--ink)]" href="#quickstart">
        Quickstart
      </a>
    </div>
    <p class="mt-6 text-[13px] text-[color:var(--muted)]">
      Start high. Choose a branch. Go down fast. Find a leaf.
    </p>
  </section>
  <section class="mx-auto mt-12 max-w-4xl fade-in">
    <div class="graph-stage">
      <svg
        class="graph-svg"
        viewBox="0 0 720 470"
        role="img"
        aria-label="Task graph animation: compare projects, descend into Site launch, discover that it needs breakdown, return to the project level, then descend through TaskGraph to the actionable leaf review notes"
      >
        <defs>
          <pattern id="graph-grid" width="24" height="24" patternUnits="userSpaceOnUse">
            <circle class="graph-grid" cx="1" cy="1" r="1" />
          </pattern>
        </defs>
        <rect width="720" height="470" fill="url(#graph-grid)" />
        <g aria-hidden="true">
          <path class="graph-edge" d="M360 96 L360 148" />
          <path class="graph-edge" d="M360 206 L360 258" />
          <path class="graph-edge" d="M360 316 L360 368" />
          <path class="graph-edge" d="M590 96 L590 148" />
          <path class="graph-edge" d="M590 206 L590 258" />
          <path class="graph-edge" d="M590 316 L590 368" />
          <path class="story-path story-path-fail-1" pathLength="1" d="M590 96 L590 148" />
          <path class="story-path story-path-fail-2" pathLength="1" d="M590 206 L590 258" />
          <path class="story-path story-path-fail-3" pathLength="1" d="M590 316 L590 368" />
          <path class="story-path story-path-return" pathLength="1" d="M590 438 L590 118 C590 108 580 108 570 108 L380 108 C370 108 360 108 360 118" />
          <path class="story-path story-path-success-1" pathLength="1" d="M360 96 L360 148" />
          <path class="story-path story-path-success-2" pathLength="1" d="M360 206 L360 258" />
          <path class="story-path story-path-success-3" pathLength="1" d="M360 316 L360 368" />
        </g>
        <g class="story-node node-writing">
          <rect class="node-box" x="45" y="38" width="170" height="58" rx="8" />
          <text class="node-kicker" x="61" y="58">project</text>
          <text class="node-label" x="61" y="81">Book proposal</text>
        </g>
        <g class="story-node node-taskgraph">
          <rect class="node-box" x="275" y="38" width="170" height="58" rx="8" />
          <text class="node-kicker" x="291" y="58">project</text>
          <text class="node-label" x="291" y="81">TaskGraph</text>
        </g>
        <g class="story-node node-site">
          <rect class="node-box" x="505" y="38" width="170" height="58" rx="8" />
          <text class="node-kicker" x="521" y="58">project</text>
          <text class="node-label" x="521" y="81">Site launch</text>
        </g>
        <g class="story-node node-graph-model">
          <rect class="node-box" x="275" y="148" width="170" height="58" rx="8" />
          <text class="node-kicker" x="291" y="168">epic</text>
          <text class="node-label" x="291" y="191">Graph model</text>
        </g>
        <g class="story-node node-review-work">
          <rect class="node-box" x="275" y="258" width="170" height="58" rx="8" />
          <text class="node-kicker" x="291" y="278">task</text>
          <text class="node-label" x="291" y="301">Review research</text>
        </g>
        <g class="story-node node-leaf">
          <rect class="node-box" x="275" y="368" width="170" height="58" rx="8" />
          <text class="node-kicker" x="291" y="388">next action</text>
          <text class="node-label" x="291" y="411">Review notes</text>
        </g>
        <g class="story-node node-launch-plan">
          <rect class="node-box" x="505" y="148" width="170" height="58" rx="8" />
          <text class="node-kicker" x="521" y="168">epic</text>
          <text class="node-label" x="521" y="191">Launch plan</text>
        </g>
        <g class="story-node node-deploy-checklist">
          <rect class="node-box" x="505" y="258" width="170" height="58" rx="8" />
          <text class="node-kicker" x="521" y="278">task</text>
          <text class="node-label" x="521" y="301">Deploy site</text>
        </g>
        <g class="story-node node-blocked">
          <rect class="node-box" x="505" y="368" width="170" height="58" rx="8" />
          <text class="node-kicker" x="521" y="388">no next action</text>
          <text class="node-label" x="521" y="411">Needs breakdown</text>
        </g>
        <circle class="story-cursor" aria-hidden="true" cx="360" cy="118" r="5" />
      </svg>
      <div class="story-status" aria-hidden="true">
        <span class="story-status-label">TaskGraph</span>
        <span class="story-caption-stack">
          <span class="story-caption story-caption-1">Orient across the work that matters.</span>
          <span class="story-caption story-caption-2">Choose a promising branch and inspect it.</span>
          <span class="story-caption story-caption-3">No actionable leaf. This branch needs breakdown.</span>
          <span class="story-caption story-caption-4">Back up instead of forcing a vague task.</span>
          <span class="story-caption story-caption-5">Choose another branch and descend quickly.</span>
          <span class="story-caption story-caption-6">Do next: review notes.</span>
        </span>
      </div>
    </div>
    <div class="mt-4 grid gap-4 text-xs text-[color:var(--muted)] sm:grid-cols-3">
      <p>High-level orientation tells you which area matters.</p>
      <p>Selective descent narrows the branch instead of dumping every leaf.</p>
      <p>The point is a real next action, not a better-organized vague project.</p>
    </div>
  </section>
  <section id="install" class="section-rule mx-auto mt-20 max-w-4xl pt-10">
    <p class="text-[15px] font-medium">## Install</p>
    <div class="mt-5 max-w-3xl">
      <p class="text-sm leading-7 text-[color:var(--muted)]">
        Install the latest release on macOS or Linux, then initialize a project and start capturing tasks.
      </p>
      <pre class="mt-5 overflow-x-auto border border-[color:var(--line)] bg-[color:var(--soft)] p-5 text-sm leading-7 text-[color:var(--ink)]"><code>curl -fsSL https://raw.githubusercontent.com/rufuspollock/taskgraph/main/scripts/install.sh | bash</code></pre>
      <p class="mt-4 text-sm leading-7 text-[color:var(--muted)]">
        Task data lives in markdown. The SQLite database is derived state only and can be rebuilt.
      </p>
    </div>
  </section>
  <section id="quickstart" class="section-rule mx-auto mt-16 max-w-4xl pt-10">
    <p class="text-[15px] font-medium">## Quickstart</p>
    <div class="quickstart-grid mt-5 grid gap-10 sm:grid-cols-[1.05fr_0.95fr]">
      <div>
        <pre class="overflow-x-auto border border-[color:var(--line)] bg-white p-5 text-sm leading-7 text-[color:var(--ink)]"><code>tg init
tg add "buy milk"
tg add "plan flower show" --labels flowershow,events
tg add "map launch dependencies" --type epic
tg create "book dentist"
tg inbox
tg list --label flowershow
tg graph --depth 3</code></pre>
      </div>
      <div class="space-y-5 text-sm leading-7 text-[color:var(--muted)]">
        <p>
          <span class="text-[color:var(--ink)]">Capture quickly.</span> `tg add` and `tg create` make it
          easy to get tasks into the system without breaking flow.
        </p>
        <p>
          <span class="text-[color:var(--ink)]">Keep markdown authoritative.</span> Inbox tasks live in
          <code class="text-[color:var(--ink)]">INBOX.md</code>, and labels stay inline as
          markdown tags.
        </p>
        <p>
          <span class="text-[color:var(--ink)]">Query the graph.</span> `tg list` and `tg graph` use a
          derived SQLite index so you can move quickly without hiding your data in a proprietary layer.
        </p>
      </div>
    </div>
  </section>
  <section id="vision" class="section-rule mx-auto mt-16 max-w-4xl pt-10">
    <p class="text-[15px] font-medium">## Vision</p>
    <div class="mt-5 max-w-3xl space-y-6 text-sm leading-7 text-[color:var(--muted)]">
      <p>
        The core problem is not capture by itself. The harder problem is deciding what to work on next
        without burning twenty minutes digging through projects, notes, and vague half-tasks.
      </p>
      <p>
        A giant flat list does not solve that. It just gives you a better organized version of overwhelm.
        Stopping at "work on this project" does not solve it either, because a project is still too vague
        to tell you what to do with the next thirty minutes.
      </p>
      <p>
        What matters is movement between levels. Start high enough to choose a direction. Then descend
        quickly until you hit a real leaf task. If a branch fails to produce a clear next action, that is
        useful information too: it may be blocked, it may need breaking down, or it may simply not be ripe
        yet.
      </p>
      <p>
        This is why the graph matters. It preserves the relationship between high-level choices and concrete
        actions. It lets you practice selective descent: choose a branch, open that branch, keep opening it,
        and surface a few real candidate leaves instead of every leaf in the universe.
      </p>
      <div class="border border-[color:var(--line)] bg-white p-5 text-[13px] leading-7 text-[color:var(--ink)]">
        <p>Start high.</p>
        <p>Choose a branch.</p>
        <p>Go down fast.</p>
        <p>Find a leaf.</p>
        <p>If there isn't a good leaf, learn why.</p>
        <p>Then cycle.</p>
      </div>
      <p>
        TaskGraph is intended to make that graph legible and traversable. AI then helps interpret, rank,
        clarify, and guide. The tool should not replace judgment. It should provide solid local structure
        for judgment to work on.
      </p>
    </div>
  </section>
  <section class="section-rule mx-auto mt-16 max-w-4xl pt-10">
    <div class="flex flex-wrap items-center gap-3 text-sm">
      <a class="border border-[color:var(--ink)] bg-[color:var(--ink)] px-4 py-2 text-white" href="#install">
        Install
      </a>
      <a
        class="border border-[color:var(--line)] px-4 py-2 text-[color:var(--ink)]"
        href="https://github.com/rufuspollock/taskgraph#usage"
        target="_blank"
        rel="noreferrer"
      >
        Read usage
      </a>
      <a
        class="border border-[color:var(--line)] px-4 py-2 text-[color:var(--ink)]"
        href="https://github.com/rufuspollock/taskgraph"
        target="_blank"
        rel="noreferrer"
      >
        GitHub
      </a>
    </div>
  </section>
</main>
