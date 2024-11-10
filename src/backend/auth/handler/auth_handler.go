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
	"github.com/laWiki/auth/model"

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
	expiration := time.Now().Add(1 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
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
	// Get state from the cookie
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		// Invalid state, possible CSRF attack
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// Exchange code for access token and retrieve user data
	data, err := getUserDataFromGoogle(r.FormValue("code"))
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

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    jwtToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	// Redirect to home page or wherever you want
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	// Create the JWT claims, which includes the user information and standard claims
	claims := jwt.MapClaims{}
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return tokenString, nil
}
