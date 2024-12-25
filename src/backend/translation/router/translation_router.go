package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/translation/handler"
)

// NewRouter configura las rutas para el servicio de traducción
func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Get("/", handler.Translate)
	})

	return r
}
