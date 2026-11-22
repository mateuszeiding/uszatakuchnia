package handlers

import (
	"strings"
	"testing"

	"uszatakuchnia/db/entities"
)

var pizza = &shareRecipe{
	ID: 8, Name: "Pizza", Servings: 2, TimeMin: 60, Status: entities.StatusPublished,
	PortionMode: entities.PortionSized, Sizes: []int{25, 30, 35, 40}, BaseD: 30,
	PieceForms: []string{"placek", "placki", "placków"},
}

var szarlotka = &shareRecipe{
	ID: 9, Name: "Szarlotka", TimeMin: 90, Status: entities.StatusPublished,
	PortionMode: entities.PortionForm, FormShape: sp(entities.FormRound), FormDCm: ip(24),
}

func TestSharePortionMeta(t *testing.T) {
	stub(t, map[uint]*shareRecipe{8: pizza, 9: szarlotka}, tplIndex)
	cases := []struct {
		path, id string
		shopping bool
		want     string
	}{
		{"/recipes/8", "8", false, `content="Przepis: 2 placki 30 cm · 1 h · Uszata Kuchnia"`},
		{"/recipes/9", "9", false, `content="Przepis: tortownica Ø 24 cm · 1 h 30 min · Uszata Kuchnia"`},
		{"/recipes/8/zakupy?sztuki=3&rozmiar=25", "8", true, `content="Lista zakupów: Pizza · 3 placki 25 cm"`},
		{"/recipes/8/zakupy?sztuki=5&rozmiar=45", "8", true, `content="Lista zakupów: Pizza · 5 placków 30 cm"`}, // size gone → recipe's
		{"/recipes/8/zakupy", "8", true, `content="Lista zakupów: Pizza · 2 placki 30 cm"`},
		{"/recipes/8/zakupy?porcje=6", "8", true, `content="Lista zakupów: Pizza · 2 placki 30 cm"`}, // porcje only in servings mode
		{"/recipes/9/zakupy?forma=o26", "9", true, `content="Lista zakupów: Szarlotka · tortownica Ø 26 cm"`},
		{"/recipes/9/zakupy?forma=p30x40", "9", true, `content="Lista zakupów: Szarlotka · blacha 30 × 40 cm"`},
		{"/recipes/9/zakupy?forma=o70", "9", true, `content="Lista zakupów: Szarlotka · tortownica Ø 24 cm"`},
		{"/recipes/9/zakupy?forma=x", "9", true, `content="Lista zakupów: Szarlotka · tortownica Ø 24 cm"`},
	}
	for _, c := range cases {
		h := RecipeSharePage
		if c.shopping {
			h = ShoppingSharePage
		}
		body := get(c.path, h, c.id)
		if !strings.Contains(body, c.want) {
			t.Errorf("%s: missing %s\n%s", c.path, c.want, body)
		}
	}
}

func TestShareServingsUnchanged(t *testing.T) {
	stub(t, map[uint]*shareRecipe{6: ragu}, tplIndex)
	body := get("/recipes/6/zakupy?porcje=4", ShoppingSharePage, "6")
	if !strings.Contains(body, `content="Lista zakupów: Ragù alla bolognese · 4 porcje"`) {
		t.Errorf("servings title changed:\n%s", body)
	}
}
