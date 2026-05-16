package users

import (
	"errors"
	"net/http"

	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
	httperror "github.com/nikitatisenko/pirksp/pkg/http/error"
	"github.com/nikitatisenko/pirksp/pkg/http/writer"
	"github.com/nikitatisenko/pirksp/pkg/logger"
	pkgmiddleware "github.com/nikitatisenko/pirksp/pkg/middleware"
)

func GetDashboard(uc usecase.ParticipationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.UserIDFromContext(r.Context())
		if !ok {
			httperror.Unauthorized(w, errs.ErrUnauthorized.Error())
			return
		}

		dashboard, err := uc.GetUserDashboard(r.Context(), userID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				httperror.NotFound(w, err.Error())
				return
			}
			logger.FromContext(r.Context()).Error("get dashboard failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, dashboard)
	}
}
