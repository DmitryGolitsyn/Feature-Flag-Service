package httpserver

import (
	"net/http"
	"time"
)

// statusRecorder позволяет узнать финальный статус ответа
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Logging — простой middleware логирования запросов
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// очень простой лог в stdout; позже заменим на structured logger
		dur := time.Since(start)
		// В проде лучше логировать в JSON; сейчас коротко и наглядно
		println(r.Method, r.URL.Path, rec.status, dur.String())
	})
}
