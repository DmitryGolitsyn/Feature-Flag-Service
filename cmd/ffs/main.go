package main

import (
	"context"
	"ffs-tutorial/internal/config"
	"ffs-tutorial/internal/httpserver"
)

func main() {
	httpCfg := config.LoadHTTP()
	s := httpserver.New(httpCfg.Addr)
	_ = s.Start(context.Background())
}
