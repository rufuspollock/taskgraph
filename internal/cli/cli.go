package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"taskgraph/internal/indexer"
	"taskgraph/internal/migrate"
	"taskgraph/internal/project"
	"taskgraph/internal/tasks"
)

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
		return runInbox(stdout, stderr)
	case "list":
		return runList(args, stdout, stderr)
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
  add <text>        Add a task to .taskgraph/issues.md
  create <text>     Alias for add
  inbox             Print inbox checklist from .taskgraph/issues.md
  list [--all]      Print indexed checklist tasks from SQLite
  index             Build SQLite index from markdown files
  migrate-beads     Import .beads/issues.jsonl into .taskgraph/issues.md
  help              Show this help

EXAMPLES
  tg init
  tg add "buy milk"
  tg create "book dentist"
  tg inbox
  tg list
  tg index
  tg migrate-beads

NOTES
  - tg add auto-initializes .taskgraph if missing
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
	if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
		fmt.Fprintln(stderr, "usage: tg add <task text>")
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

	taskText := strings.TrimSpace(args[1])
	taskFile := filepath.Join(root, ".taskgraph", "issues.md")
	prefix, err := project.ReadPrefix(root)
	if err != nil {
		return err
	}
	if err := tasks.AppendTask(taskFile, prefix, taskText); err != nil {
		return err
	}
	if _, _, err := buildAndStoreIndex(root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Added task: %s\n", taskText)
	return nil
}

func runInbox(stdout io.Writer, stderr io.Writer) error {
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
		fmt.Fprintln(stdout, line)
	}
	return nil
}

func runList(args []string, stdout io.Writer, stderr io.Writer) error {
	includeClosed, err := parseListArgs(args[1:])
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
	nodes, err := indexer.ReadChecklistNodes(dbPath, includeClosed)
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

func parseListArgs(args []string) (bool, error) {
	includeClosed := false
	for _, arg := range args {
		switch arg {
		case "--all":
			includeClosed = true
		default:
			return false, fmt.Errorf("usage: tg list [--all]")
		}
	}
	return includeClosed, nil
}
