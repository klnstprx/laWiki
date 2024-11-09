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

/*
* GET /health
* checks if the service is up
 */
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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

// Implement other CRUD operations (GetCommentByID, PutComment, DeleteComment) similarly

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
		config.App.Logger.Error().Err(err).Msg("Wiki not found")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

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

func GetCommentByContent(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")

	var comment model.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.CommentCollection.FindOne(ctx, bson.M{"content": content}).Decode(&comment)
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

func GetCommentByRating(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetCommentByDate(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
