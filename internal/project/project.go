package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const taskgraphDirName = ".taskgraph"

var builtinIssueTypes = []string{
	"idea",
	"initiative",
	"product",
	"epic",
	"feature",
	"task",
	"subtask",
	"bug",
	"chore",
	"decision",
}

// FindTaskgraphRoot walks upward from startDir until it finds a .taskgraph directory.
func FindTaskgraphRoot(startDir string) (string, bool, error) {
	if startDir == "" {
		return "", false, errors.New("start directory is required")
	}

	repoRoot, _ := findGitRepoRoot(startDir)
	cur := startDir
	for {
		candidate := filepath.Join(cur, taskgraphDirName)
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return cur, true, nil
		}
		if err != nil && !os.IsNotExist(err) {
			return "", false, err
		}

		// Don't cross git repository boundaries.
		// If we're inside a repo and reached its root with no .taskgraph, stop searching.
		if repoRoot != "" && filepath.Clean(cur) == filepath.Clean(repoRoot) {
			return "", false, nil
		}

		parent := filepath.Dir(cur)
		if parent == cur {
			return "", false, nil
		}
		cur = parent
	}
}

// InitAt creates .taskgraph/config.yml and .taskgraph/issues.md inside dir.
func InitAt(dir string) (string, bool, error) {
	if dir == "" {
		return "", false, errors.New("directory is required")
	}

	base := filepath.Join(dir, taskgraphDirName)
	already := exists(base)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", false, err
	}
	if err := ensureConfig(filepath.Join(base, "config.yml"), dir); err != nil {
		return "", false, err
	}
	if err := ensureFile(filepath.Join(base, "issues.md")); err != nil {
		return "", false, err
	}
	if err := ensureTaskgraphGitignore(filepath.Join(base, ".gitignore")); err != nil {
		return "", false, err
	}

	return dir, !already, nil
}

func ensureFile(path string) error {
	if exists(path) {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

func ensureTaskgraphGitignore(path string) error {
	const entry = "taskgraph.db"

	if !exists(path) {
		return os.WriteFile(path, []byte(entry+"\n"), 0o644)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(b)
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}

	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += entry + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureConfig(path string, rootDir string) error {
	prefix := deriveDefaultPrefix(rootDir)
	if !exists(path) {
		content := fmt.Sprintf("issue-prefix: %s\n", prefix)
		return os.WriteFile(path, []byte(content), 0o644)
	}

	// If config exists but has no usable prefix, set a default.
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if p := parsePrefix(string(b)); p != "" {
		return nil
	}
	content := fmt.Sprintf("issue-prefix: %s\n", prefix)
	return os.WriteFile(path, []byte(content), 0o644)
}

// ReadPrefix reads prefix from .taskgraph/config.yml in the given root directory.
// Falls back to default prefix derivation if config is missing/invalid.
func ReadPrefix(rootDir string) (string, error) {
	if strings.TrimSpace(rootDir) == "" {
		return "", errors.New("root directory is required")
	}
	configPath := filepath.Join(rootDir, taskgraphDirName, "config.yml")
	b, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return deriveDefaultPrefix(rootDir), nil
		}
		return "", err
	}

	prefix := parsePrefix(string(b))
	if prefix == "" {
		return deriveDefaultPrefix(rootDir), nil
	}
	return normalizePrefix(prefix), nil
}

func ReadAllowedIssueTypes(rootDir string) ([]string, error) {
	if strings.TrimSpace(rootDir) == "" {
		return nil, errors.New("root directory is required")
	}

	out := append([]string{}, builtinIssueTypes...)
	configPath := filepath.Join(rootDir, taskgraphDirName, "config.yml")
	b, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}

	custom := parseIssueTypes(string(b))
	seen := map[string]bool{}
	merged := make([]string, 0, len(out)+len(custom))
	for _, item := range append(out, custom...) {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		merged = append(merged, item)
	}
	return merged, nil
}

func parsePrefix(config string) string {
	for _, line := range strings.Split(config, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "issue-prefix:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "issue-prefix:"))
		}
		// Backward compatibility for older tg configs.
		if strings.HasPrefix(line, "prefix:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "prefix:"))
		}
	}
	return ""
}

func parseIssueTypes(config string) []string {
	for _, line := range strings.Split(config, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "issue-types:") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(line, "issue-types:"))
		if raw == "" {
			return nil
		}
		parts := strings.Split(raw, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			t := normalizeIssueType(p)
			if t != "" {
				out = append(out, t)
			}
		}
		return out
	}
	return nil
}

func deriveDefaultPrefix(dir string) string {
	base := filepath.Base(dir)
	return normalizePrefix(base)
}

func normalizePrefix(raw string) string {
	var out []rune
	for _, r := range strings.ToLower(raw) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
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

func normalizeIssueType(raw string) string {
	label := strings.TrimSpace(strings.ToLower(raw))
	if label == "" {
		return ""
	}

	var out []rune
	prevHyphen := false
	for _, r := range label {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			out = append(out, r)
			prevHyphen = false
		case r == '-' && len(out) > 0 && !prevHyphen:
			out = append(out, r)
			prevHyphen = true
		}
	}
	return strings.Trim(string(out), "-")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func findGitRepoRoot(startDir string) (string, bool) {
	cur := startDir
	for {
		gitPath := filepath.Join(cur, ".git")
		if info, err := os.Stat(gitPath); err == nil {
			// Accept both .git directory and .git file (worktree/submodule layouts).
			if info.IsDir() || !info.IsDir() {
				return cur, true
			}
		}

		parent := filepath.Dir(cur)
		if parent == cur {
			return "", false
		}
		cur = parent
	}
}
