package app

import (
	"net/http"

	"github.com/whatsapp-tts/internal/config"
	"github.com/whatsapp-tts/internal/observability"
	"github.com/whatsapp-tts/internal/webhook"
)

// App representa la aplicación principal.
type App struct {
	config *config.Config
	logger observability.Logger
	server *http.Server
}

// New crea una nueva aplicación.
func New(cfg *config.Config, logger observability.Logger) *App {
	return &App{
		config: cfg,
		logger: logger,
	}
}

// Run inicia la aplicación.
func (a *App) Run() error {
	a.logger.Info("starting_server",
		observability.F("port", a.config.ServerPort),
		observability.F("phone_number_id", a.config.WhatsAppPhoneNumberID),
	)

	// Create handler with real WhatsApp client
	handler := webhook.NewHandler(a.config, a.logger)

	// Create HTTP server
	a.server = &http.Server{
		Addr:    ":" + a.config.ServerPort,
		Handler: handler,
	}

	// Start server
	return a.server.ListenAndServe()
}

// Stop detiene la aplicación gracefully.
func (a *App) Stop() error {
	a.logger.Info("stopping_server")
	return a.server.Close()
}
