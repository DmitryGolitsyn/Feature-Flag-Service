package config

import (
	"os"
	"strconv"
	"time"
)

type HTTP struct {
	Addr              string        // вида ":8080"
	ReadHeaderTimeout time.Duration // 5s
	WriteTimeout      time.Duration // 15s
	IdleTimeout       time.Duration // 60s
}

func LoadHTTP() HTTP {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return HTTP{
		Addr:              port,
		ReadHeaderTimeout: envDuration("READ_HEADER_TIMEOUT", 5*time.Second),
		WriteTimeout:      envDuration("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:       envDuration("IDLE_TIMEOUT", 60*time.Second)}
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		// принимаем секунды как целое число: "10" => 10s
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return time.Duration(n) * time.Second
		}
		// или полноценные duration-строки: "250ms", "1.5s", "2m"
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
