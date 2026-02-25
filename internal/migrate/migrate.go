package migrate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Summary reports migration counts for one import run.
type Summary struct {
	Imported         int
	SkippedTombstone int
	SkippedInvalid   int
}

type beadsIssue struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// ImportBeadsIssues imports ./.beads/issues.jsonl into ./.taskgraph/issues.md.
func ImportBeadsIssues(cwd string) (Summary, error) {
	summary := Summary{}

	beadsDir := filepath.Join(cwd, ".beads")
	taskgraphDir := filepath.Join(cwd, ".taskgraph")
	if !dirExists(beadsDir) || !dirExists(taskgraphDir) {
		return summary, fmt.Errorf("expected .beads and .taskgraph in current directory: %s", cwd)
	}

	inputPath := filepath.Join(beadsDir, "issues.jsonl")
	if _, err := os.Stat(inputPath); err != nil {
		return summary, fmt.Errorf("missing input file %s: %w", inputPath, err)
	}

	issuesPath := filepath.Join(taskgraphDir, "issues.md")
	in, err := os.Open(inputPath)
	if err != nil {
		return summary, fmt.Errorf("open input %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.OpenFile(issuesPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return summary, fmt.Errorf("open output %s: %w", issuesPath, err)
	}
	defer out.Close()

	s := bufio.NewScanner(in)
	lineNo := 0
	for s.Scan() {
		lineNo++
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		var issue beadsIssue
		if err := json.Unmarshal([]byte(line), &issue); err != nil {
			return summary, fmt.Errorf("parse %s line %d: %w", inputPath, lineNo, err)
		}

		issue.ID = strings.TrimSpace(issue.ID)
		issue.Title = strings.TrimSpace(issue.Title)
		if issue.ID == "" || issue.Title == "" {
			summary.SkippedInvalid++
			continue
		}

		status := strings.ToLower(strings.TrimSpace(issue.Status))
		if status == "tombstone" {
			summary.SkippedTombstone++
			continue
		}

		checkbox := " "
		if status == "closed" || status == "done" || status == "resolved" {
			checkbox = "x"
		}

		if _, err := fmt.Fprintf(out, "- [%s] [beads:%s] %s\n", checkbox, issue.ID, issue.Title); err != nil {
			return summary, fmt.Errorf("append to %s: %w", issuesPath, err)
		}
		summary.Imported++
	}
	if err := s.Err(); err != nil {
		return summary, fmt.Errorf("read %s: %w", inputPath, err)
	}

	return summary, nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
