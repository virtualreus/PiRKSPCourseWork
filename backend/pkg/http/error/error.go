package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/nikitatisenko/pirksp/pkg/http/header"
)

type Response struct {
	Error RespErr `json:"error"`
}

type RespErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, status int, code, message string) {
	header.AddJSONContentType(w.Header())
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{
		Error: RespErr{Code: code, Message: message},
	})
}

func Internal(w http.ResponseWriter, message string) {
	Write(w, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

func BadRequest(w http.ResponseWriter, message string) {
	Write(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

func NotFound(w http.ResponseWriter, message string) {
	Write(w, http.StatusNotFound, "NOT_FOUND", message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Write(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Conflict(w http.ResponseWriter, message string) {
	Write(w, http.StatusConflict, "CONFLICT", message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Write(w, http.StatusForbidden, "FORBIDDEN", message)
}
