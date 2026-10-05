package handlers

import (
	"errors"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetImage(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	img, err := h.db.GetImage(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.Error(w, http.StatusNotFound, msgs.ImageNotFound)
			return
		}

		h.log.Error(
			"failed to get image from DB",
			"error", err,
		)

		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	h.JSON(w, http.StatusOK, apiImageMetaFromDB(*img))
}
