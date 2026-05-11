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
