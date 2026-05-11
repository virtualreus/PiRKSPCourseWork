package converters

import (
	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type UsersConverter struct{}

func NewUsersConverter() *UsersConverter {
	return &UsersConverter{}
}

func (c *UsersConverter) ToDTO(user entities.User) dto.User {
	return dto.User{
		ID:           user.ID.String(),
		Email:        user.Email,
		FullName:     user.FullName,
		PlatformRole: user.PlatformRole,
		CreatedAt:    user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
