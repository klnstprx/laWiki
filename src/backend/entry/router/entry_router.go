package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/entry/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Get("/", handler.GetEntries)
		r.Post("/", handler.PostEntry)
		r.Get("/title", handler.GetEntriesByTitle)
		r.Get("/author", handler.GetEntriesByAuthor)
		r.Get("/date", handler.GetEntriesByDate)
		r.Get("/wiki", handler.GetEntriesByWikiID)
		r.Route("/id", func(r chi.Router) {
			r.Get("/", handler.GetEntryByID)
			r.Put("/", handler.PutEntry)
			r.Delete("/", handler.DeleteEntry)
		})
	})

	return r
}
