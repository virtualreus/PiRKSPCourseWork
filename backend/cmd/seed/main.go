package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/nikitatisenko/pirksp/internal/adapter/repository/postgres_repo"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/seed"
	"github.com/nikitatisenko/pirksp/internal/usecase/hackathon_usecase"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	cfg, err := postgres.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := postgres.MigrateDB(db); err != nil {
		log.Fatal(err)
	}

	users := postgres_repo.NewUsersRepository(db)
	hackathonsRepo := postgres_repo.NewHackathonsRepository(db)
	hackathons := hackathon_usecase.NewHackathonUseCase(hackathonsRepo)

	seed.Run(context.Background(), users, hackathons, logger.New())
}
