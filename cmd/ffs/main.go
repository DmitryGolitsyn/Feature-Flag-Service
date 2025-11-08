package main

import (
	"context"
	"ffs-tutorial/internal/app"
	"ffs-tutorial/internal/config"
	"ffs-tutorial/internal/httpserver"
	"ffs-tutorial/internal/infra/mem"
)

func main() {
	httpCfg := config.LoadHTTP()

	repo := mem.NewFlagRepo()
	application := app.NewWithRepo(repo)

	s := httpserver.New(
		httpCfg.Addr,
		httpCfg.ReadHeaderTimeout,
		httpCfg.WriteTimeout,
		httpCfg.IdleTimeout,
		application,
	)

	if err := s.Start(context.Background()); err != nil {
		panic(err)
	}
}
