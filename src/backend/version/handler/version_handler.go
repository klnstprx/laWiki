package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	if len(versions) == 0 {
		config.App.Logger.Info().Msg("No versions found")
		http.Error(w, "No versions found", http.StatusNotFound)
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
	entryID := r.URL.Query().Get("entryID")

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

	if len(versions) == 0 {
		config.App.Logger.Info().Msg("No versions found")
		http.Error(w, "No versions found", http.StatusNotFound)
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

	filter := bson.M{
		"content": bson.M{
			"$regex":   content,
			"$options": "i",
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

	if len(versions) == 0 {
		config.App.Logger.Info().Msg("No versions found")
		http.Error(w, "No versions found", http.StatusNotFound)
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

	if len(versions) == 0 {
		config.App.Logger.Info().Msg("No versions found")
		http.Error(w, "No versions found", http.StatusNotFound)
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

	if len(versions) == 0 {
		config.App.Logger.Info().Msg("No versions found")
		http.Error(w, "No versions found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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
	versionID := r.URL.Query().Get("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete associated comments first
	commentServiceURL := fmt.Sprintf("%s/api/comments/version?versionID=%s", config.App.API_GATEWAY_URL, versionID)
	config.App.Logger.Info().Str("url", commentServiceURL).Msg("Preparing to delete associated comments")

	req, err := http.NewRequest("DELETE", commentServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to comment service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Str("url", commentServiceURL).Msg("Sending request to delete associated comments")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to comment service")
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
			Msg("Comment service returned error")
		http.Error(w, "Failed to delete associated comments", http.StatusInternalServerError)
		return
	}

	// Now proceed to delete the version document
	objID, err := primitive.ObjectIDFromHex(versionID)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid version ID")
		http.Error(w, "Invalid version ID", http.StatusBadRequest)
		return
	}

	result, err := database.VersionCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete version")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Version not found")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	config.App.Logger.Info().Str("versionID", versionID).Msg("Version and associated comments deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

// DeleteVersionsByEntryID godoc
// @Summary      Deletes all versions by the Entry ID
// @Description  Deletes all versions associated with a specific Entry ID.
// @Tags         Versions
// @Param        id    query     string  true  "Entry ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {string}  string  "EntryID is required"
// @Failure      404   {string}  string  "No versions found for the given entry ID"
// @Failure      500   {string}  string  "Internal server error"
// @Failure      500   {string}  string  "Failed to delete associated comments"
// @Router       /api/versions/entry/ [delete]
func DeleteVersionsByEntryID(w http.ResponseWriter, r *http.Request) {
	entryID := r.URL.Query().Get("entryID")
	if entryID == "" {
		config.App.Logger.Warn().Msg("Missing entryID parameter")
		http.Error(w, "EntryID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve all versions associated with the entryID
	versionsCursor, err := database.VersionCollection.Find(ctx, bson.M{"entry_id": entryID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error while fetching versions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer versionsCursor.Close(ctx)

	var versions []model.Version
	if err := versionsCursor.All(ctx, &versions); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode versions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(versions) == 0 {
		config.App.Logger.Info().Str("entryID", entryID).Msg("No versions found for the given entryID")
		http.Error(w, "No versions found for the given entry ID", http.StatusNotFound)
		return
	}

	// Collect all versionIDs
	var versionIDs []string
	for _, version := range versions {
		versionIDs = append(versionIDs, version.ID)
	}

	// Delete associated comments for each versionID
	for _, versionID := range versionIDs {
		commentServiceURL := fmt.Sprintf("%s/api/comments/version?versionID=%s", config.App.API_GATEWAY_URL, versionID)
		req, err := http.NewRequest("DELETE", commentServiceURL, nil)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to create request to comment service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		config.App.Logger.Info().Str("url", commentServiceURL).Msg("Sending delete request to comment service")
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to send delete request to comment service")
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
				Msg("Comment service returned error during deletion")
			http.Error(w, "Failed to delete associated comments", http.StatusInternalServerError)
			return
		}
	}

	config.App.Logger.Info().Str("entryID", entryID).Msg("Associated comments deleted successfully")

	// Delete all versions associated with the entryID
	deleteResult, err := database.VersionCollection.DeleteMany(ctx, bson.M{"entry_id": entryID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete versions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if deleteResult.DeletedCount == 0 {
		config.App.Logger.Info().Str("entryID", entryID).Msg("No versions found to delete for the given entryID")
		http.Error(w, "No versions found to delete for the given entryID", http.StatusNotFound)
		return
	}

	config.App.Logger.Info().Int64("deletedCount", deleteResult.DeletedCount).Str("entryID", entryID).Msg("Versions deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}
