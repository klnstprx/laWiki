package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Importamos el paquete gorilla/handlers para manejar CORS
	"github.com/laWiki/gateway/config"
	"github.com/laWiki/gateway/router"
	"github.com/rs/zerolog/log"
)

// @title           API Gateway
// @version         1.0
// @description     Combined API documentation for all services.

// @host            localhost:8000
// @BasePath        /api
func main() {
	// is the service run in docker?
	var configPath string
	if os.Getenv("DOCKER") == "true" {
		configPath = "./config.toml"
	} else {
		configPath = "../config.toml"
	}
	// config setup
	config.New()
	config.App.LoadConfig(configPath)
	config.SetupLogger(config.App.PrettyLogs, config.App.Debug)
	config.App.Logger = &log.Logger
	xlog := config.App.Logger.With().Str("service", "gateway").Logger()

	// r setup
	r := router.NewRouter()

	// CORS Configuration: Allow all origins (replace '*' with specific URLs for production)

	// O puedes restringir a dominios específicos, por ejemplo:
	// corsObj := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	// Contexto para apagado suave
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Lógica de apagado suave
	// Esto no es crucial, pero es útil tenerlo para una salida ordenada
	signalCaught := false
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		if signalCaught {
			xlog.Warn().Msg("Caught second signal, terminating immediately")
			os.Exit(1)
		}
		signalCaught = true
		xlog.Info().Msg("Caught shutdown signal")
		cancel()
	}()

	// Server setup
	httpServer := http.Server{
		Addr:    config.App.Port,
		Handler: r,
	}

	// Iniciamos el servidor HTTP en un go routine para no bloquear el hilo principal
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			xlog.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()
	xlog.Info().Str("Port", config.App.Port).Msg("HTTP server started")

	// Bloqueamos hasta que se reciba una señal de cancelación
	<-ctx.Done()

	// Lógica de apagado del servidor
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		xlog.Error().Err(err).Msg("HTTP server failed to shutdown")
	}
	xlog.Info().Msg("Server shut down successfully")
}
