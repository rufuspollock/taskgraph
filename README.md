# TaskGraph

## Vision

The tool is a local task-indexing and sense-making layer that turns a corpus of Markdown files into a single, queryable task graph. Its purpose is not task entry or manual planning, but extraction, normalization, and interpretation: given an existing body of notes, projects, and checklists, it continuously derives a coherent view of “all the work there is,” so that relevant next actions and meaningful structures can be surfaced automatically.

## Description of what is being built

The system scans one or more folders of Markdown files (e.g. an Obsidian vault) and extracts all checklist items as first-class tasks, regardless of where they appear. Each task is stored in a local database together with rich contextual metadata: source file, heading path, surrounding text, tags, links, timestamps, and any explicit task annotations (due dates, estimates, statuses). Pages themselves may also be treated as tasks when they function as containers for other tasks (e.g. files in a “Projects” folder or pages with task-heavy structure).

Rather than imposing a fixed hierarchy of projects, tasks, and subtasks, the system models relationships between tasks: containment (task appears within a section or page), sequencing or dependency (implicit or explicit), and grouping (tasks co-located or semantically linked). The result is a directed, acyclic task graph that reflects how work is actually expressed in notes, not how a task manager expects it to be entered.

On top of this graph sits a query layer. The user can ask questions such as: “What can I do in 10 minutes?”, “What is the next available task in this chain?”, or “What tasks belong to this area of work and are unblocked?” The interface—whether textual, tabular, or visual—is secondary to the core capability: a continuously updated, queryable database of tasks derived from Markdown, capable of supporting both fast next-action selection and higher-level reasoning about the flow and structure of work.

## MVP usage

```bash
pnpm install
pnpm link --global
taskgraph index fixtures/rufus-projects
taskgraph query
taskgraph query "meeting"
```
