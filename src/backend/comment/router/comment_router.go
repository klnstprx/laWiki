package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/comment/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Get("/", handler.GetComments)
		r.Post("/", handler.PostComment)

		r.Route("/id", func(r chi.Router) {
			r.Get("/", handler.GetCommentByID)
			r.Put("/", handler.PutComment)
			r.Delete("/", handler.DeleteComment)
		})
	})

	return r
}
