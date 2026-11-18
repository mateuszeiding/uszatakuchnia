// Package auth rozpoznaje, kto wysyła żądanie do API.
//
// Prawo zapisu (i podglądu szkiców) ma tylko właściciel, na dwa sposoby:
//   - przeglądarka: token dostępu Auth0 (JWT RS256), którego `sub` = AUTH0_OWNER_SUB,
//   - sesje Claude Code na komputerze właściciela: klucz API_ADMIN_KEY z `.env.local`.
//
// Oba idą w nagłówku `Authorization: Bearer …`.
package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Result — wynik sprawdzenia żądania.
type Result int

const (
	Anonymous Result = iota // brak nagłówka Authorization → 401
	Invalid                 // zły klucz, nieważny / wygasły token → 401
	Forbidden               // ważny token, ale nie właściciela → 403
	Owner                   // właściciel (przeglądarka albo klucz API)
)

// MinKeyLen — krótszy klucz API jest ignorowany (za łatwy do zgadnięcia).
const MinKeyLen = 32

type Config struct {
	APIKey   string // API_ADMIN_KEY
	Issuer   string // https://<AUTH0_DOMAIN>/
	Audience string // identyfikator API w Auth0
	OwnerSub string // AUTH0_OWNER_SUB, np. "google-oauth2|123…"
	JWKSURL  string // domyślnie Issuer + ".well-known/jwks.json"
}

// env — pierwsza niepusta zmienna. Frontend ma już VITE_AUTH0_*, więc backend przyjmuje je jako zapas.
func env(names ...string) string {
	for _, n := range names {
		if v := strings.TrimSpace(os.Getenv(n)); v != "" {
			return v
		}
	}
	return ""
}

// FromEnv — konfiguracja ze zmiennych środowiskowych (Vercel / .env.local).
func FromEnv() Config {
	c := Config{
		APIKey:   env("API_ADMIN_KEY"),
		Audience: env("AUTH0_AUDIENCE", "VITE_AUTH0_AUDIENCE"),
		OwnerSub: env("AUTH0_OWNER_SUB", "VITE_AUTH0_OWNER_SUB"),
	}
	if d := env("AUTH0_DOMAIN", "VITE_AUTH0_DOMAIN"); d != "" {
		d = strings.TrimSuffix(strings.TrimPrefix(d, "https://"), "/")
		c.Issuer = "https://" + d + "/"
	}
	return c
}

// Check — kto wysłał żądanie.
func Check(r *http.Request) Result { return FromEnv().Check(r) }

// IsOwner — skrót do warunków typu „pokaż szkic”.
func IsOwner(r *http.Request) bool { return Check(r) == Owner }

func (c Config) Check(r *http.Request) Result {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return Anonymous
	}
	scheme, token, ok := strings.Cut(h, " ")
	token = strings.TrimSpace(token)
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return Invalid
	}

	if c.keyMatches(token) {
		return Owner
	}
	if strings.Count(token, ".") != 2 {
		return Invalid // nie JWT, a klucz się nie zgadza
	}

	sub, err := c.verifyJWT(token)
	if err != nil {
		return Invalid
	}
	if c.OwnerSub == "" || sub != c.OwnerSub {
		return Forbidden
	}
	return Owner
}

func (c Config) keyMatches(token string) bool {
	if len(c.APIKey) < MinKeyLen {
		return false
	}
	// porównanie skrótów: stały czas niezależnie od długości
	a := sha256.Sum256([]byte(token))
	b := sha256.Sum256([]byte(c.APIKey))
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

var errNotConfigured = errors.New("auth0 not configured")

func (c Config) verifyJWT(token string) (string, error) {
	if c.Issuer == "" || c.Audience == "" {
		return "", errNotConfigured
	}
	jwksURL := c.JWKSURL
	if jwksURL == "" {
		jwksURL = c.Issuer + ".well-known/jwks.json"
	}

	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (any, error) {
			kid, _ := t.Header["kid"].(string)
			return publicKey(jwksURL, kid)
		},
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(c.Issuer),
		jwt.WithAudience(c.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
