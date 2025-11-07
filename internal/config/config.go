package config

import "os"

type HTTP struct {
	Addr string // вида ":8080"
}

func LoadHTTP() HTTP {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return HTTP{Addr: port}
}
