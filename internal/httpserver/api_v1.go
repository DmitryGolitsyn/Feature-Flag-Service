package httpserver

import (
	"ffs-tutorial/internal/app"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// mountAPIv1 вешает все маршруты версии v1
func mountAPIv1(r chi.Router, app *app.Application) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/ping", handlePing(app))

	})
}

type pingResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

func handlePing(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := app.Ping.Ping(r.Context())
		if err != nil {
			ErrorJSON(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		JSON(w, r, http.StatusOK, pingResponse{
			OK:        ok,
			RequestID: requestIDFromContext(r.Context()),
		})
	}
}
