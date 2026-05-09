package middleware

import (
	"log/slog"
	"net/http"

	"github.com/nikitatisenko/pirksp/pkg/logger"
)

func LoggerContext(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.ContextWithLogger(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
