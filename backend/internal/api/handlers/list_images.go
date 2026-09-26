package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/api/server"
	"github.com/P3rCh1/immersive-images/backend/internal/common"
	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/google/uuid"
)

func (h *Handlers) collectListImagesData(ctx context.Context, limit int, start *uuid.UUID) (*models.ListImagesResponse, error) {
	txCtx, err := h.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}

	defer h.db.Rollback(txCtx)

	dalImages, err := h.db.ListImages(txCtx, limit+1, start)
	if err != nil {
		return nil, err
	}

	total, err := h.db.TotalImages(txCtx)
	if err != nil {
		return nil, err
	}

	var next *uuid.UUID
	if len(dalImages) > limit {
		next = common.Ptr(dalImages[len(dalImages)-1].ID)
		dalImages = dalImages[:limit]
	}

	apiImages := make([]models.ImageMeta, len(dalImages))
	for i, dalImage := range dalImages {
		apiImages[i] = apiImageMetaFromDB(dalImage)
	}

	return &models.ListImagesResponse{
		Items: apiImages,
		Limit: limit,
		Total: total,
		Next:  next,
	}, nil
}

func (h *Handlers) ListImages(w http.ResponseWriter, r *http.Request, params server.ListImagesParams) {
	limit := config.Config.Images.DefaultLimit
	if params.Limit != nil {
		limit = max(*params.Limit, 1)
		limit = min(limit, config.Config.Images.MaxLimit)
	}

	resp, err := h.collectListImagesData(r.Context(), limit, params.Start)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	h.JSON(w, http.StatusOK, resp)
}
