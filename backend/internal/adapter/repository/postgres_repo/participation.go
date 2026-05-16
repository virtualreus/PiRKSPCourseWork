package postgres_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

type participationRepository struct {
	db *postgres.Postgres
}

func NewParticipationRepository(db *postgres.Postgres) repository.ParticipationRepository {
	return &participationRepository{db: db}
}

func (r *participationRepository) CreateRegistration(ctx context.Context, reg entities.HackathonRegistration) (entities.HackathonRegistration, error) {
	qb := r.db.Builder.Insert("hackathon_registrations").
		Columns("hackathon_id", "user_id").
		Values(reg.HackathonID, reg.UserID).
		Suffix("RETURNING id, hackathon_id, user_id, registered_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.HackathonRegistration{}, err
	}

	var created entities.HackathonRegistration
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.HackathonRegistration{}, errs.ErrAlreadyRegistered
		}
		return entities.HackathonRegistration{}, err
	}

	return created, nil
}

func (r *participationRepository) DeleteRegistration(ctx context.Context, hackathonID, userID uuid.UUID) error {
	qb := r.db.Builder.Delete("hackathon_registrations").
		Where(squirrel.Eq{"hackathon_id": hackathonID, "user_id": userID})

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotRegistered
	}

	return nil
}

func (r *participationRepository) GetRegistration(ctx context.Context, hackathonID, userID uuid.UUID) (entities.HackathonRegistration, error) {
	qb := r.db.Builder.
		Select("id", "hackathon_id", "user_id", "registered_at").
		From("hackathon_registrations").
		Where(squirrel.Eq{"hackathon_id": hackathonID, "user_id": userID})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.HackathonRegistration{}, err
	}

	var reg entities.HackathonRegistration
	if err := r.db.SqlxDB().GetContext(ctx, &reg, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.HackathonRegistration{}, errs.ErrNotRegistered
		}
		return entities.HackathonRegistration{}, err
	}

	return reg, nil
}

func (r *participationRepository) ListRegistrationsWithUsers(ctx context.Context, hackathonID uuid.UUID) ([]entities.RegistrationWithUser, error) {
	qb := r.db.Builder.
		Select(
			"r.id", "r.hackathon_id", "r.user_id", "r.registered_at",
			"u.email", "u.full_name",
		).
		From("hackathon_registrations r").
		Join("users u ON u.id = r.user_id").
		Where(squirrel.Eq{"r.hackathon_id": hackathonID}).
		OrderBy("r.registered_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entities.RegistrationWithUser
	if err := r.db.SqlxDB().SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *participationRepository) CreateTeam(ctx context.Context, team entities.Team) (entities.Team, error) {
	qb := r.db.Builder.Insert("teams").
		Columns("hackathon_id", "name", "captain_id", "track_id", "case_id").
		Values(team.HackathonID, team.Name, team.CaptainID, team.TrackID, team.CaseID).
		Suffix("RETURNING id, hackathon_id, name, captain_id, track_id, case_id, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Team{}, err
	}

	var created entities.Team
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&created); err != nil {
		return entities.Team{}, err
	}

	return created, nil
}

func (r *participationRepository) GetTeam(ctx context.Context, teamID uuid.UUID) (entities.Team, error) {
	qb := r.db.Builder.
		Select("id", "hackathon_id", "name", "captain_id", "track_id", "case_id", "created_at").
		From("teams").
		Where(squirrel.Eq{"id": teamID})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Team{}, err
	}

	var team entities.Team
	if err := r.db.SqlxDB().GetContext(ctx, &team, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Team{}, errs.ErrNotFound
		}
		return entities.Team{}, err
	}

	return team, nil
}

func (r *participationRepository) UpdateTeam(ctx context.Context, team entities.Team) (entities.Team, error) {
	qb := r.db.Builder.Update("teams").
		Set("name", team.Name).
		Set("track_id", team.TrackID).
		Set("case_id", team.CaseID).
		Where(squirrel.Eq{"id": team.ID}).
		Suffix("RETURNING id, hackathon_id, name, captain_id, track_id, case_id, created_at")

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Team{}, err
	}

	var updated entities.Team
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Team{}, errs.ErrNotFound
		}
		return entities.Team{}, err
	}

	return updated, nil
}

func (r *participationRepository) ListTeamsByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.Team, error) {
	qb := r.db.Builder.
		Select("id", "hackathon_id", "name", "captain_id", "track_id", "case_id", "created_at").
		From("teams").
		Where(squirrel.Eq{"hackathon_id": hackathonID}).
		OrderBy("created_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var teams []entities.Team
	if err := r.db.SqlxDB().SelectContext(ctx, &teams, query, args...); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *participationRepository) GetUserTeamInHackathon(ctx context.Context, hackathonID, userID uuid.UUID) (entities.Team, error) {
	qb := r.db.Builder.
		Select("t.id", "t.hackathon_id", "t.name", "t.captain_id", "t.track_id", "t.case_id", "t.created_at").
		From("teams t").
		Join("team_members m ON m.team_id = t.id").
		Where(squirrel.Eq{"t.hackathon_id": hackathonID, "m.user_id": userID})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Team{}, err
	}

	var team entities.Team
	if err := r.db.SqlxDB().GetContext(ctx, &team, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Team{}, errs.ErrNotInTeam
		}
		return entities.Team{}, err
	}

	return team, nil
}

func (r *participationRepository) AddMember(ctx context.Context, teamID, userID uuid.UUID, role string) error {
	qb := r.db.Builder.Insert("team_members").
		Columns("team_id", "user_id", "team_role").
		Values(teamID, userID, role)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.SqlxDB().ExecContext(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrAlreadyInTeam
		}
		return err
	}

	return nil
}

func (r *participationRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	qb := r.db.Builder.Delete("team_members").
		Where(squirrel.Eq{"team_id": teamID, "user_id": userID})

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotInTeam
	}

	return nil
}

func (r *participationRepository) UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role string) error {
	qb := r.db.Builder.Update("team_members").
		Set("team_role", role).
		Where(squirrel.Eq{"team_id": teamID, "user_id": userID})

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.SqlxDB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotInTeam
	}

	return nil
}

func (r *participationRepository) ListMembers(ctx context.Context, teamID uuid.UUID) ([]entities.TeamMember, error) {
	qb := r.db.Builder.
		Select("m.team_id", "m.user_id", "m.team_role", "u.full_name", "m.joined_at").
		From("team_members m").
		Join("users u ON u.id = m.user_id").
		Where(squirrel.Eq{"m.team_id": teamID}).
		OrderBy("m.joined_at ASC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var members []entities.TeamMember
	if err := r.db.SqlxDB().SelectContext(ctx, &members, query, args...); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *participationRepository) CountMembers(ctx context.Context, teamID uuid.UUID) (int, error) {
	qb := r.db.Builder.Select("COUNT(*)").From("team_members").Where(squirrel.Eq{"team_id": teamID})
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

func (r *participationRepository) UpsertSubmission(ctx context.Context, sub entities.Submission) (entities.Submission, error) {
	now := time.Now().UTC()
	qb := r.db.Builder.Insert("submissions").
		Columns(
			"team_id", "hackathon_id", "title", "summary", "repo_url",
			"demo_url", "pitch_url", "video_url", "submitted_at", "updated_at",
		).
		Values(
			sub.TeamID, sub.HackathonID, sub.Title, sub.Summary, sub.RepoURL,
			sub.DemoURL, sub.PitchURL, sub.VideoURL, sub.SubmittedAt, now,
		).
		Suffix(`
			ON CONFLICT (team_id) DO UPDATE SET
				title = EXCLUDED.title,
				summary = EXCLUDED.summary,
				repo_url = EXCLUDED.repo_url,
				demo_url = EXCLUDED.demo_url,
				pitch_url = EXCLUDED.pitch_url,
				video_url = EXCLUDED.video_url,
				submitted_at = COALESCE(submissions.submitted_at, EXCLUDED.submitted_at),
				updated_at = EXCLUDED.updated_at
			RETURNING id, team_id, hackathon_id, title, summary, repo_url, demo_url, pitch_url, video_url, submitted_at, updated_at
		`)

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Submission{}, err
	}

	var saved entities.Submission
	if err := r.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&saved); err != nil {
		logger.FromContext(ctx).Error("upsert submission", "err", err)
		return entities.Submission{}, fmt.Errorf("upsert submission: %w", err)
	}

	return saved, nil
}

func (r *participationRepository) GetSubmissionByTeam(ctx context.Context, teamID uuid.UUID) (entities.Submission, error) {
	qb := r.db.Builder.
		Select("id", "team_id", "hackathon_id", "title", "summary", "repo_url", "demo_url", "pitch_url", "video_url", "submitted_at", "updated_at").
		From("submissions").
		Where(squirrel.Eq{"team_id": teamID})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Submission{}, err
	}

	var sub entities.Submission
	if err := r.db.SqlxDB().GetContext(ctx, &sub, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Submission{}, errs.ErrNotFound
		}
		return entities.Submission{}, err
	}

	return sub, nil
}

func (r *participationRepository) ListSubmissionsByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.SubmissionListRow, error) {
	qb := r.db.Builder.
		Select(
			"s.id", "s.team_id", "s.hackathon_id", "s.title", "s.summary", "s.repo_url",
			"s.demo_url", "s.pitch_url", "s.video_url", "s.submitted_at", "s.updated_at",
			"t.name AS team_name", "c.title AS case_title", "tr.title AS track_title",
		).
		From("submissions s").
		Join("teams t ON t.id = s.team_id").
		LeftJoin("cases c ON c.id = t.case_id").
		LeftJoin("tracks tr ON tr.id = t.track_id").
		Where(squirrel.Eq{"s.hackathon_id": hackathonID}).
		OrderBy("s.updated_at DESC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entities.SubmissionListRow
	if err := r.db.SqlxDB().SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *participationRepository) ListUserParticipations(ctx context.Context, userID uuid.UUID) ([]entities.UserParticipationRow, error) {
	qb := r.db.Builder.
		Select(
			"h.id AS hackathon_id",
			"h.title",
			"h.short_description",
			"h.status",
			"h.format",
			"h.registration_opens_at",
			"h.submission_deadline_at",
			"r.registered_at",
			"t.id AS team_id",
			"t.name AS team_name",
			"t.captain_id AS captain_id",
			"tm.team_role",
			"COALESCE((SELECT COUNT(*)::int FROM team_members WHERE team_id = t.id), 0) AS team_member_count",
			"tr.title AS track_title",
			"c.title AS case_title",
			"(t.id IS NOT NULL AND t.case_id IS NOT NULL) AS has_case",
			"s.id AS submission_id",
			"s.title AS submission_title",
			"s.repo_url AS submission_repo_url",
			"s.submitted_at AS submission_submitted_at",
		).
		From("hackathon_registrations r").
		Join("hackathons h ON h.id = r.hackathon_id").
		LeftJoin("team_members tm ON tm.user_id = r.user_id").
		LeftJoin("teams t ON t.id = tm.team_id AND t.hackathon_id = r.hackathon_id").
		LeftJoin("tracks tr ON tr.id = t.track_id").
		LeftJoin("cases c ON c.id = t.case_id").
		LeftJoin("submissions s ON s.team_id = t.id").
		Where(squirrel.Eq{"r.user_id": userID}).
		OrderByClause(
			"CASE WHEN h.status IN ('registration', 'running') THEN 0 ELSE 1 END, h.submission_deadline_at DESC",
		)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entities.UserParticipationRow
	if err := r.db.SqlxDB().SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	return rows, nil
}
