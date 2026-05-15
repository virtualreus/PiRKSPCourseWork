package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Hackathon struct {
	ID                     uuid.UUID      `db:"id"`
	OrganizerID            uuid.UUID      `db:"organizer_id"`
	Title                  string         `db:"title"`
	ShortDescription       sql.NullString `db:"short_description"`
	Description            string         `db:"description"`
	Format                 string         `db:"format"`
	Status                 string         `db:"status"`
	MaxTeamSize            int            `db:"max_team_size"`
	PrizesInfo             sql.NullString `db:"prizes_info"`
	RegistrationOpensAt    time.Time      `db:"registration_opens_at"`
	RegistrationClosesAt   time.Time      `db:"registration_closes_at"`
	EventStartsAt          time.Time      `db:"event_starts_at"`
	EventEndsAt            time.Time      `db:"event_ends_at"`
	SubmissionDeadlineAt   time.Time      `db:"submission_deadline_at"`
	CreatedAt              time.Time      `db:"created_at"`
}

type Track struct {
	ID          uuid.UUID      `db:"id"`
	HackathonID uuid.UUID      `db:"hackathon_id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
}

type Case struct {
	ID           uuid.UUID      `db:"id"`
	TrackID      uuid.UUID      `db:"track_id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	CustomerName sql.NullString `db:"customer_name"`
	ResourcesURL sql.NullString `db:"resources_url"`
	CreatedAt    time.Time      `db:"created_at"`
}
