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

func Login(uc usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperror.BadRequest(w, errs.ErrInvalidJSON.Error())
			return
		}

		resp, err := uc.Login(r.Context(), req)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidCreds) {
				httperror.Unauthorized(w, err.Error())
				return
			}
			logger.FromContext(r.Context()).Error("login failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, resp)
	}
}
