package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
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

		ip := ""
		if v := r.Context().Value(realIPKey{}); v != nil {
			if s, ok := v.(string); ok {
				ip = s
			}
		}
		ua := ""
		if v := r.Context().Value(userAgentKey{}); v != nil {
			if s, ok := v.(string); ok {
				ua = s
			}
		}

		// очень простой лог в stdout; позже заменим на structured logger
		dur := time.Since(start)
		// В проде лучше логировать в JSON; сейчас коротко и наглядно
		if rid != "" {
			println(rid, r.Method, r.URL.Path, rec.status, dur.String(), ip, ua)
		} else {
			println(r.Method, r.URL.Path, rec.status, dur.String(), ip, ua)
		}
	})
}

// ===== RealIP + UserAgent =====

type realIPKey struct{}
type userAgentKey struct{}

// RealIPUA извлекает IP клиента (с учётом прокси) и User-Agent и кладёт в контекст
func RealIPUA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		ua := r.Header.Get("User-Agent")

		ctx := context.WithValue(r.Context(), realIPKey{}, ip)
		ctx = context.WithValue(ctx, userAgentKey{}, ua)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// clientIP пытается аккуратно извлечь реальный IP из заголовков прокси
func clientIP(r *http.Request) string {
	// 1) X-Forwarded-For: может быть списком IP через запятую — берём первый не-пустой
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// формат: "client, proxy1, proxy2"
		for _, p := range splitCommaTrim(xff) {
			if p != "" {
				return p
			}
		}
	}
	// 2) X-Real-Ip
	if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		return xrip
	}
	// 3) RemoteAddr (host:port)
	host, _, err := netSplitHostPort(r.RemoteAddr)
	if err != nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

// splitCommaTrim — помощник для X-Forwarded-For
func splitCommaTrim(s string) []string {
	out := make([]string, 0, 4)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, trimSpaces(s[start:i]))
			start = i + 1
		}
	}
	out = append(out, trimSpaces(s[start:]))
	return out
}

func trimSpaces(s string) string {
	// без аллокаций ради простоты можно strings.TrimSpace, но оставим компактно

	i, j := 0, len(s)-1
	for i <= j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\t') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}

// netSplitHostPort — обёртка над net.SplitHostPort, чтобы не импортить вверху
func netSplitHostPort(addr string) (host, port string, err error) {
	// локально дернем из std без глобального импорта наверху
	type splitter interface {
		SplitHostPort(hostport string) (host, port string, err error)
	}
	return (&netPkg{}).SplitHostPort(addr)
}

type netPkg struct{}

func (*netPkg) SplitHostPort(hostport string) (string, string, error) {
	return net.SplitHostPort(hostport)
}
