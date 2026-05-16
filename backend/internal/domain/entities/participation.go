package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type HackathonRegistration struct {
	ID           uuid.UUID `db:"id"`
	HackathonID  uuid.UUID `db:"hackathon_id"`
	UserID       uuid.UUID `db:"user_id"`
	RegisteredAt time.Time `db:"registered_at"`
}

type Team struct {
	ID          uuid.UUID      `db:"id"`
	HackathonID uuid.UUID      `db:"hackathon_id"`
	Name        string         `db:"name"`
	CaptainID   uuid.UUID      `db:"captain_id"`
	TrackID     sql.NullString `db:"track_id"`
	CaseID      sql.NullString `db:"case_id"`
	CreatedAt   time.Time      `db:"created_at"`
}

type TeamMember struct {
	TeamID   uuid.UUID `db:"team_id"`
	UserID   uuid.UUID `db:"user_id"`
	TeamRole string    `db:"team_role"`
	FullName string    `db:"full_name"`
	JoinedAt time.Time `db:"joined_at"`
}

type Submission struct {
	ID          uuid.UUID      `db:"id"`
	TeamID      uuid.UUID      `db:"team_id"`
	HackathonID uuid.UUID      `db:"hackathon_id"`
	Title       sql.NullString `db:"title"`
	Summary     sql.NullString `db:"summary"`
	RepoURL     string         `db:"repo_url"`
	DemoURL     sql.NullString `db:"demo_url"`
	PitchURL    sql.NullString `db:"pitch_url"`
	VideoURL    sql.NullString `db:"video_url"`
	SubmittedAt sql.NullTime   `db:"submitted_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type SubmissionListRow struct {
	Submission
	TeamName   string         `db:"team_name"`
	CaseTitle  sql.NullString `db:"case_title"`
	TrackTitle sql.NullString `db:"track_title"`
}

type RegistrationWithUser struct {
	HackathonRegistration
	Email    string `db:"email"`
	FullName string `db:"full_name"`
}
