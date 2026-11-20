// Package api is the only Vercel function. vercel.json sends all of /api/* here and the router
// dispatches by method and path (routes: router/routes.go, logic: handlers/). Helper packages
// must live outside api/: every file here would become a separate Vercel function.
package api

import (
	"net/http"

	resp "uszatakuchnia/http"
	"uszatakuchnia/router"
)

var routes = router.New(router.Routes)

// Handler serves the whole API. resp.Begin runs once, before every route (404 and 405 included):
// CORS, preflight and the write gate. Enforced by http/begin_test.go and router/router_test.go.
func Handler(w http.ResponseWriter, r *http.Request) {
	if resp.Begin(w, r) {
		return
	}
	routes.ServeHTTP(w, r)
}
