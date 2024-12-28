package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/auth/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/login", handler.Login)
		r.Get("/me", handler.UserInfo)
		r.Get("/callback", handler.Callback)
		r.Get("/logout", handler.Logout)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", handler.GetUsers)
			r.Post("/", handler.PostUser)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetUserByID)
				r.Put("/", handler.PutUser)
				r.Delete("/", handler.DeleteUser)
				r.Route("/notifications", func(r chi.Router) {
					r.Post("/", handler.AddUserNotification)
				})
			})
		})
	})

	return r
}
