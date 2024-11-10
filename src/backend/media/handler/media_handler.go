package handler

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/laWiki/media/config"
	"github.com/laWiki/media/database"
	"github.com/laWiki/media/model"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the service is up
// @Tags         Health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/media/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getImageFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	err := r.ParseMultipartForm(config.App.MB_LIMIT << 20) // 5 MB
	if err != nil {
		log.Println("Error parsing form data:", err)
		return nil, nil, err
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Error retrieving image from form data:", err)
		return nil, nil, err
	}

	if file == nil {
		log.Println("No image found in form data")
		return nil, nil, err
	}
	return file, header, nil
}

// PostMedia godoc
// @Summary      Upload a new media file
// @Description  Uploads a new media file to Cloudinary and stores media info in the database.
// @Tags         Media
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        image  formData  file     true  "Image file to upload"
// @Success      201    {object}  model.Media
// @Failure      400    {string}  string  "Bad request"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/media/ [post]
func PostMedia(w http.ResponseWriter, r *http.Request) {
	var media model.Media

	if config.App.Cld == nil {
		config.App.Logger.Error().Msg("Cloudinary not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the image from the form data
	//
	file, header, err := getImageFile(r)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to get image file")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// get the file name

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	media.PublicID = strings.Split(header.Filename, ".")[0]

	// Upload the image and set the PublicID to header.Filename
	//
	uploadResp, err := config.App.Cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: media.PublicID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Cloudinary error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	media.UploadUrl = uploadResp.SecureURL

	// Update the media object in the database

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
	// myImage, err := config.App.Cld.Image(header.Filename)
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

// GetMedia godoc
// @Summary      Get all media files
// @Description  Retrieves a list of all media files from the database.
// @Tags         Media
// @Produce      application/json
// @Success      200  {array}   model.Media
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/media/ [get]
func GetMedia(w http.ResponseWriter, r *http.Request) {
	var media []model.Media

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.MediaCollection.Find(ctx, bson.M{})
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

// GetMediaByID godoc
// @Summary      Get media by ID
// @Description  Retrieves a media file by its ID.
// @Tags         Media
// @Produce      application/json
// @Param        id   query     string  true  "Media ID"
// @Success      200  {object}  model.Media
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      404  {string}  string  "Media not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/media/id/ [get]
func GetMediaByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Info().Msgf("Retrieved ID: %s", id)
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

// GetMediaByID godoc
// @Summary      Get media by PublicID
// @Description  Retrieves a media file by its PublicID (filename).
// @Tags         Media
// @Produce      application/json
// @Param        id   query     string  true  "Media PublicID"
// @Success      200  {object}  model.Media
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      404  {string}  string  "Media not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/media/pubid/ [get]
func GetMediaByPublicID(w http.ResponseWriter, r *http.Request) {
	publicID := r.URL.Query().Get("publicId")

	var media model.Media

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.MediaCollection.FindOne(ctx, bson.M{"publicId": publicID}).Decode(&media)
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

// DeleteMedia godoc
// @Summary      Delete media by ID
// @Description  Deletes a media file by its ID from both Cloudinary and the database.
// @Tags         Media
// @Param        id  query     string  true  "Media ID"
// @Success      204 {string}  string  "No Content"
// @Failure      400 {string}  string  "Invalid ID"
// @Failure      404 {string}  string  "Media not found"
// @Failure      500 {string}  string  "Internal server error"
// @Router       /api/media/id/ [delete]
func DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the file name
	var media model.Media
	err = database.MediaCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&media)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to find media")
		http.Error(w, "Media not found", http.StatusNotFound)
		return
	}

	// Delete the image from Cloudinary
	_, err = config.App.Cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: media.PublicID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Error deleting image:")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	_, err = database.MediaCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete media")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PutMedia godoc
// @Summary      Update media by ID
// @Description  Updates a media file by its ID in Cloudinary and updates media info in the database.
// @Tags         Media
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        id     query     string  true  "Media ID"
// @Param        image  formData  file    true  "Updated image file to upload"
// @Success      204    {string}  string  "No Content"
// @Failure      400    {string}  string  "Invalid ID or bad request"
// @Failure      404    {string}  string  "Media not found"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/media/id/ [put]
func PutMedia(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var media model.Media

	// Get the image from the form data
	//
	file, header, err := getImageFile(r)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to get image file")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	media.PublicID = strings.Split(header.Filename, ".")[0]

	// Update the image in Cloudinary
	uploadResp, err := config.App.Cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: media.PublicID, Overwrite: api.Bool(true)})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Cloudinary error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	media.UploadUrl = uploadResp.SecureURL

	_, err = database.MediaCollection.ReplaceOne(ctx, bson.M{"_id": objID}, media)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to update media")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
