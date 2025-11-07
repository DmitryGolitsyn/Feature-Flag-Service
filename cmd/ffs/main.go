package main

import (
	"context"
	"ffs-tutorial/internal/app"
	"ffs-tutorial/internal/config"
	"ffs-tutorial/internal/httpserver"
)

func main() {
	application := app.New()
	httpCfg := config.LoadHTTP()
	s := httpserver.New(
		httpCfg.Addr,
		httpCfg.ReadHeaderTimeout,
		httpCfg.WriteTimeout,
		httpCfg.IdleTimeout,
		application,
	)
	_ = s.Start(context.Background())
}
