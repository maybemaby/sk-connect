package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justinas/alice"
)

var RequestIDHeader = "X-Request-ID"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type RequestLoggerContextKey string

const RequestLoggerKey RequestLoggerContextKey = "logger"

func CorsMiddleware(origin string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, X-User-Agent")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestIdMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			r.Header.Set(RequestIDHeader, uuid.New().String())
			next.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware(logger *slog.Logger) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{w, http.StatusOK}
			url := r.URL.String()
			requestId := r.Header.Get(RequestIDHeader)
			method := r.Method

			_logger := logger.WithGroup("request").With(slog.String("url", url), slog.String("method", method), slog.String("request_id", requestId))

			_logger.Info("Request received")

			context := context.WithValue(r.Context(), RequestLoggerKey, _logger)

			start := time.Now()
			next.ServeHTTP(rec, r.WithContext(context))

			_logger.Info("Request completed", slog.Duration("duration", time.Since(start)), slog.Int("status", rec.status))
		})
	}
}

func RootMiddleware(logger *slog.Logger, origin string) alice.Chain {
	return alice.New(RequestIdMiddleware(), LoggingMiddleware(logger), CorsMiddleware(origin))
}
