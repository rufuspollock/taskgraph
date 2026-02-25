package indexer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	headingPattern   = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	checklistPattern = regexp.MustCompile(`^\s*-\s*\[( |x|X)\]\s+(.*)$`)
)

// Node is one indexed markdown element.
type Node struct {
	ID         string
	Kind       string
	Title      string
	State      string
	Path       string
	Line       int
	ParentID   string
	Context    string
	SearchText string
	Source     string
}

// BuildNodes scans the root directory for markdown files and returns indexed nodes.
func BuildNodes(root string) ([]Node, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("root is required")
	}

	files, err := discoverMarkdownFiles(root)
	if err != nil {
		return nil, err
	}

	tasksPath := filepath.Join(root, ".taskgraph", "issues.md")
	if _, err := os.Stat(tasksPath); err == nil {
		files = appendUnique(files, tasksPath)
	}
	sort.Strings(files)

	var nodes []Node
	for _, absPath := range files {
		rel, err := filepath.Rel(root, absPath)
		if err != nil {
			return nil, err
		}
		rel = filepath.ToSlash(rel)
		source := "scan"
		if rel == ".taskgraph/issues.md" {
			source = "tasks_md"
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}
		fileNodes := indexMarkdown(string(content), rel, source)
		nodes = append(nodes, fileNodes...)
	}

	return nodes, nil
}

func FileNodeCount(nodes []Node) int {
	count := 0
	for _, n := range nodes {
		if n.Kind == "file" {
			count++
		}
	}
	return count
}

func discoverMarkdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if name == "node_modules" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(filepath.Ext(name), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func appendUnique(paths []string, p string) []string {
	for _, item := range paths {
		if item == p {
			return paths
		}
	}
	return append(paths, p)
}

func indexMarkdown(content, relPath, source string) []Node {
	fileTitle := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))
	fileID := buildNodeID(relPath, nil, 0, "file")
	nodes := []Node{{
		ID:         fileID,
		Kind:       "file",
		Title:      fileTitle,
		State:      "unknown",
		Path:       relPath,
		Line:       0,
		ParentID:   "",
		Context:    fileTitle,
		SearchText: normalizeSearch(fileTitle),
		Source:     source,
	}}

	type headingEntry struct {
		level int
		title string
		id    string
	}
	var stack []headingEntry

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNo := i + 1

		if m := headingPattern.FindStringSubmatch(line); len(m) == 3 {
			level := len(m[1])
			title := strings.TrimSpace(m[2])
			for len(stack) > 0 && stack[len(stack)-1].level >= level {
				stack = stack[:len(stack)-1]
			}
			parentID := fileID
			pathBits := []string{}
			for _, h := range stack {
				pathBits = append(pathBits, h.title)
				parentID = h.id
			}
			pathBits = append(pathBits, title)
			id := buildNodeID(relPath, pathBits, lineNo, "heading")
			stack = append(stack, headingEntry{level: level, title: title, id: id})
			context := buildContext(fileTitle, pathBits)
			nodes = append(nodes, Node{
				ID:         id,
				Kind:       "heading",
				Title:      title,
				State:      "unknown",
				Path:       relPath,
				Line:       lineNo,
				ParentID:   parentID,
				Context:    context,
				SearchText: normalizeSearch(context + " " + title),
				Source:     source,
			})
			continue
		}

		if m := checklistPattern.FindStringSubmatch(line); len(m) == 3 {
			checked := strings.EqualFold(m[1], "x")
			title := strings.TrimSpace(m[2])
			state := "open"
			if checked {
				state = "closed"
			}
			parentID := fileID
			pathBits := []string{}
			for _, h := range stack {
				pathBits = append(pathBits, h.title)
				parentID = h.id
			}
			pathBits = append(pathBits, title)
			id := buildNodeID(relPath, pathBits, lineNo, "checklist")
			context := buildContext(fileTitle, pathBits)
			nodes = append(nodes, Node{
				ID:         id,
				Kind:       "checklist",
				Title:      title,
				State:      state,
				Path:       relPath,
				Line:       lineNo,
				ParentID:   parentID,
				Context:    context,
				SearchText: normalizeSearch(context + " " + title),
				Source:     source,
			})
		}
	}

	return nodes
}

func buildContext(fileTitle string, pathBits []string) string {
	parts := []string{fileTitle}
	parts = append(parts, pathBits...)
	return strings.Join(parts, " > ")
}

func normalizeSearch(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func buildNodeID(path string, pathBits []string, line int, kind string) string {
	raw := path + "::" + strings.Join(pathBits, "::") + "::" + fmt.Sprint(line) + "::" + kind
	sum := sha1.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}
