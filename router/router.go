// Package router — rozdziela żądania API na handlery po metodzie i ścieżce (http.ServeMux z Go ≥ 1.22).
//
// To zwykły http.Handler: na Vercelu woła go jedyna funkcja w api/index.go (po resp.Begin), a po
// wyprowadzce z Vercela wystarczy http.ListenAndServe(addr, http.HandlerFunc(api.Handler)).
// Trasy są w routes.go — nowy endpoint to jedna linijka w tabeli Routes.
package router

import (
	"net/http"
	"strconv"
	"strings"

	resp "uszatakuchnia/http"
)

// Route — jedna metoda jednej ścieżki.
type Route struct {
	Method string // GET obejmuje też HEAD
	// Path — wzorzec http.ServeMux. {id} przepuszcza tylko liczbę ≥ 1, inaczej 404
	// (jak dawne :id(\d+) w vercel.json; handlery czytają go przez r.PathValue("id")).
	Path    string
	Handler http.HandlerFunc
}

// New — router z tabeli tras. Nieznana ścieżka → 404, znana ścieżka ze złą metodą → 405
// (z nagłówkiem Allow), oba jako JSON {"error": …} przez resp.Error. Nie woła resp.Begin —
// to robi raz funkcja w api/ przed routerem. Panikuje przy kolidujących trasach (jak ServeMux).
func New(routes []Route) http.Handler {
	mux := http.NewServeMux()
	paths := http.NewServeMux() // te same ścieżki bez metody — odróżnia 405 od 404
	allow := map[string][]string{}
	for _, rt := range routes {
		h := rt.Handler
		if strings.Contains(rt.Path, "{id}") {
			h = requireID(h)
		}
		mux.HandleFunc(rt.Method+" "+rt.Path, h)
		if _, seen := allow[rt.Path]; !seen {
			paths.HandleFunc(rt.Path, http.NotFound) // nigdy niewołany, tylko do dopasowania
		}
		allow[rt.Path] = append(allow[rt.Path], rt.Method)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ServeHTTP (nie sam handler z mux.Handler), bo dopiero on ustawia r.PathValue.
		if _, pattern := mux.Handler(r); pattern != "" {
			mux.ServeHTTP(w, r)
			return
		}
		if _, pattern := paths.Handler(r); pattern != "" {
			w.Header().Set("Allow", strings.Join(allow[pattern], ", "))
			resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		resp.Error(w, http.StatusNotFound, "not found")
	})
}

func requireID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil || id == 0 {
			resp.Error(w, http.StatusNotFound, "not found")
			return
		}
		next(w, r)
	}
}
