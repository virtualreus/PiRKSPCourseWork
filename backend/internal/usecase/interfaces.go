package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
)

type AuthUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.User, error)
}

type HackathonUseCase interface {
	ListPublic(ctx context.Context, statuses []string, query string) ([]dto.HackathonListItem, error)
	GetPublic(ctx context.Context, id uuid.UUID) (*dto.HackathonDetail, error)
	ListOrganizer(ctx context.Context, organizerID uuid.UUID) ([]dto.HackathonListItem, error)
	GetOrganizer(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error)
	Create(ctx context.Context, organizerID uuid.UUID, req dto.CreateHackathonRequest) (*dto.HackathonDetail, error)
	Update(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.UpdateHackathonRequest) (*dto.HackathonDetail, error)
	Delete(ctx context.Context, organizerID, hackathonID uuid.UUID) error
	Publish(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error)
	ListTracks(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.Track, error)
	CreateTrack(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error)
	UpdateTrack(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error)
	DeleteTrack(ctx context.Context, organizerID, trackID uuid.UUID) error
	ListCases(ctx context.Context, organizerID, trackID uuid.UUID) ([]dto.Case, error)
	CreateCase(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error)
	UpdateCase(ctx context.Context, organizerID, caseID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error)
	DeleteCase(ctx context.Context, organizerID, caseID uuid.UUID) error
}

type ParticipationUseCase interface {
	GetParticipation(ctx context.Context, userID, hackathonID uuid.UUID) (*dto.ParticipationStatus, error)
	RegisterForHackathon(ctx context.Context, userID, hackathonID uuid.UUID) (*dto.HackathonRegistration, error)
	UnregisterFromHackathon(ctx context.Context, userID, hackathonID uuid.UUID) error
	ListTeams(ctx context.Context, hackathonID uuid.UUID) ([]dto.Team, error)
	CreateTeam(ctx context.Context, userID, hackathonID uuid.UUID, req dto.CreateTeamRequest) (*dto.Team, error)
	GetTeam(ctx context.Context, teamID uuid.UUID) (*dto.Team, error)
	UpdateTeam(ctx context.Context, userID, teamID uuid.UUID, req dto.UpdateTeamRequest) (*dto.Team, error)
	JoinTeam(ctx context.Context, userID, teamID uuid.UUID, req dto.JoinTeamRequest) (*dto.Team, error)
	LeaveTeam(ctx context.Context, userID, teamID uuid.UUID) error
	UpdateTeamMemberRole(ctx context.Context, actorID, teamID, memberID uuid.UUID, req dto.UpdateTeamMemberRoleRequest) (*dto.TeamMember, error)
	GetTeamSubmission(ctx context.Context, userID, teamID uuid.UUID) (*dto.Submission, error)
	UpsertTeamSubmission(ctx context.Context, userID, teamID uuid.UUID, req dto.UpsertSubmissionRequest) (*dto.Submission, error)
	ListHackathonRegistrations(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.HackathonRegistrationWithUser, error)
	ListHackathonSubmissions(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.SubmissionWithTeam, error)
	GetUserDashboard(ctx context.Context, userID uuid.UUID) (*dto.UserDashboard, error)
}
