package health

import (
	"net/http"

	"github.com/nikitatisenko/pirksp/internal/infrastructure/database/postgres"
	"github.com/nikitatisenko/pirksp/pkg/http/writer"
)

type Response struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func Check(db *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			Status:   "ok",
			Database: "ok",
		}

		if err := db.SqlxDB().PingContext(r.Context()); err != nil {
			resp.Status = "degraded"
			resp.Database = "unavailable"
			writer.WriteJSON(w, http.StatusServiceUnavailable, resp)
			return
		}

		writer.WriteJSON(w, http.StatusOK, resp)
	}
}
