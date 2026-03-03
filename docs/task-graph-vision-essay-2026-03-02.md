# How I want to decide what to do next

I keep coming back to the same problem.

I want to know what to work on next, without burning 20 minutes digging through projects, notes, and vague half-tasks first.

Capture matters. Storage matters. Where the files live matters. Those questions are downstream.

What matters is being able to move from a large mess of possible work to a concrete thing I can do now.

I don't want a giant flat list. That just gives me a better organized version of overwhelm.

I also don't want to stop at "work on Project X". That's too vague. A project doesn't tell me what to do with the next 30 minutes.

What I want is a movement between levels.

I want to start high enough to choose a direction. Then I want to go down fast until I hit a real leaf task. If there isn't a good leaf task, I want that to be obvious too.

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

This helps me orient.

It helps me say, "yes, that's the area." But it still doesn't tell me what to do.

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

Better. Still too vague.

Projects are often containers. Even tasks are often containers. I don't want to stop at "I should work on Project B" or even "Task B2.1". I want the thing I can actually do.

## 3. Keep going until the branch yields a leaf

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

This is where the structure starts paying rent. It lets me move from "this seems like the right area" to "review notes" or "email John". Small. Concrete. Actionable.

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

This bit matters a lot.

I don't want all the leaf tasks. If there are 200 of them, showing me all 200 is just a more elaborate way to drown me.

I want selective descent. Choose a branch high up. Open that branch. Keep opening it until I hit a few good candidate leaves.

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

This is useful information.

Sometimes I choose an area and there isn't a good next task in it. Fine. I want to know that. Maybe the branch is blocked. Maybe it needs breaking down. Maybe it's just not ripe yet.

That is a real answer. It tells me not to pretend I've made progress by "choosing a project" when I still haven't found work I can actually do.

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

I want to filter high up, descend quickly, test whether the branch yields a real next action, then cycle if it doesn't. That's what I keep doing in my head anyway.

## Why the graph matters

The problem lives in the movement between scales.

At one level, I want to decide which project or epic matters. At another, I want a leaf task I can do in 10 minutes, 30 minutes, or 2 hours. I move between those levels constantly.

That's why the graph matters. It keeps the relationship between the high-level choice and the concrete action.

Dependencies matter here too. Some branches are blocked. Some need breaking down. Bigger items higher up the graph usually have more structure beneath them. Some ideas connect across branches later. Real work is messy like that.

## What I'm trying to build

Here's the process I want to support:

```text
Start high.
Choose a branch.
Go down fast.
Find a leaf.
If there isn't a good leaf, learn why.
Then cycle.
```

That's the north star for what I'm building with Task Graph.

The job of the tool is to create the graph, keep it legible, and make it easy to move through. The job of AI is different. AI helps interpret, rank, clarify, and guide. I don't want the tool to replace that. I want it to give AI, and me, something solid to think with.

If this works, then a lot of the lower-level questions about storage format, issue layout, and representation can get worked out through use. The main thing is the movement: from "what area matters?" to "what exact thing should I do now?"
