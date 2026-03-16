package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/whatsapp-tts/internal/app"
	"github.com/whatsapp-tts/internal/config"
	"github.com/whatsapp-tts/internal/observability"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create logger
	logger := observability.NewStdLogger()

	logger.Info("starting_whatsapp_tts")

	// Create application
	application := app.New(cfg, logger)

	// Start server in goroutine
	go func() {
		if err := application.Run(); err != nil {
			logger.Error("server_error", observability.F("error", err.Error()))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting_down")

	// Graceful shutdown
	if err := application.Stop(); err != nil {
		logger.Error("shutdown_error", observability.F("error", err.Error()))
	}

	// Wait for context cancellation (optional)
	<-context.Background().Done()
}
