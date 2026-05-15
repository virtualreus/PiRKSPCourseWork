package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type HackathonListFilter struct {
	Statuses    []string
	Query       string
	OrganizerID *uuid.UUID
	ExcludeDraft bool
}

type HackathonsRepository interface {
	List(ctx context.Context, filter HackathonListFilter) ([]entities.Hackathon, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Hackathon, error)
	Create(ctx context.Context, h entities.Hackathon) (entities.Hackathon, error)
	Update(ctx context.Context, h entities.Hackathon) (entities.Hackathon, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (entities.Hackathon, error)

	ListTracks(ctx context.Context, hackathonID uuid.UUID) ([]entities.Track, error)
	GetTrack(ctx context.Context, id uuid.UUID) (entities.Track, error)
	CreateTrack(ctx context.Context, t entities.Track) (entities.Track, error)
	UpdateTrack(ctx context.Context, t entities.Track) (entities.Track, error)
	DeleteTrack(ctx context.Context, id uuid.UUID) error

	ListCasesByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.Case, error)
	ListCasesByTrack(ctx context.Context, trackID uuid.UUID) ([]entities.Case, error)
	GetCase(ctx context.Context, id uuid.UUID) (entities.Case, error)
	CreateCase(ctx context.Context, c entities.Case) (entities.Case, error)
	UpdateCase(ctx context.Context, c entities.Case) (entities.Case, error)
	DeleteCase(ctx context.Context, id uuid.UUID) error
	CountCasesByHackathon(ctx context.Context, hackathonID uuid.UUID) (int, error)
}
