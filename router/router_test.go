package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"uszatakuchnia/api"
	"uszatakuchnia/router"
)

const testKey = "0123456789abcdef0123456789abcdef-router"

// isBotPage: a page for link-preview bots, the only routes outside /api/, GET only, under /recipes/.
func isBotPage(rt router.Route) bool {
	return rt.Method == http.MethodGet && strings.HasPrefix(rt.Path, "/recipes/")
}

func samplePath(rt router.Route) string { return strings.ReplaceAll(rt.Path, "{id}", "7") }

// call sends a request to the real Vercel function (api.Handler: Begin + router).
func call(method, path string, hdr map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	for k, v := range hdr {
		r.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	api.Handler(w, r)
	return w
}

func jsonError(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type %q, want application/json", ct)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Error == "" {
		t.Errorf("body %q is not {\"error\": …}", w.Body.String())
	}
	return body.Error
}

// Every route in the table (new ones are picked up automatically): a write without a token or with a bad key → 401,
// preflight from an allowed origin → 204 with CORS, CORS headers set once.
func TestEveryRouteIsGuardedByBegin(t *testing.T) {
	t.Setenv("API_ADMIN_KEY", testKey)
	if len(router.Routes) == 0 {
		t.Fatal("no routes")
	}
	for _, rt := range router.Routes {
		path := samplePath(rt)
		name := rt.Method + " " + path
		if !strings.HasPrefix(rt.Path, "/api/") && !isBotPage(rt) {
			t.Errorf("%s: route outside /api/ (only GET /recipes/… pages for link-preview bots)", name)
		}

		if rt.Method != http.MethodGet {
			w := call(rt.Method, path, nil)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("%s without token: %d, want 401", name, w.Code)
			} else {
				jsonError(t, w)
			}
			if w := call(rt.Method, path, map[string]string{"Authorization": "Bearer wrong"}); w.Code != http.StatusUnauthorized {
				t.Errorf("%s with wrong key: %d, want 401", name, w.Code)
			}
		}

		w := call(http.MethodOptions, path, map[string]string{
			"Origin":                        "http://localhost:5173",
			"Access-Control-Request-Method": rt.Method,
		})
		if w.Code != http.StatusNoContent {
			t.Errorf("OPTIONS %s: %d, want 204", path, w.Code)
		}
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("OPTIONS %s: ACAO %q", path, got)
		}
		if v := w.Header().Values("Vary"); len(v) != 1 {
			t.Errorf("OPTIONS %s: Vary %v — Begin called more than once?", path, v)
		}
	}
}

// Begin also runs before 404 and 405: a write to an unknown path or with a wrong method, no token → 401.
// With the key, the router answers 404 / 405 in JSON. None of these requests touch the database.
func TestBeginRunsBeforeNotFoundAndMethodNotAllowed(t *testing.T) {
	t.Setenv("API_ADMIN_KEY", testKey)
	withKey := map[string]string{"Authorization": "Bearer " + testKey}

	cases := []struct {
		method, path string
		hdr          map[string]string
		want         int
	}{
		{http.MethodPost, "/api/nope", nil, http.StatusUnauthorized},
		{http.MethodDelete, "/api/recipes/list", nil, http.StatusUnauthorized},
		{http.MethodPost, "/api/nope", withKey, http.StatusNotFound},
		{http.MethodGet, "/api/nope", nil, http.StatusNotFound},
		{http.MethodGet, "/api", nil, http.StatusNotFound},
		{http.MethodGet, "/api/recipes/", nil, http.StatusNotFound},
		{http.MethodGet, "/api/recipes/abc", nil, http.StatusNotFound},
		{http.MethodGet, "/api/recipes/0", nil, http.StatusNotFound},
		{http.MethodGet, "/api/recipes/1x", nil, http.StatusNotFound},
		{http.MethodGet, "/api/taxonomy/units?kind=equipment", nil, http.StatusMethodNotAllowed},
		{http.MethodGet, "/api/recipes", nil, http.StatusMethodNotAllowed},
		{http.MethodPut, "/api/recipes", withKey, http.StatusMethodNotAllowed},
		{http.MethodPost, "/api/recipes/7", withKey, http.StatusMethodNotAllowed},
		{http.MethodGet, "/api/ingredients/7", nil, http.StatusMethodNotAllowed},
	}
	for _, c := range cases {
		w := call(c.method, c.path, c.hdr)
		if w.Code != c.want {
			t.Errorf("%s %s: %d, want %d", c.method, c.path, w.Code, c.want)
			continue
		}
		jsonError(t, w)
		if c.want == http.StatusMethodNotAllowed && w.Header().Get("Allow") == "" {
			t.Errorf("%s %s: 405 without Allow", c.method, c.path)
		}
	}

	if w := call(http.MethodOptions, "/api/nope", map[string]string{"Origin": "http://localhost:5173"}); w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS unknown path: %d, want 204", w.Code)
	}
}

// Every route reaches its own handler with {id} in r.PathValue (handlers replaced with stubs,
// no database). The query string stays intact (?moveTo=, ?admin=true).
func TestEveryRouteReachesItsHandler(t *testing.T) {
	stubs := make([]router.Route, len(router.Routes))
	hit := -1
	var gotID, gotQuery string
	for i, rt := range router.Routes {
		i := i
		stubs[i] = router.Route{Method: rt.Method, Path: rt.Path, Handler: func(w http.ResponseWriter, r *http.Request) {
			hit, gotID, gotQuery = i, r.PathValue("id"), r.URL.RawQuery
			w.WriteHeader(http.StatusTeapot)
		}}
	}
	h := router.New(stubs)

	seen := map[string]bool{}
	for i, rt := range stubs {
		key := rt.Method + " " + rt.Path
		if seen[key] {
			t.Errorf("duplicate route %s", key)
		}
		seen[key] = true

		methods := []string{rt.Method}
		if rt.Method == http.MethodGet {
			methods = append(methods, http.MethodHead)
		}
		for _, m := range methods {
			hit, gotID, gotQuery = -1, "", ""
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(m, samplePath(rt)+"?moveTo=3", nil))
			if w.Code != http.StatusTeapot || hit != i {
				t.Errorf("%s %s: code %d, handler #%d, want #%d", m, samplePath(rt), w.Code, hit, i)
				continue
			}
			if strings.Contains(rt.Path, "{id}") && gotID != "7" {
				t.Errorf("%s %s: PathValue(id) = %q", m, rt.Path, gotID)
			}
			if gotQuery != "moveTo=3" {
				t.Errorf("%s %s: query %q", m, rt.Path, gotQuery)
			}
		}
	}
}

type vercelRewrite struct {
	Source, Destination string
	Has                 []struct{ Type, Key, Value string }
}

func readVercelRewrites(t *testing.T) ([]vercelRewrite, string) {
	t.Helper()
	raw, err := os.ReadFile("../vercel.json")
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct{ Rewrites []vercelRewrite }
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatal(err)
	}
	return cfg.Rewrites, string(raw)
}

// vercel.json: the whole API in one rewrite (no conditions) to the only function.
// The other rewrites to the function are bot preview pages, each conditioned on the user agent.
func TestVercelJSONHasSingleAPIRewrite(t *testing.T) {
	rewrites, raw := readVercelRewrites(t)
	api := 0
	for _, rw := range rewrites {
		if strings.HasPrefix(rw.Destination, "/api") {
			if rw.Destination != "/api/index" {
				t.Errorf("rewrite %s → %s: API must go to /api/index", rw.Source, rw.Destination)
			}
			if len(rw.Has) == 0 {
				api++
			}
		}
	}
	if api != 1 {
		t.Errorf("%d unconditional API rewrites in vercel.json, want 1", api)
	}
	if strings.Contains(string(raw), "kind=") {
		t.Error("vercel.json still has ?kind=")
	}
}

// Bot preview pages: every route outside /api/ has a rewrite to /api/index in vercel.json,
// conditioned on the user agent (people get the static index.html), placed before the SPA rewrite.
func TestBotPagesRewrittenOnlyForBots(t *testing.T) {
	rewrites, _ := readVercelRewrites(t)
	spa := -1
	for i, rw := range rewrites {
		if rw.Source == "/(.*)" {
			spa = i
		}
	}
	if spa < 0 {
		t.Fatal("no SPA rewrite /(.*)")
	}
	pages := 0
	for _, rt := range router.Routes {
		if !isBotPage(rt) {
			continue
		}
		pages++
		want := strings.ReplaceAll(rt.Path, "{id}", `:id(\d+)`)
		found := false
		for i, rw := range rewrites {
			if rw.Source != want {
				continue
			}
			found = true
			if rw.Destination != "/api/index" {
				t.Errorf("%s → %s, want /api/index", rw.Source, rw.Destination)
			}
			if i > spa {
				t.Errorf("%s: bot rewrite after the SPA rewrite never matches", rw.Source)
			}
			ua := false
			for _, h := range rw.Has {
				if h.Type == "header" && strings.EqualFold(h.Key, "user-agent") &&
					strings.Contains(h.Value, "facebookexternalhit") && strings.Contains(h.Value, "WhatsApp") {
					ua = true
				}
			}
			if !ua {
				t.Errorf("%s: rewrite without user-agent condition — humans would hit the function", rw.Source)
			}
		}
		if !found {
			t.Errorf("GET %s: no bot rewrite %q in vercel.json", rt.Path, want)
		}
	}
	if pages == 0 {
		t.Error("no bot page routes")
	}
}
