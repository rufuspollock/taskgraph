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

// FindTaskgraphRoot walks upward from startDir until it finds a .taskgraph directory.
func FindTaskgraphRoot(startDir string) (string, bool, error) {
	if startDir == "" {
		return "", false, errors.New("start directory is required")
	}

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

		parent := filepath.Dir(cur)
		if parent == cur {
			return "", false, nil
		}
		cur = parent
	}
}

// InitAt creates .taskgraph/config.yml and .taskgraph/tasks.md inside dir.
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
	if err := ensureFile(filepath.Join(base, "tasks.md")); err != nil {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
