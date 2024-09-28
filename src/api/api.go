package api

import (
	"fmt"
	"lawiki/api/wikis"
	"lawiki/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// xlog is a global logger for the api package
var xlog zerolog.Logger

func Routes() *chi.Mux {
	xlog = config.App.Logger.With().Str("component", "api").Logger()
	router := chi.NewRouter()
	router.Mount("/wikis", wikis.Routes()) // Mount the wikis routes

	return router
}

/*
* JsonError is a helper function to log errors and return a JSON response
* with the error message and status code
 */
func JsonError(w http.ResponseWriter, r *http.Request, errorMsg string, code int) {
	xlog.Error().Ctx(r.Context()).Str("error", errorMsg).Int("code", code).Str("URI", r.RequestURI).Msg("JSON API Error")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, `{"message": "`+errorMsg+`"}`)
}
