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

// AuthConfig holds the configuration specific to the auth service
type AuthConfig struct {
	Port                    int    `toml:"PORT"`
	GoogleOAuthClientID     string `toml:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `toml:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURL  string `toml:"GOOGLE_OAUTH_REDIRECT_URL"`
	JWTSecret               string `toml:"JWT_SECRET"`
	PrettyLogs              *bool  `toml:"PRETTY_LOGS"`
	Debug                   *bool  `toml:"DEBUG"`
}

// Config represents the structure of the config.toml file
type Config struct {
	Auth AuthConfig `toml:"auth"`
}

type AppConfig struct {
	Logger     *zerolog.Logger
	Port       string
	PrettyLogs bool
	Debug      bool

	GoogleOAuthConfig *oauth2.Config
	JWTSecret         string
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
	if config.Auth.PrettyLogs != nil {
		cfg.PrettyLogs = *config.Auth.PrettyLogs
	} else {
		cfg.PrettyLogs = true // Default to true
		log.Warn().Msg("PRETTY_LOGS not set in config file. Using default 'true'.")
	}

	// DEBUG with default value
	if config.Auth.Debug != nil {
		cfg.Debug = *config.Auth.Debug
	} else {
		cfg.Debug = true // Default to true
		log.Warn().Msg("DEBUG not set in config file. Using default 'true'.")
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
	if config.Auth.JWTSecret == "" {
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
	cfg.JWTSecret = config.Auth.JWTSecret
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
