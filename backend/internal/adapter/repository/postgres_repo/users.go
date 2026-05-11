package postgres_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

type usersRepository struct {
	db *postgres.Postgres
}

func NewUsersRepository(db *postgres.Postgres) repository.UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	qb := r.db.Builder.Insert("users").
		Columns("email", "password_hash", "full_name", "platform_role").
		Values(user.Email, user.PasswordHash, user.FullName, user.PlatformRole).
		Suffix("RETURNING id, email, password_hash, full_name, platform_role, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("create user sql: %w", err)
	}

	var created entities.User
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.User{}, errs.ErrEmailTaken
		}
		logger.FromContext(ctx).Error("Create user failed", "err", err)
		return entities.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (r *usersRepository) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	return r.getOne(ctx, squirrel.Eq{"email": email})
}

func (r *usersRepository) GetByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
	return r.getOne(ctx, squirrel.Eq{"id": id})
}

func (r *usersRepository) UpdateFullName(ctx context.Context, id uuid.UUID, fullName string) (entities.User, error) {
	qb := r.db.Builder.Update("users").
		Set("full_name", fullName).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, email, password_hash, full_name, platform_role, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("update user sql: %w", err)
	}

	var updated entities.User
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, errs.ErrNotFound
		}
		logger.FromContext(ctx).Error("Update user failed", "err", err)
		return entities.User{}, fmt.Errorf("update user: %w", err)
	}

	return updated, nil
}

func (r *usersRepository) getOne(ctx context.Context, where squirrel.Eq) (entities.User, error) {
	qb := r.db.Builder.
		Select("id", "email", "password_hash", "full_name", "platform_role", "created_at").
		From("users").
		Where(where)

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("get user sql: %w", err)
	}

	var user entities.User
	if err := r.db.SqlxDB().GetContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, errs.ErrNotFound
		}
		logger.FromContext(ctx).Error("Get user failed", "err", err)
		return entities.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
