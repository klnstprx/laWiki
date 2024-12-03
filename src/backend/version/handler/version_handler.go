package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/version/config"
	"github.com/laWiki/version/database"
	"github.com/laWiki/version/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		w.WriteHeader(http.StatusNoContent)
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
// @Router       /api/versions/{id} [get]
func GetVersionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

// SearchVersions godoc
// @Summary      Search versions
// @Description  Search for versions using various query parameters. You can search by content, editor, createdAt, or entryID. All parameters are optional and can be combined.
// @Tags         Versions
// @Produce      application/json
// @Param        content     query     string  false  "Partial content to search for (case-insensitive)"
// @Param        editor      query     string  false  "Editor to search for"
// @Param        createdAt   query     string  false  "Creation date (YYYY-MM-DD)"
// @Param        entryID     query     string  false  "Entry ID to search for"
// @Success      200         {array}   model.Version
// @Failure      400         {string}  string  "Bad Request"
// @Failure      500         {string}  string  "Internal Server Error"
// @Router       /api/versions/search [get]
func SearchVersions(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")
	editor := r.URL.Query().Get("editor")
	createdAtFromString := r.URL.Query().Get("createdAtFrom")
	createdAtToString := r.URL.Query().Get("createdAtTo")
	entryID := r.URL.Query().Get("entryID")

	filter := bson.M{}

	if content != "" {
		filter["content"] = bson.M{
			"$regex":   content,
			"$options": "i",
		}
	}

	if editor != "" {
		filter["editor"] = editor
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

	if entryID != "" {
		filter["entry_id"] = entryID
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	var versions []model.Version

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VersionCollection.Find(ctx, filter, opts)
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
		w.WriteHeader(http.StatusNoContent)
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

	version.CreatedAt = time.Now().UTC()

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
// @Router       /api/versions/{id} [put]
func PutVersion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var newVersion model.Version
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newVersion); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingVersion model.Version
	err = database.VersionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingVersion)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Version not found")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	// Identify media_ids to delete
	mediaIDsToDelete := difference(existingVersion.MediaIDs, newVersion.MediaIDs)

	// Delete unreferenced media files
	client := &http.Client{Timeout: 5 * time.Second}
	for _, mediaID := range mediaIDsToDelete {
		mediaServiceURL := fmt.Sprintf("%s/api/media/%s", config.App.API_GATEWAY_URL, mediaID)
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
			config.App.Logger.Error().Int("status", resp.StatusCode).Str("body", bodyString).Msg("Media service returned error")
			http.Error(w, "Failed to delete associated media files", http.StatusInternalServerError)
			return
		}
	}

	newVersion.UpdatedAt = time.Now().UTC()

	update := bson.M{
		"$set": bson.M{
			"content":    newVersion.Content,
			"editor":     newVersion.Editor,
			"updated_at": newVersion.UpdatedAt,
			"address":    newVersion.Address,
			"media_ids":  newVersion.MediaIDs,
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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Retrieve the updated document (optional)
	err = database.VersionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&newVersion)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated version")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newVersion); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Helper function to find the difference between two slices
func difference(slice1, slice2 []string) []string {
	var diff []string
	m := make(map[string]bool)

	for _, item := range slice2 {
		m[item] = true
	}

	for _, item := range slice1 {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
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
// @Router       /api/versions/{id} [delete]
func DeleteVersion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Get the MediaIDs associated with the version

	var version model.Version
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = database.VersionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&version)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Version not found")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	// Delete associated media files first
	for _, mediaID := range version.MediaIDs {
		mediaServiceURL := fmt.Sprintf("%s/api/media/%s", config.App.API_GATEWAY_URL, mediaID)
		req, err := http.NewRequest("DELETE", mediaServiceURL, nil)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to create request to media service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		config.App.Logger.Info().Str("url", mediaServiceURL).Msg("Sending delete request to media service")

		resp, err := client.Do(req)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to send delete request to media service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			config.App.Logger.Error().Int("status", resp.StatusCode).Str("body", bodyString).Msg("Media service returned error")
			http.Error(w, "Failed to delete associated media files", http.StatusInternalServerError)
			return
		}
	}

	// Delete associated comments first
	commentServiceURL := fmt.Sprintf("%s/api/comments/version?versionID=%s", config.App.API_GATEWAY_URL, id)
	config.App.Logger.Info().Str("url", commentServiceURL).Msg("Preparing to delete associated comments")

	req, err := http.NewRequest("DELETE", commentServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to comment service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Str("url", commentServiceURL).Msg("Sending request to delete associated comments")

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

	result, err := database.VersionCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete version")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Version not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Str("versionID", id).Msg("Version and associated comments deleted successfully")
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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Collect all versionIDs
	var versionIDs []string
	for _, version := range versions {
		versionIDs = append(versionIDs, version.ID)
	}

	//this is the client to send requests to the media service
	client := &http.Client{Timeout: 5 * time.Second}

	//delete associated media files
	for _, version := range versions {
		for _, mediaID := range version.MediaIDs {
			mediaServiceURL := fmt.Sprintf("%s/api/media/%s", config.App.API_GATEWAY_URL, mediaID)
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
				config.App.Logger.Error().Int("status", resp.StatusCode).Str("body", bodyString).Msg("Media service returned error")
				http.Error(w, "Failed to delete associated media files", http.StatusInternalServerError)
				return
			}
		}
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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Int64("deletedCount", deleteResult.DeletedCount).Str("entryID", entryID).Msg("Versions deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}
