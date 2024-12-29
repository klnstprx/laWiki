package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/laWiki/auth/config"
	"github.com/laWiki/auth/database"
	"github.com/laWiki/auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the service is up
// @Tags         Health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /api/auth/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// A partir de aqui son las fuciones para la coleccion usuarios(no he borrado lo anterior por si hacia falta)
// CRUD DE USUARIOS

// GetVersions godoc
// @Summary      Get all versions
// @Description  Retrieves the list of all version JSON objects from the database.
// @Tags         Versions
// @Produce      application/json
// @Success      200  {array}   model.Version
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/versions/ [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var usuarios []model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.UsuarioCollection.Find(ctx, bson.M{})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var usuario model.User
		if err := cursor.Decode(&usuario); err != nil {
			config.App.Logger.Error().Err(err).Msg("Failed to decode user")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		usuarios = append(usuarios, usuario)
	}

	if err := cursor.Err(); err != nil {
		config.App.Logger.Error().Err(err).Msg("Cursor error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usuarios) == 0 {
		config.App.Logger.Info().Msg("No users found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVersionByID godoc
// @Summary      Get a version by ID
// @Description  Retrieves a version by its ID.
// @Tags         Versions
// @Produce      application/json
// @Param        id    query     string  true  "Version ID"
// @Success      200   {object}  model.Version
// @Failure      400   {string}  string  "Invalid ID"
// @Failure      404   {string}  string  "Version not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/versions/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verifica si el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var version model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&version)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Version not found")
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// PostVersion godoc
// @Summary      Create a new version
// @Description  Creates a new version. Expects a JSON object in the request body.
// @Tags         Versions
// @Accept       application/json
// @Produce      application/json
// @Param        version  body      model.Version  true  "Version information"
// @Success      201      {object}  model.Version
// @Failure      400      {string}  string  "Invalid request body"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/versions/ [post]
func PostUser(w http.ResponseWriter, r *http.Request) {
	var usuario model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usuario); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.UsuarioCollection.InsertOne(ctx, usuario)
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
	usuario.ID = objID.Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Return 201 Created
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	config.App.Logger.Info().Interface("usuario", usuario).Msg("Added new user")
}

// PutVersion godoc
// @Summary      Update a version by ID
// @Description  Updates a version by its ID. Expects a JSON object in the request body.
// @Tags         Versions
// @Accept       application/json
// @Produce      application/json
// @Param        id      query     string          true  "Version ID"
// @Param        version body      model.Version   true  "Updated version information"
// @Success      200     {object}  model.Version
// @Failure      400     {string}  string  "Invalid ID or request body"
// @Failure      404     {string}  string  "Version not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/versions/{id} [put]
func PutUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verifica si el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser model.User
	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingUser)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fusiona los campos nuevos con los existentes
	updatedFields := bson.M{}

	// Verifica si `enable_mails` está presente en el payload
	if enableMails, ok := payload["enable_mails"]; ok {
		if enableMailsBool, isBool := enableMails.(bool); isBool {
			updatedFields["enable_mails"] = enableMailsBool
		} else {
			http.Error(w, "Invalid type for enable_mails", http.StatusBadRequest)
			return
		}
	} else {
		updatedFields["enable_mails"] = existingUser.EnableMails
	}

	if role, ok := payload["role"]; ok {
		updatedFields["role"] = role
	}

	if valoration, ok := payload["valoration"]; ok {
		updatedFields["valoration"] = valoration
	}

	if notifications, ok := payload["notifications"]; ok {
		updatedFields["notifications"] = notifications
	}

	update := bson.M{
		"$set": updatedFields,
	}

	result, err := database.UsuarioCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("User not found for update")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Helper function to find the difference between two slices
func difference(slice1, slice2 []string) []string {
	var diff []string
	m := make(map[string]bool)

	for _, item := range slice2 {
		m[item] = true
	}

	for _, item := range slice1 {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// DeleteVersion godoc
// @Summary      Delete a version by ID
// @Description  Deletes a version by its ID.
// @Tags         Versions
// @Param        id query string true "Version ID"
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Version not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/versions/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verifica si el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the MediaIDs associated with the version

	var usuario model.User

	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	result, err := database.UsuarioCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete User")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("User not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Str("usuarioID", id).Msg("User deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func AddUserNotification(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verifica si el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Notification string `json:"notification"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to decode provided request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser model.User
	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingUser)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Hago que solo se pueda actualizar el rol y la valoracion del usuario(por ahora no queremos otro)
	update := bson.M{
		"$set": bson.M{
			"notifications": append(existingUser.Notifications, body.Notification),
		},
	}

	result, err := database.UsuarioCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		config.App.Logger.Warn().Str("id", id).Msg("User not found for update")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
