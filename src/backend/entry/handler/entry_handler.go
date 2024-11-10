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

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the service is up
// @Tags         Health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/entries/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// GetEntries godoc
// @Summary      Get all entries
// @Description  Retrieves the list of all entries from the database.
// @Tags         Entries
// @Produce      application/json
// @Success      200  {array}   model.Entry
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/entries/ [get]
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

// GetEntryByID godoc
// @Summary      Get an entry by ID
// @Description  Retrieves an entry by its ID.
// @Tags         Entries
// @Produce      application/json
// @Param        id    query     string  true  "Entry ID"
// @Success      200   {object}  model.Entry
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Entry not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/entries/id/ [get]
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

// GetEntryByExactTitle godoc
// @Summary      Get entries by title
// @Description  Retrieves entry that matches exactly the given title.
// @Tags         Entries
// @Produce      application/json
// @Param        title  query     string  true  "Title to search"
// @Success      200    {object}  model.Entry
// @Failure      404    {string}  string  "Entry not found"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/entries/exactTitle [get]
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

// GetEntriesByTitle godoc
// @Summary      Get all entries by title.
// @Description  Retrieves the list of all entries that matches the title.
// @Tags         Entries
// @Produce      application/json
// @Param        title  query     string  true  "Title to search"
// @Success      200  {array}   model.Entry
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/entries/title [get]
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

// GetEntriesByAuthor godoc
// @Summary      Get entries by author
// @Description  Retrieves entries authored by the given author(s).
// @Tags         Entries
// @Produce      application/json
// @Param        author  query     string true  "Author to search"
// @Success      200     {array}   model.Entry
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/entries/author [get]
func GetEntriesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query()["author"]

	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.EntryCollection.Find(ctx, bson.M{"author": author})
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
		config.App.Logger.Warn().Strs("author", author).Msg("No entries found")
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

// GetEntriesByDate godoc
// @Summary      Get entries by date
// @Description  Retrieves entries created on the given date.
// @Tags         Entries
// @Produce      application/json
// @Param        createdAt  query     string  true  "Creation date (YYYY-MM-DD)"
// @Success      200        {array}   model.Entry
// @Failure      400        {string}  string  "Invalid date format. Expected YYYY-MM-DD"
// @Failure      500        {string}  string  "Internal server error"
// @Router       /api/entries/date [get]
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

// GetEntriesByWikiID godoc
// @Summary      Get entries by Wiki ID
// @Description  Retrieves entries associated with a specific Wiki ID.
// @Tags         Entries
// @Produce      application/json
// @Param        wikiID  query     string  true  "Wiki ID"
// @Success      200     {array}   model.Entry
// @Failure      400     {string}  string  "WikiID is required"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/entries/wiki [get]
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

// PostEntry godoc
// @Summary      Create a new entry
// @Description  Creates a new entry. Expects a JSON object in the request body.
// @Tags         Entries
// @Accept       application/json
// @Produce      application/json
// @Param        entry  body      model.Entry  true  "Entry information"
// @Success      201    {object}  model.Entry
// @Failure      400    {string}  string  "Invalid request body"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/entries/ [post]
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

// PutEntry godoc
// @Summary      Update an entry by ID
// @Description  Updates an entry by its ID. Expects a JSON object in the request body.
// @Tags         Entries
// @Accept       application/json
// @Produce      application/json
// @Param        id     query     string      true  "Entry ID"
// @Param        entry  body      model.Entry true  "Updated entry information"
// @Success      200    {object}  model.Entry
// @Failure      400    {string}  string  "Invalid ID or request body"
// @Failure      404    {string}  string  "Entry not found"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/entries/id/ [put]
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
			"author":     entry.Author,
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

// DeleteEntry godoc
// @Summary      Delete an entry by ID
// @Description  Deletes an entry by its ID.
// @Tags         Entries
// @Param        id    query     string  true  "Entry ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Entry not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/entries/id/ [delete]
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

// DeleteEntriesByWikiID godoc
// @Summary      Deletes all entries by the Wiki ID
// @Description  Deletes all entries associated with a specific Wiki ID.
// @Tags         Entries
// @Param        id    query     string  true  "Wiki ID"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {string}  string  "WikiID is required"
// @Failure      404   {string}  string  "No entries found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/entries/wiki/ [delete]
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
