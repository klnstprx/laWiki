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
* creates a new version
* expects a json object in the request body
 */
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

/*
* GET /
* gets the list of all Version json objects from db
 */
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

/*
* GET /{id}
* gets a version by id from db
 */

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

/*
* PUT /{id}
* updates a version by id
* expects a json object in the request
 */

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

/*
* DELETE /{id}
* deletes a Version
 */
func DeleteVersion(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	versionID := vars.Get("id")

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

/*
* GET /
* gets the list of all Version json objects that correspond to an Entry from db
 */
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
