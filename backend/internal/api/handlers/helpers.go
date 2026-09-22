package handlers

import (
	"encoding/json/v2"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
)

func (h *Handlers) structResponse(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.MarshalWrite(w, v, json.DefaultOptionsV2()); err != nil {
		h.log.Error(
			"failed to send response",
			"error", err,
		)
	}
}

func (h *Handlers) errorResponse(w http.ResponseWriter, statusCode int, err error) {
	apiError := models.Error{
		Error: struct {
			Message string "json:\"message\""
		}{
			Message: err.Error(),
		},
	}

	h.structResponse(w, statusCode, apiError)
}
