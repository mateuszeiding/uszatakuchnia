// Package auth identifies who sent an API request.
//
// Only the owner may write (and preview drafts), in one of two ways:
//   - browser: an Auth0 access token (JWT RS256) whose `sub` = AUTH0_OWNER_SUB,
//   - Claude Code sessions on the owner's machine: the API_ADMIN_KEY from `.env.local`.
//
// Both come in the `Authorization: Bearer …` header.
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

type Result int

const (
	Anonymous Result = iota // no Authorization header → 401
	Invalid                 // wrong key, invalid / expired token → 401
	Forbidden               // valid token, but not the owner's → 403
	Owner                   // owner (browser or API key)
)

// MinKeyLen: a shorter API key is ignored (too easy to guess).
const MinKeyLen = 32

type Config struct {
	APIKey   string // API_ADMIN_KEY
	Issuer   string // https://<AUTH0_DOMAIN>/
	Audience string // Auth0 API identifier
	OwnerSub string // AUTH0_OWNER_SUB, e.g. "google-oauth2|123…"
	JWKSURL  string // defaults to Issuer + ".well-known/jwks.json"
}

// env returns the first non-empty variable. The frontend already has VITE_AUTH0_*, so the backend accepts those as a fallback.
func env(names ...string) string {
	for _, n := range names {
		if v := strings.TrimSpace(os.Getenv(n)); v != "" {
			return v
		}
	}
	return ""
}

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

func Check(r *http.Request) Result { return FromEnv().Check(r) }

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
		return Invalid // not a JWT, and the key doesn't match
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
	// compare hashes: constant time regardless of length
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
