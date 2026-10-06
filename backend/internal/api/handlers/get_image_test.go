package handlers

import (
	"encoding/json/v2"
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

func TestGetImage_OK(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(&img, nil)

	handlers.GetImage(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	expData, err := json.Marshal(apiImageMetaFromDB(img))
	require.NoError(t, err)

	assert.JSONEq(t, string(expData), string(data))
}

func TestGetImage_NotFound(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(nil, postgres.ErrNotFound)

	handlers.GetImage(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetImage_Internal(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().GetImage(mock.Anything, img.ID).Return(nil, errors.New("internal"))

	handlers.GetImage(rec, httptest.NewRequest(http.MethodGet, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
