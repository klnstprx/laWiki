package wikis

import (
	"encoding/json"
	"fmt"
	"lawiki/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Wiki struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Catergory   string `json:"category"`
	URL         string `json:"url"`
}

// xlog is a global logger for the api package
var xlog zerolog.Logger

// Routes returns a router that will be used as a subrouter of the main router
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
	// placeholder
	var wikis []Wiki
	for i := 0; i < 10; i++ {
		wiki := Wiki{
			ID:          fmt.Sprint(i),
			Title:       "First Wiki",
			Description: "This is the first wiki",
			Content:     "This is the content of the first wiki",
			Catergory:   "General",
			URL:         "https://lawiki.net/first-wiki",
		}
		wikis = append(wikis, wiki)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonWikis, _ := json.Marshal(wikis)
	w.Write(jsonWikis)
}

/*
* GET /api/wikis/{id}
* gets a wiki by id from db
 */
func getWiki(w http.ResponseWriter, r *http.Request) {
	xlog.Info().Msg("getWiki")
	// placeholder
	wiki := Wiki{
		ID:          "1",
		Title:       "First Wiki",
		Description: "This is the first wiki",
		Content:     "This is the content of the first wiki",
		Catergory:   "General",
		URL:         "https://lawiki.net/first-wiki",
	}
	w.Header().Set("Content-Type", "application/json")
	jsonWiki, _ := json.Marshal(wiki)
	w.Write(jsonWiki)
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
