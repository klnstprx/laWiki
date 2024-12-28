package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/laWiki/gateway/config"
	"github.com/rs/zerolog"
)

type key int

const requestIDKey key = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.App.FrontendURL) // Reemplaza con el dominio del frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Si es una petición interna, omitir la autenticación
		if r.Header.Get("X-Internal-Request") == "true" {
			next.ServeHTTP(w, r)
			return
		}

		// Si es una solicitud OPTIONS, respondemos inmediatamente
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// retrieve token from "jwt_token" cookie
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		// parse/validate using HS256
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.App.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}
		role, _ := claims["role"].(string)

		// put the role in context for route handlers
		ctx := context.WithValue(r.Context(), "user_role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
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
