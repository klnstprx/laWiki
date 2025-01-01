package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/laWiki/gateway/config"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rs/zerolog"
)

type key int

const requestIDKey key = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for auth routes and health check

		w.Header().Set("Access-Control-Allow-Origin", config.App.FrontendURL) // Reemplaza con el dominio del frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Si es una solicitud OPTIONS, respondemos inmediatamente
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		config.App.Logger.Debug().Msgf("Authenticating request to: %s", r.URL.Path)

		if strings.Contains(r.URL.Path, "/api/auth") && r.Method == http.MethodPost {
			config.App.Logger.Info().Msg("Request to create new user, skipping authentication.")
			next.ServeHTTP(w, r)
			return
		}

		// Asi para que a las solicitudes get(y las que son a auth) no se les pida autenticacion
		if (strings.Contains(r.URL.Path, "/api/auth") && r.Method != http.MethodPut && r.Method != http.MethodDelete) || r.URL.Path == "/health" || (r.Method == http.MethodGet) {
			if r.URL.Path != "/api/auth" { // Para no poder ver todos los usuarios sin autenticacion
				config.App.Logger.Debug().Msg(r.URL.Path)
				config.App.Logger.Debug().Msg("Request to auth, health or a GET request. Passing without authentication.")
				next.ServeHTTP(w, r)
				return
			}
		}

		internalAuthHeader := r.Header.Get("X-Internal-Auth")
		if internalAuthHeader == config.App.JWTSecret {
			config.App.Logger.Debug().Msg("Internal request authenticated.")
			next.ServeHTTP(w, r)
			return
		}

		config.App.Logger.Debug().Str("X-Internal-Auth", internalAuthHeader)

		// Get the JWT token from the cookie
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			config.App.Logger.Error().Err(err).Msg("Missing jwt_token cookie.")
			return
		}

		cookieRole, err := r.Cookie("role")
		if err != nil {
			http.Error(w, "Unauthorized: missing role", http.StatusUnauthorized)
			config.App.Logger.Error().Err(err).Msg("Missing role cookie.")
			return
		}

		tokenString := cookie.Value
		role := cookieRole.Value
		config.App.Logger.Debug().Str("role", role)

		// Parse and validate the token using RS256 and public key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is RS256
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
			}

			// Get the public key from JWKS
			jwksURL := "https://www.googleapis.com/oauth2/v3/certs" // URL del JWKS de Google, cambia si es otro proveedor
			keySet, err := jwk.Fetch(r.Context(), jwksURL)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch public keys: %v", err)
			}

			// Find the key with the appropriate kid (Key ID) from the token header
			kid := token.Header["kid"].(string)
			key, ok := keySet.LookupKeyID(kid)
			if !ok {
				return nil, fmt.Errorf("unable to find appropriate key")
			}

			// Return the public key to verify the token
			var pubKey interface{}
			err = key.Raw(&pubKey)
			if err != nil {
				return nil, fmt.Errorf("failed to extract key: %v", err)
			}
			return pubKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user claims and add them to the request context, chequea que el rol sea el correcto
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user", claims)

			if role == "redactor" {
				if strings.Contains(r.URL.Path, "/api/entries") || strings.Contains(r.URL.Path, "/api/comments") || strings.Contains(r.URL.Path, "/api/media") || strings.Contains(r.URL.Path, "/api/versions") || strings.Contains(r.URL.Path, "/api/auth") {
					if r.Method == http.MethodPost || r.Method == http.MethodPut {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					} else {
						http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
						return
					}
				} else {
					http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
					return
				}
			} else if role == "editor" {
				if strings.Contains(r.URL.Path, "/api/media") || strings.Contains(r.URL.Path, "/api/wikis") || strings.Contains(r.URL.Path, "/api/auth") {
					if r.Method == http.MethodPost || r.Method == http.MethodPut {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					} else {
						http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
						return
					}
				} else if strings.Contains(r.URL.Path, "/api/entries") || strings.Contains(r.URL.Path, "/api/comments") || strings.Contains(r.URL.Path, "/api/versions") {
					if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					} else {
						http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
						return
					}
				} else {
					http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
					return
				}
			} else if role == "admin" {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		} else {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}
	})
}

// RequestID is a middleware that injects a request ID into the context
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-Id", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetReqID returns the request ID from the context
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// LoggerMiddleware - http logs to zerolog format
func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now().UTC()
			defer func() {
				if rec := recover(); rec != nil {
					log.Error().Timestamp().Interface("recover", rec).Bytes("stack", debug.Stack()).Msg("panic")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				fields := map[string]interface{}{
					"remote_ip":  r.RemoteAddr,
					"host":       r.Host,
					"path":       r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency":    time.Since(start).String(),
				}
				switch {
				case ww.Status() < 400:
					log.Info().Timestamp().Fields(fields).Msg("http")
				case ww.Status() < 500:
					log.Warn().Timestamp().Fields(fields).Msg("http")
				case ww.Status() > 500:
					log.Error().Timestamp().Fields(fields).Msg("http")
				}
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
