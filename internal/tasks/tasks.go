package tasks

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var idPattern = regexp.MustCompile(`\[[a-z0-9]+-[0-9a-z]{3,8}\]`)
var labelPattern = regexp.MustCompile(`(^|[\s(])#([A-Za-z0-9][A-Za-z0-9-]*)`)
var typeLabelPrefix = "t-"

// AppendTask appends one markdown checklist line to tasksFile.
func AppendTask(tasksFile, prefix, text string, labels []string, taskType string) error {
	if strings.TrimSpace(tasksFile) == "" {
		return errors.New("tasks file is required")
	}
	cleanPrefix := normalizePrefix(prefix)
	clean := strings.TrimSpace(text)
	if clean == "" {
		return errors.New("task text is required")
	}
	resolvedType, cleanInputLabels, err := ResolveTaskType(clean, labels, taskType)
	if err != nil {
		return err
	}

	existing, err := os.ReadFile(tasksFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	now := time.Now()
	existingIDs := collectExistingIDs(string(existing))
	id := generateIssueID(cleanPrefix, clean, now, existingIDs)
	mergedLabels := MergeLabels(ExtractLabels(clean), cleanInputLabels)
	if resolvedType != "" {
		mergedLabels = MergeLabels(mergedLabels, []string{TypeLabel(resolvedType)})
	}
	if len(mergedLabels) > 0 {
		clean = stripLabels(clean)
		clean = strings.TrimSpace(clean + " " + formatLabels(mergedLabels))
	}

	line := fmt.Sprintf("- [ ] ➕%s [%s] %s\n", now.Format("2006-01-02"), id, clean)
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

func CloseTask(tasksFile, id, reason string) error {
	if strings.TrimSpace(tasksFile) == "" {
		return errors.New("tasks file is required")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("task id is required")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return errors.New("close reason is required")
	}

	content, err := os.ReadFile(tasksFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	needle := "[" + id + "]"
	found := false
	for i, line := range lines {
		if !strings.Contains(line, needle) {
			continue
		}
		found = true
		if !strings.HasPrefix(line, "- [ ] ") {
			return fmt.Errorf("task already closed: %s", id)
		}
		lines[i] = line + " **✅" + time.Now().Format("2006-01-02") + " " + reason + "**"
		lines[i] = strings.Replace(lines[i], "- [ ] ", "- [x] ", 1)
		break
	}
	if !found {
		return fmt.Errorf("task not found: %s", id)
	}

	return os.WriteFile(tasksFile, []byte(strings.Join(lines, "\n")), 0o644)
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

func collectExistingIDs(content string) map[string]bool {
	out := map[string]bool{}
	for _, match := range idPattern.FindAllString(content, -1) {
		id := strings.TrimSuffix(strings.TrimPrefix(match, "["), "]")
		out[id] = true
	}
	return out
}

func generateIssueID(prefix, text string, created time.Time, existing map[string]bool) string {
	for length := 3; length <= 8; length++ {
		for nonce := 0; nonce < 10; nonce++ {
			candidate := generateHashID(prefix, text, created, length, nonce)
			if !existing[candidate] {
				return candidate
			}
		}
	}
	// Extremely unlikely fallback.
	return generateHashID(prefix, text, time.Now().Add(time.Nanosecond), 8, 999)
}

func generateHashID(prefix, text string, created time.Time, length, nonce int) string {
	content := fmt.Sprintf("%s|%d|%d", text, created.UnixNano(), nonce)
	hash := sha256.Sum256([]byte(content))
	short := encodeBase36(hash[:], length)
	return fmt.Sprintf("%s-%s", prefix, short)
}

func encodeBase36(data []byte, length int) string {
	num := new(big.Int).SetBytes(data)
	base := big.NewInt(36)
	zero := big.NewInt(0)
	mod := new(big.Int)

	chars := make([]byte, 0, length)
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		chars = append(chars, alphabet[mod.Int64()])
	}
	if len(chars) == 0 {
		chars = append(chars, '0')
	}

	var b strings.Builder
	for i := len(chars) - 1; i >= 0; i-- {
		b.WriteByte(chars[i])
	}
	s := b.String()
	if len(s) < length {
		s = strings.Repeat("0", length-len(s)) + s
	}
	if len(s) > length {
		s = s[len(s)-length:]
	}
	return s
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

func NormalizeLabelsCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		label := normalizeLabel(part)
		if label != "" {
			out = append(out, label)
		}
	}
	return MergeLabels(out)
}

func ExtractLabels(text string) []string {
	matches := labelPattern.FindAllStringSubmatch(text, -1)
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		label := normalizeLabel(match[2])
		if label != "" {
			out = append(out, label)
		}
	}
	return MergeLabels(out)
}

func MergeLabels(groups ...[]string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, group := range groups {
		for _, item := range group {
			label := normalizeLabel(item)
			if label == "" || seen[label] {
				continue
			}
			seen[label] = true
			out = append(out, label)
		}
	}
	return out
}

func normalizeLabel(raw string) string {
	label := strings.TrimSpace(strings.TrimPrefix(raw, "#"))
	label = strings.ToLower(label)
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
	label = strings.Trim(string(out), "-")
	return label
}

func NormalizeTaskType(raw string) string {
	normalized := normalizeLabel(raw)
	if strings.HasPrefix(normalized, typeLabelPrefix) {
		return strings.TrimPrefix(normalized, typeLabelPrefix)
	}
	return normalized
}

func TypeLabel(taskType string) string {
	t := NormalizeTaskType(taskType)
	if t == "" {
		return ""
	}
	return typeLabelPrefix + t
}

func ExtractTaskTypeFromText(text string) (string, error) {
	return ExtractTaskTypeFromLabels(ExtractLabels(text))
}

func ExtractTaskTypeFromLabels(labels []string) (string, error) {
	found := ""
	for _, label := range labels {
		l := normalizeLabel(label)
		if !strings.HasPrefix(l, typeLabelPrefix) {
			continue
		}
		t := strings.TrimPrefix(l, typeLabelPrefix)
		if t == "" {
			continue
		}
		if found == "" || found == t {
			found = t
			continue
		}
		return "", fmt.Errorf("multiple task types found: %s, %s", found, t)
	}
	return found, nil
}

func ResolveTaskType(text string, labels []string, taskType string) (string, []string, error) {
	inlineType, err := ExtractTaskTypeFromText(text)
	if err != nil {
		return "", nil, err
	}
	labelType, err := ExtractTaskTypeFromLabels(labels)
	if err != nil {
		return "", nil, err
	}
	flagType := NormalizeTaskType(taskType)

	selected := ""
	for _, candidate := range []string{inlineType, labelType, flagType} {
		if candidate == "" {
			continue
		}
		if selected == "" || selected == candidate {
			selected = candidate
			continue
		}
		return "", nil, fmt.Errorf("conflicting task types: %s, %s", selected, candidate)
	}

	outLabels := make([]string, 0, len(labels))
	for _, label := range labels {
		l := normalizeLabel(label)
		if strings.HasPrefix(l, typeLabelPrefix) {
			continue
		}
		if l != "" {
			outLabels = append(outLabels, l)
		}
	}

	return selected, MergeLabels(outLabels), nil
}

func stripLabels(text string) string {
	cleaned := labelPattern.ReplaceAllString(text, "$1")
	return strings.Join(strings.Fields(cleaned), " ")
}

func formatLabels(labels []string) string {
	parts := make([]string, 0, len(labels))
	for _, label := range MergeLabels(labels) {
		parts = append(parts, "#"+label)
	}
	return strings.Join(parts, " ")
}
