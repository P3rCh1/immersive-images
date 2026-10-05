package handlers

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) UpdateImage(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	var req models.UpdateImageRequest
	if err := json.UnmarshalRead(r.Body, &req, json.DefaultOptionsV2()); err != nil {
		h.Error(w, http.StatusBadRequest, msgs.InvalidBody)
		return
	}

	if len(req.Name) == 0 || len(req.Name) > maxImageNameLen {
		h.Error(
			w,
			http.StatusBadRequest,
			fmt.Sprintf(msgs.InvalidNameLenFmt, len(req.Name), 1, maxImageNameLen),
		)

		return
	}

	img := postgres.Image{
		ID:   id,
		Name: req.Name,
	}

	if err := h.db.UpdateImage(r.Context(), &img); err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.Error(w, http.StatusNotFound, msgs.ImageNotFound)
			return
		}

		h.log.Error(
			"failed to update image in DB",
			"error", err,
		)

		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	h.JSON(w, http.StatusOK, apiImageMetaFromDB(img))
}
