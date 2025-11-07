package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	ready      *int32
}

// New creates router + http.Server
func New(addr string, rht, wt, it time.Duration) *Server {
	r := chi.NewRouter()

	// атомарный флаг готовности (1 = ready, 0 = not ready)
	var ready int32
	atomic.StoreInt32(&ready, 1) // на старте считаем, что готовы; позже здесь будет логика БД/кэшей

	r.Use(RequestID)
	r.Use(RealIPUA)
	r.Use(Logging)

	//r.Use(func(next http.Handler) http.Handler { return Logging(next) })

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

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: rht,
		WriteTimeout:      wt,
		IdleTimeout:       it,
	}

	return &Server{
		httpServer: srv,
		ready:      &ready,
	}
}

// Start blocks until shutdown
func (s *Server) Start(ctx context.Context) error {
	// для отмены по сигналу
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigCh:
			fmt.Println("got signal:", sig.String())
		}
		s.Shutdown()
	}()

	fmt.Println("listening on", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	atomic.StoreInt32(s.ready, 0)
	s.httpServer.SetKeepAlivesEnabled(false)
	fmt.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Println("forced close:", err)
		_ = s.httpServer.Close()
	}

	fmt.Println("bye")
}
