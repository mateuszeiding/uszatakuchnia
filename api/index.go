// Package api — jedyna funkcja Vercela. vercel.json kieruje tu całe /api/*, a router rozdziela
// żądania po metodzie i ścieżce (trasy: router/routes.go, logika: handlers/). Pakiety pomocnicze
// muszą leżeć poza api/ — każdy plik tutaj byłby osobną funkcją Vercela.
package api

import (
	"net/http"

	resp "uszatakuchnia/http"
	"uszatakuchnia/router"
)

var routes = router.New(router.Routes)

// Handler — całe API. resp.Begin raz, przed każdą trasą (także przed 404 i 405): CORS, preflight
// i bramka zapisu. Pilnują tego testy w http/begin_test.go i router/router_test.go.
func Handler(w http.ResponseWriter, r *http.Request) {
	if resp.Begin(w, r) {
		return
	}
	routes.ServeHTTP(w, r)
}
