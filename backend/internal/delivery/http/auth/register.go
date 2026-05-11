package auth

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
)

func Register(uc usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		resp, err := uc.Register(r.Context(), req)
		if err != nil {
			switch {
			case errors.Is(err, errs.ErrEmailTaken):
				httperror.Conflict(w, err.Error())
			case errors.Is(err, errs.ErrWeakPassword),
				errors.Is(err, errs.ErrEmptyEmail),
				errors.Is(err, errs.ErrEmptyName),
				errors.Is(err, errs.ErrInvalidRole):
				httperror.BadRequest(w, err.Error())
			default:
				logger.FromContext(r.Context()).Error("register failed", "err", err)
				httperror.Internal(w, "internal error")
			}
			return
		}

		writer.WriteJSON(w, http.StatusCreated, resp)
	}
}
