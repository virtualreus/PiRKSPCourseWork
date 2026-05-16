package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authhttp "github.com/nikitatisenko/pirksp/internal/delivery/http/auth"
	hackathonshttp "github.com/nikitatisenko/pirksp/internal/delivery/http/hackathons"
	"github.com/nikitatisenko/pirksp/internal/delivery/http/health"
	organizerhttp "github.com/nikitatisenko/pirksp/internal/delivery/http/organizer"
	participationhttp "github.com/nikitatisenko/pirksp/internal/delivery/http/participation"
	usershttp "github.com/nikitatisenko/pirksp/internal/delivery/http/users"
	pkgmiddleware "github.com/nikitatisenko/pirksp/pkg/middleware"
)

func (s *Server) initRoutes() {
	s.router.Use(middleware.Recoverer)
	s.router.Use(pkgmiddleware.CORS)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Hackathon Platform API\n"))
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Use(pkgmiddleware.LoggerContext(s.logger))
		r.Use(pkgmiddleware.RequestLog)

		r.Get("/health", health.Check(s.database))

		r.Post("/auth/register", authhttp.Register(s.authUseCase))
		r.Post("/auth/login", authhttp.Login(s.authUseCase))

		r.Get("/hackathons", hackathonshttp.List(s.hackathonUseCase))
		r.Get("/hackathons/{hackathonId}", hackathonshttp.Get(s.hackathonUseCase))

		r.Get("/hackathons/{hackathonId}/teams", participationhttp.ListTeams(s.participationUseCase))
		r.Get("/teams/{teamId}", participationhttp.GetTeam(s.participationUseCase))

		r.Route("/users", func(r chi.Router) {
			r.Use(pkgmiddleware.AuthRequired(s.tokens))
			r.Get("/me", usershttp.GetMe(s.authUseCase))
			r.Get("/me/dashboard", usershttp.GetDashboard(s.participationUseCase))
			r.Patch("/me", usershttp.UpdateMe(s.authUseCase))
		})

		r.Group(func(r chi.Router) {
			r.Use(pkgmiddleware.AuthRequired(s.tokens))

			r.Get("/hackathons/{hackathonId}/participation", participationhttp.GetParticipation(s.participationUseCase))
			r.Post("/hackathons/{hackathonId}/register", participationhttp.Register(s.participationUseCase))
			r.Delete("/hackathons/{hackathonId}/register", participationhttp.Unregister(s.participationUseCase))
			r.Post("/hackathons/{hackathonId}/teams", participationhttp.CreateTeam(s.participationUseCase))

			r.Patch("/teams/{teamId}", participationhttp.UpdateTeam(s.participationUseCase))
			r.Post("/teams/{teamId}/join", participationhttp.JoinTeam(s.participationUseCase))
			r.Post("/teams/{teamId}/leave", participationhttp.LeaveTeam(s.participationUseCase))
			r.Patch("/teams/{teamId}/members/{userId}", participationhttp.UpdateMemberRole(s.participationUseCase))
			r.Get("/teams/{teamId}/submission", participationhttp.GetSubmission(s.participationUseCase))
			r.Put("/teams/{teamId}/submission", participationhttp.UpsertSubmission(s.participationUseCase))
		})

		r.Route("/organizer", func(r chi.Router) {
			r.Use(pkgmiddleware.AuthRequired(s.tokens))
			r.Use(pkgmiddleware.OrganizerRequired)

			r.Get("/hackathons", organizerhttp.ListHackathons(s.hackathonUseCase))
			r.Get("/hackathons/{hackathonId}", organizerhttp.GetHackathon(s.hackathonUseCase))
			r.Post("/hackathons", organizerhttp.CreateHackathon(s.hackathonUseCase))
			r.Patch("/hackathons/{hackathonId}", organizerhttp.UpdateHackathon(s.hackathonUseCase))
			r.Delete("/hackathons/{hackathonId}", organizerhttp.DeleteHackathon(s.hackathonUseCase))
			r.Post("/hackathons/{hackathonId}/publish", organizerhttp.PublishHackathon(s.hackathonUseCase))

			r.Get("/hackathons/{hackathonId}/tracks", organizerhttp.ListTracks(s.hackathonUseCase))
			r.Post("/hackathons/{hackathonId}/tracks", organizerhttp.CreateTrack(s.hackathonUseCase))
			r.Patch("/tracks/{trackId}", organizerhttp.UpdateTrack(s.hackathonUseCase))
			r.Delete("/tracks/{trackId}", organizerhttp.DeleteTrack(s.hackathonUseCase))

			r.Get("/tracks/{trackId}/cases", organizerhttp.ListCases(s.hackathonUseCase))
			r.Post("/tracks/{trackId}/cases", organizerhttp.CreateCase(s.hackathonUseCase))
			r.Patch("/cases/{caseId}", organizerhttp.UpdateCase(s.hackathonUseCase))
			r.Delete("/cases/{caseId}", organizerhttp.DeleteCase(s.hackathonUseCase))

			r.Get("/hackathons/{hackathonId}/registrations", participationhttp.ListRegistrations(s.participationUseCase))
			r.Get("/hackathons/{hackathonId}/submissions", participationhttp.ListSubmissions(s.participationUseCase))
		})
	})
}
