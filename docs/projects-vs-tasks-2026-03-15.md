# Projects vs Tasks: Why Scale Matters and How to Keep It Simple

What I really want from Task Graph is described in the [north star vision](https://rufuspollock.com/ref/task-graph-vision-essay-2026-03-02): start with my high-level projects, pick one to focus on, descend to a concrete task I can actually do, and if there is no good task, cycle back and pick another project. That movement between scales of work is the whole point.

But for that movement to work, I need to know which things are "high" and which things are "low." I need to be able to say: show me my projects. Show me the big things I am choosing among. And then, once I have picked one, show me the tasks inside it.

That means the system needs some notion of scale. Not a rigid six-level hierarchy. Just enough to distinguish "this is a thing I pick among" from "this is a thing I do."

## The problem with Getting Things Done

In GTD, a "project" is anything with more than one step. Booking a dentist appointment is a project. Building a product is a project. They sit in the same list.

That is fine for GTD's purposes. But it does not help me with the movement I described above. If I ask "show me my projects" and I get back a list that includes "book dentist" alongside "build TaskGraph" and "Life Itself strategy", that is not useful. I cannot orient at a high level because the levels are all mixed together.

So I need a distinction that GTD does not make. I need to separate the big things from the small things. Not in a complicated way. Just enough to support the core movement: start high, descend, find a leaf.

## What I actually need: two levels

I have gone back and forth on how many levels of "big" I need. The current task-types schema has six: initiative, product, epic, feature, task, subtask. That is probably right as a rough conceptual model. But for the actual system, for what tg needs to know to show me the right view, I think two is enough:

- **Project**: a big thing I pick among. It has real scope and identity. It might take weeks or months. It has subtasks inside it. Examples: "TaskGraph website", "Florence trip 2024", "Life Itself strategy."
- **Task**: everything else. A concrete thing I do, possibly with subtasks of its own, but not something I would choose among at the top level.

That is it. Two levels. Project and task.

I am deliberately reclaiming the word "project" here and not using it in the GTD sense. In tg, a project means something with real scope. "Book dentist" is a task, even if it has two steps. "Build a landing page for TaskGraph" is a project.

If I later need a level above project — "initiative", for things like "Life Itself" or "Datopian" that contain multiple projects — I can add that. But I do not think I need it yet for the core movement to work. When I sit down to plan my day or week, I am not usually choosing between initiatives. I am choosing between projects. The initiative level is useful context but not usually the decision point.

So: start with two. Design the system so a third level can be added later without breaking anything. But do not build for three until I actually need three.

## How do I know something is a project?

This is the real design question. If I have to manually label every item as "project" or "task" when I create it, I will not do it consistently. I know myself. I am not disciplined about labelling things. So the system needs to help.

There are three signals, in rough order of reliability:

**1. It has its own file.** This is the strongest signal I have right now. If something is a whole markdown file — not a checklist item inside another file, but its own file — that is a pretty good sign it is project-scale. My `projects/` folder is full of these. Each one is a project. The file is the project.

This is already a convention I follow naturally. I do not create a file for "buy milk." I create a file for "Florence trip 2024" or "TaskGraph landing page." The act of giving something its own file is already an implicit declaration that it is big enough to be a project.

**2. It is explicitly labelled.** When I `tg add "plan flower show" --type project`, that is a clear signal. The problem is I do not do this very often right now. But I could do it more, especially if the system prompted me.

**3. It can be inferred from the language.** If I write "build the TaskGraph website" or "plan the festival," the words themselves suggest project-scale. An AI-assisted `tg add` could notice this and prompt me: "This sounds like a project — should I mark it as one?" Similarly, `tg index` could look at items it finds and make guesses based on language and structure.

The practical plan is probably: rely on file-based inference first (it is free and already works), support explicit `--type project` for when I want to be clear, and build NLP-assisted prompting later as a nice-to-have.

## What this means for TaskGraph

If the system can reliably distinguish projects from tasks, even roughly, then the core movement starts to work:

1. `tg projects` (or `tg list --type project`) shows me my projects. These are the high-level things I choose among.
2. I pick one.
3. `tg list --project <name>` or just descending the graph shows me the tasks inside it.
4. I find a leaf task and do it.

That is the cycle. That is what I have been trying to get to.

The key insight is that I do not need a perfect taxonomy. I do not need to decide whether something is an "epic" or a "feature" or a "product." I just need to know: is this a thing I pick among, or a thing I do? Project or task. Big or small. Root or leaf.

One thing worth noting: scale is relative. TaskGraph itself is a project in my life. But within the TaskGraph repo, the "projects" are really features — things like "add project/task distinction" or "fix the graph view." If I were inside Life Itself, the projects would be "conscious coliving program" or "new website." The word "project" does not mean a fixed absolute size. It means "the big things at whatever level I am currently looking at." A tg repo is a scope. The projects within it are the roots of that scope. This is fine. It actually reinforces the binary design — the system does not need to know absolute scale, it just needs to know what is root-level and what is nested within the current context.

Everything else — the richer type system, dependencies, epics, initiatives — can layer on top of that binary distinction later. But the binary distinction is the foundation that makes the graph navigable.

## Design implications

1. **Two types for now**: `project` and `task` (with `task` as the default). Everything is a task unless marked or inferred as a project.
2. **File = project heuristic**: if `tg index` encounters a standalone markdown file (not a checklist item inside another file), it should consider it a project candidate. Especially if it lives in a known projects directory.
3. **Explicit labelling**: `tg add --type project` should work. `type: project` in frontmatter should work.
4. **Inference at add time**: when running `tg add` with AI assistance, the system could prompt: "This sounds like it might be a project. Should I mark it as one?"
5. **Inference at index time**: `tg index` could use simple heuristics (has its own file, lives in projects folder, title contains "project"/"plan"/"initiative") to guess scale. These guesses can be surfaced for confirmation rather than silently applied.
6. **Design for a third level later**: the schema should allow `initiative` or similar above `project` without restructuring. But do not implement it yet.
