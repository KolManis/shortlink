package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	infrastructurePostgres "github.com/KolManis/shortlink/internal/infrastructure/postgres"
	"github.com/KolManis/shortlink/internal/repository/cache"
	postgresRepo "github.com/KolManis/shortlink/internal/repository/postgres"
	transportHttp "github.com/KolManis/shortlink/internal/transport/http"
	httpHandlers "github.com/KolManis/shortlink/internal/transport/http/handlers"
	url "github.com/KolManis/shortlink/internal/usecase/url"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// Level: slog.LevelInfo,
	}))

	cfg := loadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := infrastructurePostgres.Open(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("open postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	urlRepo := postgresRepo.New(pool)
	urlCache := cache.NewMemoryCache()
	urlUseCase := url.NewService(urlRepo, urlCache, logger)
	urlHandler := httpHandlers.NewUrlHandler(urlUseCase)

	router := transportHttp.NewRouter(urlHandler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", "error", err)
		}
	}()

	logger.Info("http server started", "addr", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}

type config struct {
	HTTPAddr    string
	DatabaseDSN string
}

func loadConfig() config {
	cfg := config{
		HTTPAddr:    envOrDefault("HTTP_ADDR", ":8080"),
		DatabaseDSN: envOrDefault("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/urlservice?sslmode=disable"),
	}

	if cfg.DatabaseDSN == "" {
		panic(fmt.Errorf("DATABASE_DSN is required"))
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
