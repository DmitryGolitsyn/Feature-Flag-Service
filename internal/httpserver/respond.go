package httpserver

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — единый стандарт ошибки
type ErrorResponse struct {
	Error      string `json:"error"`
	RequestID  string `json:"request_id"`
	StatusCode int    `json:"status_code"`
}

// JSON — помогает отдавать обычные успешные ответы
func JSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	// проставим X-Request-Id в ответ
	if rid := requestIDFromContext(r.Context()); rid != "" {
		w.Header().Set(HeaderRequestID, rid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

// ErrorJSON — всегда отдаёт одну и ту же структуру ошибки
func ErrorJSON(w http.ResponseWriter, r *http.Request, status int, msg string) {
	rid := requestIDFromContext(r.Context())

	resp := ErrorResponse{
		Error:      msg,
		RequestID:  rid,
		StatusCode: status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}
