package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/entry/config"
	"github.com/laWiki/entry/database"
	"github.com/laWiki/entry/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mailersend/mailersend-go"
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
// @Router       /api/entries/{id} [get]
func GetEntryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

// SearchEntries godoc
// @Summary      Search entries
// @Description  Search for entries using various query parameters. You can search by title, exact_title, author, createdAt, or wikiID. All parameters are optional and can be combined.
// @Tags         Entries
// @Produce      application/json
// @Param        title        query     string  false  "Partial title to search for (case-insensitive)"
// @Param        exact_title  query     string  false  "Exact title to search for"
// @Param        author       query     string  false  "Author to search for"
// @Param        createdAt    query     string  false  "Creation date (YYYY-MM-DD)"
// @Param        wikiID       query     string  false  "Wiki ID to search for"
// @Success      200          {array}   model.Entry
// @Failure      400          {string}  string  "Bad Request"
// @Failure      500          {string}  string  "Internal Server Error"
// @Router       /api/entries/search [get]
func SearchEntries(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	title := r.URL.Query().Get("title")
	exactTitle := r.URL.Query().Get("exact_title")
	author := r.URL.Query().Get("author")
	createdAtFromString := r.URL.Query().Get("createdAtFrom")
	createdAtToString := r.URL.Query().Get("createdAtTo")
	wikiID := r.URL.Query().Get("wikiID")

	// Build the MongoDB filter dynamically
	filter := bson.M{}

	if title != "" {
		filter["title"] = bson.M{
			"$regex":   title,
			"$options": "i",
		}
	}

	if exactTitle != "" {
		filter["title"] = exactTitle
	}

	if author != "" {
		filter["author"] = author
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

	if wikiID != "" {
		filter["wiki_id"] = wikiID
	}

	// Query the database
	var entries []model.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		config.App.Logger.Info().Msg("No entries found")
		w.WriteHeader(http.StatusNoContent)
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

	if entry.Title == "" || entry.WikiID == "" {
		config.App.Logger.Error().Msg("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	entry.CreatedAt = time.Now().UTC()

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
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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
// @Router       /api/entries/{id} [put]
func PutEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

	entry.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":      entry.Title,
			"author":     entry.Author,
			"updated_at": entry.UpdatedAt,
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
		w.WriteHeader(http.StatusNoContent)
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
// @Router       /api/entries/{id} [delete]
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete associated versions first
	versionServiceURL := fmt.Sprintf("%s/api/versions/entry?entryID=%s", config.App.API_GATEWAY_URL, id)
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
	req.Header.Set("X-Internal-Request", "true")
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

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid entry ID")
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	// Retrieve the entry for email notification

	var entry model.Entry

	err = database.EntryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&entry)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Entry not found")
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Now proceed to delete the entry document

	result, err := database.EntryCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete entry")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("Entry not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Str("entryID", id).Msg("Version and associated versions deleted successfully")
	w.WriteHeader(http.StatusNoContent)

	// Retrieve the user from the user service with the author ID from the entry
	userServiceURL := fmt.Sprintf("%s/api/auth/user?id=%s", config.App.API_GATEWAY_URL, entry.Author)
	config.App.Logger.Info().Str("url", userServiceURL).Msg("Sending request to user service")

	req, err = http.NewRequest("GET", userServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to user service")
		return
	}
	req.Header.Set("X-Internal-Request", "true")
	resp, err = client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to user service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		config.App.Logger.Error().
			Int("status", resp.StatusCode).
			Str("body", bodyString).
			Msg("User service returned error")
		return
	}

	var user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode user response")
		return
	}

	// Send email notification to the entry author
	/*
		notifyEmail("Tu entrada ha sido eliminada",
			"Hola {{ nombre }},\nTu entrada \"{{ entrada }}\" ha sido eliminada.",
			"<p> Hola {{ nombre }},</p><p>Tu entrada \"{{ entrada }}\" ha sido eliminada.</p>",
			user.Name,
			user.Email,
			entry.Title)
	*/
	notifyInterno("Tu entrada "+entry.Title+" ha sido eliminada", entry.Author)
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
		w.WriteHeader(http.StatusNoContent)
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
		req.Header.Set("X-Internal-Request", "true")
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

func notifyEmail(subject string, text string, html string, destinoNombre string, destinoEmail string, entryTitle string) {
	ms := mailersend.NewMailersend("mlsn.9938f4dc11ca834ac853af3f07c9d9552a39e8007391e356dfb28d76094516c8")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  "laWiki",
		Email: "laWiki@trial-x2p0347d0v94zdrn.mlsender.net",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  destinoNombre,
			Email: destinoEmail,
		},
	}

	personalization := []mailersend.Personalization{
		{
			Email: destinoEmail,
			Data: map[string]interface{}{
				"nombre":  destinoNombre,
				"entrada": entryTitle,
			},
		},
	}

	tags := []string{}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	message.SetTags(tags)
	message.SetPersonalization(personalization)

	res, _ := ms.Email.Send(ctx, message)

	fmt.Printf(res.Header.Get("X-Message-Id"))
}

func notifyInterno(mensaje string, editor string) {
	notificationMessage := fmt.Sprintf(
		mensaje,
	)

	// Construir la URL del servicio de usuarios con query string
	userServiceURL := fmt.Sprintf("%s/api/auth/notifications?id=%s", config.App.API_GATEWAY_URL, editor)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Crear el cuerpo de la solicitud array con la cadena
	notificationPayload := map[string]string{
		"notification": notificationMessage,
	}

	payloadBytes, _ := json.Marshal(notificationPayload)

	req, err := http.NewRequest("POST", userServiceURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to user service")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Request", "true")
	// Enviar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to user service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		config.App.Logger.Error().
			Int("status", resp.StatusCode).
			Str("body", bodyString).
			Msg("User service returned error")
		return
	}

	config.App.Logger.Info().
		Str("userId", editor).
		Msg("Notification sent to user service")
}
