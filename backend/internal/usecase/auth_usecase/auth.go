package auth_usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/nikitatisenko/pirksp/internal/converters"
	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/auth"
	"github.com/nikitatisenko/pirksp/internal/usecase"
)

const bcryptCost = 12

type authUseCase struct {
	users      repository.UsersRepository
	tokens     *auth.TokenService
	converter  *converters.UsersConverter
}

func NewAuthUseCase(users repository.UsersRepository, tokens *auth.TokenService) usecase.AuthUseCase {
	return &authUseCase{
		users:     users,
		tokens:    tokens,
		converter: converters.NewUsersConverter(),
	}
}

func (u *authUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	if email == "" {
		return nil, errs.ErrEmptyEmail
	}
	if strings.TrimSpace(req.FullName) == "" {
		return nil, errs.ErrEmptyName
	}
	if len(req.Password) < 8 {
		return nil, errs.ErrWeakPassword
	}

	role := strings.TrimSpace(req.PlatformRole)
	if role == "" {
		role = "participant"
	}
	if role != "participant" && role != "organizer" {
		return nil, errs.ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	created, err := u.users.Create(ctx, entities.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     strings.TrimSpace(req.FullName),
		PlatformRole: role,
	})
	if err != nil {
		return nil, err
	}

	return u.authResponse(created)
}

func (u *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	if email == "" || req.Password == "" {
		return nil, errs.ErrInvalidCreds
	}

	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrInvalidCreds
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errs.ErrInvalidCreds
	}

	return u.authResponse(user)
}

func (u *authUseCase) GetMe(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	dtoUser := u.converter.ToDTO(user)
	return &dtoUser, nil
}

func (u *authUseCase) UpdateMe(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.User, error) {
	if req.FullName == nil || strings.TrimSpace(*req.FullName) == "" {
		return nil, errs.ErrEmptyName
	}

	name := strings.TrimSpace(*req.FullName)
	user, err := u.users.UpdateFullName(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	dtoUser := u.converter.ToDTO(user)
	return &dtoUser, nil
}

func (u *authUseCase) authResponse(user entities.User) (*dto.AuthResponse, error) {
	token, err := u.tokens.Issue(user.ID, user.PlatformRole)
	if err != nil {
		return nil, err
	}

	dtoUser := u.converter.ToDTO(user)
	return &dto.AuthResponse{
		AccessToken: token,
		User:        dtoUser,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
