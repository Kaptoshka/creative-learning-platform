package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func JSON(log *slog.Logger, w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("failed to encode response", "error", err)
	}
}

func OK(log *slog.Logger, w http.ResponseWriter, data any) {
	JSON(log, w, http.StatusOK, data)
}

func Created(log *slog.Logger, w http.ResponseWriter, data any) {
	JSON(log, w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(log *slog.Logger, w http.ResponseWriter, err error) {
	status := domainErrToStatus(err)
	JSON(log, w, status, errorResponse{Error: err.Error()})
}

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func domainErrToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUnauthenticated):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrPermissionDenied):
		return http.StatusForbidden
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrSubmissionClosed):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.ErrSubmissionNotSubmitted):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.ErrNoUpdates):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrInvalidPageToken):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
