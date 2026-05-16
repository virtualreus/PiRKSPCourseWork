package participation_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/converters"
	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
)

var validTeamRoles = map[string]struct{}{
	"team_lead": {}, "developer": {}, "designer": {},
	"data_scientist": {}, "devops_qa": {}, "other": {},
}

type participationUseCase struct {
	hackathons     repository.HackathonsRepository
	participation  repository.ParticipationRepository
	users          repository.UsersRepository
	converter      *converters.ParticipationConverter
}

func NewParticipationUseCase(
	hackathons repository.HackathonsRepository,
	participation repository.ParticipationRepository,
	users repository.UsersRepository,
) usecase.ParticipationUseCase {
	return &participationUseCase{
		hackathons:    hackathons,
		participation: participation,
		users:         users,
		converter:     converters.NewParticipationConverter(),
	}
}

func (u *participationUseCase) GetParticipation(ctx context.Context, userID, hackathonID uuid.UUID) (*dto.ParticipationStatus, error) {
	h, err := u.hackathons.GetByID(ctx, hackathonID)
	if err != nil {
		return nil, err
	}
	if h.Status == "draft" {
		return nil, errs.ErrNotFound
	}

	now := time.Now().UTC()
	reg, regErr := u.participation.GetRegistration(ctx, hackathonID, userID)
	isRegistered := regErr == nil

	var regDTO *dto.HackathonRegistration
	if isRegistered {
		dtoReg := u.converter.ToRegistration(reg)
		regDTO = &dtoReg
	}

	var teamDTO *dto.Team
	hasSubmission := false
	if team, err := u.participation.GetUserTeamInHackathon(ctx, hackathonID, userID); err == nil {
		loaded, err := u.loadTeamDTO(ctx, team)
		if err != nil {
			return nil, err
		}
		teamDTO = &loaded
		if _, err := u.participation.GetSubmissionByTeam(ctx, team.ID); err == nil {
			hasSubmission = true
		}
	}

	user, _ := u.users.GetByID(ctx, userID)
	canRegister := user.PlatformRole == "participant" && !isRegistered &&
		u.canRegisterWindow(h, now) && h.Status == "registration"

	inTeam := teamDTO != nil
	canCreateTeam := isRegistered && !inTeam && (h.Status == "registration" || h.Status == "running")

	submitBlockReason := ""
	if !inTeam {
		submitBlockReason = "no_team"
	} else if teamDTO.CaseID == nil {
		submitBlockReason = "no_case"
	} else if !now.Before(h.SubmissionDeadlineAt) {
		submitBlockReason = "deadline_passed"
	} else if h.Status == "finished" {
		submitBlockReason = "hackathon_finished"
	} else if h.Status != "registration" && h.Status != "running" {
		submitBlockReason = "hackathon_not_active"
	}
	canSubmit := submitBlockReason == ""

	return &dto.ParticipationStatus{
		HackathonID:          hackathonID.String(),
		IsRegistered:         isRegistered,
		Registration:         regDTO,
		Team:                 teamDTO,
		HasSubmission:        hasSubmission,
		CanRegister:          canRegister,
		CanCreateTeam:        canCreateTeam,
		CanSubmit:            canSubmit,
		SubmitBlockReason:    submitBlockReason,
		HackathonStatus:      h.Status,
		SubmissionDeadlineAt: h.SubmissionDeadlineAt.UTC().Format(time.RFC3339),
	}, nil
}

func (u *participationUseCase) RegisterForHackathon(ctx context.Context, userID, hackathonID uuid.UUID) (*dto.HackathonRegistration, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.PlatformRole != "participant" {
		return nil, errs.ErrParticipantsOnly
	}

	h, err := u.hackathons.GetByID(ctx, hackathonID)
	if err != nil {
		return nil, err
	}
	if h.Status != "registration" || !u.canRegisterWindow(h, time.Now().UTC()) {
		return nil, errs.ErrRegistrationClosed
	}

	if _, err := u.participation.GetUserTeamInHackathon(ctx, hackathonID, userID); err == nil {
		return nil, errs.ErrAlreadyInTeam
	}

	created, err := u.participation.CreateRegistration(ctx, entities.HackathonRegistration{
		HackathonID: hackathonID,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}

	dtoReg := u.converter.ToRegistration(created)
	return &dtoReg, nil
}

func (u *participationUseCase) UnregisterFromHackathon(ctx context.Context, userID, hackathonID uuid.UUID) error {
	if _, err := u.participation.GetUserTeamInHackathon(ctx, hackathonID, userID); err == nil {
		return errs.ErrInTeamBlockUnregister
	}

	return u.participation.DeleteRegistration(ctx, hackathonID, userID)
}

func (u *participationUseCase) ListTeams(ctx context.Context, hackathonID uuid.UUID) ([]dto.Team, error) {
	if _, err := u.hackathons.GetByID(ctx, hackathonID); err != nil {
		return nil, err
	}

	teams, err := u.participation.ListTeamsByHackathon(ctx, hackathonID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.Team, 0, len(teams))
	for _, t := range teams {
		dtoTeam, err := u.loadTeamDTO(ctx, t)
		if err != nil {
			return nil, err
		}
		out = append(out, dtoTeam)
	}

	return out, nil
}

func (u *participationUseCase) CreateTeam(ctx context.Context, userID, hackathonID uuid.UUID, req dto.CreateTeamRequest) (*dto.Team, error) {
	h, err := u.ensureRegistered(ctx, userID, hackathonID)
	if err != nil {
		return nil, err
	}

	if _, err := u.participation.GetUserTeamInHackathon(ctx, hackathonID, userID); err == nil {
		return nil, errs.ErrAlreadyInTeam
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errs.ErrValidation
	}

	role := normalizeTeamRole(req.TeamRole, "team_lead")
	trackID, err := converters.ParseOptionalUUID(req.TrackID)
	if err != nil {
		return nil, errs.ErrValidation
	}
	caseID, err := converters.ParseOptionalUUID(req.CaseID)
	if err != nil {
		return nil, errs.ErrValidation
	}

	team, err := u.participation.CreateTeam(ctx, entities.Team{
		HackathonID: hackathonID,
		Name:        name,
		CaptainID:   userID,
		TrackID:     trackID,
		CaseID:      caseID,
	})
	if err != nil {
		return nil, err
	}

	if err := u.participation.AddMember(ctx, team.ID, userID, role); err != nil {
		return nil, err
	}

	_ = h
	loaded, err := u.loadTeamDTO(ctx, team)
	if err != nil {
		return nil, err
	}

	return &loaded, nil
}

func (u *participationUseCase) GetTeam(ctx context.Context, teamID uuid.UUID) (*dto.Team, error) {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	loaded, err := u.loadTeamDTO(ctx, team)
	if err != nil {
		return nil, err
	}

	return &loaded, nil
}

func (u *participationUseCase) UpdateTeam(ctx context.Context, userID, teamID uuid.UUID, req dto.UpdateTeamRequest) (*dto.Team, error) {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team.CaptainID != userID {
		return nil, errs.ErrNotCaptain
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errs.ErrValidation
		}
		team.Name = name
	}
	if req.TrackID != nil {
		trackID, err := converters.ParseOptionalUUID(req.TrackID)
		if err != nil {
			return nil, errs.ErrValidation
		}
		team.TrackID = trackID
	}
	if req.CaseID != nil {
		caseID, err := converters.ParseOptionalUUID(req.CaseID)
		if err != nil {
			return nil, errs.ErrValidation
		}
		team.CaseID = caseID
	}

	updated, err := u.participation.UpdateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	loaded, err := u.loadTeamDTO(ctx, updated)
	if err != nil {
		return nil, err
	}

	return &loaded, nil
}

func (u *participationUseCase) JoinTeam(ctx context.Context, userID, teamID uuid.UUID, req dto.JoinTeamRequest) (*dto.Team, error) {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	if _, err := u.ensureRegistered(ctx, userID, team.HackathonID); err != nil {
		return nil, err
	}

	if _, err := u.participation.GetUserTeamInHackathon(ctx, team.HackathonID, userID); err == nil {
		return nil, errs.ErrAlreadyInTeam
	}

	h, err := u.hackathons.GetByID(ctx, team.HackathonID)
	if err != nil {
		return nil, err
	}

	count, err := u.participation.CountMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if count >= h.MaxTeamSize {
		return nil, errs.ErrTeamFull
	}

	role := normalizeTeamRole(req.TeamRole, "developer")
	if err := u.participation.AddMember(ctx, teamID, userID, role); err != nil {
		return nil, err
	}

	loaded, err := u.loadTeamDTO(ctx, team)
	if err != nil {
		return nil, err
	}

	return &loaded, nil
}

func (u *participationUseCase) LeaveTeam(ctx context.Context, userID, teamID uuid.UUID) error {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return err
	}
	if team.CaptainID == userID {
		return errs.ErrCaptainCannotLeave
	}

	return u.participation.RemoveMember(ctx, teamID, userID)
}

func (u *participationUseCase) UpdateTeamMemberRole(ctx context.Context, actorID, teamID, memberID uuid.UUID, req dto.UpdateTeamMemberRoleRequest) (*dto.TeamMember, error) {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if actorID != memberID && team.CaptainID != actorID {
		return nil, errs.ErrForbidden
	}

	role := normalizeTeamRole(req.TeamRole, "")
	if role == "" {
		return nil, errs.ErrValidation
	}

	if err := u.participation.UpdateMemberRole(ctx, teamID, memberID, role); err != nil {
		return nil, err
	}

	members, err := u.participation.ListMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		if m.UserID == memberID {
			return &dto.TeamMember{
				UserID:   m.UserID.String(),
				FullName: m.FullName,
				TeamRole: m.TeamRole,
			}, nil
		}
	}

	return nil, errs.ErrNotInTeam
}

func (u *participationUseCase) GetTeamSubmission(ctx context.Context, userID, teamID uuid.UUID) (*dto.Submission, error) {
	if err := u.ensureTeamMember(ctx, userID, teamID); err != nil {
		return nil, err
	}

	sub, err := u.participation.GetSubmissionByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	dtoSub := u.converter.ToSubmission(sub)
	return &dtoSub, nil
}

func (u *participationUseCase) UpsertTeamSubmission(ctx context.Context, userID, teamID uuid.UUID, req dto.UpsertSubmissionRequest) (*dto.Submission, error) {
	team, err := u.participation.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if err := u.ensureTeamMember(ctx, userID, teamID); err != nil {
		return nil, err
	}

	h, err := u.hackathons.GetByID(ctx, team.HackathonID)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().After(h.SubmissionDeadlineAt) {
		return nil, errs.ErrDeadlinePassed
	}
	if !team.CaseID.Valid {
		return nil, errs.ErrCaseRequired
	}

	repoURL := strings.TrimSpace(req.RepoURL)
	if repoURL == "" {
		return nil, errs.ErrValidation
	}

	submittedAt := sql.NullTime{}
	if _, err := u.participation.GetSubmissionByTeam(ctx, teamID); err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return nil, err
		}
		submittedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	}

	saved, err := u.participation.UpsertSubmission(ctx, entities.Submission{
		TeamID:      teamID,
		HackathonID: team.HackathonID,
		Title:       converters.NullString(strings.TrimSpace(req.Title)),
		Summary:     converters.NullString(strings.TrimSpace(req.Summary)),
		RepoURL:     repoURL,
		DemoURL:     converters.NullString(strings.TrimSpace(req.DemoURL)),
		PitchURL:    converters.NullString(strings.TrimSpace(req.PitchURL)),
		VideoURL:    converters.NullString(strings.TrimSpace(req.VideoURL)),
		SubmittedAt: submittedAt,
	})
	if err != nil {
		return nil, err
	}

	dtoSub := u.converter.ToSubmission(saved)
	return &dtoSub, nil
}

func (u *participationUseCase) ListHackathonRegistrations(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.HackathonRegistrationWithUser, error) {
	if err := u.ensureHackathonOwner(ctx, organizerID, hackathonID); err != nil {
		return nil, err
	}

	rows, err := u.participation.ListRegistrationsWithUsers(ctx, hackathonID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.HackathonRegistrationWithUser, 0, len(rows))
	for _, row := range rows {
		out = append(out, u.converter.ToRegistrationWithUser(row))
	}

	return out, nil
}

func (u *participationUseCase) ListHackathonSubmissions(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.SubmissionWithTeam, error) {
	if err := u.ensureHackathonOwner(ctx, organizerID, hackathonID); err != nil {
		return nil, err
	}

	rows, err := u.participation.ListSubmissionsByHackathon(ctx, hackathonID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.SubmissionWithTeam, 0, len(rows))
	for _, row := range rows {
		out = append(out, u.converter.ToSubmissionWithTeam(row))
	}

	return out, nil
}

func (u *participationUseCase) loadTeamDTO(ctx context.Context, team entities.Team) (dto.Team, error) {
	members, err := u.participation.ListMembers(ctx, team.ID)
	if err != nil {
		return dto.Team{}, err
	}
	return u.converter.ToTeam(team, members), nil
}

func (u *participationUseCase) ensureRegistered(ctx context.Context, userID, hackathonID uuid.UUID) (entities.Hackathon, error) {
	h, err := u.hackathons.GetByID(ctx, hackathonID)
	if err != nil {
		return entities.Hackathon{}, err
	}
	if _, err := u.participation.GetRegistration(ctx, hackathonID, userID); err != nil {
		return entities.Hackathon{}, err
	}
	return h, nil
}

func (u *participationUseCase) ensureTeamMember(ctx context.Context, userID, teamID uuid.UUID) error {
	members, err := u.participation.ListMembers(ctx, teamID)
	if err != nil {
		return err
	}
	for _, m := range members {
		if m.UserID == userID {
			return nil
		}
	}
	return errs.ErrNotInTeam
}

func (u *participationUseCase) ensureHackathonOwner(ctx context.Context, organizerID, hackathonID uuid.UUID) error {
	h, err := u.hackathons.GetByID(ctx, hackathonID)
	if err != nil {
		return err
	}
	if h.OrganizerID != organizerID {
		return errs.ErrNotOwner
	}
	return nil
}

func (u *participationUseCase) canRegisterWindow(h entities.Hackathon, now time.Time) bool {
	return !now.Before(h.RegistrationOpensAt) && now.Before(h.RegistrationClosesAt)
}

func (u *participationUseCase) GetUserDashboard(ctx context.Context, userID uuid.UUID) (*dto.UserDashboard, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userConv := converters.NewUsersConverter()
	hackConv := converters.NewHackathonsConverter()
	dtoUser := userConv.ToDTO(user)

	rows, err := u.participation.ListUserParticipations(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	stats := dto.UserDashboardStats{}
	participations := make([]dto.UserParticipation, 0, len(rows))

	for _, row := range rows {
		stats.RegistrationsCount++

		hackItem := dto.HackathonListItem{
			ID:                   row.HackathonID.String(),
			Title:                row.Title,
			Status:               row.Status,
			Format:               row.Format,
			RegistrationOpensAt:  row.RegistrationOpensAt.UTC().Format(time.RFC3339),
			SubmissionDeadlineAt: row.SubmissionDeadlineAt.UTC().Format(time.RFC3339),
		}
		if row.ShortDescription.Valid {
			hackItem.ShortDescription = row.ShortDescription.String
		}

		if row.Status == "registration" || row.Status == "running" {
			stats.ActiveHackathons++
		}

		var teamDTO *dto.UserParticipationTeam
		if row.TeamID.Valid {
			stats.TeamsCount++
			isCaptain := row.CaptainID.Valid && row.CaptainID.String == userID.String()
			teamDTO = &dto.UserParticipationTeam{
				ID:          row.TeamID.String,
				Name:        row.TeamName.String,
				IsCaptain:   isCaptain,
				TeamRole:    row.TeamRole.String,
				MemberCount: row.TeamMemberCount,
			}
			if row.TrackTitle.Valid {
				teamDTO.TrackTitle = row.TrackTitle.String
			}
			if row.CaseTitle.Valid {
				teamDTO.CaseTitle = row.CaseTitle.String
			}
		}

		var subDTO *dto.UserParticipationSubmission
		if row.SubmissionID.Valid {
			stats.SubmissionsCount++
			subDTO = &dto.UserParticipationSubmission{
				ID:      row.SubmissionID.String,
				RepoURL: row.SubmissionRepoURL.String,
			}
			if row.SubmissionTitle.Valid {
				subDTO.Title = row.SubmissionTitle.String
			}
			if row.SubmissionSubmittedAt.Valid {
				s := row.SubmissionSubmittedAt.Time.UTC().Format(time.RFC3339)
				subDTO.SubmittedAt = &s
			}
		}

		canSubmit, blockReason := submitBlockFromRow(row, now)

		participations = append(participations, dto.UserParticipation{
			Hackathon:         hackItem,
			RegisteredAt:      row.RegisteredAt.UTC().Format(time.RFC3339),
			Team:              teamDTO,
			Submission:        subDTO,
			CanSubmit:         canSubmit,
			SubmitBlockReason: blockReason,
		})
	}

	dashboard := &dto.UserDashboard{
		User:           dtoUser,
		Stats:          stats,
		Participations: participations,
	}

	if user.PlatformRole == "organizer" {
		organized, err := u.hackathons.List(ctx, repository.HackathonListFilter{OrganizerID: &userID})
		if err != nil {
			return nil, err
		}
		dashboard.Organized = make([]dto.HackathonListItem, 0, len(organized))
		for _, h := range organized {
			dashboard.Organized = append(dashboard.Organized, hackConv.ToListItem(h))
		}
		stats.OrganizedCount = len(dashboard.Organized)
		dashboard.Stats = stats
	}

	return dashboard, nil
}

func submitBlockFromRow(row entities.UserParticipationRow, now time.Time) (bool, string) {
	if !row.TeamID.Valid {
		return false, "no_team"
	}
	if !row.HasCase {
		return false, "no_case"
	}
	if !now.Before(row.SubmissionDeadlineAt) {
		return false, "deadline_passed"
	}
	if row.Status == "finished" {
		return false, "hackathon_finished"
	}
	if row.Status != "registration" && row.Status != "running" {
		return false, "hackathon_not_active"
	}
	return true, ""
}

func normalizeTeamRole(role, fallback string) string {
	role = strings.TrimSpace(role)
	if role == "" {
		return fallback
	}
	if _, ok := validTeamRoles[role]; ok {
		return role
	}
	return fallback
}
