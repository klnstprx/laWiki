package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/laWiki/gateway/config"
	"github.com/laWiki/gateway/handler"
	custommw "github.com/laWiki/gateway/middleware" /*custom name cuz "middleware" would be redeclared*/
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Recoverer)
	// Custom middleware for authentication, etc.
	r.Use(custommw.RequestID)
	r.Use(custommw.LoggerMiddleware(config.App.Logger))

	// Health Check
	r.Get("/health", handler.HealthCheck)

	// Define routes to backend services
	// aqui anadimos (con r.Mount()) cada microservico al gateway.
	r.Route("/api", func(r chi.Router) {
		// Wiki Service Routes
		r.Mount("/wikis", proxyHandler("WIKI_SERVICE_URL", "/api/wikis"))

		// Entry Service Routes
		r.Mount("/entries", proxyHandler("ENTRY_SERVICE_URL", "/api/entries"))

		// Other service routes...
	})

	return r
}

// proxyHandler returns a handler that proxies requests to the given service
func proxyHandler(serviceEnvVar string, prefixToStrip string) http.HandlerFunc {
	serviceURL := os.Getenv(serviceEnvVar)
	if serviceURL == "" {
		config.App.Logger.Panic().Msg("Environment variable " + serviceEnvVar + " not set")
	}
	return handler.ReverseProxy(serviceURL, prefixToStrip)
}
