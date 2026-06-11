package friendrequests_http_transport

import (
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
	core_utils "github.com/shitaiv1ck/realtime-chat/internal/core/utils"
)

type FriendRequestsHTTPTransport struct {
	service FriendRequestsService
}

type FriendRequestsService interface {
	CreateFriendRequest(request domains.FriendRequest) (domains.FriendRequest, error)
	GetFriendRequests(userID int, direction *string) ([]domains.FriendRequest, error)
	DeleteFriendRequest(userID int, requestID int) error
}

func NewTransport(service FriendRequestsService) *FriendRequestsHTTPTransport {
	return &FriendRequestsHTTPTransport{
		service: service,
	}
}

func (t *FriendRequestsHTTPTransport) CreateFriendRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateRequest handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to autorize")

			return
		}

		var request CreateFriendRequestRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		friendRequest := domains.FriendRequest{
			FromUser: domains.UserBrief{ID: *userID},
			ToUser:   domains.UserBrief{ID: request.ToUserID},
		}

		createdFriendRequest, err := t.service.CreateFriendRequest(friendRequest)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create friend request")

			return
		}

		response := CreateFriendRequestResponse{
			ID:        createdFriendRequest.ID,
			FromUser:  UserDTO{ID: createdFriendRequest.FromUser.ID, Username: createdFriendRequest.FromUser.Username},
			ToUser:    UserDTO{ID: createdFriendRequest.ToUser.ID, Username: createdFriendRequest.ToUser.Username},
			CreatedAt: createdFriendRequest.CreatedAt,
		}

		responseHandler.JsonResponse(response, http.StatusCreated)
	}
}

func (t *FriendRequestsHTTPTransport) GetFriendRequestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke GetFriendRequests handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authenticate")

			return
		}

		direction := core_request.GetStringQueryParam(r, "direction")

		friendRequests, err := t.service.GetFriendRequests(*userID, direction)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get friend requests")

			return
		}

		response := friendRequestsToResponse(friendRequests)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *FriendRequestsHTTPTransport) DeleteFriendRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke DeleteFriendRequest handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		friendRequestID, err := core_request.GetIntPathValue(r, "friend_request_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get path value")

			return
		}

		if err := t.service.DeleteFriendRequest(*userID, *friendRequestID); err != nil {
			responseHandler.ErrorResponse(err, "failed to delete friend request")

			return
		}

		responseHandler.WriteHeader(http.StatusNoContent)
	}
}
