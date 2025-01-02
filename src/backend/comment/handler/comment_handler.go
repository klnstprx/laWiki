package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/comment/config"
	"github.com/laWiki/comment/database"
	"github.com/laWiki/comment/model"
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
// @Router       /api/comments/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// GetComments godoc
// @Summary      Get all comments
// @Description  Retrieves a list of all comments.
// @Tags         Comments
// @Produce      application/json
// @Success      200  {array}   model.Comment
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/comments/ [get]
func GetComments(w http.ResponseWriter, r *http.Request) {
	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, bson.M{})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment model.Comment
		if err := cursor.Decode(&comment); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode comment")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCommentByID godoc
// @Summary      Get comment by ID
// @Description  Retrieves a comment by its ID.
// @Tags         Comments
// @Produce      application/json
// @Param        id      query     string  true  "Comment ID"
// @Success      200     {object}  model.Comment
// @Failure      400     {string}  string  "Invalid ID"
// @Failure      404     {string}  string  "Comment not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/comments/{id} [get]
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var comment model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.CommentCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&comment)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Comment not found")
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// SearchComments godoc
// @Summary      Search comments
// @Description  Search for comments using various query parameters. You can search by content, author, createdAt, rating, or versionID. All parameters are optional and can be combined.
// @Tags         Comments
// @Produce      application/json
// @Param        content     query     string  false  "Partial content to search for (case-insensitive)"
// @Param        author      query     string  false  "Author nickname to search for"
// @Param        createdAt   query     string  false  "Creation date (YYYY-MM-DD)"
// @Param        rating      query     int     false  "Rating to filter by"
// @Param        versionID   query     string  false  "Version ID to search for"
// @Success      200         {array}   model.Comment
// @Failure      400         {string}  string  "Bad Request"
// @Failure      500         {string}  string  "Internal Server Error"
// @Router       /api/comments/search [get]
func SearchComments(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	content := r.URL.Query().Get("content")
	author := r.URL.Query().Get("author")
	createdAtFromString := r.URL.Query().Get("createdAtFrom")
	createdAtToString := r.URL.Query().Get("createdAtTo")
	ratingString := r.URL.Query().Get("rating")
	versionID := r.URL.Query().Get("versionID")
	entryID := r.URL.Query().Get("entryID")

	// Build the MongoDB filter dynamically
	filter := bson.M{}

	if content != "" {
		filter["content"] = bson.M{
			"$regex":   content,
			"$options": "i",
		}
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

	if ratingString != "" {
		rating, err := strconv.Atoi(ratingString)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Invalid rating format")
			http.Error(w, "Invalid rating format", http.StatusBadRequest)
			return
		}
		filter["rating"] = rating
	}

	if versionID != "" {
		filter["version_id"] = versionID
	}

	if entryID != "" {
		filter["entry_id"] = entryID
	}

	// Query the database
	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, filter)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment model.Comment
		if err := cursor.Decode(&comment); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode comment")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// PostComment godoc
// @Summary      Create a new comment
// @Description  Creates a new comment. Expects a JSON object in the request body.
// @Tags         Comments
// @Accept       application/json
// @Produce      application/json
// @Param        comment  body      model.Comment  true  "Comment to create"
// @Success      201      {object}  model.Comment
// @Failure      400      {string}  string  "Invalid request body"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/comments/ [post]
func PostComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment.CreatedAt = time.Now().UTC()

	// Retrieve EntryID from the provided VersionID by making an HTTP request
	versionServiceURL := fmt.Sprintf("%s/api/versions/%s", config.App.API_GATEWAY_URL, comment.VersionID)
	config.App.Logger.Info().Str("url", versionServiceURL).Msg("Fetching Version to get EntryID")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", versionServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to version service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Internal-Auth", config.App.JWTSecret)
	resp, err := client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to version service")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		config.App.Logger.Error().
			Int("status", resp.StatusCode).
			Str("body", bodyString).
			Msg("Version service returned error")
		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, "Version not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve Version", http.StatusInternalServerError)
		}
		return
	}

	var version struct {
		ID      string `json:"id"`
		EntryID string `json:"entry_id"`
		Editor  string `json:"editor"`
	}

	err = json.NewDecoder(resp.Body).Decode(&version)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode Version data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the EntryID in the comment
	comment.EntryID = version.EntryID

	// Proceed to insert the comment into the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.CommentCollection.InsertOne(ctx, comment)
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
	comment.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("comment", comment).Msg("Added new comment")

	// Retrieve the entry title from the entry service with the entryID from the comment
	entryServiceURL := fmt.Sprintf("%s/api/entries/%s", config.App.API_GATEWAY_URL, comment.EntryID)
	config.App.Logger.Info().Str("url", entryServiceURL).Msg("Fetching Entry to get title")

	req, err = http.NewRequest("GET", entryServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to entry service")
		return
	}
	req.Header.Set("X-Internal-Auth", config.App.JWTSecret)
	resp, err = client.Do(req)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send request to entry service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		config.App.Logger.Error().
			Int("status", resp.StatusCode).
			Str("body", bodyString).
			Msg("Entry service returned error")
		return
	}

	var entry struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode Entry data")
		return
	}

	// Retrieve the editor's email address from the user service with the editor ID from the version
	userServiceURL := fmt.Sprintf("%s/api/auth/user?id=%s", config.App.API_GATEWAY_URL, version.Editor)
	config.App.Logger.Info().Str("url", userServiceURL).Msg("Fetching User to get email")

	req, err = http.NewRequest("GET", userServiceURL, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create request to user service")
		return
	}
	req.Header.Set("X-Internal-Auth", config.App.JWTSecret)
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
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		EnableMails bool   `json:"enable_mails"`
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode User data")
		return
	}

	if user.EnableMails {
		// email notification al editor
		notifyEmail("Nuevo comentario recibido",
			"Hola {{ nombre }},\nSe ha añadido un nuevo comentario a tu entrada \"{{ entrada }}\".",
			"<p> Hola {{ nombre }},</p><p>Se ha añadido un nuevo comentario a tu entrada \"{{ entrada }}\".</p>",
			user.Name,
			user.Email,
			entry.Title)
	} else {
		// notificacion interna al editor
		notifyInterno("Se ha añadido un nuevo comentario a tu entrada "+entry.Title, version.Editor)
	}
}

// PutComment godoc
// @Summary      Update comment by ID
// @Description  Updates a comment by its ID. Expects a JSON object in the request body.
// @Tags         Comments
// @Accept       application/json
// @Produce      application/json
// @Param        id       query     string         true  "Comment ID"
// @Param        comment  body      model.Comment  true  "Updated comment"
// @Success      200      {object}  model.Comment
// @Failure      400      {string}  string  "Invalid ID or request body"
// @Failure      404      {string}  string  "Comment not found"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/comments/{id} [put]
func PutComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var comment model.Comment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"content":    comment.Content,
			"rating":     comment.Rating,
			"author":     comment.Author,
			"updated_at": comment.UpdatedAt,
		},
	}

	result, err := database.CommentCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Comment not found for update")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Retrieve the updated document (optional)
	err = database.CommentCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&comment)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated wiki")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// DeleteComment godoc
// @Summary      Delete comment by ID
// @Description  Deletes a comment by its ID.
// @Tags         Comments
// @Param        id      query     string  true  "Comment ID"
// @Success      204     {string}  string  "No Content"
// @Failure      400     {string}  string  "Invalid ID"
// @Failure      404     {string}  string  "Comment not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/comments/{id} [delete]
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.CommentCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		config.App.Logger.Error().Msg("Comment not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCommentsByVersionID godoc
// @Summary      Delete comments by version ID
// @Description  Deletes all comments associated with a specific version.
// @Tags         Comments
// @Param        versionID      query     string  true  "Version ID"
// @Success      204     {string}  string  "No Content"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/comments/version [delete]
func DeleteCommentsByVersionID(w http.ResponseWriter, r *http.Request) {
	versionID := r.URL.Query().Get("versionID")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.CommentCollection.DeleteMany(ctx, bson.M{"version_id": versionID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete comments")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

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

	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to send email")
		//need more detailed error handling

		config.App.Logger.Info().Interface("response", res).Msg("Email sent")
		return
	}
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
	req.Header.Set("X-Internal-Auth", config.App.JWTSecret)

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
