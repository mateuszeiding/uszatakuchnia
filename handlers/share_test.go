package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"uszatakuchnia/db/entities"
)

const tplIndex = `<!doctype html>
<html lang="pl">
  <head>
    <meta charset="UTF-8" />
    <!-- meta:start -->
    <title>Uszata Kuchnia</title>
    <meta name="description" content="ogólny" />
    <!-- meta:end -->
  </head>
  <body><div id="app"></div><script type="module" src="/assets/index-abc.js"></script></body>
</html>`

func sp(s string) *string { return &s }
func ip(n int) *int       { return &n }

// stub — przepisy bez bazy i index.html bez sieci; przywraca oryginały po teście.
func stub(t *testing.T, recipes map[uint]*shareRecipe, tpl string) {
	t.Helper()
	origLoad, origFetch := loadShareRecipe, fetchIndex
	loadShareRecipe = func(_ context.Context, id uint) (*shareRecipe, error) {
		if r, ok := recipes[id]; ok {
			return r, nil
		}
		return nil, errors.New("not found")
	}
	fetchIndex = func(_ context.Context, origin string) (string, error) {
		if tpl == "" {
			return "", errors.New("offline")
		}
		return tpl, nil
	}
	indexMu.Lock()
	indexHTML = ""
	indexMu.Unlock()
	t.Setenv("SHARE_ORIGIN", "https://uszatakuchnia.test")
	t.Cleanup(func() {
		loadShareRecipe, fetchIndex = origLoad, origFetch
		indexMu.Lock()
		indexHTML = ""
		indexMu.Unlock()
	})
}

func get(path string, h http.HandlerFunc, id string) string {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	r.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h(w, r)
	return w.Body.String()
}

var ragu = &shareRecipe{ID: 6, Name: "Ragù alla bolognese", Servings: 6, TimeMin: 180, TimeMax: ip(300), Status: entities.StatusPublished}

func TestRecipePageMeta(t *testing.T) {
	stub(t, map[uint]*shareRecipe{6: ragu}, tplIndex)
	body := get("/recipes/6", RecipeSharePage, "6")
	for _, want := range []string{
		`<title>Ragù alla bolognese · Uszata Kuchnia</title>`,
		`<meta property="og:title" content="Ragù alla bolognese" />`,
		`<meta property="og:description" content="Przepis na 6 porcji · 3–5 h · Uszata Kuchnia" />`,
		`<meta name="description" content="Przepis na 6 porcji · 3–5 h · Uszata Kuchnia" />`,
		`<meta property="og:url" content="https://uszatakuchnia.test/recipes/6" />`,
		`<meta property="og:type" content="article" />`,
		`<meta name="twitter:card" content="summary" />`,
		`<script type="module" src="/assets/index-abc.js">`, // reszta index.html bez zmian
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s\n%s", want, body)
		}
	}
	if strings.Contains(body, "ogólny") || strings.Count(body, "<title>") != 1 {
		t.Errorf("generic meta not replaced:\n%s", body)
	}
	if strings.Contains(body, "og:image") {
		t.Error("og:image without a default image")
	}
}

func TestShoppingPageMeta(t *testing.T) {
	stub(t, map[uint]*shareRecipe{6: ragu}, tplIndex)
	body := get("/recipes/6/zakupy?porcje=8&zamienniki=6:21", ShoppingSharePage, "6")
	for _, want := range []string{
		`<title>Lista zakupów: Ragù alla bolognese · Uszata Kuchnia</title>`,
		`<meta property="og:title" content="Lista zakupów: Ragù alla bolognese · 8 porcji" />`,
		`<meta property="og:url" content="https://uszatakuchnia.test/recipes/6/zakupy?porcje=8&amp;zamienniki=6:21" />`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s\n%s", want, body)
		}
	}
	// bez ?porcje= (albo z bzdurą) — porcje przepisu
	for _, path := range []string{"/recipes/6/zakupy", "/recipes/6/zakupy?porcje=99"} {
		if b := get(path, ShoppingSharePage, "6"); !strings.Contains(b, `· 6 porcji" />`) {
			t.Errorf("%s: servings from recipe expected\n%s", path, b)
		}
	}
}

func TestDraftAndMissingGetGenericMeta(t *testing.T) {
	draft := &shareRecipe{ID: 9, Name: "Tajny szkic", Servings: 2, Status: entities.StatusDraft}
	stub(t, map[uint]*shareRecipe{9: draft}, tplIndex)
	for _, id := range []string{"9", "404"} {
		body := get("/recipes/"+id, RecipeSharePage, id)
		if strings.Contains(body, "Tajny szkic") {
			t.Errorf("draft title leaked:\n%s", body)
		}
		if !strings.Contains(body, `<title>Uszata Kuchnia</title>`) ||
			!strings.Contains(body, `<meta property="og:title" content="Uszata Kuchnia" />`) {
			t.Errorf("id %s: generic meta expected\n%s", id, body)
		}
	}
}

func TestTaglineDescriptionAndEscaping(t *testing.T) {
	long := strings.Repeat("Długi opis przepisu ", 20)
	recipes := map[uint]*shareRecipe{
		1: {ID: 1, Name: `Sos "ostry" <b>`, Servings: 1, Tagline: sp("Ostro & szybko"), Status: entities.StatusPublished},
		2: {ID: 2, Name: "Zupa", Servings: 4, Description: sp("Pierwsza linia.\n\nDruga  linia."), Status: entities.StatusPublished},
		3: {ID: 3, Name: "Chleb", Servings: 2, Description: sp(long), Status: entities.StatusPublished},
		4: {ID: 4, Name: "Jajko", Servings: 1, TimeMin: 10, Status: entities.StatusPublished},
	}
	stub(t, recipes, tplIndex)
	b1 := get("/recipes/1", RecipeSharePage, "1")
	if !strings.Contains(b1, `content="Sos &#34;ostry&#34; &lt;b&gt;"`) || !strings.Contains(b1, `content="Ostro &amp; szybko"`) {
		t.Errorf("escaping / tagline:\n%s", b1)
	}
	if b := get("/recipes/2", RecipeSharePage, "2"); !strings.Contains(b, `content="Pierwsza linia. Druga linia."`) {
		t.Errorf("description whitespace:\n%s", b)
	}
	b3 := get("/recipes/3", RecipeSharePage, "3")
	i := strings.Index(b3, `og:description" content="`) + len(`og:description" content="`)
	desc := b3[i : i+strings.Index(b3[i:], `"`)]
	if !strings.HasSuffix(desc, "…") || len([]rune(desc)) > descMaxRunes+1 || strings.HasSuffix(desc, " …") {
		t.Errorf("truncated description %q", desc)
	}
	if b := get("/recipes/4", RecipeSharePage, "4"); !strings.Contains(b, `content="Przepis na 1 porcję · 10 min · Uszata Kuchnia"`) {
		t.Errorf("fallback description:\n%s", b)
	}
}

func TestMinimalPageWhenIndexUnavailable(t *testing.T) {
	stub(t, map[uint]*shareRecipe{6: ragu}, "")
	body := get("/recipes/6", RecipeSharePage, "6")
	for _, want := range []string{`<html lang="pl">`, `<meta property="og:title" content="Ragù alla bolognese" />`, `<h1>Ragù alla bolognese</h1>`} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s\n%s", want, body)
		}
	}
}

func TestFormatTimeMatchesFrontend(t *testing.T) {
	cases := []struct {
		min  int
		max  *int
		want string
	}{
		{0, nil, ""},
		{40, nil, "40 min"},
		{90, nil, "1 h 30 min"},
		{120, nil, "2 h"},
		{20, ip(30), "20–30 min"},
		{180, ip(300), "3–5 h"},
		{150, ip(270), "2,5–4,5 h"},
		{45, ip(90), "45 min – 1 h 30 min"},
	}
	for _, c := range cases {
		if got := formatTime(c.min, c.max); got != c.want {
			t.Errorf("formatTime(%d, %v) = %q, want %q", c.min, c.max, got, c.want)
		}
	}
}
