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

func CreateFileNameForSave(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	// minimal sanitation
	repl := []string{":", "-", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "_")
	}

	repl = []string{"."}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "")
	}

	// reduce spaces
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, " ", "_")
	if len(s) > 120 {
		s = s[:120]
	}
	return s + fmt.Sprintf("_%d", time.Now().UnixNano())
}
