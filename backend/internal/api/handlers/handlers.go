package handlers

import (
	"context"
	"database/sql"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"mime"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/common"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handlers struct {
	log *slog.Logger
	db  DB
	s3  S3
}

func New(log *slog.Logger, db DB, s3 S3) *Handlers {
	return &Handlers{
		log: log,
		db:  db,
		s3:  s3,
	}
}

type DB interface {
	Beginx(ctx context.Context) (context.Context, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (context.Context, error)
	Rollback(ctx context.Context)
	Commit(ctx context.Context) error

	ListImages(ctx context.Context, limit int, start *uuid.UUID) ([]postgres.Image, error)
	TotalImages(ctx context.Context) (int, error)
	CreateImage(ctx context.Context, image *postgres.Image) error
	GetImage(ctx context.Context, id uuid.UUID) (*postgres.Image, error)
	UpdateImage(ctx context.Context, image *postgres.Image) error
	SetUploaded(ctx context.Context, id uuid.UUID) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type S3 interface {
	Upload(ctx context.Context, objectKey string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

func (h *Handlers) JSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.MarshalWrite(w, v, json.DefaultOptionsV2()); err != nil {
		h.log.Error(
			"failed to send response",
			"error", err,
		)
	}
}

func (h *Handlers) PNG(w http.ResponseWriter, statusCode int, meta postgres.Image, data []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprint(meta.Size))
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{
		"filename": meta.Name + ".png",
	}))
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		h.log.Error(
			"failed to send response",
			"error", err,
		)
	}
}

func (h *Handlers) Error(w http.ResponseWriter, statusCode int, errMsg string) {
	apiError := models.Error{
		Error: struct {
			Message string "json:\"message\""
		}{
			Message: errMsg,
		},
	}

	h.JSON(w, statusCode, apiError)
}

func apiImageMetaFromDB(img postgres.Image) models.ImageMeta {
	var additional *models.Additional
	if img.Additional != nil {
		additional = common.Ptr(models.Additional(*img.Additional))
	}

	return models.ImageMeta{
		Uuid:      img.ID,
		Name:      img.Name,
		CreatedAt: img.CreatedAt,
		Seed:      img.Seed,
		Settings: models.GenerationSettingsResponse{
			Style:      models.Geometry(img.Style),
			Palette:    models.Palette(img.Palette),
			Additional: additional,
			Width:      img.Width,
			Height:     img.Height,
			Scale:      img.Scale,
		},
	}
}

func (h *Handlers) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ValidationErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	h.Error(w, http.StatusBadRequest, err.Error())
}

func (h *Handlers) DeleteImage(w http.ResponseWriter, _ *http.Request, _ openapi_types.UUID) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
}

func (h *Handlers) GetImage(w http.ResponseWriter, _ *http.Request, _ openapi_types.UUID) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
}

func (h *Handlers) UpdateImage(w http.ResponseWriter, _ *http.Request, _ openapi_types.UUID) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
}
