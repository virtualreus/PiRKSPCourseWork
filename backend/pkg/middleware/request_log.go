package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nikitatisenko/pirksp/pkg/logger"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.FromContext(r.Context()).
			Info(fmt.Sprintf("%s %s %s | %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start)))
	})
}
