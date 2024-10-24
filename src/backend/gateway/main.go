package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laWiki/gateway/config"
	"github.com/laWiki/gateway/router"

	"github.com/rs/zerolog/log"
)

func main() {
	// config setup
	config.New()
	config.App.LoadEnv()
	config.SetupLogger(config.App.PrettyLogs, config.App.Debug)
	config.App.Logger = &log.Logger
	xlog := config.App.Logger.With().Str("service", "gateway").Logger()

	// r setup
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
			log.Warn().Msg("Caught second signal, terminating immediately")
			os.Exit(1)
		}
		signalCaught = true
		log.Info().Msg("Caught shutdown signal")
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
