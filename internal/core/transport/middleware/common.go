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
	requestID = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestID)
			if id == "" {
				id = uuid.NewString()
				r.Header.Set(requestID, id)
			}

			w.Header().Set(requestID, id)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestID)

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
