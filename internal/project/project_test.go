package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindTaskgraphRootFindsNearestAncestor(t *testing.T) {
	root := t.TempDir()
	nearest := filepath.Join(root, "a", "b")
	deep := filepath.Join(nearest, "c", "d")

	mustMkdirAll(t, filepath.Join(root, ".taskgraph"))
	mustMkdirAll(t, filepath.Join(nearest, ".taskgraph"))
	mustMkdirAll(t, deep)

	got, found, err := FindTaskgraphRoot(deep)
	if err != nil {
		t.Fatalf("FindTaskgraphRoot returned err: %v", err)
	}
	if !found {
		t.Fatalf("expected found=true")
	}
	if got != nearest {
		t.Fatalf("expected nearest root %q, got %q", nearest, got)
	}
}

func TestFindTaskgraphRootNotFound(t *testing.T) {
	root := t.TempDir()
	deep := filepath.Join(root, "x", "y")
	mustMkdirAll(t, deep)

	got, found, err := FindTaskgraphRoot(deep)
	if err != nil {
		t.Fatalf("FindTaskgraphRoot returned err: %v", err)
	}
	if found {
		t.Fatalf("expected found=false, got root=%q", got)
	}
	if got != "" {
		t.Fatalf("expected empty root when not found, got %q", got)
	}
}

func TestInitAtCreatesFilesAndIsIdempotent(t *testing.T) {
	root := t.TempDir()

	gotRoot, created, err := InitAt(root)
	if err != nil {
		t.Fatalf("InitAt returned err: %v", err)
	}
	if !created {
		t.Fatalf("expected created=true on first init")
	}
	if gotRoot != root {
		t.Fatalf("expected root %q, got %q", root, gotRoot)
	}

	assertExists(t, filepath.Join(root, ".taskgraph", "config.yml"))
	assertExists(t, filepath.Join(root, ".taskgraph", "tasks.md"))

	gotRoot2, created2, err := InitAt(root)
	if err != nil {
		t.Fatalf("second InitAt returned err: %v", err)
	}
	if created2 {
		t.Fatalf("expected created=false on second init")
	}
	if gotRoot2 != root {
		t.Fatalf("expected root %q, got %q", root, gotRoot2)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %q failed: %v", path, err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %q to exist: %v", path, err)
	}
}
