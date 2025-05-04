package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/config"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting server")

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to create config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	srvr, err := server.New(logger, cfg)
	if err != nil {
		logger.Error("failed to create server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info(fmt.Sprintf("server is running on port %s", cfg.ServerPort))
		if err := srvr.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srvr.Server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", slog.String("error", err.Error()))
	}

	wg.Wait()
}
