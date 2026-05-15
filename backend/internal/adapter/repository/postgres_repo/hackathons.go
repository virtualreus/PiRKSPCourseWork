package postgres_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

var hackathonColumns = []string{
	"id", "organizer_id", "title", "short_description", "description", "format", "status",
	"max_team_size", "prizes_info",
	"registration_opens_at", "registration_closes_at", "event_starts_at", "event_ends_at",
	"submission_deadline_at", "created_at",
}

type hackathonsRepository struct {
	db *postgres.Postgres
}

func NewHackathonsRepository(db *postgres.Postgres) repository.HackathonsRepository {
	return &hackathonsRepository{db: db}
}

func (r *hackathonsRepository) List(ctx context.Context, filter repository.HackathonListFilter) ([]entities.Hackathon, error) {
	qb := r.db.Builder.Select(hackathonColumns...).From("hackathons").OrderBy("registration_opens_at DESC")

	if filter.OrganizerID != nil {
		qb = qb.Where(squirrel.Eq{"organizer_id": *filter.OrganizerID})
	}

	if filter.ExcludeDraft {
		qb = qb.Where(squirrel.NotEq{"status": "draft"})
	}

	if len(filter.Statuses) > 0 {
		qb = qb.Where(squirrel.Eq{"status": filter.Statuses})
	}

	if filter.Query != "" {
		qb = qb.Where(squirrel.ILike{"title": "%" + filter.Query + "%"})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("list hackathons sql: %w", err)
	}

	var items []entities.Hackathon
	if err := r.db.SqlxDB().SelectContext(ctx, &items, query, args...); err != nil {
		logger.FromContext(ctx).Error("List hackathons failed", "err", err)
		return nil, fmt.Errorf("list hackathons: %w", err)
	}

	return items, nil
}

func (r *hackathonsRepository) GetByID(ctx context.Context, id uuid.UUID) (entities.Hackathon, error) {
	qb := r.db.Builder.Select(hackathonColumns...).From("hackathons").Where(squirrel.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Hackathon{}, fmt.Errorf("get hackathon sql: %w", err)
	}

	var h entities.Hackathon
	if err := r.db.SqlxDB().GetContext(ctx, &h, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Hackathon{}, errs.ErrNotFound
		}
		logger.FromContext(ctx).Error("Get hackathon failed", "err", err)
		return entities.Hackathon{}, fmt.Errorf("get hackathon: %w", err)
	}

	return h, nil
}

func (r *hackathonsRepository) Create(ctx context.Context, h entities.Hackathon) (entities.Hackathon, error) {
	qb := r.db.Builder.Insert("hackathons").
		Columns(
			"organizer_id", "title", "short_description", "description", "format", "status",
			"max_team_size", "prizes_info",
			"registration_opens_at", "registration_closes_at", "event_starts_at", "event_ends_at",
			"submission_deadline_at",
		).
		Values(
			h.OrganizerID, h.Title, h.ShortDescription, h.Description, h.Format, h.Status,
			h.MaxTeamSize, h.PrizesInfo,
			h.RegistrationOpensAt, h.RegistrationClosesAt, h.EventStartsAt, h.EventEndsAt,
			h.SubmissionDeadlineAt,
		).
		Suffix("RETURNING " + joinColumns(hackathonColumns))

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Hackathon{}, fmt.Errorf("create hackathon sql: %w", err)
	}

	var created entities.Hackathon
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		logger.FromContext(ctx).Error("Create hackathon failed", "err", err)
		return entities.Hackathon{}, fmt.Errorf("create hackathon: %w", err)
	}

	return created, nil
}

func (r *hackathonsRepository) Update(ctx context.Context, h entities.Hackathon) (entities.Hackathon, error) {
	qb := r.db.Builder.Update("hackathons").
		Set("title", h.Title).
		Set("short_description", h.ShortDescription).
		Set("description", h.Description).
		Set("format", h.Format).
		Set("status", h.Status).
		Set("max_team_size", h.MaxTeamSize).
		Set("prizes_info", h.PrizesInfo).
		Set("registration_opens_at", h.RegistrationOpensAt).
		Set("registration_closes_at", h.RegistrationClosesAt).
		Set("event_starts_at", h.EventStartsAt).
		Set("event_ends_at", h.EventEndsAt).
		Set("submission_deadline_at", h.SubmissionDeadlineAt).
		Where(squirrel.Eq{"id": h.ID}).
		Suffix("RETURNING " + joinColumns(hackathonColumns))

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Hackathon{}, fmt.Errorf("update hackathon sql: %w", err)
	}

	var updated entities.Hackathon
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Hackathon{}, errs.ErrNotFound
		}
		logger.FromContext(ctx).Error("Update hackathon failed", "err", err)
		return entities.Hackathon{}, fmt.Errorf("update hackathon: %w", err)
	}

	return updated, nil
}

func (r *hackathonsRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (entities.Hackathon, error) {
	qb := r.db.Builder.Update("hackathons").
		Set("status", status).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING " + joinColumns(hackathonColumns))

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Hackathon{}, fmt.Errorf("update hackathon status sql: %w", err)
	}

	var updated entities.Hackathon
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Hackathon{}, errs.ErrNotFound
		}
		return entities.Hackathon{}, fmt.Errorf("update hackathon status: %w", err)
	}

	return updated, nil
}

func (r *hackathonsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	qb := r.db.Builder.Delete("hackathons").Where(squirrel.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("delete hackathon sql: %w", err)
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete hackathon: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *hackathonsRepository) ListTracks(ctx context.Context, hackathonID uuid.UUID) ([]entities.Track, error) {
	qb := r.db.Builder.
		Select("id", "hackathon_id", "title", "description", "created_at").
		From("tracks").
		Where(squirrel.Eq{"hackathon_id": hackathonID}).
		OrderBy("created_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("list tracks sql: %w", err)
	}

	var tracks []entities.Track
	if err := r.db.SqlxDB().SelectContext(ctx, &tracks, query, args...); err != nil {
		return nil, fmt.Errorf("list tracks: %w", err)
	}

	return tracks, nil
}

func (r *hackathonsRepository) GetTrack(ctx context.Context, id uuid.UUID) (entities.Track, error) {
	qb := r.db.Builder.
		Select("id", "hackathon_id", "title", "description", "created_at").
		From("tracks").
		Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Track{}, fmt.Errorf("get track sql: %w", err)
	}

	var t entities.Track
	if err := r.db.SqlxDB().GetContext(ctx, &t, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Track{}, errs.ErrNotFound
		}
		return entities.Track{}, fmt.Errorf("get track: %w", err)
	}

	return t, nil
}

func (r *hackathonsRepository) CreateTrack(ctx context.Context, t entities.Track) (entities.Track, error) {
	qb := r.db.Builder.Insert("tracks").
		Columns("hackathon_id", "title", "description").
		Values(t.HackathonID, t.Title, t.Description).
		Suffix("RETURNING id, hackathon_id, title, description, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Track{}, fmt.Errorf("create track sql: %w", err)
	}

	var created entities.Track
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		return entities.Track{}, fmt.Errorf("create track: %w", err)
	}

	return created, nil
}

func (r *hackathonsRepository) UpdateTrack(ctx context.Context, t entities.Track) (entities.Track, error) {
	qb := r.db.Builder.Update("tracks").
		Set("title", t.Title).
		Set("description", t.Description).
		Where(squirrel.Eq{"id": t.ID}).
		Suffix("RETURNING id, hackathon_id, title, description, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Track{}, fmt.Errorf("update track sql: %w", err)
	}

	var updated entities.Track
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Track{}, errs.ErrNotFound
		}
		return entities.Track{}, fmt.Errorf("update track: %w", err)
	}

	return updated, nil
}

func (r *hackathonsRepository) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	qb := r.db.Builder.Delete("tracks").Where(squirrel.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *hackathonsRepository) ListCasesByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.Case, error) {
	qb := r.db.Builder.
		Select("c.id", "c.track_id", "c.title", "c.description", "c.customer_name", "c.resources_url", "c.created_at").
		From("cases c").
		Join("tracks t ON t.id = c.track_id").
		Where(squirrel.Eq{"t.hackathon_id": hackathonID}).
		OrderBy("c.created_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var cases []entities.Case
	if err := r.db.SqlxDB().SelectContext(ctx, &cases, query, args...); err != nil {
		return nil, err
	}

	return cases, nil
}

func (r *hackathonsRepository) ListCasesByTrack(ctx context.Context, trackID uuid.UUID) ([]entities.Case, error) {
	qb := r.db.Builder.
		Select("id", "track_id", "title", "description", "customer_name", "resources_url", "created_at").
		From("cases").
		Where(squirrel.Eq{"track_id": trackID}).
		OrderBy("created_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var cases []entities.Case
	if err := r.db.SqlxDB().SelectContext(ctx, &cases, query, args...); err != nil {
		return nil, err
	}

	return cases, nil
}

func (r *hackathonsRepository) GetCase(ctx context.Context, id uuid.UUID) (entities.Case, error) {
	qb := r.db.Builder.
		Select("id", "track_id", "title", "description", "customer_name", "resources_url", "created_at").
		From("cases").
		Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Case{}, err
	}

	var c entities.Case
	if err := r.db.SqlxDB().GetContext(ctx, &c, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Case{}, errs.ErrNotFound
		}
		return entities.Case{}, err
	}

	return c, nil
}

func (r *hackathonsRepository) CreateCase(ctx context.Context, c entities.Case) (entities.Case, error) {
	qb := r.db.Builder.Insert("cases").
		Columns("track_id", "title", "description", "customer_name", "resources_url").
		Values(c.TrackID, c.Title, c.Description, c.CustomerName, c.ResourcesURL).
		Suffix("RETURNING id, track_id, title, description, customer_name, resources_url, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Case{}, err
	}

	var created entities.Case
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		return entities.Case{}, err
	}

	return created, nil
}

func (r *hackathonsRepository) UpdateCase(ctx context.Context, c entities.Case) (entities.Case, error) {
	qb := r.db.Builder.Update("cases").
		Set("title", c.Title).
		Set("description", c.Description).
		Set("customer_name", c.CustomerName).
		Set("resources_url", c.ResourcesURL).
		Where(squirrel.Eq{"id": c.ID}).
		Suffix("RETURNING id, track_id, title, description, customer_name, resources_url, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Case{}, err
	}

	var updated entities.Case
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Case{}, errs.ErrNotFound
		}
		return entities.Case{}, err
	}

	return updated, nil
}

func (r *hackathonsRepository) DeleteCase(ctx context.Context, id uuid.UUID) error {
	qb := r.db.Builder.Delete("cases").Where(squirrel.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *hackathonsRepository) CountCasesByHackathon(ctx context.Context, hackathonID uuid.UUID) (int, error) {
	qb := r.db.Builder.
		Select("COUNT(*)").
		From("cases c").
		Join("tracks t ON t.id = c.track_id").
		Where(squirrel.Eq{"t.hackathon_id": hackathonID})

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	if err := r.db.SqlxDB().GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

func joinColumns(cols []string) string {
	out := ""
	for i, c := range cols {
		if i > 0 {
			out += ", "
		}
		out += c
	}
	return out
}
