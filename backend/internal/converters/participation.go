package converters

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type ParticipationConverter struct{}

func NewParticipationConverter() *ParticipationConverter {
	return &ParticipationConverter{}
}

func (c *ParticipationConverter) ToRegistration(reg entities.HackathonRegistration) dto.HackathonRegistration {
	return dto.HackathonRegistration{
		ID:           reg.ID.String(),
		HackathonID:  reg.HackathonID.String(),
		UserID:       reg.UserID.String(),
		RegisteredAt: reg.RegisteredAt.UTC().Format(time.RFC3339),
	}
}

func (c *ParticipationConverter) ToRegistrationWithUser(row entities.RegistrationWithUser) dto.HackathonRegistrationWithUser {
	return dto.HackathonRegistrationWithUser{
		HackathonRegistration: c.ToRegistration(row.HackathonRegistration),
		User: dto.User{
			ID:           row.UserID.String(),
			Email:        row.Email,
			FullName:     row.FullName,
			PlatformRole: "participant",
		},
	}
}

func (c *ParticipationConverter) ToTeam(team entities.Team, members []entities.TeamMember) dto.Team {
	out := dto.Team{
		ID:          team.ID.String(),
		HackathonID: team.HackathonID.String(),
		Name:        team.Name,
		CaptainID:   team.CaptainID.String(),
		Members:     make([]dto.TeamMember, 0, len(members)),
	}
	if team.TrackID.Valid {
		s := team.TrackID.String
		out.TrackID = &s
	}
	if team.CaseID.Valid {
		s := team.CaseID.String
		out.CaseID = &s
	}
	for _, m := range members {
		out.Members = append(out.Members, dto.TeamMember{
			UserID:   m.UserID.String(),
			FullName: m.FullName,
			TeamRole: m.TeamRole,
		})
	}
	return out
}

func (c *ParticipationConverter) ToSubmission(sub entities.Submission) dto.Submission {
	out := dto.Submission{
		ID:          sub.ID.String(),
		TeamID:      sub.TeamID.String(),
		HackathonID: sub.HackathonID.String(),
		RepoURL:     sub.RepoURL,
		UpdatedAt:   sub.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if sub.Title.Valid {
		out.Title = sub.Title.String
	}
	if sub.Summary.Valid {
		out.Summary = sub.Summary.String
	}
	if sub.DemoURL.Valid {
		out.DemoURL = sub.DemoURL.String
	}
	if sub.PitchURL.Valid {
		out.PitchURL = sub.PitchURL.String
	}
	if sub.VideoURL.Valid {
		out.VideoURL = sub.VideoURL.String
	}
	if sub.SubmittedAt.Valid {
		s := sub.SubmittedAt.Time.UTC().Format(time.RFC3339)
		out.SubmittedAt = &s
	}
	return out
}

func (c *ParticipationConverter) ToSubmissionWithTeam(row entities.SubmissionListRow) dto.SubmissionWithTeam {
	out := dto.SubmissionWithTeam{
		Submission: c.ToSubmission(row.Submission),
		TeamName:   row.TeamName,
	}
	if row.CaseTitle.Valid {
		out.CaseTitle = row.CaseTitle.String
	}
	if row.TrackTitle.Valid {
		out.TrackTitle = row.TrackTitle.String
	}
	return out
}

func ParseOptionalUUID(s *string) (sql.NullString, error) {
	if s == nil || *s == "" {
		return sql.NullString{}, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return sql.NullString{}, err
	}
	return sql.NullString{String: id.String(), Valid: true}, nil
}

func NullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}
