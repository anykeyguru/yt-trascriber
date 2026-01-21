package sanitize

import (
	"os"
	"strings"
)

var repl = []string{":", "-", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t", "@"}

const maxLenTitle = 120

func SanitizeString(input string) string {

	input = clean(input)

	out := make([]rune, 0, len(input))

	for _, r := range input {
		if isAllowedRune(r) {
			out = append(out, r)
		}
	}

	return string(out)
}

func clean(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, string(os.PathSeparator), "-")
	// minimal sanitation
	for _, r := range repl {
		input = strings.ReplaceAll(input, r, "-")
	}

	repl = []string{".", "."}
	for _, r := range repl {
		input = strings.ReplaceAll(input, r, "")
	}

	// reduce spaces
	input = strings.Join(strings.Fields(input), " ")
	if len(input) > maxLenTitle {
		input = input[:maxLenTitle]
	}
	return input
}

func isAllowedRune(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}

	if r >= '0' && r <= '9' {
		return true
	}

	if r >= 'а' && r <= 'я' {
		return true
	}
	if r >= 'А' && r <= 'Я' {
		return true
	}
	if r == 'ё' || r == 'Ё' {
		return true
	}
	if r == '_' {
		return true
	}

	return false
}
