package sessions_http_transport

import (
	"context"
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
)

type SessionsHTTPTransport struct {
	service SessionsService
}

type SessionsService interface {
	CreateSession(ctx context.Context, user domains.User) (domains.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

func NewHTTPTransport(service SessionsService) *SessionsHTTPTransport {
	return &SessionsHTTPTransport{
		service: service,
	}
}

func (s *SessionsHTTPTransport) CreateSessionHandler() http.HandlerFunc {
	type CreateSessionRequest UserDTORequest

	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateSession handler")

		var request CreateSessionRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		user := domains.User{
			Username: request.Username,
			Password: request.Password,
		}

		session, err := s.service.CreateSession(r.Context(), user)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create session")

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.SessionToken,
			Path:     "/",
			Expires:  session.ExpiredAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    session.CSRFToken,
			Path:     "/",
			Expires:  session.ExpiredAt,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})

		responseHandler.WriteHeader(http.StatusCreated)
	}
}

func (s *SessionsHTTPTransport) DeleteSessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke DeleteSession handler")

		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrCoockie, "failed to get cookie")

			return
		}

		if err := s.service.DeleteSession(r.Context(), sessionToken.Value); err != nil {
			responseHandler.ErrorResponse(err, "failed to delete session")

			return
		}

		responseHandler.WriteHeader(http.StatusNoContent)
	}
}
