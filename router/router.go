// Package router dispatches API requests to handlers by method and path (http.ServeMux, Go ≥ 1.22).
//
// It's a plain http.Handler: on Vercel the only function in api/index.go calls it (after resp.Begin);
// off Vercel, http.ListenAndServe(addr, http.HandlerFunc(api.Handler)) is enough.
// Routes live in routes.go: a new endpoint is one line in the Routes table.
package router

import (
	"net/http"
	"strconv"
	"strings"

	resp "uszatakuchnia/http"
)

type Route struct {
	Method string // GET also covers HEAD
	// Path is an http.ServeMux pattern. {id} matches only numbers ≥ 1, otherwise 404
	// (handlers read it via r.PathValue("id")).
	Path    string
	Handler http.HandlerFunc
}

// New builds the router from the route table. Unknown path → 404, known path with a wrong method → 405
// (with an Allow header), both as JSON {"error": …} via resp.Error. It doesn't call resp.Begin:
// the function in api/ does that once, before the router. Panics on conflicting routes (like ServeMux).
func New(routes []Route) http.Handler {
	mux := http.NewServeMux()
	paths := http.NewServeMux() // the same paths without a method: tells 405 from 404
	allow := map[string][]string{}
	for _, rt := range routes {
		h := rt.Handler
		if strings.Contains(rt.Path, "{id}") {
			h = requireID(h)
		}
		mux.HandleFunc(rt.Method+" "+rt.Path, h)
		if _, seen := allow[rt.Path]; !seen {
			paths.HandleFunc(rt.Path, http.NotFound) // never called, only used for matching
		}
		allow[rt.Path] = append(allow[rt.Path], rt.Method)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ServeHTTP (not the handler from mux.Handler), because only it sets r.PathValue.
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
