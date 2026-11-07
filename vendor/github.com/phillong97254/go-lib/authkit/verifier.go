package authkit

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

	"github.com/golang-jwt/jwt/v5"
)

// Verifier kiểm tra chữ ký JWT bằng public key lấy từ JWKS endpoint của auth server
// (Keycloak: {realm}/protocol/openid-connect/certs), rồi kiểm tra iss, aud, exp.
type Verifier struct {
	JWKSURL  string
	Issuer   string
	Audience string

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
}

func NewVerifier(jwksURL, issuer, audience string) *Verifier {
	return &Verifier{JWKSURL: jwksURL, Issuer: issuer, Audience: audience, keys: map[string]*rsa.PublicKey{}}
}

func (v *Verifier) Verify(raw string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(raw, claims, v.keyFunc,
		jwt.WithValidMethods([]string{"RS256"}), // chặn tấn công đổi alg sang "none"/HS256
		jwt.WithIssuer(v.Issuer),
		jwt.WithAudience(v.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (v *Verifier) keyFunc(t *jwt.Token) (any, error) {
	kid, _ := t.Header["kid"].(string)
	if k := v.cached(kid); k != nil {
		return k, nil
	}
	// kid lạ => có thể auth server vừa xoay key, tải lại JWKS (tối đa 1 lần / 10s).
	if err := v.refresh(); err != nil {
		return nil, err
	}
	if k := v.cached(kid); k != nil {
		return k, nil
	}
	return nil, fmt.Errorf("unknown kid %q", kid)
}

func (v *Verifier) cached(kid string) *rsa.PublicKey {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.keys[kid]
}

func (v *Verifier) refresh() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if time.Since(v.fetchedAt) < 10*time.Second && len(v.keys) > 0 {
		return nil
	}
	resp, err := http.Get(v.JWKSURL)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()
	var set JWKS
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}
	keys := map[string]*rsa.PublicKey{}
	for _, k := range set.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pk, err := k.RSAPublicKey()
		if err != nil {
			return err
		}
		keys[k.Kid] = pk
	}
	v.keys, v.fetchedAt = keys, time.Now()
	return nil
}

// ---- JWKS (RFC 7517) ----

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewJWK(kid string, pub *rsa.PublicKey) JWK {
	return JWK{
		Kty: "RSA", Kid: kid, Use: "sig", Alg: "RS256",
		N: base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E: base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}
}

func (k JWK) RSAPublicKey() (*rsa.PublicKey, error) {
	n, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}
	e, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}
	if len(n) == 0 || len(e) == 0 {
		return nil, errors.New("invalid jwk")
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}, nil
}
