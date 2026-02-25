package tasks

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestAppendTaskWritesChecklistLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")

	if err := AppendTask(path, "tg", "first task"); err != nil {
		t.Fatalf("AppendTask returned err: %v", err)
	}

	got := readFile(t, path)
	assertTaskLineFormat(t, got, "tg", "first task")
}

func TestAppendTaskUsesBracketedIDAndDate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")

	if err := AppendTask(path, "taskgraph", "first task"); err != nil {
		t.Fatalf("AppendTask returned err: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "➕"+todayISO()) || !strings.Contains(got, "[task-") {
		t.Fatalf("expected date and bracketed ID, got: %q", got)
	}
}

func TestAppendTaskPreservesExistingLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	mustWrite(t, path, "- [ ] existing")

	if err := AppendTask(path, "demo", "second"); err != nil {
		t.Fatalf("AppendTask returned err: %v", err)
	}

	got := readFile(t, path)
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two lines, got %q", got)
	}
	if lines[0] != "- [ ] existing" {
		t.Fatalf("expected first line preserved, got %q", lines[0])
	}
	assertTaskLineFormat(t, lines[1]+"\n", "demo", "second")
}

func TestReadChecklistLinesReturnsInOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	mustWrite(t, path, "- [ ] a\n\n- [x] done\nnot a task\n- [ ] b\n")

	got, err := ReadChecklistLines(path)
	if err != nil {
		t.Fatalf("ReadChecklistLines returned err: %v", err)
	}

	want := []string{"- [ ] a", "- [x] done", "- [ ] b"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d mismatch: want %q got %q", i, want[i], got[i])
		}
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	return string(b)
}

func todayISO() string {
	return time.Now().Format("2006-01-02")
}

func assertTaskLineFormat(t *testing.T, line, prefix, text string) {
	t.Helper()
	pattern := `^\- \[ \] ➕` + regexp.QuoteMeta(todayISO()) + ` \[` + regexp.QuoteMeta(prefix) + `\-[0-9a-z]{3,8}\] ` + regexp.QuoteMeta(text) + `\n$`
	ok, err := regexp.MatchString(pattern, line)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !ok {
		t.Fatalf("line did not match expected format.\nline: %q\npattern: %q", line, pattern)
	}
}
