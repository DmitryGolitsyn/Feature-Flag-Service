package httpserver

import (
	"encoding/json"
	"ffs-tutorial/internal/app"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// mountAPIv1 вешает все маршруты версии v1
func mountAPIv1(r chi.Router, application *app.Application) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/ping", handlePing(application))
		r.Post("/echo", handleEcho(application))
		r.Post("/flags", handleFlagUpsert(application))
		r.Post("/evaluate", handleEvaluate(application))
	})
}

type pingResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

func handlePing(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := application.Ping.Ping(r.Context())
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

func handleEcho(application *app.Application) http.HandlerFunc {
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
		respMsg, err := application.Echo.Do(r.Context(), req.Message)
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

type flagUpsertRequest struct {
	Tenant  string `json:"tenant"`
	Key     string `json:"key"`
	Rollout int    `json:"rollout"` // 0..100
}

type flagUpsertResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

func handleFlagUpsert(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req flagUpsertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ErrorJSON(w, r, http.StatusBadRequest, "invalid json")
			return
		}
		f := app.Flag{
			Tenant: app.TenantID(req.Tenant),
			Key:    app.FlagKey(req.Key),
			Rule: app.Rule{
				Rollout: app.Percentage(req.Rollout),
			},
		}
		if err := application.Flags.Upsert.Do(r.Context(), f); err != nil {
			ErrorJSON(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		JSON(w, r, http.StatusOK, flagUpsertResponse{
			OK:        true,
			RequestID: requestIDFromContext(r.Context()),
		})
	}
}

type evaluateRequest struct {
	Tenant string `json:"tenant"`
	Key    string `json:"key"`
	UserID string `json:"user_id"`
}
type evaluateResponse struct {
	Enabled   bool   `json:"enabled"`
	RequestID string `json:"request_id"`
}

func handleEvaluate(application *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req evaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ErrorJSON(w, r, http.StatusBadRequest, "invalid json")
			return
		}
		out, err := application.Flags.Evaluate.Do(r.Context(), app.EvaluateInput{
			Tenant: app.TenantID(req.Tenant),
			Key:    app.FlagKey(req.Key),
			UserID: req.UserID,
		})
		if err != nil {
			status := httpStatusFor(err)
			ErrorJSON(w, r, status, err.Error())
			return
		}
		JSON(w, r, http.StatusOK, evaluateResponse{
			Enabled:   out.Enabled,
			RequestID: requestIDFromContext(r.Context()),
		})
	}
}
