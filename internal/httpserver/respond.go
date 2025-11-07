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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

// ErrorJSON — всегда отдаёт одну и ту же структуру ошибки
func ErrorJSON(w http.ResponseWriter, r *http.Request, status int, msg string) {
	var rid string
	if v := r.Context().Value(reqIDKey{}); v != nil {
		rid, _ = v.(string)
	}

	resp := ErrorResponse{
		Error:      msg,
		RequestID:  rid,
		StatusCode: status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}
