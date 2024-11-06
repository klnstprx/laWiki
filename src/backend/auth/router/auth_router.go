package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/auth/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/login", handler.Login)

	// Handle callback from Google
	r.Get("/callback", handler.Callback)

	// Protected endpoint example
	r.Get("/protected", handler.ProtectedEndpoint)

	return r
}
