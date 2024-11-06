package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/version/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetVersions)
		r.Post("/", handler.PostVersion)

		/*
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetWikiByID)
				r.Put("/", handler.PutWiki)
				r.Delete("/", handler.DeleteWiki)
			})
		*/
	})

	return r
}
