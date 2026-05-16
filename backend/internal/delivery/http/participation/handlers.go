package participation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
	httperror "github.com/nikitatisenko/pirksp/pkg/http/error"
	"github.com/nikitatisenko/pirksp/pkg/http/writer"
	"github.com/nikitatisenko/pirksp/pkg/logger"
	pkgmiddleware "github.com/nikitatisenko/pirksp/pkg/middleware"
)

func GetParticipation(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		status, err := uc.GetParticipation(r.Context(), userID, hackathonID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, status)
	}
}

func Register(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		reg, err := uc.RegisterForHackathon(r.Context(), userID, hackathonID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusCreated, reg)
	}
}

func Unregister(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		if err := uc.UnregisterFromHackathon(r.Context(), userID, hackathonID); err != nil {
			handleErr(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListTeams(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		items, err := uc.ListTeams(r.Context(), hackathonID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func CreateTeam(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		var req dto.CreateTeamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		team, err := uc.CreateTeam(r.Context(), userID, hackathonID, req)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusCreated, team)
	}
}

func GetTeam(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		team, err := uc.GetTeam(r.Context(), teamID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, team)
	}
}

func UpdateTeam(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		var req dto.UpdateTeamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		team, err := uc.UpdateTeam(r.Context(), userID, teamID, req)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, team)
	}
}

func JoinTeam(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		var req dto.JoinTeamRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		team, err := uc.JoinTeam(r.Context(), userID, teamID, req)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, team)
	}
}

func LeaveTeam(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		if err := uc.LeaveTeam(r.Context(), userID, teamID); err != nil {
			handleErr(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateMemberRole(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		memberID, err := uuid.Parse(chi.URLParam(r, "userId"))
		if err != nil {
			httperror.BadRequest(w, "invalid user id")
			return
		}

		var req dto.UpdateTeamMemberRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		member, err := uc.UpdateTeamMemberRole(r.Context(), userID, teamID, memberID, req)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, member)
	}
}

func GetSubmission(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		sub, err := uc.GetTeamSubmission(r.Context(), userID, teamID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, sub)
	}
}

func UpsertSubmission(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		teamID, err := uuid.Parse(chi.URLParam(r, "teamId"))
		if err != nil {
			httperror.BadRequest(w, "invalid team id")
			return
		}

		var req dto.UpsertSubmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		sub, err := uc.UpsertTeamSubmission(r.Context(), userID, teamID, req)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, sub)
	}
}

func ListRegistrations(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		items, err := uc.ListHackathonRegistrations(r.Context(), organizerID, hackathonID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func ListSubmissions(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		hackathonID, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		items, err := uc.ListHackathonSubmissions(r.Context(), organizerID, hackathonID)
		if err != nil {
			handleErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func handleErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		httperror.NotFound(w, err.Error())
	case errors.Is(err, errs.ErrForbidden), errors.Is(err, errs.ErrNotOwner),
		errors.Is(err, errs.ErrNotCaptain), errors.Is(err, errs.ErrParticipantsOnly):
		httperror.Forbidden(w, err.Error())
	case errors.Is(err, errs.ErrValidation), errors.Is(err, errs.ErrAlreadyRegistered),
		errors.Is(err, errs.ErrNotRegistered), errors.Is(err, errs.ErrAlreadyInTeam),
		errors.Is(err, errs.ErrNotInTeam), errors.Is(err, errs.ErrTeamFull),
		errors.Is(err, errs.ErrCaptainCannotLeave), errors.Is(err, errs.ErrDeadlinePassed),
		errors.Is(err, errs.ErrRegistrationClosed), errors.Is(err, errs.ErrCaseRequired),
		errors.Is(err, errs.ErrInTeamBlockUnregister):
		httperror.BadRequest(w, err.Error())
	default:
		logger.FromContext(r.Context()).Error("participation handler", "err", err)
		httperror.Internal(w, "internal error")
	}
}
