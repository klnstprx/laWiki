package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MediaConfig holds the configuration specific to the Media service
type MediaConfig struct {
	Port                int    `toml:"PORT"`
	CLOUDIFY_CLOUD_NAME string `toml:"CLOUDIFY_CLOUD_NAME"`
	CLOUDIFY_API_KEY    string `toml:"CLOUDIFY_API_KEY"`
	CLOUDIFY_API_SECRET string `toml:"CLOUDIFY_API_SECRET"`
	MongoDBURI          string `toml:"MONGODB_URI"`
	DBCollectionName    string `toml:"DB_COLLECTION_NAME"`
	DBName              string `toml:"DB_NAME"`
	PrettyLogs          *bool  `toml:"PRETTY_LOGS"`
	Debug               *bool  `toml:"DEBUG"`
	MB_LIMIT            int64  `toml:"MB_LIMIT"`
}

// Config represents the structure of the config.toml file
type Config struct {
	Media MediaConfig `toml:"Media"`
}
type AppConfig struct {
	Logger           *zerolog.Logger
	Cld              *cloudinary.Cloudinary
	Port             string
	PrettyLogs       bool
	Debug            bool
	MongoDBURI       string
	DBCollectionName string
	DBName           string
	MB_LIMIT         int64
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

	// PORT with default value
	if config.Media.Port == 0 {
		cfg.Port = ":8002" // Default port
		log.Warn().Msg("PORT not set in config file. Using default ':8002'.")
	} else {
		cfg.Port = fmt.Sprintf(":%d", config.Media.Port)
	}

	// PRETTY_LOGS with default value
	if config.Media.PrettyLogs != nil {
		cfg.PrettyLogs = *config.Media.PrettyLogs
	} else {
		cfg.PrettyLogs = true // Default to true
		log.Warn().Msg("PRETTY_LOGS not set in config file. Using default 'true'.")
	}

	// DEBUG with default value
	if config.Media.Debug != nil {
		cfg.Debug = *config.Media.Debug
	} else {
		cfg.Debug = true // Default to true
		log.Warn().Msg("DEBUG not set in config file. Using default 'true'.")
	}

	// DBNAME with default value
	if config.Media.DBName != "" {
		cfg.DBName = config.Media.DBName
	} else {
		cfg.DBName = "laWiki" // Default to "laWiki"
		log.Warn().Msg("DBNAME not set in config file. Using default 'laWiki'.")
	}
	// DBCOLLECTIONNAME with default value
	if config.Media.DBCollectionName != "" {
		cfg.DBCollectionName = config.Media.DBCollectionName
	} else {
		cfg.DBCollectionName = "versiones" // Default to "wikis"
		log.Warn().Msg("DBCOLLECTIONNAME not set in config file. Using default 'wiki'.")
	}

	// MONGODB_URI is required
	if config.Media.MongoDBURI != "" {
		cfg.MongoDBURI = config.Media.MongoDBURI
	} else {
		cfg.MongoDBURI = "mongodb://localhost:27017" // Default to locally hosted DB
		log.Warn().Msg("DMONGODB_URI not set in config file. Using default 'mongodb://localhost:27017'.")
	}

	// Initialize Cloudinary
	if config.Media.CLOUDIFY_CLOUD_NAME != "" && config.Media.CLOUDIFY_API_KEY != "" && config.Media.CLOUDIFY_API_SECRET != "" {
		var err error
		cfg.Cld, err = cloudinary.NewFromParams(config.Media.CLOUDIFY_CLOUD_NAME, config.Media.CLOUDIFY_API_KEY, config.Media.CLOUDIFY_API_SECRET)
		if err != nil {
			log.Error().Msgf("Failed to initialize Cloudinary: %v", err)
		}
	} else {

		// CLOUDIFY_CLOUD_NAME is required
		if config.Media.CLOUDIFY_CLOUD_NAME == "" {
			missingVars = append(missingVars, "CLOUDIFY_CLOUD_NAME")
			log.Warn().Msg("CLOUDIFY_CLOUD_NAME not set in config file.")
		}
		// CLOUDIFY_API_KEY is required
		if config.Media.CLOUDIFY_API_KEY == "" {
			missingVars = append(missingVars, "CLOUDIFY_API_KEY")
			log.Warn().Msg("CLOUDIFY_API_KEY not set in config file.")
		}
		// CLOUDIFY_API_SECRET is required
		if config.Media.CLOUDIFY_API_SECRET == "" {
			missingVars = append(missingVars, "CLOUDIFY_API_SECRET")
			log.Warn().Msg("CLOUDIFY_API_SECRET not set in config file.")
		}
	}

	if config.Media.MB_LIMIT == 0 {
		cfg.MB_LIMIT = 5 // Default to 5 MB
		log.Warn().Msg("MB_LIMIT not set in config file. Using default '5'.")
	} else {
		cfg.MB_LIMIT = config.Media.MB_LIMIT
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
