package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	hmock "github.com/P3rCh1/immersive-images/backend/internal/api/handlers/mocks"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetImageContent_OK(t *testing.T) {
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	img := getTestImage()
	png := []byte{0x89, 0x50, 0x4e, 0x47}

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(&img, nil)
	s3.EXPECT().Get(mock.Anything, img.ID.String()).Return(png, nil)

	handlers.GetImageContent(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "image/png")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.Equal(t, png, data)
}

func TestGetImageContent_NotFound(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(nil, postgres.ErrNotFound)

	handlers.GetImageContent(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetImageContent_Internal(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(nil, errors.New("internal"))

	handlers.GetImageContent(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetImageContent_Unavailable(t *testing.T) {
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(&img, nil)
	s3.EXPECT().Get(mock.Anything, img.ID.String()).Return(nil, errors.New("internal"))

	handlers.GetImageContent(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
