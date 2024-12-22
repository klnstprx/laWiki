package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/auth/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetUsers)
		r.Post("/", handler.PostUser)

		r.Route("/user", func(r chi.Router) {
			r.Get("/", handler.GetUserByID)
			r.Put("/", handler.PutUser)
			r.Delete("/", handler.DeleteUser)
		})
	})

	return r
}
