package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// AppendTask appends one markdown checklist line to tasksFile.
func AppendTask(tasksFile, text string) error {
	if strings.TrimSpace(tasksFile) == "" {
		return errors.New("tasks file is required")
	}
	clean := strings.TrimSpace(text)
	if clean == "" {
		return errors.New("task text is required")
	}

	existing, err := os.ReadFile(tasksFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	line := fmt.Sprintf("- [ ] ➕%s %s\n", time.Now().Format("2006-01-02"), clean)
	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		line = "\n" + line
	}

	f, err := os.OpenFile(tasksFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line)
	return err
}

// ReadChecklistLines reads markdown checklist lines in file order.
func ReadChecklistLines(tasksFile string) ([]string, error) {
	if strings.TrimSpace(tasksFile) == "" {
		return nil, errors.New("tasks file is required")
	}

	f, err := os.Open(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	out := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "- [ ] ") || strings.HasPrefix(line, "- [x] ") {
			out = append(out, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
