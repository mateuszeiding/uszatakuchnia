package services

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// Step text embeds refs as {{<kind>:<id>}} or {{<kind>:<id>|<arg>}}. Two kinds: "ing" (an ingredient row, by its
// key; arg = the part of the row) and "term" (a glossary term, by its id; arg = the author's word it wraps).
// Unknown kinds are rejected instead of being shown as text. Hints take only "term".
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

var termIDRe = regexp.MustCompile(`^[1-9][0-9]{0,9}$`)

// MaxTermWordLen is the longest word a glossary ref may wrap (characters).
const MaxTermWordLen = 60

func termID(r StepRef) (uint, bool) {
	if !termIDRe.MatchString(r.ID) {
		return 0, false
	}
	n, err := strconv.ParseUint(r.ID, 10, 32)
	return uint(n), err == nil
}

// TermRefIDs returns the glossary term ids referenced in the texts (each once, in order of appearance).
func TermRefIDs(texts ...string) []uint {
	seen := map[uint]bool{}
	out := []uint{}
	for _, t := range texts {
		for _, r := range ParseStepRefs(t) {
			if id, ok := termID(r); r.Kind == "term" && ok && !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}
	return out
}

// recipeTermIDs: term ids in the request's steps and hints.
func recipeTermIDs(req dtos.UpsertRecipeRequest) []uint {
	texts := []string{}
	for _, s := range req.Steps {
		texts = append(texts, s.Text)
		texts = append(texts, s.Hints...)
	}
	return TermRefIDs(texts...)
}

// VisibleTermText is the text with glossary refs replaced by their word, as the reader sees it. Limits count
// this, so marking a word as a term doesn't eat into them.
func VisibleTermText(text string) string {
	return stepRefRe.ReplaceAllStringFunc(text, func(m string) string {
		sub := stepRefRe.FindStringSubmatch(m)
		if sub[1] == "term" {
			return sub[3]
		}
		return m
	})
}

// ReplaceTermRefs rewrites refs to the term `id`: to the term `with`, or (with = 0) to the bare word.
func ReplaceTermRefs(text string, id, with uint) string {
	return stepRefRe.ReplaceAllStringFunc(text, func(m string) string {
		sub := stepRefRe.FindStringSubmatch(m)
		if n, ok := termID(StepRef{ID: sub[2]}); sub[1] != "term" || !ok || n != id {
			return m
		}
		if with == 0 {
			return sub[3]
		}
		return fmt.Sprintf("{{term:%d|%s}}", with, sub[3])
	})
}

// strayBraces: "{{" or "}}" left after the valid refs — a broken or nested ref ({{term:1|mąkę {{ing:x}}}}).
func strayBraces(text string) bool {
	rest := stepRefRe.ReplaceAllString(text, "")
	return strings.Contains(rest, "{{") || strings.Contains(rest, "}}")
}

// checkTermRef: a malformed glossary ref is always an error; a term that isn't in the glossary (deleted, or an id
// from another database) only blocks publishing — the page shows the bare word anyway.
func checkTermRef(r StepRef, where, f string, publish bool, terms map[uint]bool, add func(field, msg string)) {
	id, ok := termID(r)
	word := strings.TrimSpace(r.Arg)
	switch {
	case !ok:
		add(f, fmt.Sprintf("%s: nieprawidłowe pojęcie „%s” (id ze słowniczka)", where, r.ID))
	case word == "":
		add(f, fmt.Sprintf("%s: pojęcie bez oznaczonego słowa", where))
	case strings.Contains(r.Arg, "|"):
		add(f, fmt.Sprintf("%s: oznaczone słowo nie może zawierać „|”", where))
	case utf8.RuneCountInString(r.Arg) > MaxTermWordLen:
		add(f, fmt.Sprintf("%s: oznaczone słowo może mieć najwyżej %d znaków", where, MaxTermWordLen))
	case publish && !terms[id]:
		add(f, fmt.Sprintf("%s: pojęcia nie ma w słowniczku — wybierz inne albo usuń oznaczenie", where))
	}
}

// knownTerms: which of the request's term ids exist in the glossary.
func knownTerms(tx *gorm.DB, req dtos.UpsertRecipeRequest) (map[uint]bool, error) {
	out := map[uint]bool{}
	ids := recipeTermIDs(req)
	if len(ids) == 0 {
		return out, nil
	}
	var found []uint
	if err := tx.Model(&entities.GlossaryTerm{}).Where("id IN ?", ids).Pluck("id", &found).Error; err != nil {
		return nil, err
	}
	for _, id := range found {
		out[id] = true
	}
	return out, nil
}

// validateStepRefs: a malformed ref is always an error (a typo in an API body); a ref to a row that isn't in the
// recipe, or to a term that isn't in the glossary, only blocks publishing, so a draft can still be saved and fixed.
func validateStepRefs(req dtos.UpsertRecipeRequest, publish bool, terms map[uint]bool, add func(field, msg string)) {
	keys := map[string]bool{}
	for _, it := range req.Ingredients {
		if k := strings.TrimSpace(it.Key); k != "" {
			keys[k] = true
		}
	}
	for i, s := range req.Steps {
		for hi, h := range s.Hints {
			f := fmt.Sprintf("steps[%d].hints[%d]", i, hi)
			where := fmt.Sprintf("Krok %d, wskazówka %d", i+1, hi+1)
			if strayBraces(h) {
				add(f, where+": nieprawidłowy znacznik „{{…}}”")
			}
			for _, r := range ParseStepRefs(h) {
				if r.Kind == "term" {
					checkTermRef(r, where, f, publish, terms, add)
				} else {
					add(f, fmt.Sprintf("%s: we wskazówce można oznaczać tylko pojęcia, nie „{{%s:…}}”", where, r.Kind))
				}
			}
		}
		f := fmt.Sprintf("steps[%d].text", i)
		if strayBraces(s.Text) {
			add(f, fmt.Sprintf("Krok %d: nieprawidłowy znacznik „{{…}}” (odwołania nie mogą się zagnieżdżać)", i+1))
		}
		for _, r := range ParseStepRefs(s.Text) {
			switch {
			case r.Kind == "term":
				checkTermRef(r, fmt.Sprintf("Krok %d", i+1), f, publish, terms, add)
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
