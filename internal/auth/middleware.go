package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type ctxKey int

const claimsKey ctxKey = iota

// ClaimsFromContext extracts JWT claims from the request context.
func ClaimsFromContext(ctx context.Context) *Claims {
	c, _ := ctx.Value(claimsKey).(*Claims)
	return c
}

// Middleware validates the Authorization header and injects claims into context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints.
		if strings.HasPrefix(r.URL.Path, "/api/auth/") {
			next.ServeHTTP(w, r)
			return
		}
		// Allow metrics and static assets.
		if r.URL.Path == "/metrics" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		// Non-API paths serve the frontend, no auth needed.
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			writeAuthError(w, "missing Authorization header")
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			writeAuthError(w, "invalid Authorization format, expected Bearer token")
			return
		}

		claims, err := Verify(token)
		if err != nil {
			writeAuthError(w, "authentication failed: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		slog.Debug("auth", "user", claims.Username, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeAuthError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("WWW-Authenticate", `Bearer realm="OpenResty Admin"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + msg + `"}`))
}