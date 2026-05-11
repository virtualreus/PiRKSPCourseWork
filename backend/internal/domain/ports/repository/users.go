package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type UsersRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.User, error)
	UpdateFullName(ctx context.Context, id uuid.UUID, fullName string) (entities.User, error)
}
