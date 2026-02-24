package project

import (
	"errors"
	"os"
	"path/filepath"
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
	if err := ensureFile(filepath.Join(base, "config.yml")); err != nil {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
