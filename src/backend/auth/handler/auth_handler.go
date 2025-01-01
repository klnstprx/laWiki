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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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

	// Check if a user with the same email already exists
	var existingUser model.User
	err := database.UsuarioCollection.FindOne(ctx, bson.M{"email": usuario.Email}).Decode(&existingUser)
	if err == nil {
		// User already exists
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

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

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing user email", http.StatusBadRequest)
		return
	}

	var user model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.UsuarioCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the cookie

	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		config.App.Logger.Error().Err(err).Msg("Missing jwt_token cookie.")
		return
	}

	//return it via response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cookie.Value))
}
