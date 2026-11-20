package handlers

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
)

// Podgląd linku w komunikatorach (Messenger, WhatsApp, Slack…). Boty podglądu nie uruchamiają JS,
// więc meta przepisu muszą przyjść w HTML z serwera. vercel.json kieruje tu TYLKO ich user-agenty
// (rewrite z `has`); ludzie dostają statyczny index.html bez funkcji.
//
// Strona = index.html aplikacji z tego samego wdrożenia (skrypty z hashami pasują do buildu),
// w którym blok między <!-- meta:start --> a <!-- meta:end --> zastępujemy meta przepisu.
// Gdy pobranie się nie uda — minimalny HTML z samymi meta (bot i tak nie uruchamia JS).

// ogDefaultImage — obraz podglądu (adres bezwzględny albo ścieżka od korzenia). Pusty = bez og:image.
// Jeden obraz od designera dla wszystkich stron, dopóki przepisy nie mają zdjęć (web/public/og-image.png,
// ten sam w web/index.html).
const (
	ogDefaultImage       = "/og-image.png"
	ogDefaultImageWidth  = 1200
	ogDefaultImageHeight = 630
	ogDefaultImageAlt    = "uszatakuchnia — przepisy z domowej kuchni"
)

const (
	siteName        = "Uszata Kuchnia"
	siteDescription = "Przepisy Uszatej Kuchni: składniki z zamiennikami, lista zakupów i tryb gotowania krok po kroku."
	descMaxRunes    = 150
	indexCacheTTL   = 10 * time.Minute
)

// RecipeSharePage — GET /recipes/{id} (tylko boty podglądu linków).
func RecipeSharePage(w http.ResponseWriter, r *http.Request) { sharePage(w, r, false) }

// ShoppingSharePage — GET /recipes/{id}/zakupy (tylko boty podglądu linków).
func ShoppingSharePage(w http.ResponseWriter, r *http.Request) { sharePage(w, r, true) }

type pageMeta struct {
	Title       string // <title>
	OGTitle     string
	Description string
	URL         string // og:url; "" = bez
	Image       string // og:image; "" = bez
	Type        string // og:type
}

// shareRecipe — to, czego potrzebują meta (bez kroków i składników).
type shareRecipe struct {
	ID          uint
	Name        string
	Servings    int
	Description *string
	Tagline     *string
	TimeMin     int
	TimeMax     *int
	Status      string
}

// loadShareRecipe — podmieniana w testach (bez bazy).
var loadShareRecipe = func(ctx context.Context, id uint) (*shareRecipe, error) {
	var rec shareRecipe
	err := db.DB().WithContext(ctx).Model(&entities.Recipe{}).
		Select("id, name, servings, description, tagline, time_min, time_max, status").
		Where("id = ?", id).Take(&rec).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func sharePage(w http.ResponseWriter, r *http.Request, shopping bool) {
	origin := shareOrigin()
	meta := genericMeta(origin)

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()
	// szkic, brak przepisu albo błąd bazy → meta ogólne (bez tytułu szkicu), dalej 200
	if rec, err := loadShareRecipe(ctx, pathID(r)); err == nil && rec.Status == entities.StatusPublished {
		path := r.URL.Path
		if shopping && r.URL.RawQuery != "" {
			path += "?" + r.URL.RawQuery // porcje i zamienniki: każdy wariant listy ma swój podgląd
		}
		meta = recipeMeta(rec, shopping, servingsParam(r, rec.Servings), origin, path)
	}

	page := renderPage(indexTemplate(ctx, origin), meta)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// krótko w CDN: poprawka przepisu w panelu dociera do podglądów w ciągu kilku minut
	w.Header().Set("Cache-Control", "public, max-age=0, s-maxage=300")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, page)
}

// shareOrigin — adres strony do pobrania index.html i og:url. SHARE_ORIGIN (testy lokalne),
// a na Vercelu domena produkcyjna (zmienna systemowa). Nie bierzemy Host z żądania: index.html
// jest trzymany w pamięci funkcji, więc obcy Host podmieniłby stronę innym botom.
func shareOrigin() string {
	if o := strings.TrimRight(os.Getenv("SHARE_ORIGIN"), "/"); o != "" {
		return o
	}
	if h := os.Getenv("VERCEL_PROJECT_PRODUCTION_URL"); h != "" {
		return "https://" + h
	}
	return ""
}

func genericMeta(origin string) pageMeta {
	return pageMeta{
		Title:       siteName,
		OGTitle:     siteName,
		Description: siteDescription,
		URL:         origin,
		Image:       absURL(origin, ogDefaultImage),
		Type:        "website",
	}
}

func recipeMeta(rec *shareRecipe, shopping bool, servings int, origin, path string) pageMeta {
	m := pageMeta{
		Title:       rec.Name + " · " + siteName,
		OGTitle:     rec.Name,
		Description: recipeDescription(rec),
		Image:       absURL(origin, ogDefaultImage),
		Type:        "article",
	}
	if origin != "" {
		m.URL = origin + path
	}
	if shopping {
		m.Title = "Lista zakupów: " + rec.Name + " · " + siteName
		m.OGTitle = fmt.Sprintf("Lista zakupów: %s · %d %s", rec.Name, servings, servingsNoun(servings))
	}
	return m
}

// recipeDescription — tagline, a bez niego początek opisu (~150 znaków, na granicy słowa);
// bez obu: „Przepis na 6 porcji · 3–5 h · Uszata Kuchnia”.
func recipeDescription(rec *shareRecipe) string {
	if rec.Tagline != nil && strings.TrimSpace(*rec.Tagline) != "" {
		return collapse(*rec.Tagline)
	}
	if rec.Description != nil && strings.TrimSpace(*rec.Description) != "" {
		return truncateWords(collapse(*rec.Description), descMaxRunes)
	}
	parts := []string{}
	if rec.Servings > 0 {
		parts = append(parts, fmt.Sprintf("Przepis na %d %s", rec.Servings, servingsAcc(rec.Servings)))
	} else {
		parts = append(parts, "Przepis")
	}
	if t := formatTime(rec.TimeMin, rec.TimeMax); t != "" {
		parts = append(parts, t)
	}
	return strings.Join(append(parts, siteName), " · ")
}

// servingsParam — ?porcje= z linku (1–20, jak na froncie), inaczej porcje przepisu.
func servingsParam(r *http.Request, def int) int {
	var n int
	if _, err := fmt.Sscanf(r.URL.Query().Get("porcje"), "%d", &n); err == nil && n >= 1 && n <= 20 {
		return n
	}
	return def
}

func plural(n int, one, few, many string) string {
	if n == 1 {
		return one
	}
	m10, m100 := n%10, n%100
	if m10 >= 2 && m10 <= 4 && (m100 < 10 || m100 >= 20) {
		return few
	}
	return many
}

// „6 porcji” (mianownik, jak licznik porcji na stronie) i „na 6 porcji” (biernik: porcję / porcje / porcji).
func servingsNoun(n int) string { return plural(n, "porcja", "porcje", "porcji") }
func servingsAcc(n int) string  { return plural(n, "porcję", "porcje", "porcji") }

// formatTime — jak web/src/shared/utils/time.ts (formatTime): „40 min”, „20–30 min”, „3–5 h”,
// „45 min – 1 h 30 min”; "" gdy czasu nie podano.
func formatTime(min int, max *int) string {
	one := func(m int) string {
		if m < 60 {
			return fmt.Sprintf("%d min", m)
		}
		if m%60 == 0 {
			return fmt.Sprintf("%d h", m/60)
		}
		return fmt.Sprintf("%d h %d min", m/60, m%60)
	}
	hours := func(m int) string {
		s := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", float64(m)/60), "0"), ".")
		return strings.ReplaceAll(s, ".", ",")
	}
	if min <= 0 {
		return ""
	}
	if max == nil || *max <= min {
		return one(min)
	}
	if *max < 60 {
		return fmt.Sprintf("%d–%d min", min, *max)
	}
	if min >= 60 && min%15 == 0 && *max%15 == 0 {
		return hours(min) + "–" + hours(*max) + " h"
	}
	return one(min) + " – " + one(*max)
}

var spaces = regexp.MustCompile(`\s+`)

func collapse(s string) string { return strings.TrimSpace(spaces.ReplaceAllString(s, " ")) }

func truncateWords(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	cut := string([]rune(s)[:max])
	if i := strings.LastIndex(cut, " "); i > max/2 {
		cut = cut[:i]
	}
	return strings.TrimRight(cut, " ,.;:–-") + "…"
}

func absURL(origin, p string) string {
	if p == "" || strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	if origin == "" {
		return ""
	}
	return origin + p
}

// ─── index.html ──────────────────────────────────────────────

var (
	indexMu      sync.Mutex
	indexHTML    string
	indexFetched time.Time
)

// fetchIndex — podmieniana w testach.
var fetchIndex = func(ctx context.Context, origin string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, origin+"/index.html", nil)
	if err != nil {
		return "", err
	}
	res, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 512<<10))
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), metaStart) {
		return "", fmt.Errorf("index.html: %d", res.StatusCode)
	}
	return string(body), nil
}

// indexTemplate — index.html z pamięci funkcji (10 min), "" = nie udało się pobrać.
func indexTemplate(ctx context.Context, origin string) string {
	if origin == "" {
		return ""
	}
	indexMu.Lock()
	defer indexMu.Unlock()
	if indexHTML != "" && time.Since(indexFetched) < indexCacheTTL {
		return indexHTML
	}
	if s, err := fetchIndex(ctx, origin); err == nil {
		indexHTML, indexFetched = s, time.Now()
	}
	return indexHTML // przy błędzie: ostatnia dobra kopia albo ""
}

const (
	metaStart = "<!-- meta:start -->"
	metaEnd   = "<!-- meta:end -->"
)

func metaBlock(m pageMeta) string {
	e := html.EscapeString
	var b strings.Builder
	b.WriteString(metaStart + "\n")
	fmt.Fprintf(&b, "    <title>%s</title>\n", e(m.Title))
	fmt.Fprintf(&b, "    <meta name=\"description\" content=\"%s\" />\n", e(m.Description))
	fmt.Fprintf(&b, "    <meta property=\"og:site_name\" content=\"%s\" />\n", e(siteName))
	fmt.Fprintf(&b, "    <meta property=\"og:type\" content=\"%s\" />\n", e(m.Type))
	fmt.Fprintf(&b, "    <meta property=\"og:title\" content=\"%s\" />\n", e(m.OGTitle))
	fmt.Fprintf(&b, "    <meta property=\"og:description\" content=\"%s\" />\n", e(m.Description))
	fmt.Fprintf(&b, "    <meta property=\"og:locale\" content=\"pl_PL\" />\n")
	if m.URL != "" {
		fmt.Fprintf(&b, "    <meta property=\"og:url\" content=\"%s\" />\n", e(m.URL))
	}
	card := "summary"
	if m.Image != "" {
		fmt.Fprintf(&b, "    <meta property=\"og:image\" content=\"%s\" />\n", e(m.Image))
		fmt.Fprintf(&b, "    <meta property=\"og:image:width\" content=\"%d\" />\n", ogDefaultImageWidth)
		fmt.Fprintf(&b, "    <meta property=\"og:image:height\" content=\"%d\" />\n", ogDefaultImageHeight)
		fmt.Fprintf(&b, "    <meta property=\"og:image:alt\" content=\"%s\" />\n", e(ogDefaultImageAlt))
		card = "summary_large_image"
	}
	fmt.Fprintf(&b, "    <meta name=\"twitter:card\" content=\"%s\" />\n", card)
	b.WriteString("    " + metaEnd)
	return b.String()
}

// renderPage — index.html z podmienionym blokiem meta; bez szablonu minimalny dokument.
func renderPage(tpl string, m pageMeta) string {
	block := metaBlock(m)
	if i := strings.Index(tpl, metaStart); i >= 0 {
		if j := strings.Index(tpl[i:], metaEnd); j >= 0 {
			return tpl[:i] + block + tpl[i+j+len(metaEnd):]
		}
	}
	e := html.EscapeString
	link := ""
	if m.URL != "" {
		link = fmt.Sprintf("\n    <p><a href=\"%s\">%s</a></p>", e(m.URL), e(m.URL))
	}
	return "<!doctype html>\n<html lang=\"pl\">\n  <head>\n    <meta charset=\"UTF-8\" />\n    " + block +
		"\n  </head>\n  <body>\n    <h1>" + e(m.OGTitle) + "</h1>\n    <p>" + e(m.Description) + "</p>" + link +
		"\n  </body>\n</html>\n"
}
