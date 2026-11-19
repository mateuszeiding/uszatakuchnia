// Package handlers — logika endpointów API. Każdy handler obsługuje jedną metodę jednej trasy;
// trasy, metody i 404/405 są w router/routes.go, a CORS i bramkę zapisu (resp.Begin) robi raz
// funkcja Vercela w api/index.go przed routerem. Handlery NIE wołają resp.Begin (podwójne nagłówki
// CORS) — pilnuje tego test w router/.
package handlers

import (
	"net/http"
	"strconv"
)

// pathID — {id} z adresu. Router przepuszcza do handlera tylko liczby ≥ 1 (inaczej 404).
func pathID(r *http.Request) uint {
	id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
	return uint(id)
}
