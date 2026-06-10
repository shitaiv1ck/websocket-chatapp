package core_middleware

import (
	"context"
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
)

type SessionsService interface {
	GetSession(sessionToken string) (domains.Session, error)
}

const (
	csrfTokenHeader = "X-CSRF-Token"
)

func Authorization(service SessionsService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseHandler := core_repsponse.NewResponseWriter(w)

			sessionToken, err := r.Cookie("session_token")
			if err != nil {
				responseHandler.ErrorResponse(core_errors.ErrCoockie, "failed to get cookie")

				return
			}

			session, err := service.GetSession(sessionToken.Value)
			if err != nil {
				responseHandler.ErrorResponse(core_errors.ErrCoockie, "check cookie")

				return
			}

			if !core_request.IsMethodSafe(r.Method) {
				csrfToken := r.Header.Get(csrfTokenHeader)

				if csrfToken != session.CSRFToken {
					responseHandler.ErrorResponse(core_errors.ErrCoockie, "check cookie")

					return
				}
			}

			ctx := context.WithValue(r.Context(), "user_id", session.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
