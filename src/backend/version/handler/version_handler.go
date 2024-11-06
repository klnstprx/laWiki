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
		version = append(versions, version)
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

/*
* GET /{id}
* gets a wiki by id from db

	func GetWikiByID(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Invalid ID format")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var wiki model.Wiki

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Wiki not found")
			http.Error(w, "Wiki not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to encode response")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

* PUT /{id}
* updates a wiki by id
* expects a json object in the request

	func PutWiki(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Invalid ID format")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var wiki model.Wiki
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"title":       wiki.Title,
				"description": wiki.Description,
				"category":    wiki.Category,
				"url":         wiki.URL,
			},
		}

		result, err := database.WikiCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Database error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if result.MatchedCount == 0 {
			config.App.Logger.Warn().Str("id", id).Msg("Wiki not found for update")
			http.Error(w, "Wiki not found", http.StatusNotFound)
			return
		}

		// Retrieve the updated document (optional)
		err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated wiki")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(wiki); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to encode response")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}


* DELETE /{id}
* deletes a wiki

func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.WikiCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("Wiki not found for deletion")
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

*/
