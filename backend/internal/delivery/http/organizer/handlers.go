package organizer

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

func GetHackathon(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		item, err := uc.GetOrganizer(r.Context(), organizerID, hackathonID)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func ListHackathons(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		items, err := uc.ListOrganizer(r.Context(), organizerID)
		if err != nil {
			logger.FromContext(r.Context()).Error("list organizer hackathons failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func CreateHackathon(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		var req dto.CreateHackathonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.Create(r.Context(), organizerID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusCreated, item)
	}
}

func UpdateHackathon(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		var req dto.UpdateHackathonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.Update(r.Context(), organizerID, hackathonID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func DeleteHackathon(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		if err := uc.Delete(r.Context(), organizerID, hackathonID); err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func PublishHackathon(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		item, err := uc.Publish(r.Context(), organizerID, hackathonID)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func ListTracks(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		items, err := uc.ListTracks(r.Context(), organizerID, hackathonID)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func CreateTrack(uc usecase.HackathonUseCase) http.HandlerFunc {
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

		var req dto.CreateTrackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.CreateTrack(r.Context(), organizerID, hackathonID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusCreated, item)
	}
}

func UpdateTrack(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		trackID, err := uuid.Parse(chi.URLParam(r, "trackId"))
		if err != nil {
			httperror.BadRequest(w, "invalid track id")
			return
		}

		var req dto.CreateTrackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.UpdateTrack(r.Context(), organizerID, trackID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func DeleteTrack(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		trackID, err := uuid.Parse(chi.URLParam(r, "trackId"))
		if err != nil {
			httperror.BadRequest(w, "invalid track id")
			return
		}

		if err := uc.DeleteTrack(r.Context(), organizerID, trackID); err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListCases(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		trackID, err := uuid.Parse(chi.URLParam(r, "trackId"))
		if err != nil {
			httperror.BadRequest(w, "invalid track id")
			return
		}

		items, err := uc.ListCases(r.Context(), organizerID, trackID)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func CreateCase(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		trackID, err := uuid.Parse(chi.URLParam(r, "trackId"))
		if err != nil {
			httperror.BadRequest(w, "invalid track id")
			return
		}

		var req dto.CreateCaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.CreateCase(r.Context(), organizerID, trackID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusCreated, item)
	}
}

func UpdateCase(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		caseID, err := uuid.Parse(chi.URLParam(r, "caseId"))
		if err != nil {
			httperror.BadRequest(w, "invalid case id")
			return
		}

		var req dto.CreateCaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		item, err := uc.UpdateCase(r.Context(), organizerID, caseID, req)
		if err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func DeleteCase(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizerID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, "unauthorized")
			return
		}

		caseID, err := uuid.Parse(chi.URLParam(r, "caseId"))
		if err != nil {
			httperror.BadRequest(w, "invalid case id")
			return
		}

		if err := uc.DeleteCase(r.Context(), organizerID, caseID); err != nil {
			handleHackathonErr(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleHackathonErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		httperror.NotFound(w, err.Error())
	case errors.Is(err, errs.ErrNotOwner):
		httperror.Forbidden(w, err.Error())
	case errors.Is(err, errs.ErrValidation),
		errors.Is(err, errs.ErrCannotDelete),
		errors.Is(err, errs.ErrPublishNotReady):
		httperror.BadRequest(w, err.Error())
	default:
		logger.FromContext(r.Context()).Error("hackathon handler failed", "err", err)
		httperror.Internal(w, "internal error")
	}
}
