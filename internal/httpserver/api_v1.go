package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// mountAPIv1 вешает все маршруты версии v1
func mountAPIv1(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/ping", handlePing)

	})
}

type pingResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	var rid string
	if v := r.Context().Value(reqIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			rid = s
		}
	}
	JSON(w, r, http.StatusOK, pingResponse{
		OK:        true,
		RequestID: rid,
	})
}
