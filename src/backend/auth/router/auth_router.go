package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/auth/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.HealthCheck)

	r.Get("/login", handler.Login)

	// Handle callback from Google
	r.Get("/callback", handler.Callback)
	return r
}
