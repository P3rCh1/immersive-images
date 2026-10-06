package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	hmock "github.com/P3rCh1/immersive-images/backend/internal/api/handlers/mocks"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImage_OK(t *testing.T) {
	db := hmock.NewMockDB(t)
	s3 := hmock.NewMockS3(t)
	handlers := New(slog.New(slog.DiscardHandler), db, s3)
	rec := httptest.NewRecorder()

	img := getTestImage()
	deleted := make(chan struct{})

	db.EXPECT().DeleteImage(mock.Anything, img.ID).Return(nil)
	s3.EXPECT().Delete(mock.Anything, img.ID.String()).Run(func(context.Context, string) {
		close(deleted)
	}).Return(nil)

	handlers.DeleteImage(rec, httptest.NewRequest(http.MethodDelete, "/", nil), openapi_types.UUID(img.ID))

	<-deleted

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteImage_NotFound(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().DeleteImage(mock.Anything, img.ID).Return(postgres.ErrNotFound)

	handlers.DeleteImage(rec, httptest.NewRequest(http.MethodDelete, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteImage_Internal(t *testing.T) {
	db := hmock.NewMockDB(t)
	handlers := New(slog.New(slog.DiscardHandler), db, nil)
	rec := httptest.NewRecorder()

	img := getTestImage()

	db.EXPECT().DeleteImage(mock.Anything, img.ID).Return(errors.New("internal"))

	handlers.DeleteImage(rec, httptest.NewRequest(http.MethodDelete, "/", nil), openapi_types.UUID(img.ID))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
