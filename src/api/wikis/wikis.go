package wikis

import (
	"lawiki/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// xlog is a global logger for the api package
var xlog zerolog.Logger

func Routes() *chi.Mux {
	xlog = config.App.Logger.With().Str("component", "wikis-api").Logger()
	router := chi.NewRouter()
	router.Get("/", getWikis)
	router.Get("/{id}", getWiki)
	router.Post("/", postWiki)
	router.Put("/{id}", putWiki)
	router.Delete("/{id}", deleteWiki)
	return router
}

/*
* GET /api/wikis
* gets the list of all wiki json objects from db
 */
func getWikis(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("getWikis")
}

/*
* GET /api/wikis/{id}
* gets a wiki by id from db
 */
func getWiki(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("getWiki")
}

/*
* POST /api/wikis
* creates a new wiki
* expects a json object in the request body
 */
func postWiki(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("postWiki")
}

/*
* PUT /api/wikis/{id}
* updates a wiki by id
* expects a json object in the request
 */
func putWiki(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("updateWiki")
}

/*
* DELETE /api/wikis/{id}
* deletes a wiki
 */
func deleteWiki(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("deleteWiki")
}
