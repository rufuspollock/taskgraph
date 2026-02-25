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

// AppendTask appends one markdown checklist line to tasksFile.
func AppendTask(tasksFile, prefix, text string) error {
	if strings.TrimSpace(tasksFile) == "" {
		return errors.New("tasks file is required")
	}
	cleanPrefix := normalizePrefix(prefix)
	clean := strings.TrimSpace(text)
	if clean == "" {
		return errors.New("task text is required")
	}

	existing, err := os.ReadFile(tasksFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	now := time.Now()
	existingIDs := collectExistingIDs(string(existing))
	id := generateIssueID(cleanPrefix, clean, now, existingIDs)

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
