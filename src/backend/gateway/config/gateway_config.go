package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GatewayConfig holds the configuration specific to the gateway service
type GatewayConfig struct {
	Port              int    `toml:"PORT"`
	WikiServiceURL    string `toml:"WIKI_SERVICE_URL"`
	EntryServiceURL   string `toml:"ENTRY_SERVICE_URL"`
	AuthServiceURL    string `toml:"AUTH_SERVICE_URL"`
	VersionServiceURL string `toml:"VERSION_SERVICE_URL"`
	PrettyLogs        *bool  `toml:"PRETTY_LOGS"`
	Debug             *bool  `toml:"DEBUG"`
}

// Config represents the structure of the config.toml file
type Config struct {
	Gateway GatewayConfig `toml:"gateway"`
}

type AppConfig struct {
	Logger            *zerolog.Logger
	Port              string
	PrettyLogs        bool
	Debug             bool
	WikiServiceURL    string
	EntryServiceURL   string
	AuthServiceURL    string
	VersionServiceURL string
}

// App holds app configuration
var App AppConfig

// Creates global AppConfig
func New() {
	App = AppConfig{}
}

// LoadConfig reads the configuration from config.toml and populates AppConfig
func (cfg *AppConfig) LoadConfig(configPath string) {
	var config Config
	// Check if the config.toml file exists
	_, err := os.Stat(configPath)
	if err != nil {
		log.Error().Msgf("Config file '%s' not found.", configPath)
		os.Exit(1)
	}

	// Decode the TOML file into the Config struct
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Error().Err(err).Msg("Error decoding config file.")
		os.Exit(1)
	}

	missingVars := []string{}

	// PORT with default value
	if config.Gateway.Port == 0 {
		cfg.Port = ":8000" // Default port
		log.Warn().Msg("PORT not set in config file. Using default ':8000'.")
	} else {
		cfg.Port = fmt.Sprintf(":%d", config.Gateway.Port)
	}

	// PRETTY_LOGS with default value
	if config.Gateway.PrettyLogs != nil {
		cfg.PrettyLogs = *config.Gateway.PrettyLogs
	} else {
		cfg.PrettyLogs = true // Default to true
		log.Warn().Msg("PRETTY_LOGS not set in config file. Using default 'true'.")
	}

	// DEBUG with default value
	if config.Gateway.Debug != nil {
		cfg.Debug = *config.Gateway.Debug
	} else {
		cfg.Debug = true // Default to true
		log.Warn().Msg("DEBUG not set in config file. Using default 'true'.")
	}

	// WIKI_SERVICE_URL is required
	if config.Gateway.WikiServiceURL == "" {
		missingVars = append(missingVars, "WIKI_SERVICE_URL")
	} else {
		cfg.WikiServiceURL = config.Gateway.WikiServiceURL
	}

	// ENTRY_SERVICE_URL is required
	if config.Gateway.EntryServiceURL == "" {
		missingVars = append(missingVars, "ENTRY_SERVICE_URL")
	} else {
		cfg.EntryServiceURL = config.Gateway.EntryServiceURL
	}

	// AUTH_SERVICE_URL is required
	if config.Gateway.AuthServiceURL == "" {
		missingVars = append(missingVars, "AUTH_SERVICE_URL")
	} else {
		cfg.AuthServiceURL = config.Gateway.AuthServiceURL
	}

	// VERSION_SERVICE_URL is required
	if config.Gateway.VersionServiceURL == "" {
		missingVars = append(missingVars, "ENTRY_SERVICE_URL")
	} else {
		cfg.VersionServiceURL = config.Gateway.VersionServiceURL
	}

	// If there are missing required variables, log them and exit
	if len(missingVars) > 0 {
		for _, v := range missingVars {
			log.Error().Msgf("Missing required configuration variable: %s", v)
		}
		os.Exit(1)
	}
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
