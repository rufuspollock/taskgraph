package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoArgsPrintsHelp(t *testing.T) {
	stdout, stderr, err := run([]string{})
	if err != nil {
		t.Fatalf("expected nil error, got %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "TaskGraph") || !strings.Contains(stdout, "USAGE") {
		t.Fatalf("expected help output, got %q", stdout)
	}
}

func TestHelpFlagsPrintHelp(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		stdout, stderr, err := run([]string{arg})
		if err != nil {
			t.Fatalf("expected nil error for %q, got %v stderr=%q", arg, err, stderr)
		}
		if !strings.Contains(stdout, "COMMANDS") || !strings.Contains(stdout, "tg add") {
			t.Fatalf("expected help output for %q, got %q", arg, stdout)
		}
	}
}

func TestUnknownCommandIncludesHelpHint(t *testing.T) {
	_, stderr, err := run([]string{"bogus"})
	if err == nil {
		t.Fatalf("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "Run 'tg --help'") {
		t.Fatalf("expected help hint in error, got %q", err.Error())
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
}

func TestInitCreatesTaskgraphFiles(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	stdout, stderr, err := run([]string{"init"})
	if err != nil {
		t.Fatalf("init returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "Initialized .taskgraph") {
		t.Fatalf("expected init output, got %q", stdout)
	}
	assertExists(t, filepath.Join(dir, ".taskgraph", "config.yml"))
	assertExists(t, filepath.Join(dir, ".taskgraph", "tasks.md"))
}

func TestAddAutoInitsWhenMissing(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	stdout, stderr, err := run([]string{"add", "first task"})
	if err != nil {
		t.Fatalf("add returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "Initialized .taskgraph") {
		t.Fatalf("expected auto-init notice, got %q", stdout)
	}
	if !strings.Contains(stdout, "Added task: first task") {
		t.Fatalf("expected add confirmation, got %q", stdout)
	}

	content := readFile(t, filepath.Join(dir, ".taskgraph", "tasks.md"))
	if content != "- [ ] first task\n" {
		t.Fatalf("unexpected tasks.md content: %q", content)
	}
}

func TestAddUsesNearestAncestorTaskgraph(t *testing.T) {
	root := t.TempDir()
	_, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(root, ".taskgraph"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "config.yml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "tasks.md"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	chdir(t, nested)

	stdout, stderr, err := run([]string{"add", "from nested"})
	if err != nil {
		t.Fatalf("add returned err: %v stderr=%q", err, stderr)
	}
	if strings.Contains(stdout, "Initialized .taskgraph") {
		t.Fatalf("did not expect init output when ancestor exists: %q", stdout)
	}

	content := readFile(t, filepath.Join(root, ".taskgraph", "tasks.md"))
	if content != "- [ ] from nested\n" {
		t.Fatalf("unexpected root tasks.md content: %q", content)
	}
}

func TestCreateIsAliasForAdd(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"create", "alias task"})
	if err != nil {
		t.Fatalf("create returned err: %v stderr=%q", err, stderr)
	}

	content := readFile(t, filepath.Join(dir, ".taskgraph", "tasks.md"))
	if content != "- [ ] alias task\n" {
		t.Fatalf("unexpected tasks.md content: %q", content)
	}
}

func TestListPrintsRawChecklistLines(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "config.yml"), "")
	mustWrite(t, filepath.Join(dir, ".taskgraph", "tasks.md"), "- [ ] a\n- [x] done\n")

	stdout, stderr, err := run([]string{"list"})
	if err != nil {
		t.Fatalf("list returned err: %v stderr=%q", err, stderr)
	}
	if stdout != "- [ ] a\n- [x] done\n" {
		t.Fatalf("unexpected list output: %q", stdout)
	}
}

func TestAddRequiresTaskText(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"add"})
	if err == nil {
		t.Fatalf("expected error for missing task text")
	}
	if !strings.Contains(stderr, "usage: tg add <task text>") {
		t.Fatalf("expected usage message, got %q", stderr)
	}
}

func run(args []string) (string, string, error) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	err := Run(args, &out, &errOut)
	return out.String(), errOut.String(), err
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(old)
	})
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %q to exist: %v", path, err)
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
