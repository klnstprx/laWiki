package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/laWiki/entry/config"
	"github.com/laWiki/entry/database"
	"github.com/laWiki/entry/model"
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

func PostEntry(w http.ResponseWriter, r *http.Request) {
	var entry model.Entry
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	entry.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.EntryCollection.InsertOne(ctx, entry)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		config.App.Logger.Error().Msg("Failed to convert InsertedID to ObjectID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	entry.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("entry", entry).Msg("Added new entry")
}

func GetEntries(w http.ResponseWriter, r *http.Request) {
	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.EntryCollection.Find(ctx, bson.M{})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntryByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var entry model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.EntryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&entry)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Entry not found")
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	entryID := r.URL.Query().Get("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete associated versions first
	versionServiceURL := fmt.Sprintf("%s/api/versions/entry?entryID=%s", config.App.API_GATEWAY_URL, entryID)
	config.App.Logger.Info().Str("url", versionServiceURL).Msg("Preparing to delete associated versions")

	req, err := http.NewRequest("DELETE", versionServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to version service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Str("url", versionServiceURL).Msg("Sending request to delete associated versions")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to version service")
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
			Msg("version service returned error")
		http.Error(w, "Failed to delete associated versions", http.StatusInternalServerError)
		return
	}

	// Now proceed to delete the entry document
	objID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid entry ID")
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	result, err := database.EntryCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete entry")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Version not found")
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	config.App.Logger.Info().Str("entryID", entryID).Msg("Version and associated versions deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func PutEntry(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var entry model.Entry
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":      entry.Title,
			"authors":    entry.Authors,
			"created_at": entry.CreatedAt,
		},
	}

	result, err := database.EntryCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Entry not found for update")
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Retrieve the updated document (optional)
	err = database.EntryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&entry)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated entry")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntryByExactTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.EntryCollection.Find(ctx, bson.M{"title": title})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Info().Str("title", title).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntriesByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}

	cursor, err := database.EntryCollection.Find(ctx, filter)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Info().Str("title", title).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntriesByAuthors(w http.ResponseWriter, r *http.Request) {
	authors := r.URL.Query()["author"]

	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.EntryCollection.Find(ctx, bson.M{"authors": authors})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Warn().Strs("authors", authors).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntriesByDate(w http.ResponseWriter, r *http.Request) {
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

	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build the query to find entries within the date range
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cursor, err := database.EntryCollection.Find(ctx, filter)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Warn().Str("createdAt", createdAtString).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetEntriesByWikiID(w http.ResponseWriter, r *http.Request) {
	wikiID := r.URL.Query().Get("wikiID")

	if wikiID == "" {
		config.App.Logger.Error().Msg("WikiID is required")
		http.Error(w, "WikiID is required", http.StatusBadRequest)
		return
	}

	var entries []model.Entry
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.EntryCollection.Find(ctx, bson.M{"wiki_id": wikiID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry model.Entry
		if err := cursor.Decode(&entry); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode entry")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Warn().Str("wikiID", wikiID).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteEntriesByWikiID(w http.ResponseWriter, r *http.Request) {
	wikiID := r.URL.Query().Get("wikiID")

	if wikiID == "" {
		config.App.Logger.Error().Msg("WikiID is required")
		http.Error(w, "WikiID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve all entries with the given wikiID
	var entries []model.Entry
	cursor, err := database.EntryCollection.Find(ctx, bson.M{"wiki_id": wikiID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error while fetching entries")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &entries); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode entries")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		config.App.Logger.Info().Str("wikiID", wikiID).Msg("No entries found")
		http.Error(w, "No entries found", http.StatusNotFound)
		return
	}

	// For each entryID, delete associated versions
	for _, entry := range entries {
		entryID := entry.ID

		// Send request to version service to delete versions associated with entryID
		versionServiceURL := fmt.Sprintf("%s/api/versions/entry?entryID=%s", config.App.API_GATEWAY_URL, entryID)
		config.App.Logger.Info().Str("url", versionServiceURL).Msg("Preparing to delete associated versions")

		req, err := http.NewRequest("DELETE", versionServiceURL, nil)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to create request to version service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to send request to version service")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			config.App.Logger.Error().
				Int("status", resp.StatusCode).
				Str("body", bodyString).
				Msg("Version service returned error")
			http.Error(w, "Failed to delete associated versions", http.StatusInternalServerError)
			return
		}

		config.App.Logger.Info().Str("entryID", entryID).Msg("Associated versions deleted successfully")
	}

	// Proceed to delete entries associated with the wikiID
	deleteResult, err := database.EntryCollection.DeleteMany(ctx, bson.M{"wiki_id": wikiID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete entries")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().
		Str("wikiID", wikiID).
		Int64("deletedCount", deleteResult.DeletedCount).
		Msg("Entries and their associated versions deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}
