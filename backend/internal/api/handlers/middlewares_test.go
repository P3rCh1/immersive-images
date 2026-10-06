package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	handlers := New(slog.New(slog.DiscardHandler), nil, nil)
	router := chi.NewRouter()
	router.Use(handlers.Recover())

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("some error")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
