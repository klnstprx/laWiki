package handler

// Import the required packages for upload and admin.

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/laWiki/media/config"
	"github.com/laWiki/media/database"
	"github.com/laWiki/media/model"
)

// Initialize Cloudinary in the init function

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

var cld = config.App.Cld

//IMPORTANTE, HAY QUE CAMBIAR EL NOMBRE DE LA NUBE, LA CLAVE DE LA API Y EL SECRETO DE API.

//IMPORTANTE, HAY QUE CAMBIAR EL NOMBRE DE LA IMAGEN QUE SE SUBE A CLOUDINARY.

func PostMedia(w http.ResponseWriter, r *http.Request) {
	var media model.Media

	// Decode the request body into the media object
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&media); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Upload the image and set the PublicID to "my_image"
	uploadResp, err := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{PublicID: "my_image"})
	if err != nil {
		log.Println("Error uploading image:", err)
		return
	}
	media.UploadUrl = uploadResp.SecureURL

	// Get details about the image
	assetResp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "my_image"})
	if err != nil {
		log.Println("Error fetching asset details:", err)
		return
	}
	media.AssetUrl = assetResp.SecureURL

	//Update the media object in the database

	result, err := database.MediaCollection.InsertOne(ctx, media)
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
	media.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(media); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("media", media).Msg("Added new media")

	// // Create a transformation for the image
	// myImage, err := cld.Image("my_image")
	// if err != nil {
	// 	log.Println("Error creating image object:", err)
	// 	return
	// }
	// myImage.Transformation = "c_fill,h_250,w_250"

	// Generate and log the transformed image URL
	// url, err := myImage.String()
	// if err != nil {
	// 	log.Println("Error generating URL:", err)
	// 	return
	// }
	// log.Println("Transformed image URL:", url)
}

func GetMedia(w http.ResponseWriter, r *http.Request) {
	var media []model.Media

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.MediaCollection.Find(ctx, nil)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var m model.Media
		if err := cursor.Decode(&m); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode media")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		media = append(media, m)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(media); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetMediaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var media model.Media

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.MediaCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&media)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to find media")
		http.Error(w, "Media not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(media); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.MediaCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete media")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Delete the image from Cloudinary
	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: id})
	if err != nil {
		log.Println("Error deleting image:", err)
		return
	}
}

func PutMedia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var media model.Media

	// Decode the request body into the media object
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&media); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.MediaCollection.ReplaceOne(ctx, bson.M{"_id": objID}, media)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to update media")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Update the image in Cloudinary
	_, err = cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{PublicID: id})
	if err != nil {
		log.Println("Error updating image:", err)
		return
	}
}
