package httpserver

import (
	"encoding/json"
	"ffs-tutorial/internal/app"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// mountAPIv1 вешает все маршруты версии v1
func mountAPIv1(r chi.Router, app *app.Application) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/ping", handlePing(app))
		r.Post("/echo", handleEcho(app))
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
			status := httpStatusFor(err)
			ErrorJSON(w, r, status, err.Error())
			return
		}
		JSON(w, r, http.StatusOK, pingResponse{
			OK:        ok,
			RequestID: requestIDFromContext(r.Context()),
		})
	}
}

type echoRequest struct {
	Message string `json:"message"`
}

type echoResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func handleEcho(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req echoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ErrorJSON(w, r, http.StatusBadRequest, "invalid json")
			return
		}
		// минимальная валидация (это HTTP-слой)
		if req.Message == "" {
			ErrorJSON(w, r, http.StatusBadRequest, "message is required")
			return
		}
		// вызываем usecase
		respMsg, err := app.Echo.Do(r.Context(), req.Message)
		if err != nil {
			status := httpStatusFor(err)
			ErrorJSON(w, r, status, err.Error())
			return
		}

		JSON(w, r, http.StatusOK, echoResponse{
			Message:   respMsg,
			RequestID: requestIDFromContext(r.Context()),
		})
	}
}
