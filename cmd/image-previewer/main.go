package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devv4n/image-previewer/internal/api/rest"
	"github.com/devv4n/image-previewer/internal/cache"
	"github.com/devv4n/image-previewer/internal/config"
	"github.com/devv4n/image-previewer/internal/service"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "c", ".config.json", "local config file path")
	flag.Parse()

	lvl := new(slog.LevelVar)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	slog.SetDefault(logger)

	slog.Info("service starting")

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("error loading configuration", "error", err)
	}

	lvl.Set(cfg.LogLevel)

	cch := cache.NewLRUCache(cfg.CacheSize)

	svc := service.NewService(cch)

	srv := rest.NewServer(svc, cfg)

	go srv.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("stopping service")

	if err = srv.Shutdown(context.Background()); err != nil {
		slog.Error("error stopping server", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
