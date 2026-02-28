## v1 - late Dec 2025 - see old README

## v2 2026-01-17

See docs/vision-v2-2026-01-16.md

## v3 reboot note - 2026-02-24

### Overall Vision

Build Task Graph as a practical, daily-use system that helps decide what to do next and capture work fast, while building a durable task graph underneath. Use it directly while building it (dogfooding), first in this repo and then in the planning repo, so decisions are grounded in real day-to-day usage.

### Jobs To Be Done

1. `tg create` (alias: `tg add`, future alias `tg capture`)
Capture tasks and ideas instantly in a GTD-style flow.
Status: `Implemented` (`tg add`/`tg create` append to `.taskgraph/issues.md` and refresh index DB).

2. `tg inbox`
See and process captured inbox items from `.taskgraph/issues.md`.
Status: `Implemented` (raw inbox checklist view).

3. `tg list`
See task graph checklist items across indexed markdown in one queryable view.
Status: `Implemented (v1)` (SQLite-backed checklist listing, open-only default, `--all` to include closed).

4. `tg graph`
Support richer hierarchical/relationship-oriented task views (projects, subprojects, tasks).
Status: `Partial` (hierarchy is indexed in DB; dedicated graph UX/command not implemented).

5. `tg next`
Help decide what to do next quickly with low overhead.
Status: `Not implemented` (planned to build on DB-backed list + ranking heuristics).

### JTBD Review Loop

- Keep this section as the source of truth for command intent and support level.
- When command behavior changes, update JTBD status in the same PR/commit.
- Prefer explicit status labels: `Not implemented`, `Partial`, `Implemented (vN)`.

### Key Principle: AI Uses The Tool

Task Graph should not try to automate every part of planning by itself. Instead, it should provide strong primitives that make an AI workflow effective:

- A reliable task/backlog data layer
- Fast CLI commands for capture, logging, and structure
- Clear graph relationships and hierarchy

This is primarily an "AI uses the tool" architecture, rather than "the tool uses a little AI." The software should do what software is best at: persistence, structure, and fast operations. The AI should do what AI is best at: interpretation, prioritization support, clarification, and guidance.

### Principle: use markdown for issue items

We use markdown for issue storage and the checklist format for "tasks".

- Try and find a format that will allow us to push and pull from github or simliar

### Principles of development

- Build by use, agile iteration of a prototype

### Ideas

- [ ] ➕2026-02-24 want to integrate with existing issue trackers rather than replacing. e.g. if a project is using github issues or linear lets use that. can still keep a local cache in .local for offline. what would be local is the graph.
