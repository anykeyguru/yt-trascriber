package pipeline

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func normalizeURL(s string) string {
	s = strings.TrimSpace(s)

	// убираем обрамляющие кавычки/скобки
	s = strings.Trim(s, "[]()<>\"'")

	return s
}

const maxLenTitle = 120

var repl = []string{":", "-", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t", "@"}

func CreateFileNameForSave(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	// minimal sanitation
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "_")
	}

	repl = []string{".", "."}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "")
	}

	// reduce spaces
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, " ", "_")
	if len(s) > maxLenTitle {
		s = s[:maxLenTitle]
	}
	return fmt.Sprintf("%s_%d", s, time.Now().UnixNano())
}

func Sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, string(os.PathSeparator), "-")
	// minimal sanitation
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "-")
	}

	repl = []string{".", "."}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "")
	}

	// reduce spaces
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > maxLenTitle {
		s = s[:maxLenTitle]
	}
	return fmt.Sprintf("%s_%d", s, time.Now().UnixNano())
}
