package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/common"
	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var teststruct = struct {
	Integer int    `json:"integer"`
	String  string `json:"string"`
}{
	Integer: 5,
	String:  "somestr",
}

func TestJSON(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)
	rec := httptest.NewRecorder()

	handlers.JSON(rec, http.StatusTeapot, &teststruct)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"integer":5,"string":"somestr"}`, string(data))
}

func TestHealthz(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)

	rec := httptest.NewRecorder()

	handlers.Healthz(rec, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPNG(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)

	rec := httptest.NewRecorder()

	meta := postgres.Image{
		Name: "testimg",
		Size: 4,
	}
	data := []byte{0x89, 0x50, 0x4e, 0x47}

	handlers.PNG(rec, http.StatusOK, meta, data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "image/png")
	assert.Equal(t, rec.Header().Get("Content-Length"), "4")
	assert.Equal(t, rec.Header().Get("Cache-Control"), "public, max-age=31536000, immutable")
	assert.Equal(t, rec.Header().Get("Content-Disposition"), "inline; filename=testimg.png")

	sent, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.Equal(t, data, sent)
}

func TestError(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)

	rec := httptest.NewRecorder()

	handlers.Error(rec, http.StatusBadRequest, "some error")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error":{"message":"some error"}}`, string(data))
}

func TestValidationErrorHandler(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)

	rec := httptest.NewRecorder()

	handlers.ValidationErrorHandler(rec, nil, errors.New("internal"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error":{"message":"internal"}}`, string(data))
}

func TestAPIImageMetaFromDB(t *testing.T) {
	createdAt := time.Date(2026, time.October, 6, 12, 0, 0, 0, time.UTC)
	additional := "glow"

	img := postgres.Image{
		ID:         uuid.MustParse("d1c07a5e-8c6b-4d9a-9f3e-2a1b0c9d8e7f"),
		Name:       "someimg",
		Size:       4,
		Seed:       42,
		Style:      "circles",
		Palette:    "neon",
		Additional: &additional,
		Width:      1920,
		Height:     1080,
		Scale:      1.5,
		CreatedAt:  createdAt,
		Uploaded:   true,
	}

	meta := apiImageMetaFromDB(img)

	assert.Equal(t, models.ImageMeta{
		Uuid:      img.ID,
		Name:      img.Name,
		CreatedAt: createdAt,
		Seed:      42,
		Settings: models.GenerationSettingsResponse{
			Style:      models.Geometry("circles"),
			Palette:    models.Palette("neon"),
			Additional: common.Ptr(models.Additional("glow")),
			Width:      1920,
			Height:     1080,
			Scale:      1.5,
		},
	}, meta)
}

func getTestImage() postgres.Image {
	return postgres.Image{
		ID:         uuid.New(),
		Name:       "testimg",
		Size:       4,
		Seed:       42,
		Style:      "circles",
		Palette:    "neon",
		Additional: common.Ptr("pixel"),
		Width:      1920,
		Height:     1080,
		Scale:      1.5,
		CreatedAt:  time.Now(),
		Uploaded:   true,
	}
}

func initImagesConfig() {
	config.Config.Images = config.Images{
		DefaultLimit:   10,
		MaxLimit:       100,
		DefaultStyle:   "PATCHES",
		DefaultPalette: "NEON",
		DefaultWidth:   256,
		DefaultHeight:  256,
		MinWidth:       64,
		MaxWidth:       1024,
		MinHeight:      64,
		MaxHeight:      1024,
		DefaultScale:   3,
		MinScale:       0.5,
		MaxScale:       10,
	}
}
