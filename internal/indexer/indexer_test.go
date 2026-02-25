package indexer

import (
	"os"
	"path/filepath"
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
