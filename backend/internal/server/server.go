package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/nikitatisenko/pirksp/internal/adapter/repository/postgres_repo"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/auth"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/seed"
	"github.com/nikitatisenko/pirksp/internal/usecase"
	"github.com/nikitatisenko/pirksp/internal/usecase/auth_usecase"
	"github.com/nikitatisenko/pirksp/internal/usecase/hackathon_usecase"
	"github.com/nikitatisenko/pirksp/internal/usecase/participation_usecase"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

type Server struct {
	logger   *slog.Logger
	database *postgres.Postgres
	tokens   *auth.TokenService

	usersRepository     repository.UsersRepository
	hackathonsRepository     repository.HackathonsRepository
	participationRepository  repository.ParticipationRepository
	authUseCase              usecase.AuthUseCase
	hackathonUseCase         usecase.HackathonUseCase
	participationUseCase     usecase.ParticipationUseCase

	router *chi.Mux
	server *http.Server
}

func NewServer() (*Server, error) {
	s := &Server{}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) init() error {
	s.logger = logger.New()
	s.router = chi.NewRouter()

	if err := s.initDB(); err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	if err := postgres.MigrateDB(s.database); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if err := s.initAuth(); err != nil {
		return fmt.Errorf("init auth: %w", err)
	}

	s.initRepo()
	s.initUseCases()
	seed.Run(context.Background(), s.usersRepository, s.hackathonUseCase, s.logger)
	s.initRoutes()
	s.initHTTPServer()

	return nil
}

func (s *Server) initDB() error {
	cfg, err := postgres.NewConfig()
	if err != nil {
		return err
	}

	pg, err := postgres.NewPostgres(cfg)
	if err != nil {
		return err
	}

	s.database = pg
	s.logger.Info("storage: postgres")
	return nil
}

func (s *Server) initAuth() error {
	tokens, err := auth.NewTokenService()
	if err != nil {
		return err
	}
	s.tokens = tokens
	return nil
}

func (s *Server) initRepo() {
	s.usersRepository = postgres_repo.NewUsersRepository(s.database)
	s.hackathonsRepository = postgres_repo.NewHackathonsRepository(s.database)
	s.participationRepository = postgres_repo.NewParticipationRepository(s.database)
}

func (s *Server) initUseCases() {
	s.authUseCase = auth_usecase.NewAuthUseCase(s.usersRepository, s.tokens)
	s.hackathonUseCase = hackathon_usecase.NewHackathonUseCase(s.hackathonsRepository)
	s.participationUseCase = participation_usecase.NewParticipationUseCase(
		s.hackathonsRepository,
		s.participationRepository,
		s.usersRepository,
	)
}

func (s *Server) initHTTPServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Run() {
	go func() {
		s.logger.Info("server started", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.database != nil {
		_ = s.database.Close()
	}

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("shutdown", "err", err)
		return
	}

	s.logger.Info("server stopped")
}
