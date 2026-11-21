package services

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"

	"uszatakuchnia/dtos"
)

// Step text embeds refs as {{<kind>:<id>}} or {{<kind>:<id>|<arg>}}. Today the only kind is "ing" (an ingredient
// row, by its key; arg = the part of the row). A glossary ref ({{term:<id>}}) can come later without touching
// stored data, so the grammar is generic and unknown kinds are rejected instead of being shown as text.
var stepRefRe = regexp.MustCompile(`\{\{([^{}:|]*):([^{}|]*)(?:\|([^{}]*))?\}\}`)

var rowKeyRe = regexp.MustCompile(`^[a-z0-9-]{1,32}$`)

// RefParts are the allowed parts of an ingredient row; "rest" = the whole minus this row's fractions in other steps.
var RefParts = map[string]bool{"1/2": true, "1/3": true, "1/4": true, "2/3": true, "3/4": true, "rest": true}

type StepRef struct {
	Kind, ID, Arg string
}

func ParseStepRefs(text string) []StepRef {
	out := []StepRef{}
	for _, m := range stepRefRe.FindAllStringSubmatch(text, -1) {
		out = append(out, StepRef{Kind: m[1], ID: m[2], Arg: m[3]})
	}
	return out
}

func validateIngredientKeys(req dtos.UpsertRecipeRequest, add func(field, msg string)) {
	seen := map[string]bool{}
	for i, it := range req.Ingredients {
		k := strings.TrimSpace(it.Key)
		if k == "" {
			continue
		}
		f := fmt.Sprintf("ingredients[%d].key", i)
		if !rowKeyRe.MatchString(k) {
			add(f, "Klucz wiersza: małe litery, cyfry i „-”, najwyżej 32 znaki")
		} else if seen[k] {
			add(f, fmt.Sprintf("Klucz wiersza „%s” powtarza się w przepisie", k))
		}
		seen[k] = true
	}
}

// validateStepRefs: a malformed ref is always an error (a typo in an API body); a ref to a row that isn't in the
// recipe only blocks publishing, so a draft with a removed row can still be saved and fixed later.
func validateStepRefs(req dtos.UpsertRecipeRequest, publish bool, add func(field, msg string)) {
	keys := map[string]bool{}
	for _, it := range req.Ingredients {
		if k := strings.TrimSpace(it.Key); k != "" {
			keys[k] = true
		}
	}
	for i, s := range req.Steps {
		f := fmt.Sprintf("steps[%d].text", i)
		for _, r := range ParseStepRefs(s.Text) {
			switch {
			case r.Kind != "ing":
				add(f, fmt.Sprintf("Krok %d: nieznany znacznik „{{%s:…}}”", i+1, r.Kind))
			case r.ID == "":
				add(f, fmt.Sprintf("Krok %d: odwołanie bez składnika", i+1))
			case !rowKeyRe.MatchString(r.ID):
				// no row can have such a key, so it's a typo, not a removed row
				add(f, fmt.Sprintf("Krok %d: nieprawidłowy klucz wiersza „%s” (małe litery, cyfry i „-”, najwyżej 32 znaki)", i+1, r.ID))
			case r.Arg != "" && !RefParts[r.Arg]:
				add(f, fmt.Sprintf("Krok %d: nieznana część „%s” (1/2, 1/3, 1/4, 2/3, 3/4, rest)", i+1, r.Arg))
			case publish && !keys[r.ID]:
				add(f, fmt.Sprintf("Krok %d: odwołanie do składnika, którego nie ma w przepisie — wybierz inny albo usuń odwołanie", i+1))
			}
		}
	}
}

// ingredientKeys returns the key for each row: the given one, or a new random one for rows without it.
func ingredientKeys(items []dtos.UpsertIngredientRequest) []string {
	used := map[string]bool{}
	for _, it := range items {
		used[strings.TrimSpace(it.Key)] = true
	}
	out := make([]string, len(items))
	for i, it := range items {
		k := strings.TrimSpace(it.Key)
		for k == "" || (it.Key == "" && used[k]) {
			k = randomKey()
		}
		used[k] = true
		out[i] = k
	}
	return out
}

const keyAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomKey() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = keyAlphabet[int(b[i])%len(keyAlphabet)]
	}
	return string(b)
}
