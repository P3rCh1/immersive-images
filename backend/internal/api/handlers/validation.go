package handlers

import "net/http"

func (h *Handlers) ValidationErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	h.errorResponse(w, http.StatusBadRequest, err)
}
