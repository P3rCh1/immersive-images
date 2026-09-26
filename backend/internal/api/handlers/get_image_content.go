package handlers

import (
	"errors"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetImageContent(w http.ResponseWriter, r *http.Request, UUID openapi_types.UUID) {
	imageMeta, err := h.db.GetImage(r.Context(), UUID)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.Error(w, http.StatusNotFound, msgs.ImageNotFound)
			return
		}

		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	data, err := h.s3.Get(r.Context(), UUID.String())
	if err != nil {
		h.Error(w, http.StatusServiceUnavailable, msgs.Unavailable)
		return
	}

	h.PNG(w, http.StatusOK, *imageMeta, data)
}
