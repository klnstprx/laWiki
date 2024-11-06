package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/entry/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetEntries)
		r.Post("/", handler.PostEntry)

		// TODO: Implement the following routes:
		/*r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetEntryByID)
			r.Put("/", handler.PutEntry)
			r.Delete("/", handler.DeleteEntry)
		})*/
	})

	return r
}
