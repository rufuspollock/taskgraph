package project

import (
	"os"
	"path/filepath"
	"strings"
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

func TestFindTaskgraphRootStopsAtGitRepoBoundary(t *testing.T) {
	root := t.TempDir()
	outerTaskgraph := filepath.Join(root, ".taskgraph")
	mustMkdirAll(t, outerTaskgraph)

	repoRoot := filepath.Join(root, "other-repo")
	mustMkdirAll(t, filepath.Join(repoRoot, ".git"))
	deep := filepath.Join(repoRoot, "a", "b")
	mustMkdirAll(t, deep)

	got, found, err := FindTaskgraphRoot(deep)
	if err != nil {
		t.Fatalf("FindTaskgraphRoot returned err: %v", err)
	}
	if found {
		t.Fatalf("expected found=false (should not cross git repo root), got root=%q", got)
	}
}

func TestFindTaskgraphRootRespectsTaskgraphInsideRepo(t *testing.T) {
	root := t.TempDir()
	repoRoot := filepath.Join(root, "repo")
	mustMkdirAll(t, filepath.Join(repoRoot, ".git"))
	mustMkdirAll(t, filepath.Join(repoRoot, ".taskgraph"))
	deep := filepath.Join(repoRoot, "x", "y")
	mustMkdirAll(t, deep)

	got, found, err := FindTaskgraphRoot(deep)
	if err != nil {
		t.Fatalf("FindTaskgraphRoot returned err: %v", err)
	}
	if !found {
		t.Fatalf("expected found=true")
	}
	if got != repoRoot {
		t.Fatalf("expected root %q, got %q", repoRoot, got)
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
	prefix, err := ReadPrefix(root)
	if err != nil {
		t.Fatalf("ReadPrefix returned err: %v", err)
	}
	if len(prefix) == 0 || len(prefix) > 4 {
		t.Fatalf("expected prefix length 1..4, got %q", prefix)
	}

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

func TestNormalizePrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"TaskGraph", "task"},
		{"My Project!", "mypr"},
		{"a", "a"},
		{"", "tg"},
		{"---", "tg"},
		{"AB12Z", "ab12"},
	}

	for _, tt := range tests {
		got := normalizePrefix(tt.in)
		if got != tt.want {
			t.Fatalf("normalizePrefix(%q)=%q want %q", tt.in, got, tt.want)
		}
	}
}

func TestReadPrefixFromConfig(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, ".taskgraph"))
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "config.yml"), []byte("issue-prefix: demo\n"), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	got, err := ReadPrefix(root)
	if err != nil {
		t.Fatalf("ReadPrefix returned err: %v", err)
	}
	if got != "demo" {
		t.Fatalf("expected demo prefix, got %q", got)
	}
}

func TestReadPrefixBackwardCompatibleWithOldKey(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, ".taskgraph"))
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "config.yml"), []byte("prefix: old\n"), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	got, err := ReadPrefix(root)
	if err != nil {
		t.Fatalf("ReadPrefix returned err: %v", err)
	}
	if got != "old" {
		t.Fatalf("expected old prefix, got %q", got)
	}
}

func TestEnsureConfigBackfillsMissingPrefix(t *testing.T) {
	root := t.TempDir()
	taskgraphDir := filepath.Join(root, ".taskgraph")
	mustMkdirAll(t, taskgraphDir)
	configPath := filepath.Join(taskgraphDir, "config.yml")
	if err := os.WriteFile(configPath, []byte("# empty\n"), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	if err := ensureConfig(configPath, root); err != nil {
		t.Fatalf("ensureConfig returned err: %v", err)
	}
	b, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config failed: %v", err)
	}
	if !strings.Contains(string(b), "prefix: ") {
		t.Fatalf("expected config to contain prefix, got %q", string(b))
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
