package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/server"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handlers struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Handlers {
	return &Handlers{
		log: log,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ListImages(w http.ResponseWriter, _ *http.Request, _ server.ListImagesParams) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
}

func (h *Handlers) CreateImage(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
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

func (h *Handlers) GetImageContent(w http.ResponseWriter, _ *http.Request, _ openapi_types.UUID) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "Unimplemented")
}
