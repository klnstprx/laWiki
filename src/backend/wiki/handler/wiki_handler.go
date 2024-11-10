package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/laWiki/wiki/config"
	"github.com/laWiki/wiki/database"
	"github.com/laWiki/wiki/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
* GET /health
* checks if the service is up
 */
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

/*
* POST /
* creates a new wiki
* expects a json object in the request body
 */
func PostWiki(w http.ResponseWriter, r *http.Request) {
	var wiki model.Wiki
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.WikiCollection.InsertOne(ctx, wiki)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ensure InsertedID is of type primitive.ObjectID
	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		config.App.Logger.Error().Msg("Failed to convert InsertedID to ObjectID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	wiki.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Return 201 Created
	if err := json.NewEncoder(w).Encode(wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("wiki", wiki).Msg("Added new wiki")
}

/*
* GET /
* gets the list of all wiki json objects from db
 */
func GetWikis(w http.ResponseWriter, r *http.Request) {
	var wikis []model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.WikiCollection.Find(ctx, bson.M{})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var wiki model.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode wiki")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		wikis = append(wikis, wiki)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

/*
* GET /{id}
* gets a wiki by id from db
 */
func GetWikiByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var wiki model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Wiki not found")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

/*
* PUT /{id}
* updates a wiki by id
* expects a json object in the request
 */
func PutWiki(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var wiki model.Wiki
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       wiki.Title,
			"description": wiki.Description,
			"category":    wiki.Category,
		},
	}

	result, err := database.WikiCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Wiki not found for update")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	// Retrieve the updated document (optional)
	err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated wiki")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

/*
* DELETE /{id}
* deletes a wiki
 */
func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	wikiID := r.URL.Query().Get("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete associated entries first
	entryServiceURL := fmt.Sprintf("%s/api/entries/wiki?wikiID=%s", config.App.API_GATEWAY_URL, wikiID)
	config.App.Logger.Info().Str("url", entryServiceURL).Msg("Preparing to delete associated entries")

	req, err := http.NewRequest("DELETE", entryServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to entry service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Str("url", entryServiceURL).Msg("Sending request to delete associated entries")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to entry service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		config.App.Logger.Error().
			Int("status", resp.StatusCode).
			Str("body", bodyString).
			Msg("Entry service returned error")
		http.Error(w, "Failed to delete associated entries", http.StatusInternalServerError)
		return
	}

	// Now proceed to delete the wiki document
	objID, err := primitive.ObjectIDFromHex(wikiID)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid wiki ID")
		http.Error(w, "Invalid wiki ID", http.StatusBadRequest)
		return
	}

	result, err := database.WikiCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete wiki")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Wiki not found")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	config.App.Logger.Info().Str("wikiID", wikiID).Msg("Wiki and associated entries deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func GetWikisByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var wikis []model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.WikiCollection.Find(ctx, bson.M{"title": title})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var wiki model.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode wiki")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		wikis = append(wikis, wiki)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetWikisByDescription(w http.ResponseWriter, r *http.Request) {
	description := r.URL.Query().Get("description")

	var wikis []model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.WikiCollection.Find(ctx, bson.M{"description": description})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var wiki model.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode wiki")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		wikis = append(wikis, wiki)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetWikisByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var wikis []model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.WikiCollection.Find(ctx, bson.M{"category": category})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var wiki model.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode wiki")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		wikis = append(wikis, wiki)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
