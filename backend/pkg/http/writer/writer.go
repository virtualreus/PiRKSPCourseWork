package writer

import (
	"encoding/json"
	"net/http"

	"github.com/nikitatisenko/pirksp/pkg/http/header"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	header.AddJSONContentType(w.Header())
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
