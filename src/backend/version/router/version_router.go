package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/version/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Get("/", handler.GetVersions)
		r.Post("/", handler.PostVersion)
		r.Get("/entry/{id}", handler.GetVersionsByEntryID)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetVersionByID)
			r.Put("/", handler.PutVersion)
			r.Delete("/", handler.DeleteVersion)
		})
	})

	return r
}
