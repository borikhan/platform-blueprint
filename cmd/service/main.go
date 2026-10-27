package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct{ Port, Environment, Release string }

func load() (config, error) {
	c := config{Port: get("PORT", "8080"), Environment: os.Getenv("APP_ENV"), Release: os.Getenv("RELEASE_SHA")}
	if c.Environment == "" || c.Release == "" {
		return c, errors.New("APP_ENV and RELEASE_SHA are required")
	}
	return c, nil
}
func main() {
	c, err := load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]string{"status": "alive"}) })
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]string{"status": "ready", "release": c.Release})
	})
	mux.HandleFunc("GET /v1/info", func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]string{"environment": c.Environment, "release": c.Release})
	})
	srv := http.Server{Addr: ":" + c.Port, Handler: mux, ReadHeaderTimeout: 3 * time.Second}
	go func() {
		slog.Info("service listening", "port", c.Port, "release", c.Release)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	slog.Info("shutdown complete")
}
func get(k, f string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return f
}
func write(w http.ResponseWriter, s int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(v)
}
