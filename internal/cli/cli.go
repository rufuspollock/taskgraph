package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"taskgraph/internal/project"
	"taskgraph/internal/tasks"
)

// Run dispatches CLI commands.
func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: tg <command>")
		return errors.New("missing command")
	}

	switch args[0] {
	case "init":
		return runInit(stdout)
	case "add", "create":
		return runAdd(args, stdout, stderr)
	case "list":
		return runList(stdout, stderr)
	default:
		return fmt.Errorf("command not implemented: %s", args[0])
	}
}

func runInit(stdout io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	root, _, err := project.InitAt(cwd)
	if err != nil {
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

	cwd, err := os.Getwd()
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
	taskFile := filepath.Join(root, ".taskgraph", "tasks.md")
	if err := tasks.AppendTask(taskFile, taskText); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Added task: %s\n", taskText)
	return nil
}

func runList(stdout io.Writer, stderr io.Writer) error {
	cwd, err := os.Getwd()
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

	lines, err := tasks.ReadChecklistLines(filepath.Join(root, ".taskgraph", "tasks.md"))
	if err != nil {
		return err
	}
	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}
	return nil
}
