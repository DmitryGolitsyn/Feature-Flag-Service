package main

import (
	"context"
	"ffs-tutorial/internal/httpserver"
)

func main() {
	s := httpserver.New()
	_ = s.Start(context.Background())
}
