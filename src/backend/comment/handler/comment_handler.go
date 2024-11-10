package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/laWiki/comment/config"
	"github.com/laWiki/comment/database"
	"github.com/laWiki/comment/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Router       /api/comments/id/ [get]
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

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

// GetCommentsByContent godoc
// @Summary      Get comment by content
// @Description  Retrieves a comment by its content.
// @Tags         Comments
// @Produce      application/json
// @Param        content  query     string  true  "Content to search"
// @Success      200      {object}  model.Comment
// @Failure      404      {string}  string  "Comment not found"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/comments/content [get]
func GetCommentsByContent(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")

	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"content": bson.M{
			"$regex":   content,
			"$options": "i",
		},
	}

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
		config.App.Logger.Warn().Str("content", content).Msg("No comments found")
		http.Error(w, "No comments found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCommentsByRating godoc
// @Summary      Get comments by rating
// @Description  Retrieves comments with a specific rating.
// @Tags         Comments
// @Produce      application/json
// @Param        rating  query     int     true  "Rating to filter"
// @Success      200     {array}   model.Comment
// @Failure      400     {string}  string  "Invalid rating format"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/comments/rating [get]
func GetCommentsByRating(w http.ResponseWriter, r *http.Request) {
	ratingString := r.URL.Query().Get("rating")

	// cast rating to int
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid rating format")
		http.Error(w, "Invalid rating format", http.StatusBadRequest)
		return
	}

	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, bson.M{"rating": rating})
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
		config.App.Logger.Warn().Int("rating", rating).Msg("No comments found")
		http.Error(w, "No comments found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCommentsByDate godoc
// @Summary      Get comments by date
// @Description  Retrieves comments created on a specific date.
// @Tags         Comments
// @Produce      application/json
// @Param        createdAt  query     string  true  "Creation date (YYYY-MM-DD)"
// @Success      200        {array}   model.Comment
// @Failure      400        {string}  string  "Invalid date format. Expected YYYY-MM-DD"
// @Failure      500        {string}  string  "Internal server error"
// @Router       /api/comments/date [get]
func GetCommentsByDate(w http.ResponseWriter, r *http.Request) {
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

	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build the query to find comments within the date range
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

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
		config.App.Logger.Warn().Str("createdAt", createdAtString).Msg("No comments found")
		http.Error(w, "No comments found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCommentsByAuthor godoc
// @Summary      Get comments by author
// @Description  Retrieves comments created by a specific author.
// @Tags         Comments
// @Produce      application/json
// @Param        author  query     string  true  "Author nickname"
// @Success      200        {array}   model.Comment
// @Failure      500        {string}  string  "Internal server error"
// @Router       /api/comments/author [get]
func GetCommentsByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, bson.M{"author": author})
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
		config.App.Logger.Warn().Str("author", author).Msg("No comments found")
		http.Error(w, "No comments found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetCommentsByVersionID godoc
// @Summary      Get comments by the id of a version its associated to
// @Description  Retrieves comments associated with a specific version.
// @Tags         Comments
// @Produce      application/json
// @Param        versionID  query     string  true  "Version ID"
// @Success      200        {array}   model.Comment
// @Failure      500        {string}  string  "Internal server error"
// @Router       /api/comments/version [get]
func GetCommentsByVersionID(w http.ResponseWriter, r *http.Request) {
	versionID := r.URL.Query().Get("versionID")

	var comments []model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, bson.M{"versionID": versionID})
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
		config.App.Logger.Warn().Str("versionID", versionID).Msg("No comments found")
		http.Error(w, "No comments found", http.StatusNotFound)
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

	comment.CreatedAt = time.Now()

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
// @Router       /api/comments/id/ [put]
func PutComment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"content":    comment.Content,
			"rating":     comment.Rating,
			"created_at": comment.CreatedAt,
			"author":     comment.Author,
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
		http.Error(w, "Comment not found", http.StatusNotFound)
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
// @Router       /api/comments/id/ [delete]
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

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
		http.Error(w, "Comment not found", http.StatusNotFound)
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
