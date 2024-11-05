package config

import (
	"io"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Logger     *zerolog.Logger
	Port       string
	PrettyLogs bool
	Debug      bool
}

// App holds app configuration
var App AppConfig

// Creates global AppConfig
func New() {
	App = AppConfig{}
}

func (cfg *AppConfig) LoadEnv() {
	_, err := os.Stat(".env")
	if err != nil {
		log.Logger.Warn().Msg("No .env file found, setting defaults.")
		godotenv.Load(".env.default")
	} else {
		godotenv.Load()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	cfg.Port = ":" + port
	cfg.PrettyLogs = os.Getenv("PRETTY_LOGS") == "true"
	cfg.Debug = os.Getenv("DEBUG") == "true"
}

// Setups pretty logs and debug level
func SetupLogger(prettylogs bool, debug bool) {
	var writers []io.Writer
	if prettylogs {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		writers = append(writers, os.Stderr)
	}
	finalWriter := io.MultiWriter(writers...)
	log.Logger = zerolog.New(finalWriter).With().Timestamp().Logger()
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	App.Logger = &log.Logger
}
