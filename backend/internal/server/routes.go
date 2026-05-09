package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nikitatisenko/pirksp/internal/delivery/http/health"
	pkgmiddleware "github.com/nikitatisenko/pirksp/pkg/middleware"
)

func (s *Server) initRoutes() {
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Hackathon Platform API\n"))
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Use(pkgmiddleware.LoggerContext(s.logger))
		r.Use(pkgmiddleware.RequestLog)

		r.Get("/health", health.Check(s.database))
	})
}
