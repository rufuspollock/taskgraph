# North Star: A Task Graph That Helps Me Decide What To Do Next

What I really want is not just somewhere to capture tasks.

I have kind of got that now. I can capture files, tasks, issues, inbox items. That matters, but it is not the main thing I am struggling with.

What I keep struggling with is: how do I see what I should work on next?

That has always been the real "next issue" question. That is where having a task graph is actually useful. The graph part is not just a nice representation. It is the thing that can help me move from a large landscape of work to a real next action without getting overwhelmed.

This is why some of the architectural questions do not feel primary yet. Should issues live in `.taskgraph/` or in the main repo? Should there be one file per issue? Those questions matter, but they are not the north star. It is fine to try one approach for actual usage and see. The more important question is what process I want the system to support.

The process I want to support is something like this:

- I want to look at fairly high-level things first.
- I want to choose among projects, epics, or major areas of work.
- Then I want to go down quickly to find an actual small task I can do.
- If there is no clear small task, I want that to become visible too.

That movement up and down the graph is the core of it.

I do not want to be shown all the leaf tasks at once. That is overwhelming. There might be tens or hundreds or thousands of them.

But conversely, if I just say "I want to work on Project X", that is not enough either. A project is too high-level. It does not tell me what to do next. I need to find the leaf tasks under it, or discover that I need to create them.

So the useful thing is not just hierarchy by itself, and not just a giant flat backlog. The useful thing is being able to cycle between levels:

1. Start high up, at the level of projects, epics, or source nodes.
2. Choose an area that feels important or alive.
3. Descend quickly through the graph.
4. Find the actionable leaf nodes.
5. If the branch does not yield a good next task, go back up and choose another branch.

That is the behavior I want Task Graph to support.

Issue types and dependencies matter here because they help shape that movement. Larger items higher up the graph should tend to have richer internal structure. Epics should have more subitems. Dependencies should help explain why something is not yet actionable, or why one leaf task is the real next move before another.

The graph is useful because it can help answer two related but different questions:

- What area of work do I want to focus on?
- What is the concrete thing I can actually do next inside that area?

Those are not the same question. In practice I am often cycling between them.

## The core movement

At the top level, I may want something like this:

```text
High-level landscape

  [Project A] -----
                   \
  [Project B] ------>  "What do I want to work on?"
                   /
  [Project C] -----
```

This is already helpful. It lets me orient. It lets me say, "yes, I want to work on that project." But it still does not tell me what to do.

So then I want to descend:

```text
Zoom into one branch

  [Project B]
      |
      +--> [Epic B1]
      |        |
      |        +--> [Task B1.1]
      |        +--> [Task B1.2]
      |
      +--> [Epic B2]
               |
               +--> [Task B2.1]
               +--> [Task B2.2]
```

This gets closer, but it still may be too coarse. Some of those tasks may still be containers. Some may not be ready. Some may not have a genuine next action inside them.

So I want to keep going until I get to leaf nodes:

```text
Find the real leaf tasks

  [Project B]
      |
      +--> [Epic B2]
               |
               +--> [Task B2.1]
               |        |
               |        +--> [leaf: draft outline]
               |        +--> [leaf: review notes]
               |
               +--> [Task B2.2]
                        |
                        +--> [leaf: email John]
                        +--> [leaf: book venue]
```

Now I am somewhere useful. I have specific, concrete actions. But even here I do not necessarily want every leaf task. I want the relevant leaf tasks for the branch I currently care about.

That suggests another move: selective descent rather than total expansion.

```text
Selective descent

  [Project A]   [Project B]   [Project C]
                     *
                     |
                  [Epic B2]
                     |
         +-----------+-----------+
         |                       |
   [Task B2.1]             [Task B2.2]
         |                       |
   [leaf: review notes]    [leaf: email John]
```

This is closer to how I actually want to work. I choose at a high level, then selectively open up a branch until I hit actionable leaves.

But there is a further twist: sometimes I choose a branch and discover that it does not yet contain a good next task.

```text
Branch with no clear next action

  [Project C]
      |
      +--> [Epic C1]
               |
               +--> [Task C1.1]
                        |
                        +--> [? no small actionable leaf]
                        +--> [? blocked by dependency]
                        +--> [? needs breakdown]
```

That is not a failure. That is useful information. It tells me one of three things:

- this branch is blocked
- this branch needs decomposition
- this branch is not the place to work right now

And then I want to cycle back:

```text
Cycle back up, then down another branch

  [Project A]   [Project B]   [Project C]
                     ^              x
                     |
               return here
                     |
                 choose B
                     |
                  [Epic B2]
                     |
            [leaf: review notes]
```

That cycling back and forth is a big part of the point.

## Why a graph matters

A flat task list loses too much information. It can show tasks, but not why they matter, how they relate, what they unblock, or which larger effort they belong to.

A pure hierarchy is better, but still not enough. Real work is messier than a strict tree. Some tasks depend on other tasks across branches. Some issues are really project-level. Some are tiny inbox items. Some larger items should have many subitems beneath them. Some ideas are orthogonal and later attach to multiple places.

That is why I keep coming back to the graph.

By constructing the graph, I can support a process that combines:

- orientation at a high level
- filtering by which project or epic feels important
- descent to concrete leaf tasks
- visibility into dependencies and blockedness
- rapid cycling when one branch does not yield a good next action

The graph is not there for abstract modeling. It is there to support this movement.

## The north star

The north star is not "store all my issues neatly."

The north star is not "show me every task."

The north star is not even just "have a nice hierarchy."

The north star is:

```text
Help me move from "what area matters?" to "what exact thing should I do now?"
without overwhelm, and with enough structure to see why.
```

Or said another way:

```text
Start high.
Choose a branch.
Go down fast.
Find a leaf.
If there is no good leaf, learn why.
Then cycle.
```

That feels like the real promise of Task Graph.

If it can do that well, then many of the lower-level questions about storage format, file layout, and representation can be worked out through use. Those questions are important, but they are downstream of the main thing.

The main thing is helping me decide what to do next by moving well between scales of work.
