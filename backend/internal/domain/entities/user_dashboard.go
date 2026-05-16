package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserParticipationRow struct {
	HackathonID            uuid.UUID      `db:"hackathon_id"`
	Title                  string         `db:"title"`
	ShortDescription       sql.NullString `db:"short_description"`
	Status                 string         `db:"status"`
	Format                 string         `db:"format"`
	RegistrationOpensAt    time.Time      `db:"registration_opens_at"`
	SubmissionDeadlineAt   time.Time      `db:"submission_deadline_at"`
	RegisteredAt           time.Time      `db:"registered_at"`
	TeamID                 sql.NullString `db:"team_id"`
	TeamName               sql.NullString `db:"team_name"`
	CaptainID              sql.NullString `db:"captain_id"`
	TeamRole               sql.NullString `db:"team_role"`
	TeamMemberCount        int            `db:"team_member_count"`
	TrackTitle             sql.NullString `db:"track_title"`
	CaseTitle              sql.NullString `db:"case_title"`
	HasCase                bool           `db:"has_case"`
	SubmissionID           sql.NullString `db:"submission_id"`
	SubmissionTitle        sql.NullString `db:"submission_title"`
	SubmissionRepoURL      sql.NullString `db:"submission_repo_url"`
	SubmissionSubmittedAt  sql.NullTime   `db:"submission_submitted_at"`
}
