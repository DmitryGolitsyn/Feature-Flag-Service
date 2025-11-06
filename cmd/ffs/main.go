package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// атомарный флаг готовности (1 = ready, 0 = not ready)
	var ready int32
	atomic.StoreInt32(&ready, 1) // на старте считаем, что готовы; позже здесь будет логика БД/кэшей

	// liveness curl -i http://localhost:8080/healthz
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// readiness curl -i http://localhost:8080/readyz
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&ready) == 1 {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ready"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	})

	// http.Server отдельно — чтобы иметь Shutdown()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Контекст, который завершится при SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// фатальная ошибка старта/работы
			fmt.Println("http error:", err)
			stop() // просим главный контекст завершиться
		}
	}()

	// Ждём сигнала
	<-ctx.Done()

	// Перед выключением помечаем "не готовы"
	atomic.StoreInt32(&ready, 0)
	time.Sleep(3 * time.Second)

	// Даём до 10 секунд додать активные запросы
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Если не успели корректно — жёстко закрываем
		_ = srv.Close()
	}

	fmt.Println("bye")
}
