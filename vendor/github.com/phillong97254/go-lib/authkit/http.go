package authkit

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Authenticate là HTTP middleware: đọc "Authorization: Bearer <jwt>", verify,
// gắn claims vào context. Sai/thiếu token => 401.
func Authenticate(v *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || raw == "" {
				w.Header().Set("WWW-Authenticate", `Bearer`)
				WriteError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			claims, err := v.Verify(raw)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
				WriteError(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}
			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims, raw)))
		})
	}
}

// Require là bước authorization: token hợp lệ nhưng thiếu scope/role => 403.
func Require(p Policy, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, ok := ClaimsFrom(r.Context())
		if !ok {
			WriteError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}
		if !p.Allows(c) {
			w.Header().Set("WWW-Authenticate", `Bearer error="insufficient_scope"`)
			WriteError(w, http.StatusForbidden, "forbidden: insufficient scope/role")
			return
		}
		next(w, r)
	}
}

func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, map[string]string{"error": msg})
}
