# Task Graph: choosing what to do next

I keep circling back to the same problem: I want to know what to work on next, without spending 20 minutes digging through projects, notes, and half-formed tasks first.

Capture matters. Storage matters. Whether issues live in `.taskgraph/` or the main repo probably matters too. Same for whether there should be one file per issue. But those aren't the live wire for me right now.

The live wire is this: I want a system that helps me move from a big, messy set of projects to a concrete thing I can do now.

That sounds small. It isn't. It's the whole game.

What I usually want is a movement between levels. I want to start fairly high up, around projects, epics, or major areas of work. Then I want to go down fast until I hit an actual leaf task. If there isn't a good leaf task, I want that to be obvious too.

The diagrams are the point.

## 1. Start high enough to choose a direction

```text
High-level view

  [Project A] -----
                   \
  [Project B] ------>  "What do I want to work on?"
                   /
  [Project C] -----
```

This is useful because it helps me orient. It helps me say, "yes, that's the area I want." But it still doesn't tell me what to do.

## 2. Go down into a branch

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

Better. Still not enough.

Projects are usually too vague. Even tasks are often still containers. I don't want to stop at "I should work on Project B" or even "I should work on Task B2.1". I want the thing I can actually do.

## 3. Keep going until the graph yields a leaf

```text
Find the leaf tasks

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

Now we're somewhere real.

This is where the graph starts paying rent. It lets me move from "this seems like the right area" to "review notes" or "email John". Small, concrete, actionable.

## 4. Don't show me every leaf in the universe

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

I don't want all the leaves. That's the trap.

If there are 200 leaf tasks, showing me all 200 is just a more structured version of overwhelm. What I want is selective descent: choose a branch at a high level, then open that branch until I hit good candidate leaves.

That's a very different interaction from a flat backlog dump.

## 5. Sometimes the branch fails, and that's useful

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

This matters a lot.

Sometimes I choose an area and there isn't a good next task in it. That isn't useless. That's exactly the kind of thing I want the graph to show me. Maybe the branch is blocked. Maybe it needs to be broken down. Maybe it's just not ripe.

That's good information. It stops me kidding myself that I've "chosen a project" when I still haven't found work I can actually do.

## 6. Cycle back, then try another branch

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

This back-and-forth is the real motion.

I want to filter high up, descend quickly, test whether the branch yields a real next action, then cycle if it doesn't. That's what I keep doing in my head anyway. I'd like the system to support it instead of forcing me into either a giant flat list or a vague project tree.

## Why the graph matters

The graph matters because the problem lives in the movement between scales.

At one level, I want to decide which project or epic matters. At another, I want a leaf task I can do in 10 minutes, or 30, or 2 hours. I move back and forth between those levels constantly.

Dependencies matter here. Issue types matter. Bigger items higher up the graph should usually house more subitems. Some branches should be visibly blocked. Some should clearly need decomposition. Some ideas are orthogonal and later connect across branches. Real work is messy like that.

So yes, the graph is a data structure. But more importantly, it's a way of seeing.

## North star

Here's the thing I want:

```text
Start high.
Choose a branch.
Go down fast.
Find a leaf.
If there isn't a good leaf, learn why.
Then cycle.
```

If Task Graph can support that process well, a lot of the storage and representation questions can get worked out through use.

That's the north star I care about. A system that helps me decide what to do next by moving well between scales of work.
