package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testKey  = "0123456789abcdef0123456789abcdef-test"
	ownerSub = "auth0|owner"
	audience = "https://api.test"
)

type fixture struct {
	cfg  Config
	priv *rsa.PrivateKey
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kid": "k1", "kty": "RSA",
			"n": base64.RawURLEncoding.EncodeToString(priv.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.E)).Bytes()),
		}}})
	}))
	t.Cleanup(srv.Close)
	return fixture{
		priv: priv,
		cfg: Config{
			APIKey:   testKey,
			Issuer:   srv.URL + "/",
			Audience: audience,
			OwnerSub: ownerSub,
			JWKSURL:  srv.URL + "/.well-known/jwks.json",
		},
	}
}

func (f fixture) token(t *testing.T, sub string, mutate func(*jwt.RegisteredClaims)) string {
	t.Helper()
	c := jwt.RegisteredClaims{
		Issuer:    f.cfg.Issuer,
		Subject:   sub,
		Audience:  jwt.ClaimStrings{audience, f.cfg.Issuer + "userinfo"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	if mutate != nil {
		mutate(&c)
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	tok.Header["kid"] = "k1"
	s, err := tok.SignedString(f.priv)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func req(auth string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/api/recipes", nil)
	if auth != "" {
		r.Header.Set("Authorization", auth)
	}
	return r
}

func TestCheck(t *testing.T) {
	f := newFixture(t)
	other, _ := rsa.GenerateKey(rand.Reader, 2048)
	forged := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer: f.cfg.Issuer, Subject: ownerSub, Audience: jwt.ClaimStrings{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	forged.Header["kid"] = "k1"
	forgedStr, _ := forged.SignedString(other)

	cases := []struct {
		name string
		auth string
		want Result
	}{
		{"brak nagłówka", "", Anonymous},
		{"klucz API", "Bearer " + testKey, Owner},
		{"klucz API, małe bearer", "bearer " + testKey, Owner},
		{"zły klucz", "Bearer " + testKey + "x", Invalid},
		{"nie Bearer", "Basic " + testKey, Invalid},
		{"token właściciela", "Bearer " + f.token(t, ownerSub, nil), Owner},
		{"token innego konta", "Bearer " + f.token(t, "auth0|someone", nil), Forbidden},
		{"wygasły token", "Bearer " + f.token(t, ownerSub, func(c *jwt.RegisteredClaims) {
			c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
		}), Invalid},
		{"obca audience", "Bearer " + f.token(t, ownerSub, func(c *jwt.RegisteredClaims) {
			c.Audience = jwt.ClaimStrings{"https://other"}
		}), Invalid},
		{"obcy issuer", "Bearer " + f.token(t, ownerSub, func(c *jwt.RegisteredClaims) {
			c.Issuer = "https://evil.example/"
		}), Invalid},
		{"podpis innym kluczem", "Bearer " + forgedStr, Invalid},
		{"alg none", "Bearer eyJhbGciOiJub25lIn0.eyJzdWIiOiJhdXRoMHxvd25lciJ9.", Invalid},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := f.cfg.Check(req(tc.auth)); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShortKeyIgnored(t *testing.T) {
	c := Config{APIKey: "short"}
	if got := c.Check(req("Bearer short")); got != Invalid {
		t.Fatalf("got %v, want Invalid", got)
	}
}

func TestNoOwnerConfigured(t *testing.T) {
	f := newFixture(t)
	f.cfg.OwnerSub = ""
	if got := f.cfg.Check(req("Bearer " + f.token(t, ownerSub, nil))); got != Forbidden {
		t.Fatalf("got %v, want Forbidden", got)
	}
}

func TestFromEnvFallsBackToVite(t *testing.T) {
	t.Setenv("AUTH0_DOMAIN", "")
	t.Setenv("VITE_AUTH0_DOMAIN", "example.eu.auth0.com")
	t.Setenv("AUTH0_AUDIENCE", "")
	t.Setenv("VITE_AUTH0_AUDIENCE", "https://api.example")
	t.Setenv("AUTH0_OWNER_SUB", "")
	t.Setenv("VITE_AUTH0_OWNER_SUB", "auth0|x")
	c := FromEnv()
	if c.Issuer != "https://example.eu.auth0.com/" || c.Audience != "https://api.example" || c.OwnerSub != "auth0|x" {
		t.Fatalf("unexpected config %+v", c)
	}
}
