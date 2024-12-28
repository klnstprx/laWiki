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

	"github.com/laWiki/auth/config"
	"github.com/laWiki/auth/database"
	"github.com/laWiki/auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt"
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
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
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
	// Generate a random state to prevent CSRF
	state := generateState()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// Build the Google login URL
	oauthConfig := config.App.GoogleOAuthConfig
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
// @Failure      400    {string}  string  "Missing state or code"
// @Failure      401    {string}  string  "Invalid OAuth state"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /api/auth/callback [get]
func Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state == "" || code == "" {
		http.Error(w, "Missing state or code", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("oauthstate")
	if err != nil || cookie.Value != state {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// get user info from Google
	userData, err := getUserDataFromGoogle(code)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to get user data from Google")
		http.Error(w, "Failed to retrieve Google user info", http.StatusInternalServerError)
		return
	}

	var googleUser model.GoogleUser

	if err := json.Unmarshal(userData, &googleUser); err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to parse Google user JSON")
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Look for existing user in DB or create new with a default role
	user, err := upsertUser(googleUser)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to upsert user in DB")
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// create local JWT with user’s role from DB
	tokenString, err := createLocalJWT(user)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to create local JWT")
		http.Error(w, "JWT creation failed", http.StatusInternalServerError)
		return
	}

	// set the JWT in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// redirect user to the frontend
	http.Redirect(w, r, config.App.FrontendURL, http.StatusSeeOther)
}

func upsertUser(g model.GoogleUser) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing model.User
	err := database.UsuarioCollection.FindOne(ctx, bson.M{"email": g.Email}).Decode(&existing)
	if err == nil {
		// user found
		return existing, nil
	}

	newUser := model.User{
		Email:         g.Email,
		Name:          g.Name,
		GivenName:     g.GivenName,
		FamilyName:    g.FamilyName,
		Picture:       g.Picture,
		Locale:        g.Locale,
		EmailVerified: g.EmailVerified,
		Role:          "editor",
		Valoration:    []float64{},
		Notifications: []string{},
		EnableMails:   false,
	}

	// Insert the new user and capture the InsertedID
	result, err := database.UsuarioCollection.InsertOne(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}

	// Set the assigned ObjectID to newUser.ID
	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return model.User{}, fmt.Errorf("Failed to convert InsertedID to ObjectID")
	}
	newUser.ID = objID.Hex()

	return newUser, nil
}

// createLocalJWT creates a local JWT using HS256 with user information
func createLocalJWT(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "auth_service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.App.JWTSecret))
}

// getUserDataFromGoogle exchanges code for token and retrieves user info from Google
func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := config.App.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("exchange error: %v", err)
	}
	client := config.App.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Logout godoc
// @Summary      Logout
// @Description  Clears the JWT token cookie
// @Tags         Authentication
// @Success      302  {string}  string  "Redirect after logout"
// @Router       /api/auth/logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, config.App.FrontendURL, http.StatusSeeOther)
}

// UserInfo godoc
// @Summary      Get Authenticated User Info
// @Description  Returns user information for the authenticated user
// @Tags         Authentication
// @Produce      application/json
// @Success      200  {object}  model.User
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /api/auth/me [get]
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
		userID, _ := claims["user_id"].(string)
		// Fetch user from database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user model.User
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("Invalid user ID in token")
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			config.App.Logger.Error().Err(err).Msg("User not found in DB")
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Return user info as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
		return
	}
}

// GetUsers godoc
// @Summary      Get All Users
// @Description  Retrieves the list of all user objects from the database
// @Tags         Users
// @Produce      application/json
// @Success      200  {array}   model.User
// @Success      204  {string}  string  "No Content"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/auth [get]
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

// GetUserByID godoc
// @Summary      Get a User by ID
// @Description  Retrieves a user by its ID
// @Tags         Users
// @Produce      application/json
// @Param        id    path      string  true  "User ID"
// @Success      200   {object}  model.User
// @Failure      400   {string}  string  "Invalid user ID"
// @Failure      404   {string}  string  "User not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/auth/users/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verify if the ID is a valid ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
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

// PostUser godoc
// @Summary      Create a New User
// @Description  Creates a new user. Expects a JSON object in the request body.
// @Tags         Users
// @Accept       application/json
// @Produce      application/json
// @Param        user  body      model.User  true  "User information"
// @Success      201   {object}  model.User
// @Failure      400   {string}  string  "Invalid request body"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/auth [post]
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

// PutUser godoc
// @Summary      Update a User by ID
// @Description  Updates a user by its ID. Expects a JSON object in the request body.
// @Tags         Users
// @Accept       application/json
// @Produce      application/json
// @Param        id      path      string          true  "User ID"
// @Param        user    body      map[string]interface{}  true  "Updated user information"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {string}  string  "Invalid user ID or request body"
// @Failure      404     {string}  string  "User not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/auth/{id} [put]
func PutUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verify if the ID is a valid ObjectID
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

	// Merge new fields with existing fields
	updatedFields := bson.M{}

	// Check if `enable_mails` is present in the payload
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

// DeleteUser godoc
// @Summary      Delete a User by ID
// @Description  Deletes a user by its ID
// @Tags         Users
// @Param        id      path      string  true  "User ID"
// @Success      204     {string}  string  "No Content"
// @Failure      400     {string}  string  "Invalid user ID"
// @Failure      404     {string}  string  "User not found"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/auth/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verify if the ID is a valid ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the User to be deleted
	var usuario model.User

	err = database.UsuarioCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	result, err := database.UsuarioCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Failed to delete user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		config.App.Logger.Info().Msg("User not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config.App.Logger.Info().Str("userID", id).Msg("User deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

// AddUserNotification godoc
// @Summary      Add a Notification to a User
// @Description  Adds a notification to the specified user
// @Tags         Users
// @Accept       application/json
// @Produce      application/json
// @Param        id             path      string                    true  "User ID"
// @Param        notification   body      string                    true  "Notification to add"
// @Success      200            {string}  string                    "Notification added successfully"
// @Failure      400            {string}  string                    "Invalid user ID or request body"
// @Failure      404            {string}  string                    "User not found"
// @Failure      500            {string}  string                    "Internal server error"
// @Router       /api/auth/{id}/notifications [post]
func AddUserNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Verify if the ID is a valid ObjectID
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

	// Update only the notifications field
	update := bson.M{
		"$set": bson.M{
			"notifications": append(existingUser.Notifications, body.Notification),
		},
	}

	_, err = database.UsuarioCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		config.App.Logger.Error().Err(err).Msg("Database error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification added successfully"))
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
