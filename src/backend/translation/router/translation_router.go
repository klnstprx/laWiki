package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/translation/config"
	"github.com/laWiki/translation/handler"
)

// NewRouter configura las rutas para el servicio de traducción
func NewRouter() http.Handler {
	r := chi.NewRouter()

	translateHandler := handler.NewTranslateHandler(config.LoadConfig())

	r.Route("/", func(r chi.Router) {
		r.Post("/translate", translateHandler.Translate)
	})

	return r
}
