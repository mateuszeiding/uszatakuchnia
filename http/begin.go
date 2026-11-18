package http

import (
	"net/http"
	"os"
	"strings"

	"uszatakuchnia/auth"
)

// Begin — pierwsza instrukcja KAŻDEGO handlera w api/ (pilnuje tego test w begin_test.go):
//
//	if resp.Begin(w, r) {
//		return
//	}
//
// Ustawia CORS, odpowiada na preflight i odrzuca zapis (każda metoda poza GET/HEAD)
// bez uprawnień właściciela: 401 bez tokenu / ze złym tokenem, 403 z tokenem kogoś innego.
// Zwraca true, gdy odpowiedź jest już wysłana. Nowy endpoint zapisujący jest więc
// chroniony domyślnie — nie trzeba o tym pamiętać w handlerze.
func Begin(w http.ResponseWriter, r *http.Request) bool {
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return false
	}
	switch auth.Check(r) {
	case auth.Owner:
		return false
	case auth.Forbidden:
		Error(w, http.StatusForbidden, "forbidden")
	default:
		w.Header().Set("WWW-Authenticate", `Bearer realm="uszatakuchnia"`)
		Error(w, http.StatusUnauthorized, "unauthorized")
	}
	return true
}

// IsOwner — do odczytów, które pokazują więcej właścicielowi (szkice).
func IsOwner(r *http.Request) bool { return auth.IsOwner(r) }

// Produkcja i podglądy Vercela serwują frontend i API z jednej domeny, więc CORS
// potrzebny jest tylko lokalnie (Vite na :5173 woła API na :3000). Dodatkowe domeny:
// CORS_ORIGINS="https://a.pl,https://b.pl". Klienci spoza przeglądarki (curl, sesje
// Claude Code) CORS nie dotyczy.
var devOrigins = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
	"http://localhost:4173",
	"http://127.0.0.1:4173",
}

func originAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range devOrigins {
		if o == origin {
			return true
		}
	}
	for _, o := range strings.Split(os.Getenv("CORS_ORIGINS"), ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Origin")
	origin := r.Header.Get("Origin")
	if !originAllowed(origin) {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")
}
