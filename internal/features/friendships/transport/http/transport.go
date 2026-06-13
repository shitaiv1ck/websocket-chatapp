package friendships_http_transport

import (
	"net/http"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_request "github.com/shitaiv1ck/realtime-chat/internal/core/transport/request"
	core_utils "github.com/shitaiv1ck/realtime-chat/internal/core/utils"
)

type FriendshipsHTTPTransport struct {
	service FriendshipsService
}

type FriendshipsService interface {
	CreateFriendship(userID int, requestID int) (domains.Friendship, error)
	GetFriendships(userID int, limit *int, offset *int) ([]domains.Friendship, error)
	DeleteFriendship(userID int, friendshipID int) error
}

func NewHTTPTransport(service FriendshipsService) *FriendshipsHTTPTransport {
	return &FriendshipsHTTPTransport{
		service: service,
	}
}

func (t *FriendshipsHTTPTransport) CreateFriendshipHandler() http.HandlerFunc {
	type CreateFriendshipResponse FriendshipDTOResponse

	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateFriendship handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")
		}

		var request CreateFriendshipRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		createdFriendship, err := t.service.CreateFriendship(userID, request.FriendRequestID)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create frienship")

			return
		}

		response := CreateFriendshipResponse{
			ID: createdFriendship.ID,
			FirstUser: UserDTOResponse{
				ID:       createdFriendship.FirstUser.ID,
				Username: createdFriendship.FirstUser.Username,
				IsOnline: createdFriendship.FirstUser.IsOnline,
			},
			SecondUser: UserDTOResponse{
				ID:       createdFriendship.SecondUser.ID,
				Username: createdFriendship.FirstUser.Username,
				IsOnline: createdFriendship.SecondUser.IsOnline,
			},
		}

		responseHandler.JsonResponse(response, http.StatusCreated)
	}
}

func (t *FriendshipsHTTPTransport) GetFriendshipsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke GetFriendships handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to authorize")

			return
		}

		limit, err := core_request.GetIntQueryParam(r, "limit")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get 'limit'")

			return
		}

		offset, err := core_request.GetIntQueryParam(r, "offset")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get 'offset'")

			return
		}

		friendships, err := t.service.GetFriendships(userID, limit, offset)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get friendships")

			return
		}

		response := ToDTOResponse(friendships)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *FriendshipsHTTPTransport) DeleteFriendshipHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke DeleteFriendship handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		friendshipID, err := core_request.GetIntPathValue(r, "friendship_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get friendship id")

			return
		}

		if err := t.service.DeleteFriendship(userID, *friendshipID); err != nil {
			responseHandler.ErrorResponse(err, "failed to delete friendship")

			return
		}

		responseHandler.WriteHeader(http.StatusNoContent)
	}
}
