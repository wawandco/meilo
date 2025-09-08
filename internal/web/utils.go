package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// responseWrapper captures the status code for logging
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode error response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
