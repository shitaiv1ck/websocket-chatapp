package core_middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
	originHeader    = "Origin"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]bool{
				"http://localhost:8080": true,
				"null":                  true,
			}

			origin := r.Header.Get(originHeader)

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestIDHeader)
			if id == "" {
				id = uuid.NewString()
				r.Header.Set(requestIDHeader, id)
			}

			w.Header().Set(requestIDHeader, id)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestIDHeader)

			logger := log.With(
				zap.String("request-id", id),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rw := core_repsponse.NewResponseWriter(w)

			log.Debug(">>> incoming http request")
			next.ServeHTTP(rw, r)
			log.Debug(
				"<<< done http request",
				zap.Int("status-code", rw.GetStatusCode()),
			)
		})
	}
}
