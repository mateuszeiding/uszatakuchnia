package services

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// MaxCautionLen — najdłuższa uwaga przepisu / wyzwania (znaki, nie bajty).
const MaxCautionLen = 220

// CleanCautions przycina uwagi (białe znaki w środku → jedna spacja) i pomija puste.
// Zawsze zwraca niepustą referencję (zapis jako []).
func CleanCautions(in []string) []string {
	out := make([]string, 0, len(in))
	for _, c := range in {
		if c = strings.Join(strings.Fields(c), " "); c != "" {
			out = append(out, c)
		}
	}
	return out
}

// validateCautions — limit długości; pole cautions[i] liczy tylko niepuste uwagi (puste są pomijane przy zapisie).
func validateCautions(in []string, add func(field, msg string)) {
	for i, c := range CleanCautions(in) {
		if utf8.RuneCountInString(c) > MaxCautionLen {
			add(fmt.Sprintf("cautions[%d]", i), fmt.Sprintf("Uwaga może mieć najwyżej %d znaków", MaxCautionLen))
		}
	}
}
