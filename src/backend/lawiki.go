package main

import (
	"context"
	"errors"
	"lawiki/api"
	"lawiki/config"
	"lawiki/utils"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
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
	router.Use(middleware.Recoverer)

	// cors for all origins
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	router.Use(utils.LoggerMiddleware(config.App.Logger))
	router.Use(cors.Handler(corsOptions))

	// mount the api routes
	router.Mount("/api", api.Routes())

	// serve static files
	workDir, _ := os.Getwd()
	buildPath := filepath.Join(workDir, "frontend", "dist")
	FileServer(router, "/", http.Dir(buildPath), xlog)

	// context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown logic
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
	httpServer := &http.Server{
		Addr:    config.App.ListenAddr,
		Handler: router,
	}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			xlog.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()
	xlog.Info().Str("Listen address", config.App.ListenAddr).Msg("HTTP server started")
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

// FileServer serves static files and handles client-side routing
func FileServer(router chi.Router, path string, root http.FileSystem, xlog zerolog.Logger) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	indexFile := "index.html"
	fs := http.FileServer(root)

	// Serve static files
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		requestPath = strings.TrimPrefix(requestPath, path)

		f, err := root.Open(requestPath)
		if err != nil {
			if os.IsNotExist(err) || requestPath == "/" || requestPath == "" {
				f, err = root.Open(indexFile)
				if err != nil {
					xlog.Error().Err(err).Msg("Failed to open index file")
					http.Error(w, "Index file not found", http.StatusNotFound)
				}
				defer f.Close()
				fi, _ := f.Stat()
				http.ServeContent(w, r, indexFile, fi.ModTime(), f)
				return
			}
			// for other errors, return the error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// if the request file exists, serve it
		fs.ServeHTTP(w, r)
	})
}
