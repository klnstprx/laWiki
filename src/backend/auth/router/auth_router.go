package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/auth/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Public Routes
	r.Get("/health", handler.HealthCheck)
	r.Get("/login", handler.Login)
	r.Get("/callback", handler.Callback)
	r.Get("/logout", handler.Logout)
	r.Get("/userinfo", handler.UserInfo)

	//Otras rutas para el nuevo manejo de la coleccion ususarios
	r.Get("/users", handler.GetUsers)
	r.Get("/users/{id}", handler.GetUserByID)
	r.Post("/users", handler.PostUser)
	r.Put("/users/{id}", handler.PutUser)

	return r
}
