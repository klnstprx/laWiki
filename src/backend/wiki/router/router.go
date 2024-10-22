package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/wiki/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/wikis", func(r chi.Router) {
		r.Get("/", handler.GetWikis)
		r.Post("/", handler.PostWiki)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetWikiByID)
			r.Put("/", handler.PutWiki)
			r.Delete("/", handler.DeleteWiki)
		})
	})

	return r
}
