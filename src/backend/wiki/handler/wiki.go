package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/laWiki/wiki/database"
	"github.com/laWiki/wiki/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
* POST /api/wikis
* creates a new wiki
* expects a json object in the request body
 */
func PostWiki(w http.ResponseWriter, r *http.Request) {
	var wiki model.Wiki
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wiki); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.WikiCollection.InsertOne(ctx, wiki)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wiki.ID = result.InsertedID.(primitive.ObjectID).Hex()
	json.NewEncoder(w).Encode(wiki)
}

/*
* GET /api/wikis
* gets the list of all wiki json objects from db
 */
func GetWikis(w http.ResponseWriter, r *http.Request) {
	var wikis []model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.WikiCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var wiki model.Wiki
		cursor.Decode(&wiki)
		wikis = append(wikis, wiki)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(wikis)
}

/*
* GET /api/wikis/{id}
* gets a wiki by id from db
 */
func GetWikiByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var wiki model.Wiki

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.WikiCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&wiki)
	if err != nil {
		http.Error(w, "Wiki not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wiki)
}

/*
* PUT /api/wikis/{id}
* updates a wiki by id
* expects a json object in the request
 */
func PutWiki(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var wiki model.Wiki
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wiki); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wiki.ID = id

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
	if err != nil || result.MatchedCount == 0 {
		http.Error(w, "Failed to update wiki", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(wiki)
}

/*
* DELETE /api/wikis/{id}
* deletes a wiki
 */
func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.WikiCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Failed to delete wiki", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
