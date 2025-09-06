package logger

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type contextKey string

const loggerKey contextKey = "logger"

// Middleware that injects a request logger into context
func RequestLogger(baseLogger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a request ID
			reqID := uuid.New().String()

			// Create a child logger with request-specific fields
			reqLogger := baseLogger.With().
				Str("request_id", reqID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()

			// Add logger to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, loggerKey, &reqLogger)

			reqLogger.Info().Msg("incoming request")

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))

			reqLogger.Info().Msg("request completed")
		})
	}
}

// Helper to get logger from context
func LoggerFromContext(ctx context.Context) *zerolog.Logger {
	logger, ok := ctx.Value(loggerKey).(*zerolog.Logger)
	if !ok {
		// fallback to global logger
		l := zerolog.Nop()
		return &l
	}
	return logger
}
