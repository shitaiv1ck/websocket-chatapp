package users_http_transport

import (
	"context"
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
	core_utils "github.com/shitaiv1ck/realtime-chat/internal/core/utils"
)

//go:generate mockgen -source=transport.go -destination=mocks/mock.go

type UsersHTTPTransport struct {
	service UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user domains.User) (domains.User, error)
	GetUsers(ctx context.Context, search *string, limit *int, offset *int) ([]domains.User, error)
	GetUser(ctx context.Context, userID int) (domains.User, error)
	PatchUser(ctx context.Context, userID int, patch domains.UserPatch) (domains.User, error)
}

func NewHTTPTransport(service UsersService) *UsersHTTPTransport {
	return &UsersHTTPTransport{
		service: service,
	}
}

func (t *UsersHTTPTransport) CreateUserHandler() http.HandlerFunc {
	type CreateUserRequest UserDTORequest
	type CreateUserResponse UserDTOResponse

	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		if log != nil {
			log.Debug("invoke CreateUser handler")
		}

		var request CreateUserRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		user := domains.User{
			Username: request.Username,
			Password: request.Password,
		}

		createdUser, err := t.service.CreateUser(r.Context(), user)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create user")

			return
		}

		response := CreateUserResponse{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			IsOnline: false,
		}

		responseHandler.JsonResponse(response, http.StatusCreated)
	}
}

func (t *UsersHTTPTransport) GetMeHandler() http.HandlerFunc {
	type GetMeResponse UserDTOResponse

	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		if log != nil {
			log.Debug("invoke GetMe handler")
		}

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		user, err := t.service.GetUser(r.Context(), userID)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		response := GetMeResponse{
			ID:       user.ID,
			Username: user.Username,
			IsOnline: true,
		}

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *UsersHTTPTransport) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		if log != nil {
			log.Debug("invoke CreateUser handler")
		}

		search := core_request.GetStringQueryParam(r, "search")

		limit, err := core_request.GetIntQueryParam(r, "limit")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get query param")

			return
		}

		offset, err := core_request.GetIntQueryParam(r, "offset")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get query param")

			return
		}

		users, err := t.service.GetUsers(r.Context(), search, limit, offset)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get users")

			return
		}

		response := ToDTOResponse(users)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *UsersHTTPTransport) PatchUserHandler() http.HandlerFunc {
	type PatchUserResponse UserDTOResponse

	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		if log != nil {
			log.Debug("invoke CreateUser handler")
		}

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		var request PatchUserRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		patch := domains.UserPatch{
			Username:    request.Username,
			OldPassword: request.OldPassword,
			NewPassword: request.NewPassword,
		}

		patchedUser, err := t.service.PatchUser(r.Context(), userID, patch)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to patch user")

			return
		}

		response := PatchUserResponse{
			ID:       patchedUser.ID,
			Username: patchedUser.Username,
			IsOnline: patchedUser.IsOnline,
		}

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}
