package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
	httperror "github.com/nikitatisenko/pirksp/pkg/http/error"
	"github.com/nikitatisenko/pirksp/pkg/http/writer"
	"github.com/nikitatisenko/pirksp/pkg/logger"
	pkgmiddleware "github.com/nikitatisenko/pirksp/pkg/middleware"
)

func GetMe(uc usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, errs.ErrUnauthorized.Error())
			return
		}

		user, err := uc.GetMe(r.Context(), userID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				httperror.NotFound(w, err.Error())
				return
			}
			logger.FromContext(r.Context()).Error("get me failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, user)
	}
}

func UpdateMe(uc usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, errs.ErrUnauthorized.Error())
			return
		}

		var req dto.UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		user, err := uc.UpdateMe(r.Context(), userID, req)
		if err != nil {
			if errors.Is(err, errs.ErrEmptyName) {
				httperror.BadRequest(w, err.Error())
				return
			}
			if errors.Is(err, errs.ErrNotFound) {
				httperror.NotFound(w, err.Error())
				return
			}
			logger.FromContext(r.Context()).Error("update me failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, user)
	}
}
