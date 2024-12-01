package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/wiki/config"
	"github.com/laWiki/wiki/database"
	"github.com/laWiki/wiki/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the service is up
// @Tags         Health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/wikis/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// GetWikis godoc
// @Summary      Get all wikis
// @Description  Retrieves the list of all wiki JSON objects from the database.
// @Tags         Wikis
// @Produce      application/json
// @Success      200  {array}   model.Wiki
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/wikis/ [get]
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

// GetWikiByID godoc
// @Summary      Get a wiki by ID
// @Description  Retrieves a wiki by its ID.
// @Tags         Wikis
// @Produce      application/json
// @Param        id    query     string  true  "Wiki ID"
// @Success      200   {object}  model.Wiki
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Wiki not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/wikis/{id} [get]
func GetWikiByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// SearchWikis godoc
// @Summary      Search wikis
// @Description  Search for wikis using various query parameters. You can search by title, exact_title, description, or category. All parameters are optional and can be combined.
// @Tags         Wikis
// @Produce      application/json
// @Param        title        query     string  false  "Partial title to search for (case-insensitive)"
// @Param        exact_title  query     string  false  "Exact title to search for"
// @Param        description  query     string  false  "Description to search for (case-insensitive)"
// @Param        category     query     string  false  "Category to search for"
// @Success      200          {array}   model.Wiki
// @Failure      400          {string}  string  "Bad Request"
// @Failure      500          {string}  string  "Internal Server Error"
// @Router       /api/wikis/search [get]
func SearchWikis(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	title := r.URL.Query().Get("title")
	exactTitle := r.URL.Query().Get("exact_title")
	description := r.URL.Query().Get("description")
	category := r.URL.Query().Get("category")
	createdAtFromString := r.URL.Query().Get("createdAtFrom")
	createdAtToString := r.URL.Query().Get("createdAtTo")

	// Build the MongoDB filter dynamically
	filter := bson.M{}
	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if exactTitle != "" {
		filter["title"] = exactTitle
	}
	if description != "" {
		filter["description"] = bson.M{"$regex": description, "$options": "i"}
	}
	if category != "" {
		filter["category"] = category
	}

	if createdAtFromString != "" || createdAtToString != "" {
		dateFilter := bson.M{}

		if createdAtFromString != "" {
			createdAtFrom, err := time.Parse(time.RFC3339, createdAtFromString)
			if err != nil {
				config.App.Logger.Error().Err(err).Msg("Invalid 'createdAtFrom' date format. Expected ISO8601 format.")
				http.Error(w, "Invalid 'createdAtFrom' date format. Expected ISO8601 format.", http.StatusBadRequest)
				return
			}
			dateFilter["$gte"] = createdAtFrom
		}

		if createdAtToString != "" {
			createdAtTo, err := time.Parse(time.RFC3339, createdAtToString)
			if err != nil {
				config.App.Logger.Error().Err(err).Msg("Invalid 'createdAtTo' date format. Expected ISO8601 format.")
				http.Error(w, "Invalid 'createdAtTo' date format. Expected ISO8601 format.", http.StatusBadRequest)
				return
			}
			dateFilter["$lte"] = createdAtTo
		}

		filter["created_at"] = dateFilter
	}

	// Query the database
	var wikis []model.Wiki
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := database.WikiCollection.Find(ctx, filter)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var wiki model.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
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
	if len(wikis) == 0 {
		config.App.Logger.Info().Str("title", title).Str("exact_title", exactTitle).Str("description", description).Str("category", category).Msg("No wikis found")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// PostWiki godoc
// @Summary      Create a new wiki
// @Description  Creates a new wiki. Expects a JSON object in the request body.
// @Tags         Wikis
// @Accept       application/json
// @Produce      application/json
// @Param        wiki  body      model.Wiki  true  "Wiki information"
// @Success      201   {object}  model.Wiki
// @Failure      400   {string}  string  "Invalid request body"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/wikis/ [post]
func PostWiki(w http.ResponseWriter, r *http.Request) {
	var wiki model.Wiki
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wiki); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wiki.CreatedAt = time.Now().UTC()

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

// PutWiki godoc
// @Summary      Update a wiki by ID
// @Description  Updates a wiki by its ID. Expects a JSON object in the request.
// @Tags         Wikis
// @Accept       application/json
// @Produce      application/json
// @Param        id    query     string  true  "Wiki ID"
// @Param        wiki  body      model.Wiki  true  "Updated wiki information"
// @Success      200   {object}  model.Wiki
// @Failure      400   {string}  string  "Invalid ID or request body"
// @Failure      404   {string}  string  "Wiki not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/wikis/{id} [put]
func PutWiki(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

	wiki.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       wiki.Title,
			"description": wiki.Description,
			"category":    wiki.Category,
			"updated_at":  wiki.UpdatedAt,
			"media_id":    wiki.MediaID,
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
		w.WriteHeader(http.StatusNoContent)
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

// DeleteWiki godoc
// @Summary      Delete a wiki by ID
// @Description  Deletes a wiki by its ID.
// @Tags         Wikis
// @Param        id    query     string  true  "Wiki ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Wiki not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/wikis/{id} [delete]
func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	wikiID := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	//get the media ID, found in the wiki object
	var wiki model.Wiki
	objID, err := primitive.ObjectIDFromHex(wikiID)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve wiki")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete associated media first
	if wiki.MediaID != "" {
		mediaServiceURL := fmt.Sprintf("%s/api/media/%s", config.App.API_GATEWAY_URL, wiki.MediaID)
		config.App.Logger.Info().Str("url", mediaServiceURL).Msg("Preparing to delete associated media")

		req, err := http.NewRequest("DELETE", mediaServiceURL, nil)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to create request to media service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to send request to media service")
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
				Msg("Media service returned error")
			http.Error(w, "Failed to delete associated media", http.StatusInternalServerError)
			return
		}
	}

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

	result, err := database.WikiCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete wiki")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Wiki not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Str("wikiID", wikiID).Msg("Wiki and associated entries deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}
