package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GlobalConfig holds the configuration for the application
type GlobalConfig struct {
	PrettyLogs    *bool  `toml:"PRETTY_LOGS"`
	Debug         *bool  `toml:"DEBUG"`
	JWTSecret     string `toml:"JWT_SECRET"`
	FrontendURL   string `toml:"FRONTEND_URL"`
	ApiGatewayURL string `toml:"API_GATEWAY_URL"`
	MongoDBURI    string `toml:"MONGODB_URI"`
	DBName        string `toml:"DB_NAME"`
	UseHTTPS      bool   `toml:"USE_HTTPS"`
}

// AuthConfig holds the configuration specific to the auth service
type AuthConfig struct {
	Port                    int    `toml:"PORT"`
	GoogleOAuthClientID     string `toml:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `toml:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURL  string `toml:"GOOGLE_OAUTH_REDIRECT_URL"`
	DBCollectionName        string `toml:"DB_COLLECTION_NAME"`
}

// Config represents the structure of the config.toml file
type Config struct {
	Auth   AuthConfig   `toml:"auth"`
	Global GlobalConfig `toml:"global"`
}

type AppConfig struct {
	Logger      *zerolog.Logger
	Port        string
	PrettyLogs  bool
	Debug       bool
	FrontendURL string

	GoogleOAuthConfig *oauth2.Config
	JWTSecret         string
	UseHTTPS          bool

	MongoDBURI       string
	DBCollectionName string
	DBName           string
	ApiGatewayURL    string
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
	if config.Auth.Port == 0 {
		cfg.Port = ":8002" // Default port
		log.Warn().Msg("PORT not set in config file. Using default ':8002'.")
	} else {
		cfg.Port = fmt.Sprintf(":%d", config.Auth.Port)
	}

	// PRETTY_LOGS with default value
	if config.Global.PrettyLogs != nil {
		cfg.PrettyLogs = *config.Global.PrettyLogs
	} else {
		cfg.PrettyLogs = true // Default to true
		log.Warn().Msg("PRETTY_LOGS not set in config file. Using default 'true'.")
	}

	// DEBUG with default value
	if config.Global.Debug != nil {
		cfg.Debug = *config.Global.Debug
	} else {
		cfg.Debug = true // Default to true
		log.Warn().Msg("DEBUG not set in config file. Using default 'true'.")
	}

	// DBNAME with default value
	if config.Global.DBName != "" {
		cfg.DBName = config.Global.DBName
	} else {
		cfg.DBName = "laWiki" // Default to "laWiki"
		log.Warn().Msg("DBNAME not set in config file. Using default 'laWiki'.")
	}
	// DBCOLLECTIONNAME with default value
	if config.Auth.DBCollectionName != "" {
		cfg.DBCollectionName = config.Auth.DBCollectionName
	} else {
		cfg.DBCollectionName = "usuarios" // Default to "wikis"
		log.Warn().Msg("DBCOLLECTIONNAME not set in config file. Using default 'usuarios'.")
	}

	if config.Global.FrontendURL != "" {
		cfg.FrontendURL = config.Global.FrontendURL
	} else {
		cfg.FrontendURL = "localhost:5173"
		log.Warn().Msg("FRONTEND_URL not set in config file. Using default 'localhost:5173'.")
	}

	// MONGODB_URI is required
	if config.Global.MongoDBURI != "" {
		cfg.MongoDBURI = config.Global.MongoDBURI
	} else {
		cfg.MongoDBURI = "mongodb://localhost:27017" // Default to locally hosted DB
		log.Warn().Msg("DMONGODB_URI not set in config file. Using default 'mongodb://localhost:27017'.")
	}

	// Required variables
	if config.Auth.GoogleOAuthClientID == "" {
		missingVars = append(missingVars, "GOOGLE_OAUTH_CLIENT_ID")
	}
	if config.Auth.GoogleOAuthClientSecret == "" {
		missingVars = append(missingVars, "GOOGLE_OAUTH_CLIENT_SECRET")
	}
	if config.Auth.GoogleOAuthRedirectURL == "" {
		missingVars = append(missingVars, "GOOGLE_OAUTH_REDIRECT_URL")
	}
	if config.Auth.GoogleOAuthRedirectURL == "" {
		missingVars = append(missingVars, "JWT_SECRET")
	}

	// If there are missing required variables, log them and exit
	if len(missingVars) > 0 {
		for _, v := range missingVars {
			log.Error().Msgf("Missing required configuration variable: %s", v)
		}
		os.Exit(1)
	}

	// OAuth2 Configuration
	cfg.GoogleOAuthConfig = &oauth2.Config{
		ClientID:     config.Auth.GoogleOAuthClientID,
		ClientSecret: config.Auth.GoogleOAuthClientSecret,
		RedirectURL:  config.Auth.GoogleOAuthRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "openid", "profile"},
		Endpoint:     google.Endpoint,
	}

	// JWT Secret
	cfg.JWTSecret = config.Global.JWTSecret
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
