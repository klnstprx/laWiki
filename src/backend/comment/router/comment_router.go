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
		r.Get("/content", handler.GetCommentsByContent)
		r.Get("/rating", handler.GetCommentsByRating)
		r.Get("/date", handler.GetCommentsByDate)
		r.Get("/author", handler.GetCommentsByAuthor)
		r.Get("/version", handler.GetCommentsByVersionID)
		r.Delete("/version", handler.DeleteCommentsByVersionID)

		r.Route("/id", func(r chi.Router) {
			r.Get("/", handler.GetCommentByID)
			r.Put("/", handler.PutComment)
			r.Delete("/", handler.DeleteComment)
		})
	})

	return r
}
