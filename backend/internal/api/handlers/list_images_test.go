package handlers

import (
	"context"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	hmock "github.com/P3rCh1/immersive-images/backend/internal/api/handlers/mocks"
	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/api/server"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListImages_OK(t *testing.T) {
	initImagesConfig()
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	imgs := []postgres.Image{getTestImage(), getTestImage()}

	db.EXPECT().BeginTxx(mock.Anything, mock.Anything).Return(context.Background(), nil)
	db.EXPECT().Rollback(mock.Anything)
	db.EXPECT().ListImages(mock.Anything, 11, mock.Anything).Return(imgs, nil)
	db.EXPECT().TotalImages(mock.Anything).Return(len(imgs), nil)

	handlers.ListImages(rec, httptest.NewRequest(http.MethodGet, "/", nil), server.ListImagesParams{})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get("Content-Type"), "application/json")

	data, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	expData, err := json.Marshal(models.ListImagesResponse{
		Items: []models.ImageMeta{
			apiImageMetaFromDB(imgs[0]),
			apiImageMetaFromDB(imgs[1]),
		},
		Limit: 10,
		Total: 2,
	})
	require.NoError(t, err)

	assert.JSONEq(t, string(expData), string(data))
}

func TestListImages_Internal(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	db.EXPECT().BeginTxx(mock.Anything, mock.Anything).Return(nil, errors.New("internal"))

	handlers.ListImages(rec, httptest.NewRequest(http.MethodGet, "/", nil), server.ListImagesParams{})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
