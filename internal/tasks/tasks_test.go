package tasks

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendTaskWritesChecklistLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")

	if err := AppendTask(path, "first task"); err != nil {
		t.Fatalf("AppendTask returned err: %v", err)
	}

	got := readFile(t, path)
	want := "- [ ] ➕" + todayISO() + " first task\n"
	if got != want {
		t.Fatalf("unexpected file contents\nwant: %q\ngot:  %q", want, got)
	}
}

func TestAppendTaskPreservesExistingLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	mustWrite(t, path, "- [ ] existing")

	if err := AppendTask(path, "second"); err != nil {
		t.Fatalf("AppendTask returned err: %v", err)
	}

	got := readFile(t, path)
	want := "- [ ] existing\n- [ ] ➕" + todayISO() + " second\n"
	if got != want {
		t.Fatalf("unexpected file contents\nwant: %q\ngot:  %q", want, got)
	}
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
