package utils

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// LoggerMiddleware - http logs to zerolog format
func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
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
