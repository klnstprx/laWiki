package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laWiki/translation/config"
	"github.com/laWiki/translation/router"

	"github.com/rs/zerolog/log"
)

/* THIS SERVICE IS WIP - its not included in the gateway router*/

// @title           Auth Service API
// @version         1.0
// @description     API documentation for the Auth Service. !!THIS SERVICE IS WIP!!

// @host            localhost:8080
// @BasePath        /api/auth

// @securityDefinitions.oauth2.accessCode OAuth2
// @tokenUrl         https://oauth2.googleapis.com/token
// @authorizationUrl https://accounts.google.com/o/oauth2/auth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	xlog := config.App.Logger.With().Str("service", "auth").Logger()

	// router setup
	r := router.NewRouter()

	// context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown logic
	// esto no es muy importante, pero es bueno tenerlo
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

	// server starts here
	// starts in a go routine so it doesn't block the main thread
	httpServer := http.Server{
		Addr:    config.App.Port,
		Handler: r,
	}

	// lo arrancamos en un go routine para que no bloquee el main thread
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			xlog.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()
	xlog.Info().Str("Port", config.App.Port).Msg("HTTP server started")
	// Block until context is canceled (waiting for the shutdown signal).
	<-ctx.Done()
	// Shutdown logic
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		xlog.Error().Err(err).Msg("HTTP server failed to shutdown")
	}
	xlog.Info().Msg("Server shut down successfully")
}
