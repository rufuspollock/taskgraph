package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"taskgraph/internal/indexer"
	"taskgraph/internal/migrate"
	"taskgraph/internal/project"
	"taskgraph/internal/tasks"
)

const addUsage = "usage: tg add <task text> [--labels a,b] [--type name]"

var graphLabelPattern = regexp.MustCompile(`(^|[\s(])#([A-Za-z0-9][A-Za-z0-9-]*)`)

// Run dispatches CLI commands.
func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprint(stdout, helpText())
		return nil
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, helpText())
		return nil
	case "init":
		return runInit(stdout)
	case "add", "create":
		return runAdd(args, stdout, stderr)
	case "inbox":
		return runInbox(args[1:], stdout, stderr)
	case "close":
		return runClose(args[1:], stdout, stderr)
	case "list":
		return runList(args, stdout, stderr)
	case "graph":
		return runGraph(args[1:], stdout, stderr)
	case "index":
		return runIndex(stdout, stderr)
	case "migrate-beads":
		return runMigrateBeads(stdout, stderr)
	default:
		return fmt.Errorf("unknown command: %s. Run 'tg --help'", args[0])
	}
}

func helpText() string {
	return `TaskGraph
  _____         _     ____                   _     
 |_   _|_ _ ___| | __/ ___|_ __ __ _ _ __  | |__  
   | |/ _` + "`" + ` / __| |/ / |  _| '__/ _` + "`" + ` | '_ \ | '_ \ 
   | | (_| \__ \   <| |_| | | | (_| | |_) || | | |
   |_|\__,_|___/_|\_\\____|_|  \__,_| .__/ |_| |_|
                                     |_|           

Local-first, AI-friendly task graph CLI.

USAGE
  tg <command> [args]
  tg -h | --help

COMMANDS
  init              Initialize .taskgraph in current directory
  add <text>        Add a task to .taskgraph/issues.md (supports --labels, --type)
  create <text>     Alias for add (supports --labels, --type)
  inbox [--all] [--label name]
                    Print inbox checklist from .taskgraph/issues.md
  close <id> [reason]
                    Close an inbox task in .taskgraph/issues.md
  list [--all] [--label name]
                    Print indexed checklist tasks from SQLite
  graph [--depth N] [--max-children N] [--all]
                    Print a compact graph overview from root nodes
  index             Build SQLite index from markdown files
  migrate-beads     Import .beads/issues.jsonl into .taskgraph/issues.md
  help              Show this help

EXAMPLES
  tg init
  tg add "buy milk"
  tg add "buy milk" --labels errands,home
  tg add "plan launch" --type epic
  tg create "book dentist"
  tg inbox
  tg inbox --label home
  tg close tg-abc "done on phone"
  tg list
  tg list --label errands
  tg graph
  tg graph --depth 3 --max-children 4
  tg index
  tg migrate-beads

NOTES
  - tg add auto-initializes .taskgraph if missing
  - use --type with one allowed task type per task
  - inbox is stored in .taskgraph/issues.md
  - index DB is stored in .taskgraph/taskgraph.db
`
}

func runInit(stdout io.Writer) error {
	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, _, err := project.InitAt(cwd)
	if err != nil {
		return err
	}
	if _, _, err := buildAndStoreIndex(root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Initialized .taskgraph in %s\n", root)
	return nil
}

func runAdd(args []string, stdout io.Writer, stderr io.Writer) error {
	taskText, labels, taskType, err := parseAddArgs(args[1:])
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	if strings.TrimSpace(taskText) == "" {
		fmt.Fprintln(stderr, addUsage)
		return errors.New("missing task text")
	}

	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		var created bool
		root, created, err = project.InitAt(cwd)
		if err != nil {
			return err
		}
		if created {
			fmt.Fprintf(stdout, "Initialized .taskgraph in %s\n", root)
		}
	}

	taskFile := filepath.Join(root, ".taskgraph", "issues.md")
	prefix, err := project.ReadPrefix(root)
	if err != nil {
		return err
	}
	resolvedType, cleanLabels, err := tasks.ResolveTaskType(taskText, labels, taskType)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	if resolvedType != "" {
		allowed, err := project.ReadAllowedIssueTypes(root)
		if err != nil {
			return err
		}
		if !containsString(allowed, resolvedType) {
			sort.Strings(allowed)
			msg := fmt.Sprintf("unknown task type: %s (allowed: %s)", resolvedType, strings.Join(allowed, ", "))
			fmt.Fprintln(stderr, msg)
			return errors.New(msg)
		}
	}
	if err := tasks.AppendTask(taskFile, prefix, taskText, cleanLabels, resolvedType); err != nil {
		return err
	}
	if _, _, err := buildAndStoreIndex(root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Added task: %s\n", taskText)
	return nil
}

func runInbox(args []string, stdout io.Writer, stderr io.Writer) error {
	includeClosed, requiredLabels, err := parseInboxArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}

	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	lines, err := tasks.ReadChecklistLines(filepath.Join(root, ".taskgraph", "issues.md"))
	if err != nil {
		return err
	}
	for _, line := range lines {
		if !includeClosed && strings.HasPrefix(line, "- [x] ") {
			continue
		}
		if !hasAllLabels(tasks.ExtractLabels(line), requiredLabels) {
			continue
		}
		fmt.Fprintln(stdout, line)
	}
	return nil
}

func runClose(args []string, stdout io.Writer, stderr io.Writer) error {
	id, reason, err := parseCloseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}

	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	taskFile := filepath.Join(root, ".taskgraph", "issues.md")
	if err := tasks.CloseTask(taskFile, id, reason); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	if _, _, err := buildAndStoreIndex(root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Closed task: %s\n", id)
	return nil
}

func runList(args []string, stdout io.Writer, stderr io.Writer) error {
	includeClosed, requiredLabels, err := parseListArgs(args[1:])
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}

	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	dbPath := filepath.Join(root, ".taskgraph", "taskgraph.db")
	nodes, err := indexer.ReadChecklistNodes(dbPath, includeClosed, requiredLabels)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		mark := " "
		if n.State == "closed" {
			mark = "x"
		}
		fmt.Fprintf(stdout, "- [%s] %s (%s:%d)\n", mark, n.Title, n.Path, n.Line)
	}
	return nil
}

func runIndex(stdout io.Writer, stderr io.Writer) error {
	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	fileCount, nodeCount, err := buildAndStoreIndex(root)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		stdout,
		"Indexed %d files, %d nodes into %s\n",
		fileCount,
		nodeCount,
		filepath.Join(root, ".taskgraph", "taskgraph.db"),
	)
	return nil
}

func runGraph(args []string, stdout io.Writer, stderr io.Writer) error {
	depth, maxChildren, includeClosed, err := parseGraphArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}

	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	dbPath := filepath.Join(root, ".taskgraph", "taskgraph.db")
	nodes, err := indexer.ReadGraphNodes(dbPath)
	if err != nil {
		return err
	}

	selectedRoots := selectGraphRoots(nodes)
	children := graphChildren(nodes)
	byID := make(map[string]indexer.Node, len(nodes))
	for _, node := range nodes {
		byID[node.ID] = node
	}
	visible := graphVisibility(byID, children, selectedRoots, includeClosed)
	for _, node := range nodes {
		if !selectedRoots[node.ID] {
			continue
		}
		if !visible[node.ID] {
			continue
		}
		fmt.Fprintln(stdout, formatGraphNode(node))
		renderGraphChildren(stdout, children, selectedRoots, visible, node.ID, 1, depth, maxChildren)
	}
	return nil
}

func runMigrateBeads(stdout io.Writer, stderr io.Writer) error {
	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}

	summary, err := migrate.ImportBeadsIssues(cwd)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}

	fmt.Fprintf(
		stdout,
		"Imported %d issues (%d tombstones skipped, %d invalid skipped)\n",
		summary.Imported,
		summary.SkippedTombstone,
		summary.SkippedInvalid,
	)
	return nil
}

func buildAndStoreIndex(root string) (int, int, error) {
	nodes, err := indexer.BuildNodes(root)
	if err != nil {
		return 0, 0, err
	}
	dbPath := filepath.Join(root, ".taskgraph", "taskgraph.db")
	if err := indexer.RebuildSQLite(dbPath, nodes); err != nil {
		return 0, 0, err
	}
	return indexer.FileNodeCount(nodes), len(nodes), nil
}

func effectiveCWD() (string, error) {
	if v := strings.TrimSpace(os.Getenv("TG_CWD")); v != "" {
		return v, nil
	}
	return os.Getwd()
}

func parseListArgs(args []string) (bool, []string, error) {
	includeClosed := false
	var labels []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all":
			includeClosed = true
		case "--label":
			if i+1 >= len(args) {
				return false, nil, fmt.Errorf("usage: tg list [--all] [--label name]")
			}
			label := tasks.NormalizeLabelsCSV(args[i+1])
			if len(label) == 0 {
				return false, nil, fmt.Errorf("usage: tg list [--all] [--label name]")
			}
			labels = append(labels, label...)
			i++
		default:
			return false, nil, fmt.Errorf("usage: tg list [--all] [--label name]")
		}
	}
	return includeClosed, tasks.MergeLabels(labels), nil
}

func parseInboxArgs(args []string) (bool, []string, error) {
	includeClosed := false
	var labels []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all":
			includeClosed = true
		case "--label":
			if i+1 >= len(args) {
				return false, nil, fmt.Errorf("usage: tg inbox [--all] [--label name]")
			}
			label := tasks.NormalizeLabelsCSV(args[i+1])
			if len(label) == 0 {
				return false, nil, fmt.Errorf("usage: tg inbox [--all] [--label name]")
			}
			labels = append(labels, label...)
			i++
		default:
			return false, nil, fmt.Errorf("usage: tg inbox [--all] [--label name]")
		}
	}
	return includeClosed, tasks.MergeLabels(labels), nil
}

func parseCloseArgs(args []string) (string, string, error) {
	if len(args) < 1 {
		return "", "", fmt.Errorf("usage: tg close <id> [reason]")
	}

	id := strings.TrimSpace(args[0])
	if id == "" {
		return "", "", fmt.Errorf("usage: tg close <id> [reason]")
	}

	reason := strings.TrimSpace(strings.Join(args[1:], " "))
	if strings.EqualFold(reason, "null") {
		reason = ""
	}
	return id, reason, nil
}

func parseGraphArgs(args []string) (int, int, bool, error) {
	depth := 4
	maxChildren := 5
	includeClosed := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all":
			includeClosed = true
		case "--depth":
			if i+1 >= len(args) {
				return 0, 0, false, fmt.Errorf("usage: tg graph [--depth N] [--max-children N] [--all]")
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil || value < 1 {
				return 0, 0, false, fmt.Errorf("usage: tg graph [--depth N] [--max-children N] [--all]")
			}
			depth = value
			i++
		case "--max-children":
			if i+1 >= len(args) {
				return 0, 0, false, fmt.Errorf("usage: tg graph [--depth N] [--max-children N] [--all]")
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil || value < 1 {
				return 0, 0, false, fmt.Errorf("usage: tg graph [--depth N] [--max-children N] [--all]")
			}
			maxChildren = value
			i++
		default:
			return 0, 0, false, fmt.Errorf("usage: tg graph [--depth N] [--max-children N] [--all]")
		}
	}

	return depth, maxChildren, includeClosed, nil
}

var graphRootTypes = map[string]bool{
	"idea":       true,
	"initiative": true,
	"project":    true,
	"product":    true,
	"epic":       true,
}

func selectGraphRoots(nodes []indexer.Node) map[string]bool {
	byID := make(map[string]indexer.Node, len(nodes))
	for _, node := range nodes {
		byID[node.ID] = node
	}
	children := graphChildren(nodes)

	selected := make(map[string]bool)
	for _, node := range nodes {
		if !isTypedGraphRoot(node) || hasTypedGraphRootAncestor(byID, node) {
			continue
		}
		selected[node.ID] = true
	}

	for _, node := range nodes {
		if node.Kind != "file" {
			continue
		}
		roots := collectStructuralRoots(node, children, selected)
		for _, root := range roots {
			selected[root.ID] = true
		}
	}

	return selected
}

func collectStructuralRoots(node indexer.Node, children map[string][]indexer.Node, typedRoots map[string]bool) []indexer.Node {
	if node.Kind == "file" {
		var out []indexer.Node
		for _, child := range children[node.ID] {
			out = append(out, collectStructuralRoots(child, children, typedRoots)...)
		}
		return out
	}
	if typedRoots[node.ID] {
		return nil
	}

	remainingChildren := make([]indexer.Node, 0)
	for _, child := range children[node.ID] {
		if typedRoots[child.ID] {
			continue
		}
		remainingChildren = append(remainingChildren, child)
	}

	descendantRoots := make([]indexer.Node, 0)
	for _, child := range remainingChildren {
		descendantRoots = append(descendantRoots, collectStructuralRoots(child, children, typedRoots)...)
	}
	if len(remainingChildren) == 0 {
		return nil
	}
	if len(remainingChildren) == 1 && len(descendantRoots) == 0 {
		if len(nonTypedChildren(children, remainingChildren[0].ID, typedRoots)) == 0 {
			return nil
		}
	}
	allHeadings := true
	for _, child := range remainingChildren {
		if child.Kind != "heading" {
			allHeadings = false
			break
		}
	}
	if node.Kind == "heading" && allHeadings {
		return descendantRoots
	}
	return []indexer.Node{node}
}

func nonTypedChildren(children map[string][]indexer.Node, parentID string, typedRoots map[string]bool) []indexer.Node {
	out := make([]indexer.Node, 0)
	for _, child := range children[parentID] {
		if typedRoots[child.ID] {
			continue
		}
		out = append(out, child)
	}
	return out
}

func graphChildren(nodes []indexer.Node) map[string][]indexer.Node {
	out := make(map[string][]indexer.Node)
	for _, node := range nodes {
		out[node.ParentID] = append(out[node.ParentID], node)
	}
	return out
}

func hasTypedGraphRootAncestor(byID map[string]indexer.Node, node indexer.Node) bool {
	parentID := node.ParentID
	for parentID != "" {
		parent, ok := byID[parentID]
		if !ok {
			return false
		}
		if isTypedGraphRoot(parent) {
			return true
		}
		parentID = parent.ParentID
	}
	return false
}

func isTypedGraphRoot(node indexer.Node) bool {
	taskType, err := tasks.ExtractTaskTypeFromLabels(node.Labels)
	if err != nil {
		return false
	}
	return graphRootTypes[taskType]
}

func renderGraphChildren(stdout io.Writer, children map[string][]indexer.Node, roots map[string]bool, visible map[string]bool, parentID string, level int, maxDepth int, maxChildren int) {
	if level > maxDepth {
		return
	}
	visibleChildren := make([]indexer.Node, 0)
	for _, child := range children[parentID] {
		if roots[child.ID] {
			continue
		}
		if !visible[child.ID] {
			continue
		}
		visibleChildren = append(visibleChildren, child)
	}
	hiddenCount := 0
	if len(visibleChildren) > maxChildren {
		hiddenCount = len(visibleChildren) - maxChildren
		visibleChildren = visibleChildren[:maxChildren]
	}
	for _, child := range visibleChildren {
		fmt.Fprintf(stdout, "%s%s\n", strings.Repeat("  ", level), formatGraphNode(child))
		renderGraphChildren(stdout, children, roots, visible, child.ID, level+1, maxDepth, maxChildren)
	}
	if hiddenCount > 0 {
		fmt.Fprintf(stdout, "%s... %d more\n", strings.Repeat("  ", level), hiddenCount)
	}
}

func graphVisibility(byID map[string]indexer.Node, children map[string][]indexer.Node, roots map[string]bool, includeClosed bool) map[string]bool {
	visible := make(map[string]bool, len(byID))
	visiting := make(map[string]bool, len(byID))

	var visit func(id string) bool
	visit = func(id string) bool {
		if v, ok := visible[id]; ok {
			return v
		}
		if visiting[id] {
			return false
		}
		visiting[id] = true
		defer delete(visiting, id)

		node, ok := byID[id]
		if !ok {
			return false
		}

		hasVisibleChild := false
		for _, child := range children[id] {
			if roots[child.ID] {
				continue
			}
			if visit(child.ID) {
				hasVisibleChild = true
			}
		}

		result := false
		switch node.Kind {
		case "file", "heading":
			result = hasVisibleChild
		case "checklist":
			if len(children[id]) > 0 {
				result = hasVisibleChild
			} else {
				result = includeClosed || node.State != "closed"
			}
			if roots[id] && isTypedGraphRoot(node) {
				result = true
			}
		default:
			result = hasVisibleChild
		}

		visible[id] = result
		return result
	}

	for id := range byID {
		visit(id)
	}
	return visible
}

func formatGraphNode(node indexer.Node) string {
	title := cleanGraphTitle(node.Title)
	if taskType, err := tasks.ExtractTaskTypeFromLabels(node.Labels); err == nil && graphRootTypes[taskType] {
		return fmt.Sprintf("[%s] %s", taskType, title)
	}
	return title
}

func cleanGraphTitle(title string) string {
	cleaned := graphLabelPattern.ReplaceAllString(title, "$1")
	return strings.Join(strings.Fields(cleaned), " ")
}

func hasAllLabels(actual, required []string) bool {
	if len(required) == 0 {
		return true
	}
	seen := map[string]bool{}
	for _, label := range actual {
		seen[label] = true
	}
	for _, label := range required {
		if !seen[label] {
			return false
		}
	}
	return true
}

func parseAddArgs(args []string) (string, []string, string, error) {
	var textParts []string
	var labels []string
	taskType := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--labels", "-l":
			if i+1 >= len(args) {
				return "", nil, "", fmt.Errorf(addUsage)
			}
			labels = append(labels, tasks.NormalizeLabelsCSV(args[i+1])...)
			i++
		case "--type", "-t":
			if i+1 >= len(args) {
				return "", nil, "", fmt.Errorf(addUsage)
			}
			if taskType != "" {
				return "", nil, "", fmt.Errorf("multiple --type values are not allowed")
			}
			taskType = args[i+1]
			i++
		default:
			textParts = append(textParts, args[i])
		}
	}

	return strings.TrimSpace(strings.Join(textParts, " ")), tasks.MergeLabels(labels), taskType, nil
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
