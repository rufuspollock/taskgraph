package indexer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildNodesScansMarkdownAndParsesHierarchy(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, ".taskgraph"))
	mustWrite(t, filepath.Join(root, ".taskgraph", "issues.md"), "- [ ] Captured task\n")
	mustWrite(t, filepath.Join(root, "notes.md"), "# Project\n\n## Build\n- [ ] Ship\n- [x] Done\n")
	mustMkdirAll(t, filepath.Join(root, ".hidden"))
	mustWrite(t, filepath.Join(root, ".hidden", "ignored.md"), "- [ ] Ignore hidden\n")
	mustMkdirAll(t, filepath.Join(root, "node_modules"))
	mustWrite(t, filepath.Join(root, "node_modules", "ignored.md"), "- [ ] Ignore deps\n")

	nodes, err := BuildNodes(root)
	if err != nil {
		t.Fatalf("BuildNodes returned error: %v", err)
	}
	if len(nodes) == 0 {
		t.Fatalf("expected nodes")
	}

	// Included sources.
	assertHasNodePath(t, nodes, "notes.md")
	assertHasNodePath(t, nodes, ".taskgraph/issues.md")

	// Excluded paths.
	assertNoNodePath(t, nodes, ".hidden/ignored.md")
	assertNoNodePath(t, nodes, "node_modules/ignored.md")

	// Checklist states from notes.md.
	assertHasChecklist(t, nodes, "notes.md", "Ship", "open")
	assertHasChecklist(t, nodes, "notes.md", "Done", "closed")
}

func TestBuildNodesExtractsChecklistLabels(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "notes.md"), "- [ ] Ship launch #flowershow #abc\n")

	nodes, err := BuildNodes(root)
	if err != nil {
		t.Fatalf("BuildNodes returned error: %v", err)
	}

	for _, n := range nodes {
		if n.Kind == "checklist" && n.Path == "notes.md" && n.Title == "Ship launch #flowershow #abc" {
			want := []string{"flowershow", "abc"}
			if !reflect.DeepEqual(n.Labels, want) {
				t.Fatalf("got labels %v want %v", n.Labels, want)
			}
			return
		}
	}
	t.Fatalf("expected checklist node with labels")
}

func TestBuildNodesCreatesChecklistParentChildRelationships(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "notes.md"), "# Project\n- [ ] Parent\n  - [ ] Child\n    - [ ] Grandchild\n- [ ] Sibling\n")

	nodes, err := BuildNodes(root)
	if err != nil {
		t.Fatalf("BuildNodes returned error: %v", err)
	}

	project := findNodeByKindAndTitle(t, nodes, "heading", "Project")
	parent := findNodeByKindAndTitle(t, nodes, "checklist", "Parent")
	child := findNodeByKindAndTitle(t, nodes, "checklist", "Child")
	grandchild := findNodeByKindAndTitle(t, nodes, "checklist", "Grandchild")
	sibling := findNodeByKindAndTitle(t, nodes, "checklist", "Sibling")

	if parent.ParentID != project.ID {
		t.Fatalf("parent parent_id = %q, want project heading id %q", parent.ParentID, project.ID)
	}
	if child.ParentID != parent.ID {
		t.Fatalf("child parent_id = %q, want parent checklist id %q", child.ParentID, parent.ID)
	}
	if grandchild.ParentID != child.ID {
		t.Fatalf("grandchild parent_id = %q, want child checklist id %q", grandchild.ParentID, child.ID)
	}
	if sibling.ParentID != project.ID {
		t.Fatalf("sibling parent_id = %q, want project heading id %q", sibling.ParentID, project.ID)
	}
}

func TestBuildNodesDoesNotCarryChecklistParentAcrossHeadings(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "notes.md"), "## One\n- [ ] Parent\n## Two\n  - [ ] Child\n")

	nodes, err := BuildNodes(root)
	if err != nil {
		t.Fatalf("BuildNodes returned error: %v", err)
	}

	two := findNodeByKindAndTitle(t, nodes, "heading", "Two")
	child := findNodeByKindAndTitle(t, nodes, "checklist", "Child")
	if child.ParentID != two.ID {
		t.Fatalf("child parent_id = %q, want heading Two id %q", child.ParentID, two.ID)
	}
}

func assertHasNodePath(t *testing.T, nodes []Node, want string) {
	t.Helper()
	for _, n := range nodes {
		if n.Path == want {
			return
		}
	}
	t.Fatalf("expected a node with path %q", want)
}

func assertNoNodePath(t *testing.T, nodes []Node, path string) {
	t.Helper()
	for _, n := range nodes {
		if n.Path == path {
			t.Fatalf("unexpected node with path %q", path)
		}
	}
}

func assertHasChecklist(t *testing.T, nodes []Node, path, title, state string) {
	t.Helper()
	for _, n := range nodes {
		if n.Kind == "checklist" && n.Path == path && n.Title == title && n.State == state {
			return
		}
	}
	t.Fatalf("expected checklist path=%q title=%q state=%q", path, title, state)
}

func findNodeByKindAndTitle(t *testing.T, nodes []Node, kind, title string) Node {
	t.Helper()
	for _, n := range nodes {
		if n.Kind == kind && n.Title == title {
			return n
		}
	}
	t.Fatalf("expected node kind=%q title=%q", kind, title)
	return Node{}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}
