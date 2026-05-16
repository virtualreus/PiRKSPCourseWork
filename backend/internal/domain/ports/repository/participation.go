package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type ParticipationRepository interface {
	CreateRegistration(ctx context.Context, reg entities.HackathonRegistration) (entities.HackathonRegistration, error)
	DeleteRegistration(ctx context.Context, hackathonID, userID uuid.UUID) error
	GetRegistration(ctx context.Context, hackathonID, userID uuid.UUID) (entities.HackathonRegistration, error)

	ListRegistrationsWithUsers(ctx context.Context, hackathonID uuid.UUID) ([]entities.RegistrationWithUser, error)

	CreateTeam(ctx context.Context, team entities.Team) (entities.Team, error)
	GetTeam(ctx context.Context, teamID uuid.UUID) (entities.Team, error)
	UpdateTeam(ctx context.Context, team entities.Team) (entities.Team, error)
	ListTeamsByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.Team, error)
	GetUserTeamInHackathon(ctx context.Context, hackathonID, userID uuid.UUID) (entities.Team, error)

	AddMember(ctx context.Context, teamID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role string) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]entities.TeamMember, error)
	CountMembers(ctx context.Context, teamID uuid.UUID) (int, error)

	UpsertSubmission(ctx context.Context, sub entities.Submission) (entities.Submission, error)
	GetSubmissionByTeam(ctx context.Context, teamID uuid.UUID) (entities.Submission, error)
	ListSubmissionsByHackathon(ctx context.Context, hackathonID uuid.UUID) ([]entities.SubmissionListRow, error)

	ListUserParticipations(ctx context.Context, userID uuid.UUID) ([]entities.UserParticipationRow, error)
}
