package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/platform"
)

type contextKey string

const claimsKey contextKey = "claims"

func ClaimsFromContext(ctx context.Context) map[string]any {
	claims, _ := ctx.Value(claimsKey).(map[string]any)
	return claims
}

func RequireLogin(secret string, requireAdmin bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)
			if len(parts) != 2 {
				platform.WriteError(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "Please login to access this info."})
				return
			}
			claims, err := platform.VerifyJWT(parts[1], secret)
			if err != nil {
				platform.WriteError(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "Please login to access this info."})
				return
			}
			if requireAdmin {
				if role, _ := claims["role"].(string); !strings.EqualFold(role, "admin") {
					platform.WriteError(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "You do not have permission to perform this action"})
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey, claims)))
		})
	}
}

func WithCORS(next http.Handler) http.Handler {
	allowedOrigins := parseAllowedOrigins(os.Getenv("ORCS_ALLOWED_ORIGINS"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if len(allowedOrigins) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseAllowedOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{"http://localhost:3000", "http://localhost:3001"}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:3000", "http://localhost:3001"}
	}
	return out
}

func isAllowedOrigin(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, allowedOrigin := range allowed {
		if strings.EqualFold(origin, allowedOrigin) {
			return true
		}
	}
	return false
}
