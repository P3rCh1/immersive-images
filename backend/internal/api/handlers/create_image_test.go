package handlers

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateImage_OK_RequiredParams(t *testing.T) {
	initImagesConfig()
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"testimg","seed":42}`))

	var created postgres.Image
	db.EXPECT().CreateImage(mock.Anything, mock.Anything).Run(func(_ context.Context, image *postgres.Image) {
		image.ID = uuid.New()
		created = *image
	}).Return(nil)
	s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	db.EXPECT().SetUploaded(mock.Anything, mock.Anything).Return(nil)

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	assert.Equal(t, "testimg", created.Name)
	assert.Equal(t, 42, created.Seed)
	assert.Equal(t, "PATCHES", created.Style)
	assert.Equal(t, "NEON", created.Palette)
	assert.Equal(t, 256, created.Width)
	assert.Equal(t, 256, created.Height)
	assert.Equal(t, float64(3), created.Scale)
	assert.Nil(t, created.Additional)

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	expData, err := json.Marshal(apiImageMetaFromDB(created))
	require.NoError(t, err)

	assert.JSONEq(t, string(expData), string(data))
}

func TestCreateImage_OK_AllParams(t *testing.T) {
	initImagesConfig()
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{
			"name":"testimg",
			"seed":42,
			"settings":{
				"style":"rings",
				"palette":"toxic",
				"width":64,
				"height":64,
				"scale":1.5,
				"additional":"pixel"
			}
		}`,
	))

	var created postgres.Image
	db.EXPECT().CreateImage(mock.Anything, mock.Anything).Run(func(_ context.Context, image *postgres.Image) {
		image.ID = uuid.New()
		created = *image
	}).Return(nil)
	s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	db.EXPECT().SetUploaded(mock.Anything, mock.Anything).Return(nil)

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	assert.Equal(t, "testimg", created.Name)
	assert.Equal(t, 42, created.Seed)
	assert.Equal(t, "RINGS", created.Style)
	assert.Equal(t, "TOXIC", created.Palette)
	assert.Equal(t, 64, created.Width)
	assert.Equal(t, 64, created.Height)
	assert.Equal(t, 1.5, created.Scale)
	require.NotNil(t, created.Additional)
	assert.Equal(t, "PIXEL", *created.Additional)

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	expData, err := json.Marshal(apiImageMetaFromDB(created))
	require.NoError(t, err)

	assert.JSONEq(t, string(expData), string(data))
}

func TestCreateImage_InvalidBody(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not a json"))

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateImage_InvalidName(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":""}`))

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateImage_Internal(t *testing.T) {
	initImagesConfig()
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"testimg","seed":42}`))

	db.EXPECT().CreateImage(mock.Anything, mock.Anything).Return(errors.New("internal"))

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateImage_Unavailable(t *testing.T) {
	initImagesConfig()
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"testimg","seed":42}`))

	db.EXPECT().CreateImage(mock.Anything, mock.Anything).Return(nil)
	s3.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("internal"))
	db.EXPECT().DeleteImageForce(mock.Anything, mock.Anything).Return(nil)

	handlers.CreateImage(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
