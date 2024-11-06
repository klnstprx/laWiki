package main

import (
	"crypto/internal/boring/sig"
	"net/http"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/entry/config"
	"github.com/rs/zerolog/log"
)

// main is the entrypoint for the entry service
func main() {
	config.New()
	config.App.LoadEnv()
	config.SetupLogger(config.App.PrettyLogs, config.App.Debug)
	config.App.Logger = &log.Logger
	xlog := config.App.Logger.With().Str("service", "entry").Logger()

	xlog.Info().Msg("Connecting to the database...")
	database.Connect()

	r := chi.NewRouter()
	r.Mount(/entries, router.NewRouter())

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown logic
	signalCaught := false;
	signalChannel := make(chan os.signal,1)
	signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)
	go func(){
		<-signalChannel
		if signalCaught {
			xlog.Warn().Msg("Caught second signal, terminating immediately")
			os.Exit(1)
		}
		signalCaught = true
		xlog.Info().Msg("Caught shutdown signal")
		cancel()
	}()


	//server starup
	httpServer := http.Server{
		Addr: config.App.Port
		Handler: r,
	}

	go func() {
		err:=httpServer.ListenAndServe()
		if err !=nil&&errores.Is(err,http.ErrServerClosed){
			xlog.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}
	xlog.Info().Str("port",config.App.Port).Msg("HTTP Server started")

	//wait for shutdown signal
	<-ctx.Done()

	//shutdown logic
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err:=httpServer.Shutdown(shutdownCtx);err!=nil{
		xlog.Fatal().Err(err).Msg("Failed to gracefully shutdown HTTP server")
	}
	xlog.Info().Msg("HTTP server shut down successfully")
}
