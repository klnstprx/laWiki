package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/laWiki/gateway/config"
	"github.com/laWiki/gateway/handler"
	custommw "github.com/laWiki/gateway/middleware" /*custom name cuz "middleware" would be redeclared*/
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Recoverer)
	// Custom middleware for authentication, etc.
	r.Use(custommw.RequestID)
	r.Use(custommw.LoggerMiddleware(config.App.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://lawiki.mooo.com"}, // Use this to allow specific origin hosts
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// Health Check
	r.Get("/health", handler.HealthCheck)

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Specifies the combined Swagger JSON
	))

	// Serve the Swagger JSON file
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json") // Ensure the path is correct
	})

	// Define routes to backend services
	// aqui anadimos (con r.Mount()) cada microservico al gateway.
	r.Route("/api", func(r chi.Router) {
		// Wiki Service Routes
		r.Mount("/wikis", proxyHandler(config.App.WikiServiceURL, "/api/wikis"))

		// Entry Service Routes
		r.Mount("/entries", proxyHandler(config.App.EntryServiceURL, "/api/entries"))

		// Comment Service Routes
		r.Mount("/comments", proxyHandler(config.App.CommentServiceURL, "/api/comments"))

		// Auth Service Routes
		r.Mount("/auth", proxyHandler(config.App.AuthServiceURL, "/api/auth"))

		// Version Service Routes
		r.Mount("/versions", proxyHandler(config.App.VersionServiceURL, "/api/versions"))

		// Media Service Routes
		r.Mount("/media", proxyHandler(config.App.MediaServiceURL, "/api/media"))
	})

	return r
}

// proxyHandler returns a handler that proxies requests to the given service
func proxyHandler(serviceURL string, prefixToStrip string) http.HandlerFunc {
	return handler.ReverseProxy(serviceURL, prefixToStrip)
}
