package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportBeadsIssuesRequiresBeadsAndTaskgraphDirs(t *testing.T) {
	dir := t.TempDir()
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))

	_, err := ImportBeadsIssues(dir)
	if err == nil {
		t.Fatalf("expected error when .beads is missing")
	}
	if !strings.Contains(err.Error(), "expected .beads and .taskgraph") {
		t.Fatalf("expected missing dirs message, got: %v", err)
	}
}

func TestImportBeadsIssuesRequiresIssuesJSONL(t *testing.T) {
	dir := t.TempDir()
	mustMkdirAll(t, filepath.Join(dir, ".beads"))
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "")

	_, err := ImportBeadsIssues(dir)
	if err == nil {
		t.Fatalf("expected error when issues.jsonl is missing")
	}
	if !strings.Contains(err.Error(), ".beads/issues.jsonl") {
		t.Fatalf("expected missing issues.jsonl path in error, got: %v", err)
	}
}

func TestImportBeadsIssuesMapsStatusesAndSkipsTombstones(t *testing.T) {
	dir := t.TempDir()
	mustMkdirAll(t, filepath.Join(dir, ".beads"))
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "- [ ] existing\n")
	mustWrite(t, filepath.Join(dir, ".beads", "issues.jsonl"), strings.Join([]string{
		`{"id":"pl-1","title":"Open item","status":"open"}`,
		`{"id":"pl-2","title":"Closed item","status":"closed"}`,
		`{"id":"pl-3","title":"Gone","status":"tombstone"}`,
	}, "\n")+"\n")

	summary, err := ImportBeadsIssues(dir)
	if err != nil {
		t.Fatalf("ImportBeadsIssues returned err: %v", err)
	}
	if summary.Imported != 2 {
		t.Fatalf("expected 2 imported, got %d", summary.Imported)
	}
	if summary.SkippedTombstone != 1 {
		t.Fatalf("expected 1 tombstone skipped, got %d", summary.SkippedTombstone)
	}

	got := mustRead(t, filepath.Join(dir, ".taskgraph", "issues.md"))
	if !strings.Contains(got, "- [ ] [beads:pl-1] Open item\n") {
		t.Fatalf("expected open issue line, got: %q", got)
	}
	if !strings.Contains(got, "- [x] [beads:pl-2] Closed item\n") {
		t.Fatalf("expected closed issue line, got: %q", got)
	}
	if strings.Contains(got, "pl-3") {
		t.Fatalf("did not expect tombstone issue in output: %q", got)
	}
}

func TestImportBeadsIssuesErrorsOnMalformedJSONWithLineNumber(t *testing.T) {
	dir := t.TempDir()
	mustMkdirAll(t, filepath.Join(dir, ".beads"))
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "")
	mustWrite(t, filepath.Join(dir, ".beads", "issues.jsonl"), "{bad json}\n")

	_, err := ImportBeadsIssues(dir)
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Fatalf("expected line number in error, got: %v", err)
	}
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

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	return string(b)
}
