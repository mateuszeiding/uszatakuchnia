package http

import (
	"net/http"
	"os"
	"strings"

	"uszatakuchnia/auth"
)

// Begin is the first statement of EVERY handler in api/ (enforced by begin_test.go):
//
//	if resp.Begin(w, r) {
//		return
//	}
//
// It sets CORS, answers preflight and rejects writes (any method but GET/HEAD)
// without owner rights: 401 without a token / with a bad one, 403 with someone else's token.
// Returns true when the response has already been sent. A new write endpoint is therefore
// protected by default; the handler doesn't need to remember it.
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

// IsOwner is for reads that show the owner more (drafts).
func IsOwner(r *http.Request) bool { return auth.IsOwner(r) }

// Production and Vercel previews serve the frontend and the API from one domain, so CORS
// is needed only locally (Vite on :5173 calls the API on :3000). Extra domains:
// CORS_ORIGINS="https://a.pl,https://b.pl". Non-browser clients (curl, Claude Code
// sessions) aren't affected by CORS.
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
