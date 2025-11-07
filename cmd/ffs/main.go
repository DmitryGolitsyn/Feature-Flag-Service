package main

import (
	"context"
	"ffs-tutorial/internal/config"
	"ffs-tutorial/internal/httpserver"
)

func main() {
	httpCfg := config.LoadHTTP()
	s := httpserver.New(
		httpCfg.Addr,
		httpCfg.ReadHeaderTimeout,
		httpCfg.WriteTimeout,
		httpCfg.IdleTimeout,
	)
	_ = s.Start(context.Background())
}
