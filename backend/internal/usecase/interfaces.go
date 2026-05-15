package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
)

type AuthUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.User, error)
}

type HackathonUseCase interface {
	ListPublic(ctx context.Context, statuses []string, query string) ([]dto.HackathonListItem, error)
	GetPublic(ctx context.Context, id uuid.UUID) (*dto.HackathonDetail, error)
	ListOrganizer(ctx context.Context, organizerID uuid.UUID) ([]dto.HackathonListItem, error)
	GetOrganizer(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error)
	Create(ctx context.Context, organizerID uuid.UUID, req dto.CreateHackathonRequest) (*dto.HackathonDetail, error)
	Update(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.UpdateHackathonRequest) (*dto.HackathonDetail, error)
	Delete(ctx context.Context, organizerID, hackathonID uuid.UUID) error
	Publish(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error)
	ListTracks(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.Track, error)
	CreateTrack(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error)
	UpdateTrack(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error)
	DeleteTrack(ctx context.Context, organizerID, trackID uuid.UUID) error
	ListCases(ctx context.Context, organizerID, trackID uuid.UUID) ([]dto.Case, error)
	CreateCase(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error)
	UpdateCase(ctx context.Context, organizerID, caseID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error)
	DeleteCase(ctx context.Context, organizerID, caseID uuid.UUID) error
}
