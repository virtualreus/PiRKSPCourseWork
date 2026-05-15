package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/infrastructure/auth"
	httperror "github.com/nikitatisenko/pirksp/pkg/http/error"
)

type ctxKey int

const (
	userIDKey       ctxKey = 1
	platformRoleKey ctxKey = 2
)

func AuthRequired(tokens *auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				httperror.Unauthorized(w, "missing bearer token")
				return
			}

			raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			userID, role, err := tokens.Parse(raw)
			if err != nil {
				httperror.Unauthorized(w, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, platformRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func PlatformRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(platformRoleKey).(string)
	return role, ok
}

func OrganizerRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := PlatformRoleFromContext(r.Context())
		if !ok || role != "organizer" {
			httperror.Forbidden(w, "organizer role required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
