## v1 - late Dec 2025 - see old README

## v2 2026-01-17

See docs/vision-v2-2026-01-16.md

## v3 reboot note - 2026-02-24

### Overall Vision

Build Task Graph as a practical, daily-use system that helps decide what to do next and capture work fast, while building a durable task graph underneath. Use it directly while building it (dogfooding), first in this repo and then in the planning repo, so decisions are grounded in real day-to-day usage.

### Jobs To Be Done

- `tg next` Help decide what to do next each day, quickly and with low overhead (Surface useful "next" actions quickly from the task graph)
- `tg create` (alias: `tg add`, `tg capture`. Capture tasks and ideas instantly in a GTD-style flow
- `tg graph` Support hierarchical tasks (projects, subprojects, tasks)
- `tg inbox`: do GTD conversion of captured items (in "inbox") to clarify to backlog with increasing automation over time
- `tg list` or `tg graph`: Serve as a shared operational substrate for both human use and AI-assisted planning

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
