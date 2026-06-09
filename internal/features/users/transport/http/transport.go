package users_http_transport

import (
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
	"go.uber.org/zap"
)

type UsersHTTPTransport struct {
	service     UsersService
	wsTransport UsersWSTransport
}

type UsersWSTransport interface {
	BroadcastEvent(event string, content any) error
}

type UsersService interface {
	CreateUser(user domains.User) (domains.User, error)
	GetUsers() ([]domains.User, error)
	PatchUser(id int, patch domains.UserPatch) (domains.User, error)
}

func NewHTTPTransport(service UsersService, wsTransport UsersWSTransport) *UsersHTTPTransport {
	return &UsersHTTPTransport{
		service:     service,
		wsTransport: wsTransport,
	}
}

func (t *UsersHTTPTransport) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateUser handler")

		var request CreateUserRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		user := domains.User{
			Username: request.Username,
			Password: request.Password,
		}

		createdUser, err := t.service.CreateUser(user)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create user")

			return
		}

		response := CreateUserResponse{
			ID:       createdUser.ID,
			Username: createdUser.Username,
		}

		if err := t.wsTransport.BroadcastEvent("create_user", response); err != nil {
			log.Error("failed to broadcast", zap.Error(err))
		}

		responseHandler.JsonResponse(response, http.StatusCreated)
	}
}

func (t *UsersHTTPTransport) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke GetUsers handler")

		users, err := t.service.GetUsers()
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get users")

			return
		}

		response := usersToResponse(users)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *UsersHTTPTransport) PatchUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke PatchUser handler")

		// userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		// if err != nil {
		// 	responseHandler.ErrorResponse(core_errors.ErrCoockie, "failed to authenticate")

		// 	return
		// }

		userID, err := core_request.GetIntPathValue(r, "id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get path value")

			return
		}

		var request PatchUserRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		patch := domains.UserPatch{
			Username: request.Username,
			IsOnline: request.IsOnline,
		}

		patchedUser, err := t.service.PatchUser(*userID, patch)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to patch user")

			return
		}

		response := PatchUserResponse{
			ID:       patchedUser.ID,
			Username: patchedUser.Username,
			IsOnline: patchedUser.IsOnline,
		}

		if err := t.wsTransport.BroadcastEvent("patch_user", response); err != nil {
			log.Error("failed to broadcast", zap.Error(err))
		}

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}
