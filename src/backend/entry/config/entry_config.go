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

// GlobalConfig holds the configuration for the application
type GlobalConfig struct {
	API_GATEWAY_URL string `toml:"API_GATEWAY_URL"`
}

// EntryConfig holds the configuration specific to the entry service
type EntryConfig struct {
	Port             int    `toml:"PORT"`
	MongoDBURI       string `toml:"MONGODB_URI"`
	DBCollectionName string `toml:"DB_COLLECTION_NAME"`
	DBName           string `toml:"DB_NAME"`
	PrettyLogs       *bool  `toml:"PRETTY_LOGS"`
	Debug            *bool  `toml:"DEBUG"`
}

// Config represents the structure of the config.toml file
type Config struct {
	Entry  EntryConfig  `toml:"entry"`
	Global GlobalConfig `toml:"global"`
}
type AppConfig struct {
	Logger           *zerolog.Logger
	Port             string
	PrettyLogs       bool
	Debug            bool
	MongoDBURI       string
	DBCollectionName string
	DBName           string
	API_GATEWAY_URL  string
}

// App holds app configuration
var App AppConfig

// Creates global AppConfig
func New() {
	App = AppConfig{}
}

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

	if config.Global.ApiGatewayURL == "" {
		missingVars = append(missingVars, "API_GATEWAY_URL")
	} else {
		cfg.ApiGatewayURL = config.Global.ApiGatewayURL
	}

	// PORT with default value
	if config.Entry.Port == 0 {
		cfg.Port = ":8002" // Default port
		log.Warn().Msg("PORT not set in config file. Using default ':8002'.")
	} else {
		cfg.Port = fmt.Sprintf(":%d", config.Entry.Port)
	}

	// PRETTY_LOGS with default value
	if config.Entry.PrettyLogs != nil {
		cfg.PrettyLogs = *config.Entry.PrettyLogs
	} else {
		cfg.PrettyLogs = true // Default to true
		log.Warn().Msg("PRETTY_LOGS not set in config file. Using default 'true'.")
	}

	// DEBUG with default value
	if config.Entry.Debug != nil {
		cfg.Debug = *config.Entry.Debug
	} else {
		cfg.Debug = true // Default to true
		log.Warn().Msg("DEBUG not set in config file. Using default 'true'.")
	}

	// DBNAME with default value
	if config.Entry.DBName != "" {
		cfg.DBName = config.Entry.DBName
	} else {
		cfg.DBName = "laWiki" // Default to "laWiki"
		log.Warn().Msg("DBNAME not set in config file. Using default 'laWiki'.")
	}
	// DBCOLLECTIONNAME with default value
	if config.Entry.DBCollectionName != "" {
		cfg.DBCollectionName = config.Entry.DBCollectionName
	} else {
		cfg.DBCollectionName = "versiones" // Default to "wikis"
		log.Warn().Msg("DBCOLLECTIONNAME not set in config file. Using default 'wiki'.")
	}

	// MONGODB_URI is required
	if config.Entry.MongoDBURI != "" {
		cfg.MongoDBURI = config.Entry.MongoDBURI
	} else {
		cfg.MongoDBURI = "mongodb://localhost:27017" // Default to locally hosted DB
		log.Warn().Msg("DMONGODB_URI not set in config file. Using default 'mongodb://localhost:27017'.")
	}

	// API_GATEWAY_URL is required
	if config.Global.API_GATEWAY_URL != "" {
		cfg.API_GATEWAY_URL = config.Global.API_GATEWAY_URL
	} else {
		missingVars = append(missingVars, "API_GATEWAY_URL")
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
