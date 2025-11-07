package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

// ===== Request-ID =====

const HeaderRequestID = "X-Request-Id"

type reqIDKey struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if id == "" {
			id = genReqID()
		}
		// прокинем в контекст и в ответ
		ctx := context.WithValue(r.Context(), reqIDKey{}, id)
		w.Header().Set(HeaderRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// genReqID — 16 байт криптослучайных, hex (32 символа)
func genReqID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// ===== statusRecorder для логирования статуса =====

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

		// достанем request-id из контекста (если есть)
		var rid string
		if v := r.Context().Value(reqIDKey{}); v != nil {
			if s, ok := v.(string); ok {
				rid = s
			}
		}

		// очень простой лог в stdout; позже заменим на structured logger
		dur := time.Since(start)
		// В проде лучше логировать в JSON; сейчас коротко и наглядно
		if rid != "" {
			println(rid, r.Method, r.URL.Path, rec.status, dur.String())
		} else {
			println(r.Method, r.URL.Path, rec.status, dur.String())
		}
	})
}
