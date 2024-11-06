package router

import (
	"net/http"

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
		r.Mount("/wikis", proxyHandler(config.App.WikiServiceURL, "/api/wikis"))

		// Entry Service Routes
		r.Mount("/entries", proxyHandler(config.App.EntryServiceURL, "/api/entries"))

		// Auth Service Routes
		r.Mount("/auth", proxyHandler(config.App.AuthServiceURL, "/api/auth"))

		// Version Service Routes
		r.Mount("/versions", proxyHandler(config.App.VersionServiceURL, "/api/versions"))
	})

	return r
}

// proxyHandler returns a handler that proxies requests to the given service
func proxyHandler(serviceURL string, prefixToStrip string) http.HandlerFunc {
	return handler.ReverseProxy(serviceURL, prefixToStrip)
}
