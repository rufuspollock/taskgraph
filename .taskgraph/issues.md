- [ ] ➕2026-02-25 Make logo fancier by making it graph-y or animated
  Example:
         o____         _     ____o
         |_   _|_ _ ___| | __/ ___|_ __ __ _ _ __o
           | |/ _` / __| |/ / |  _| '__/ _` | '_ \
           o | (_| \__ \   <| |_| | | | (_| | |_) |
           |_|\__,_|___/_|\_\\__o_|_|  \__,_| .__/
                                             |_o|
  or like this ...
    ██████╗  █████╗ ██████╗ ███████╗██████╗  ██████╗██╗     ██╗██████╗
    ██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔════╝██║     ██║██╔══██╗
    ██████╔╝███████║██████╔╝█████╗  ██████╔╝██║     ██║     ██║██████╔╝
    ██╔═══╝ ██╔══██║██╔═══╝ ██╔══╝  ██╔══██╗██║     ██║     ██║██╔═══╝
    ██║     ██║  ██║██║     ███████╗██║  ██║╚██████╗███████╗██║██║
    ╚═╝     ╚═╝  ╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝ ╚═════╝╚══════╝╚═╝╚═╝
- [x] ➕2026-02-25 [task-zx9] Have task ids like in beads. Bonus is supporting prefix
- [x] ➕2026-02-25 [task-k1w] beads migration or integration in some form. e.g. we either migrate items or we at least load beads items into central db. This relates to multi-source support idea. **✅2026-03-15 done weeks ago**
- [ ] ➕2026-02-28 [task-zc5] Need to store a version of our format and/or db so that we can do migrations easily
- [ ] ➕2026-02-28 [task-umi] Support migrations now that is in live use
- [ ] ➕2026-02-28 [task-imb] switch to one md file per issue approach b/c better for conflicts and allows me to immediately add subtasks by opening the file
- [ ] ➕2026-02-28 [task-mwm] handle case where we don't have ids for tasks e.g. i have hand created etc. why do we need ids btw?
- [ ] ➕2026-02-28 [task-c98] Support task dependency and subtasks - already done implicitly i think or at least in db
- [x] ➕2026-02-28 [task-rfa] Support for labels **✅2026-03-03 in 8d8fc89808d89786ab956274e1b376c3edbc52c4**
- [x] ➕2026-02-28 [task-fch] Support for task types **✅2026-03-09 fixed in b07b4dcb242ebe49f3e192d7f162fe4bb8eb6c8a**
- [ ] ➕2026-02-28 [task-rps] Support for running tg from anywhere and having a default project and being to specify the project/repo this is on
- [x] ➕2026-02-28 [task-zm8] tg close command to close issues **✅2026-03-03 in c97dfffd1b356b94dd7fa4086727dd05190b5416**
- [ ] ➕2026-02-28 [task-pyv] taskgraph website and landing page
- [x] ➕2026-02-28 [task-vai] Register taskgraph.dev **✅2026-03-03**
- [ ] ➕2026-03-02 [task-x7g] I want to parse out time estimates which have a timer emoji next to them, and then time. I also want to pass out things like the next emoji, which means next. Basically, I have a bunch of emojis that I'm using in their meaning, which I should also pause out at some point.
- [ ] ➕2026-03-03 [task-w2m] tg day command: when starting my day i want to quickly (less than 5-10m) generate a backlog of tasks for my day starting with the *deep work* and then moving to more day to day items
- [ ] ➕2026-03-03 [task-zki] document patterns of having projects in projects folder and project being GTD project i.e. anything with multiple subtasks
- [ ] ➕2026-03-03 [task-wbv] reflect on how we doc patterns vs implement stuff in actual tooling
- [ ] ➕2026-03-03 [tgcl-vpf] support closing non-inbox tasks by indexed location with safe write-back semantics
- [ ] ➕2026-03-06 [task-u6e] support for task status in `[ ]` section as per https://publish.obsidian.md/tasks/Other+Plugins/Dataview
- [x] ➕2026-03-15 [task-lei] Think through task types: write mini essay on projects vs products. GTD 'projects' range from 'book hair appointment' to full products - we need a useful distinction between these scales #design #concepts **✅2026-03-15 see projects-vs-tasks-2026-03-15.md**
- [ ] ➕2026-03-15 [task-qnq] tg index and tg graph don't handle projects-as-whole-files well. A project file (e.g. fixtures/rufus-projects/*.md) should appear as a node but currently doesn't. Need to diagnose whether this is an index issue or graph issue #bug #graph #index
- [ ] ➕2026-03-15 [task-drm] Create test examples for tg graph to systematically work out what's wrong - build a small test project tree and verify graph output makes sense #graph #testing
- [ ] ➕2026-03-15 [task-hu6] Feature: support ignoring subdirectories from tg index (e.g. .tgignore or config option to exclude paths like node_modules, attic, etc.) #feature #index
- [ ] ➕2026-03-15 [task-ctb] Meta: tg is good for capture but not yet for planning/structuring work. What's missing to go from inbox to actionable plans? Probably needs AI integration - e.g. 'tg plan' that uses AI to help prioritize, sequence, and structure tasks from the graph #meta #planning #ai
- [ ] ➕2026-03-15 [task-g8h] tg close supports search so i can be typing and when i find the task i want i can close it
- [ ] ➕2026-03-15 [task-vhe] Feature: support exclude paths in config.yml for tg index (e.g. exclude: [fixtures/, sandbox/, docs/plans/]) so irrelevant markdown files are not indexed. Separate from project inference - this is about controlling what gets indexed at all. #feature #index #config
