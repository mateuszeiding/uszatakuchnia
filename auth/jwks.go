package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// Auth0 public keys (JWKS) cached in function memory. A warm Vercel function doesn't
// fetch them on every request; an unknown `kid` (key rotation in Auth0) forces
// a refetch, at most once a minute.

const jwksRefetchEvery = time.Minute

type jwksCache struct {
	mu      sync.Mutex
	keys    map[string]*rsa.PublicKey
	fetched time.Time
}

var (
	cachesMu sync.Mutex
	caches   = map[string]*jwksCache{}
	jwksHTTP = &http.Client{Timeout: 5 * time.Second}
)

func publicKey(url, kid string) (*rsa.PublicKey, error) {
	cachesMu.Lock()
	c := caches[url]
	if c == nil {
		c = &jwksCache{}
		caches[url] = c
	}
	cachesMu.Unlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if k := c.keys[kid]; k != nil {
		return k, nil
	}
	if time.Since(c.fetched) < jwksRefetchEvery && c.keys != nil {
		return nil, fmt.Errorf("unknown kid %q", kid)
	}
	keys, err := fetchJWKS(url)
	c.fetched = time.Now()
	if err != nil {
		return nil, err
	}
	c.keys = keys
	if k := keys[kid]; k != nil {
		return k, nil
	}
	return nil, fmt.Errorf("unknown kid %q", kid)
}

func fetchJWKS(url string) (map[string]*rsa.PublicKey, error) {
	res, err := jwksHTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks: HTTP %d", res.StatusCode)
	}
	var body struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	keys := map[string]*rsa.PublicKey{}
	for _, k := range body.Keys {
		if k.Kty != "RSA" {
			continue
		}
		n, err1 := base64.RawURLEncoding.DecodeString(k.N)
		e, err2 := base64.RawURLEncoding.DecodeString(k.E)
		if err1 != nil || err2 != nil {
			continue
		}
		keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
	}
	if len(keys) == 0 {
		return nil, errors.New("jwks: no RSA keys")
	}
	return keys, nil
}
