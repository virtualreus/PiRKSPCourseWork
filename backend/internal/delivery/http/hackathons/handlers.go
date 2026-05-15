package hackathons

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
	httperror "github.com/nikitatisenko/pirksp/pkg/http/error"
	"github.com/nikitatisenko/pirksp/pkg/http/writer"
	"github.com/nikitatisenko/pirksp/pkg/logger"
)

func List(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := uc.ListPublic(r.Context(), r.URL.Query()["status"], r.URL.Query().Get("q"))
		if err != nil {
			logger.FromContext(r.Context()).Error("list hackathons failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}

func Get(uc usecase.HackathonUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "hackathonId"))
		if err != nil {
			httperror.BadRequest(w, "invalid hackathon id")
			return
		}

		item, err := uc.GetPublic(r.Context(), id)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				httperror.NotFound(w, err.Error())
				return
			}
			logger.FromContext(r.Context()).Error("get hackathon failed", "err", err)
			httperror.Internal(w, "internal error")
			return
		}

		writer.WriteJSON(w, http.StatusOK, item)
	}
}

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
