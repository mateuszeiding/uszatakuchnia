package services

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// MaxCautionLen is the longest recipe / challenge note (characters, not bytes).
const MaxCautionLen = 220

// CleanCautions trims notes (inner whitespace → one space) and drops empty ones.
// Always returns a non-nil slice (stored as []).
func CleanCautions(in []string) []string {
	out := make([]string, 0, len(in))
	for _, c := range in {
		if c = strings.Join(strings.Fields(c), " "); c != "" {
			out = append(out, c)
		}
	}
	return out
}

// validateCautions checks the length limit; the cautions[i] field counts only non-empty notes (empty ones are dropped on save).
func validateCautions(in []string, add func(field, msg string)) {
	for i, c := range CleanCautions(in) {
		if utf8.RuneCountInString(c) > MaxCautionLen {
			add(fmt.Sprintf("cautions[%d]", i), fmt.Sprintf("Uwaga może mieć najwyżej %d znaków", MaxCautionLen))
		}
	}
}
