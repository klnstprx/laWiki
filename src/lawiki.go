package main

import (
	"lawiki/api"
	"lawiki/config"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	// config setup
	config.New()
	configPath := "config.toml"
	if _, err := toml.DecodeFile(configPath, &config.App); err != nil {
		panic("Error reading config file")
	}
	config.SetupLogger(config.App.PrettyLogs, config.App.Debug)
	config.App.LogConfig()
	config.App.Logger = &log.Logger
	xlog := config.App.Logger.With().Str("app", "main").Logger()

	// router setup
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	router.Mount("/gui", api.Routes())

	// server starts here
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		xlog.Fatal().Err(err).Msg("Failed to start server")
	}
}
