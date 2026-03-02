# tg next as graph navigation

## Purpose

`tg next` should help decide what to work on next without flooding the user with all available leaf tasks.

The key idea is that "next" is not just ranking a pile of tasks. It is a navigation problem across a task graph. The user often wants to choose a promising branch high up the graph, then descend quickly to a concrete actionable leaf. If the branch does not yield a good leaf, the system should make that visible and support cycling to another branch.

## Core product idea

`tg next` should support two linked decisions:

1. What area of work should I focus on right now?
2. What exact task should I do inside that area?

Those two decisions happen at different levels of the graph. A good `next` flow should let the user move between them quickly rather than collapsing everything into one flat ranked list.

## User loop

The intended loop looks like this:

1. Show a small number of promising top-level or mid-level branches.
2. Let the user pick or implicitly choose one branch.
3. Descend into that branch until the system finds candidate actionable leaves.
4. If no good leaf exists, explain why:
   blocked by dependency
   needs breakdown
   no open leaf tasks
5. Let the user cycle back and try another branch.

This gives `tg next` a shape closer to narrowing and traversal than to generic listing.

## Why a flat list is insufficient

A flat list of all open leaf tasks fails in 3 ways:

1. Volume
Large systems may contain tens, hundreds, or thousands of leaves.

2. Loss of context
A leaf task without its enclosing branch often does not explain why it matters.

3. Poor project selection
Sometimes the real decision is not "which task?" but "which project or epic deserves focus?"

`tg next` should preserve enough graph context to answer both.

## Required graph concepts

### 1. Node types

At minimum, the system should distinguish between:

- idea
- initiative or project
- epic or feature
- task
- subtask

These do not need to be rigidly enforced, but they help drive expectations about where branching, decomposition, and actionable leaves are likely to live.

### 2. Parent-child containment

The system needs to know how items nest. This provides the basic descent path from project to leaf.

### 3. Dependency edges

Dependencies matter for deciding whether a node is actionable. A leaf that is blocked should not be treated as a good next task. A branch whose leaves are mostly blocked should probably be surfaced as blocked or weak.

### 4. Actionable leaf detection

`tg next` needs a practical definition of an actionable leaf. For now, a node is a candidate leaf if it:

- is open
- has no open children beneath it
- is not blocked by known dependencies
- is small enough to be plausibly actionable

"Small enough" can be heuristic at first.

### 5. Breakdown-needed detection

A node should be marked as needing breakdown when it is open but still behaves like a container without yielding concrete child actions. This is important because "needs breakdown" is itself a useful answer.

## Desired interaction shape

### Option A: branch-first

`tg next` first shows a few promising branches:

```text
Promising branches

1. website redesign
2. task graph
3. hiring process
```

Then it expands one:

```text
task graph
  -> inbox improvements
  -> next-action selection
  -> graph model cleanup
```

Then it yields leaves:

```text
next-action selection
  -> define actionable leaf heuristic
  -> write blocked-branch rule
```

This best matches the north-star vision.

### Option B: leaf-first with branch context

`tg next` returns leaves directly, but grouped under a few parent branches:

```text
task graph
  - define actionable leaf heuristic
  - write blocked-branch rule

website redesign
  - draft homepage hero copy
```

This is simpler, and may be a good early version.

### Recommendation

Start with `Option B`, because it is easier to build on top of current list/index behavior, but shape the data model and output toward `Option A`.

That means even early `next` output should preserve branch context and support future multi-step descent rather than pretending ranking alone is the answer.

## Ranking implications

Ranking should happen at 2 levels:

1. Branch ranking
Which project, epic, or area seems worth attention now?

2. Leaf ranking within branch
Which actionable leaves inside that branch are best candidates?

Signals may include:

- explicit priority
- dependency readiness
- recency
- effort or size
- status signals like blocked or waiting
- whether the branch has clear leaves at all

The important point is that branch quality and leaf quality are separate.

## Output principles

`tg next` should:

- show a small set
- preserve branch context
- prefer actionable leaves
- surface blockedness and missing breakdown explicitly
- support cycling to another branch

`tg next` should not:

- dump every leaf task
- hide why a leaf matters
- pretend a blocked branch is ready
- force the user to fully navigate the whole graph manually

## Early implementation shape

An early useful version could work like this:

1. Identify candidate branch nodes above the leaf level.
2. Score them roughly.
3. For each top branch, search downward for open actionable leaves.
4. Return:
   best leaves found
   or a branch-state explanation such as `blocked` or `needs-breakdown`

This gives useful behavior before building a richer interactive traversal UX.

## Open questions

- How much should item type be explicit versus inferred?
- Should "needs breakdown" become a first-class status?
- How should inbox items attach into branches, if at all?
- When should dependencies be explicit links versus inferred from structure?
- Should `tg next` optimize for 1 best answer or 3-5 branch-scoped candidates?

## Design takeaway

`tg next` should be built around guided descent through a task graph.

The product goal is to help the user move from "what area should I work on?" to "what exact thing can I do now?" quickly, with enough structure to explain why a branch is promising, blocked, or under-specified.
