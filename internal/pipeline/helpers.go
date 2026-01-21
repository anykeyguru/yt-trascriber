package pipeline

import (
	"strings"
)

func normalizeURL(s string) string {
	s = strings.TrimSpace(s)

	// убираем обрамляющие кавычки/скобки
	s = strings.Trim(s, "[]()<>\"'")

	return s
}
