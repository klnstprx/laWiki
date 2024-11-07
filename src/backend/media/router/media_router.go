package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/media/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Get("/", handler.GetMedia)
		r.Post("/", handler.PostMedia)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetMediaByID)
			r.Put("/", handler.PutMedia)
			r.Delete("/", handler.DeleteMedia)
		})
	})

	return r
}
