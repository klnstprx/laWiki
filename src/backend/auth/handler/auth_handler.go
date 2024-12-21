package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/laWiki/auth/config"
	"github.com/laWiki/auth/database"
	"github.com/laWiki/auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt"
)

const FRONTEND_URL = "http://localhost:5173"

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

// generateStateOauthCookie generates a random state string and sets it as a cookie
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().UTC().Add(1 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	return state
}

// Login godoc
// @Summary      Initiate OAuth2 Login
// @Description  Initiates the OAuth2 flow with Google
// @Tags         Authentication
// @Success      302  {string}  string  "Redirect to Google OAuth2 login"
// @Router       /api/auth/login [get]
func Login(w http.ResponseWriter, r *http.Request) {
	// Generate a random state parameter to prevent CSRF attacks.
	state := generateStateOauthCookie(w)

	// Get the OAuth2 Config
	oauthConfig := config.App.GoogleOAuthConfig

	// Redirect user to Google's OAuth consent page
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback godoc
// @Summary      OAuth2 Callback
// @Description  Handles the OAuth2 callback from Google
// @Tags         Authentication
// @Param        state  query     string  true  "OAuth state"
// @Param        code   query     string  true  "Authorization code"
// @Success      302    {string}  string  "Redirect after login"
// @Failure      401    {string}  string  "Invalid OAuth state"
// @Failure      500    {string}  string  "Could not create JWT"
// @Router       /api/auth/callback [get]
func Callback(w http.ResponseWriter, r *http.Request) {
	// Validate 'state' and 'code' parameters
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state == "" || code == "" {
		http.Error(w, "Missing 'state' or 'code' parameters", http.StatusBadRequest)
		return
	}

	// Get state from the cookie
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusUnauthorized)
		return
	}

	if state != oauthState.Value {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// Exchange code for access token and retrieve user data
	data, err := getUserDataFromGoogle(code)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to get user data from Google")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Generate JWT token
	jwtToken, err := createJWTToken(data)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create JWT token")
		http.Error(w, "Could not create JWT", http.StatusInternalServerError)
		return
	}

	// Set the JWT token in a secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    jwtToken,
		Expires:  time.Now().UTC().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to the frontend application
	http.Redirect(w, r, FRONTEND_URL, http.StatusSeeOther)
}

// getUserDataFromGoogle exchanges the code for a token and gets user info from Google
func getUserDataFromGoogle(code string) ([]byte, error) {
	oauthConfig := config.App.GoogleOAuthConfig

	// Exchange the code for a token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	// Create a new HTTP client using the token
	client := oauthConfig.Client(context.TODO(), token)

	// Get the user's info from the Google API
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %v", err)
	}
	return data, nil
}

// createJWTToken creates a JWT token with the user's information
func createJWTToken(data []byte) (string, error) {
	var user model.User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	// Validate or set user role
	if user.Role == "" {
		user.Role = "user" // Set a default role
	}

	// Create the JWT claims
	claims := jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
		"exp":   time.Now().UTC().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"iat":   time.Now().UTC().Unix(),
		"iss":   "auth_service", // Issuer identifier
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return tokenString, nil
}

// Logout godoc
// @Summary      Logout
// @Description  Clears the JWT token cookie
// @Tags         Authentication
// @Success      302  {string}  string  "Redirect after logout"
// @Router       /api/auth/logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the jwt_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().UTC().Add(-1 * time.Hour), // Set expiration in the past
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to the login page or home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UserInfo godoc
// @Summary      Get User Info
// @Description  Returns user information for the authenticated user
// @Tags         Authentication
// @Success      200  {object}  model.User
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /api/auth/userinfo [get]
func UserInfo(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the cookie
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.App.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	// Extract user claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := model.User{
			Email: claims["email"].(string),
			Name:  claims["name"].(string),
			Role:  claims["role"].(string),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
		return
	}

}

//A partir de aqui son las fuciones para la coleccion usuarios(no he borrado lo anterior por si hacia falta)
//CRUD DE USUARIOS

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
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var usuario model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
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
	id := chi.URLParam(r, "id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var newUser model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
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

	/*
		Email         string
		Name          string
		GivenName     string
		FamilyName    string
		Picture       string
		Role          string
		Valoration    double
	*/
	//Hago que solo se pueda actualizar el rol y la valoracion del usuario(por ahora no queremos otro)
	update := bson.M{
		"$set": bson.M{
			"role":       newUser.Role,
			"valoration": newUser.Valoration,
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

	// Retrieve the updated document (optional)
	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&newUser)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to retrieve updated user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
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
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the MediaIDs associated with the version

	var usuario model.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Invalid ID format")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

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
