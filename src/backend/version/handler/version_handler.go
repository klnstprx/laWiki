package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/laWiki/version/config"
	"github.com/laWiki/version/database"
	"github.com/laWiki/version/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the service is up
// @Tags         Health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/versions/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// PostVersion godoc
// @Summary      Create a new version
// @Description  Creates a new version. Expects a JSON object in the request body.
// @Tags         Versions
// @Accept       application/json
// @Produce      application/json
// @Param        version  body      model.Version  true  "Version information"
// @Success      201      {object}  model.Version
// @Failure      400      {string}  string  "Invalid request body"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/versions/ [post]
func PostVersion(w http.ResponseWriter, r *http.Request) {
	var version model.Version
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.VersionCollection.InsertOne(ctx, version)
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
	version.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Return 201 Created
	if err := json.NewEncoder(w).Encode(version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("version", version).Msg("Added new version")
}

// GetVersions godoc
// @Summary      Get all versions
// @Description  Retrieves the list of all version JSON objects from the database.
// @Tags         Versions
// @Produce      application/json
// @Success      200  {array}   model.Version
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/versions/ [get]
func GetVersions(w http.ResponseWriter, r *http.Request) {
	var versions []model.Version

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VersionCollection.Find(ctx, bson.M{})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var version model.Version
		if err := cursor.Decode(&version); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode version")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		versions = append(versions, version)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVersionByID godoc
// @Summary      Get a version by ID
// @Description  Retrieves a version by its ID.
// @Tags         Versions
// @Produce      application/json
// @Param        id    query     string  true  "Version ID"
// @Success      200   {object}  model.Version
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Version not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/versions/id/ [get]
func GetVersionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var version model.Version

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.VersionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&version)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Version not found")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// PutVersion godoc
// @Summary      Update a version by ID
// @Description  Updates a version by its ID. Expects a JSON object in the request body.
// @Tags         Versions
// @Accept       application/json
// @Produce      application/json
// @Param        id      query     string          true  "Version ID"
// @Param        version body      model.Version   true  "Updated version information"
// @Success      200     {object}  model.Version
// @Failure      400     {string}  string  "Invalid ID or request body"
// @Failure      404     {string}  string  "Version not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/versions/id/ [put]
func PutVersion(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var version model.Version
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"content":    version.Content,
			"editor":     version.Editor,
			"created_at": version.CreatedAt,
		},
	}

	result, err := database.VersionCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Version not found for update")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	// Retrieve the updated document (optional)
	err = database.VersionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&version)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated version")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// DeleteVersion godoc
// @Summary      Delete a version by ID
// @Description  Deletes a version by its ID.
// @Tags         Versions
// @Param        id query string true "Version ID"
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Version not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/versions/id/ [delete]
func DeleteVersion(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.VersionCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Version not found for deletion")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// GetVersionsByEntryID godoc
// @Summary      Get versions by Entry ID
// @Description  Retrieves versions that correspond to a specific Entry ID.
// @Tags         Versions
// @Produce      application/json
// @Param        entryId query string true "Entry ID"
// @Success      200  {array}   model.Version
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/versions/entry [get]
func GetVersionsByEntryID(w http.ResponseWriter, r *http.Request) {
	var versions []model.Version
	entryID := r.URL.Query().Get("entryId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VersionCollection.Find(ctx, bson.M{"entry_id": entryID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var version model.Version
		if err := cursor.Decode(&version); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode version")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		versions = append(versions, version)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVersionsByContent godoc
// @Summary      Get versions by content
// @Description  Retrieves versions that match the given content.
// @Tags         Versions
// @Produce      application/json
// @Param        content query string true "Content to search"
// @Success      200     {array}   model.Version
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/versions/content [get]
func GetVersionsByContent(w http.ResponseWriter, r *http.Request) {
	var versions []model.Version
	content := r.URL.Query().Get("content")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VersionCollection.Find(ctx, bson.M{"content": content})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var version model.Version
		if err := cursor.Decode(&version); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode version")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		versions = append(versions, version)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVersionsByEditor godoc
// @Summary      Get versions by editor
// @Description  Retrieves versions edited by the given editor.
// @Tags         Versions
// @Produce      application/json
// @Param        editor query string true "Editor to search"
// @Success      200    {array}   model.Version
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/versions/editor [get]
func GetVersionsByEditor(w http.ResponseWriter, r *http.Request) {
	var versions []model.Version
	editor := r.URL.Query().Get("editor")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VersionCollection.Find(ctx, bson.M{"editor": editor})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var version model.Version
		if err := cursor.Decode(&version); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode version")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		versions = append(versions, version)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVersionsByDate godoc
// @Summary      Get versions by date
// @Description  Retrieves versions created on the given date.
// @Tags         Versions
// @Produce      application/json
// @Param        createdAt query string true "Creation date (YYYY-MM-DD)"
// @Success      200       {array}   model.Version
// @Failure      400       {string}  string  "Invalid date format. Expected YYYY-MM-DD"
// @Failure      500       {string}  string  "Internal server error"
// @Router       /api/versions/date [get]
func GetVersionsByDate(w http.ResponseWriter, r *http.Request) {
	createdAtString := r.URL.Query().Get("createdAt")

	// Parse the date (expected format: YYYY-MM-DD)
	createdAt, err := time.Parse("2006-01-02", createdAtString)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid date format. Expected YYYY-MM-DD")
		http.Error(w, "Invalid date format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Define the start and end of the day
	startOfDay := createdAt
	endOfDay := createdAt.AddDate(0, 0, 1)

	var versions []model.Version

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the query to find versions within the date range
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cursor, err := database.VersionCollection.Find(ctx, filter)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var version model.Version
		if err := cursor.Decode(&version); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode version")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		versions = append(versions, version)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
