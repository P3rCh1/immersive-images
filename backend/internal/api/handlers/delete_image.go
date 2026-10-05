package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) DeleteImage(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	if err := h.db.DeleteImage(r.Context(), id); err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.Error(w, http.StatusNotFound, msgs.ImageNotFound)
			return
		}

		h.log.Error(
			"failed to delete image from DB",
			"error", err,
		)

		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	go func() { //nolint:gosec // G118: r.Context() is cancelled
		if err := h.s3.Delete(context.Background(), id.String()); err != nil {
			h.log.Error(
				"failed to delete from S3",
				"need_cleanup", "S3",
				"image_id", id.String(),
				"error", err,
			)
		}
	}()
}
