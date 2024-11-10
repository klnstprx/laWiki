package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
// @Router       /api/wikis/id/ [get]
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
// @Router       /api/wikis/id/ [put]
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

// DeleteWiki godoc
// @Summary      Delete a wiki by ID
// @Description  Deletes a wiki by its ID.
// @Tags         Wikis
// @Param        id    query     string  true  "Wiki ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Wiki not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/wikis/id/ [delete]
func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.WikiCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Wiki not found for deletion")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// GetWikisByTitle godoc
// @Summary      Get all wikis by title.
// @Description  Retrieves the list of all wikis that matches the title.
// @Tags         Wikis
// @Produce      application/json
// @Param        title  query     string  true  "Title to search"
// @Success      200  {array}   model.Wiki
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/wikis/title [get]
func GetWikisByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	var wikis []model.Wiki
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}
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
		config.App.Logger.Info().Str("title", title).Msg("No wikis found")
		http.Error(w, "No wikis found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wikis); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetWikiByExactTitle godoc
// @Summary      Get wiki by exact title
// @Description  Retrieves wiki that matches the given title.
// @Tags         Wikis
// @Produce      application/json
// @Param        title   query     string  true  "Title to search"
// @Success      200     {object}   model.Wiki
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/wikis/exactTitle [get]
func GetWikiByExactTitle(w http.ResponseWriter, r *http.Request) {
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

// GetWikisByDescription godoc
// @Summary      Get wikis by description
// @Description  Retrieves wikis that match the given description.
// @Tags         Wikis
// @Produce      application/json
// @Param        description  query     string  true  "Description to search"
// @Success      200          {array}   model.Wiki
// @Failure      500          {string}  string  "Internal server error"
// @Router       /api/wikis/description [get]
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

// GetWikisByCategory godoc
// @Summary      Get wikis by category
// @Description  Retrieves wikis under the given category.
// @Tags         Wikis
// @Produce      application/json
// @Param        category  query     string  true  "Category to search"
// @Success      200       {array}   model.Wiki
// @Failure      500       {string}  string  "Internal server error"
// @Router       /api/wikis/category [get]
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
