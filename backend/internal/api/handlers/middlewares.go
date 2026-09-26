package handlers

import (
	"net/http"

	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/api/server"
)

func (h *Handlers) Recover() server.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					h.log.Error(
						"recovered",
						"error", err,
					)
					h.Error(w, http.StatusInternalServerError, msgs.Internal)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
