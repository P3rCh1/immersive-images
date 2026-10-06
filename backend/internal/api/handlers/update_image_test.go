package handlers

import (
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	hmock "github.com/P3rCh1/immersive-images/backend/internal/api/handlers/mocks"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateImage_OK(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"newname"}`))

	db.EXPECT().UpdateImage(mock.Anything, mock.Anything).Return(nil)

	handlers.UpdateImage(rec, req, openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	expData, err := json.Marshal(apiImageMetaFromDB(postgres.Image{
		ID:   img.ID,
		Name: "newname",
	}))
	require.NoError(t, err)

	assert.JSONEq(t, string(expData), string(data))
}

func TestUpdateImage_InvalidBody(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader("not a json"))

	handlers.UpdateImage(rec, req, openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImage_InvalidName(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":""}`))

	handlers.UpdateImage(rec, req, openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImage_NotFound(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"newname"}`))

	db.EXPECT().UpdateImage(mock.Anything, mock.Anything).Return(postgres.ErrNotFound)

	handlers.UpdateImage(rec, req, openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateImage_Internal(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"newname"}`))

	db.EXPECT().UpdateImage(mock.Anything, mock.Anything).Return(errors.New("internal"))

	handlers.UpdateImage(rec, req, openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
