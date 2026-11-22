package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"uszatakuchnia/db/entities"
)

// Portion in link previews: mirrors web/src/shared/utils/portion.ts (URL params, names, fallbacks).

// pieceForms: the piece unit's forms (1 / 2–4 / 5+); without a unit, generic pieces.
func (rec *shareRecipe) pieceForms() [3]string {
	if len(rec.PieceForms) == 3 {
		return [3]string{rec.PieceForms[0], rec.PieceForms[1], rec.PieceForms[2]}
	}
	return [3]string{"sztuka", "sztuki", "sztuk"}
}

type shareForm struct {
	round   bool
	d, a, b int
}

func (f shareForm) name() string {
	if f.round {
		return fmt.Sprintf("tortownica Ø %d cm", f.d)
	}
	return fmt.Sprintf("blacha %d × %d cm", f.a, f.b)
}

func (rec *shareRecipe) recipeForm() (shareForm, bool) {
	if rec.FormShape == nil {
		return shareForm{}, false
	}
	switch *rec.FormShape {
	case entities.FormRound:
		if rec.FormDCm != nil {
			return shareForm{round: true, d: *rec.FormDCm}, true
		}
	case entities.FormRect:
		if rec.FormACm != nil && rec.FormBCm != nil {
			return shareForm{a: *rec.FormACm, b: *rec.FormBCm}, true
		}
	}
	return shareForm{}, false
}

var formParam = regexp.MustCompile(`^(?:o(\d{2})|p(\d{2})x(\d{2}))$`)

func cmOK(n int) bool { return n >= entities.PortionMinCm && n <= entities.PortionMaxCm }

// parseFormParam: ?forma=o26 / p30x40 (10–60 cm); false = missing or invalid.
func parseFormParam(s string) (shareForm, bool) {
	m := formParam.FindStringSubmatch(s)
	if m == nil {
		return shareForm{}, false
	}
	if m[1] != "" {
		d, _ := strconv.Atoi(m[1])
		return shareForm{round: true, d: d}, cmOK(d)
	}
	a, _ := strconv.Atoi(m[2])
	b, _ := strconv.Atoi(m[3])
	return shareForm{a: a, b: b}, cmOK(a) && cmOK(b)
}

func intParam(r *http.Request, name string, min, max int) (int, bool) {
	n, err := strconv.Atoi(r.URL.Query().Get(name))
	return n, err == nil && n >= min && n <= max
}

// shoppingPortion: "6 porcji", "3 placki 25 cm", "tortownica Ø 26 cm" for the shopping list preview title;
// "" = nothing to say (a form recipe without a form).
func shoppingPortion(rec *shareRecipe, r *http.Request) string {
	switch rec.PortionMode {
	case entities.PortionSized:
		n, ok := intParam(r, "sztuki", 1, 20)
		if !ok {
			n = rec.Servings
		}
		f := rec.pieceForms()
		noun := plural(n, f[0], f[1], f[2])
		d, ok := intParam(r, "rozmiar", entities.PortionMinCm, entities.PortionMaxCm)
		if !ok || !containsInt(rec.Sizes, d) {
			d = rec.BaseD // a size the author no longer offers → the recipe's size
		}
		if d == 0 {
			return fmt.Sprintf("%d %s", n, noun)
		}
		return fmt.Sprintf("%d %s %d cm", n, noun, d)
	case entities.PortionForm:
		if f, ok := parseFormParam(r.URL.Query().Get("forma")); ok {
			return f.name()
		}
		if f, ok := rec.recipeForm(); ok {
			return f.name()
		}
		return ""
	default:
		n := servingsParam(r, rec.Servings)
		return fmt.Sprintf("%d %s", n, servingsNoun(n))
	}
}

// recipePortion: "Przepis na 6 porcji" / "Przepis: 2 placki 30 cm" / "Przepis: tortownica Ø 24 cm".
// Pieces and forms use the nominative ("Przepis:"): the author's unit has no accusative form.
func recipePortion(rec *shareRecipe) string {
	switch rec.PortionMode {
	case entities.PortionSized:
		if rec.Servings <= 0 {
			return "Przepis"
		}
		f := rec.pieceForms()
		s := fmt.Sprintf("Przepis: %d %s", rec.Servings, plural(rec.Servings, f[0], f[1], f[2]))
		if rec.BaseD > 0 {
			s += fmt.Sprintf(" %d cm", rec.BaseD)
		}
		return s
	case entities.PortionForm:
		if f, ok := rec.recipeForm(); ok {
			return "Przepis: " + f.name()
		}
		return "Przepis"
	default:
		if rec.Servings <= 0 {
			return "Przepis"
		}
		return fmt.Sprintf("Przepis na %d %s", rec.Servings, servingsAcc(rec.Servings))
	}
}

func containsInt(a []int, x int) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
