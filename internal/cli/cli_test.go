package cli

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
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
		if !strings.Contains(stdout, "TaskGraph") {
			t.Fatalf("expected help header for %q, got %q", arg, stdout)
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
	assertExists(t, filepath.Join(dir, ".taskgraph", "issues.md"))
	assertExists(t, filepath.Join(dir, ".taskgraph", ".gitignore"))
	assertExists(t, filepath.Join(dir, ".taskgraph", "taskgraph.db"))
	if got := readFile(t, filepath.Join(dir, ".taskgraph", ".gitignore")); !strings.Contains(got, "taskgraph.db\n") {
		t.Fatalf("expected .gitignore to include taskgraph.db, got %q", got)
	}
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

	content := readFile(t, filepath.Join(dir, ".taskgraph", "issues.md"))
	prefix := expectedPrefixForDir(dir)
	if !matchesTaskLine(content, prefix, "first task") {
		t.Fatalf("unexpected issues.md content: %q", content)
	}
	assertExists(t, filepath.Join(dir, ".taskgraph", "taskgraph.db"))
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
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "config.yml"), []byte("prefix: root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".taskgraph", "issues.md"), []byte(""), 0o644); err != nil {
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

	content := readFile(t, filepath.Join(root, ".taskgraph", "issues.md"))
	if !matchesTaskLine(content, "root", "from nested") {
		t.Fatalf("unexpected root issues.md content: %q", content)
	}
}

func TestCreateIsAliasForAdd(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"create", "alias task"})
	if err != nil {
		t.Fatalf("create returned err: %v stderr=%q", err, stderr)
	}

	content := readFile(t, filepath.Join(dir, ".taskgraph", "issues.md"))
	prefix := expectedPrefixForDir(dir)
	if !matchesTaskLine(content, prefix, "alias task") {
		t.Fatalf("unexpected issues.md content: %q", content)
	}
}

func TestInboxPrintsRawChecklistLines(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "config.yml"), "")
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "- [ ] a\n- [x] done\n")

	stdout, stderr, err := run([]string{"inbox"})
	if err != nil {
		t.Fatalf("inbox returned err: %v stderr=%q", err, stderr)
	}
	if stdout != "- [ ] a\n- [x] done\n" {
		t.Fatalf("unexpected inbox output: %q", stdout)
	}
}

func TestListReadsChecklistFromDatabase(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"init"})
	if err != nil {
		t.Fatalf("init returned err: %v stderr=%q", err, stderr)
	}

	mustWrite(t, filepath.Join(dir, "alpha.md"), "# Alpha\n\n- [ ] Alpha task\n")
	mustWrite(t, filepath.Join(dir, "beta.md"), "# Beta\n\n- [x] Beta done\n")
	now := time.Now()
	mustChtimes(t, filepath.Join(dir, "alpha.md"), now.Add(-2*time.Hour))
	mustChtimes(t, filepath.Join(dir, "beta.md"), now.Add(-1*time.Hour))

	_, stderr, err = run([]string{"index"})
	if err != nil {
		t.Fatalf("index returned err: %v stderr=%q", err, stderr)
	}

	stdout, stderr, err := run([]string{"list"})
	if err != nil {
		t.Fatalf("list returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "[ ] Alpha task") {
		t.Fatalf("expected open task in list output, got %q", stdout)
	}
	if strings.Contains(stdout, "Beta done") {
		t.Fatalf("did not expect closed task in default list output, got %q", stdout)
	}

	stdoutAll, stderr, err := run([]string{"list", "--all"})
	if err != nil {
		t.Fatalf("list --all returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdoutAll, "[x] Beta done") {
		t.Fatalf("expected closed task in list --all output, got %q", stdoutAll)
	}

	if strings.Index(stdoutAll, "Beta done") > strings.Index(stdoutAll, "Alpha task") {
		t.Fatalf("expected newest file tasks first, got %q", stdoutAll)
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

func TestIndexBuildsTaskgraphDatabase(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"init"})
	if err != nil {
		t.Fatalf("init returned err: %v stderr=%q", err, stderr)
	}

	mustWrite(t, filepath.Join(dir, "project.md"), "# Project\n\n- [ ] Ship v1\n")

	stdout, stderr, err := run([]string{"index"})
	if err != nil {
		t.Fatalf("index returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "Indexed") {
		t.Fatalf("expected index summary, got %q", stdout)
	}
	assertExists(t, filepath.Join(dir, ".taskgraph", "taskgraph.db"))
}

func TestAddUsesTGCWDOverride(t *testing.T) {
	targetDir := t.TempDir()
	otherDir := t.TempDir()

	chdir(t, otherDir)
	t.Setenv("TG_CWD", targetDir)

	_, stderr, err := run([]string{"add", "from override"})
	if err != nil {
		t.Fatalf("add returned err: %v stderr=%q", err, stderr)
	}

	content := readFile(t, filepath.Join(targetDir, ".taskgraph", "issues.md"))
	if !strings.Contains(content, "from override") {
		t.Fatalf("expected task in TG_CWD directory, got %q", content)
	}
}

func TestAddUpdatesSQLiteIndex(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"add", "indexed task"})
	if err != nil {
		t.Fatalf("add returned err: %v stderr=%q", err, stderr)
	}

	dbPath := filepath.Join(dir, ".taskgraph", "taskgraph.db")
	assertExists(t, dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow(
		"SELECT COUNT(*) FROM index_nodes WHERE kind = 'checklist' AND path = '.taskgraph/issues.md'",
	).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected at least 1 checklist node for .taskgraph/issues.md, got %d", count)
	}
}

func TestMigrateBeadsImportsIntoIssuesMarkdown(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustMkdirAll(t, filepath.Join(dir, ".beads"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "- [ ] existing\n")
	mustWrite(t, filepath.Join(dir, ".beads", "issues.jsonl"), strings.Join([]string{
		`{"id":"pl-1","title":"Open item","status":"open"}`,
		`{"id":"pl-2","title":"Closed item","status":"closed"}`,
		`{"id":"pl-3","title":"Deleted item","status":"tombstone"}`,
	}, "\n")+"\n")

	stdout, stderr, err := run([]string{"migrate-beads"})
	if err != nil {
		t.Fatalf("migrate-beads returned err: %v stderr=%q", err, stderr)
	}
	if !strings.Contains(stdout, "Imported 2 issues") {
		t.Fatalf("expected import summary, got %q", stdout)
	}

	content := readFile(t, filepath.Join(dir, ".taskgraph", "issues.md"))
	if !strings.Contains(content, "- [ ] [beads:pl-1] Open item\n") {
		t.Fatalf("missing open import line: %q", content)
	}
	if !strings.Contains(content, "- [x] [beads:pl-2] Closed item\n") {
		t.Fatalf("missing closed import line: %q", content)
	}
	if strings.Contains(content, "pl-3") {
		t.Fatalf("unexpected tombstone import: %q", content)
	}
}

func TestMigrateBeadsRequiresLocalBeadsAndTaskgraphDirs(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	_, stderr, err := run([]string{"migrate-beads"})
	if err == nil {
		t.Fatalf("expected error for missing directories")
	}
	if !strings.Contains(stderr, "expected .beads and .taskgraph in current directory") {
		t.Fatalf("expected clear missing-directory message, got %q", stderr)
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

func mustChtimes(t *testing.T, path string, mod time.Time) {
	t.Helper()
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatalf("chtimes failed: %v", err)
	}
}

func expectedTaskLine(text string) string {
	return "- [ ] ➕" + time.Now().Format("2006-01-02") + " [tg-abc] " + text + "\n"
}

func expectedPrefixForDir(dir string) string {
	base := strings.ToLower(filepath.Base(dir))
	var out []rune
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out = append(out, r)
		}
	}
	if len(out) == 0 {
		return "tg"
	}
	if len(out) > 4 {
		out = out[:4]
	}
	return string(out)
}

func matchesTaskLine(line, prefix, text string) bool {
	pattern := `^\- \[ \] ➕` + regexp.QuoteMeta(time.Now().Format("2006-01-02")) + ` \[` + regexp.QuoteMeta(prefix) + `\-[0-9a-z]{3,8}\] ` + regexp.QuoteMeta(text) + `\n$`
	ok, _ := regexp.MatchString(pattern, line)
	return ok
}
