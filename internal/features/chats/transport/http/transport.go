package chats_http_transport

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

type ChatsHTTPTransport struct {
	service ChatsService
}

type ChatsService interface {
	CreateOrGetChat(ctx context.Context, userID int, friendID int) (domains.Chat, error)
	GetChats(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Chat, error)
	DeleteChat(ctx context.Context, userID int, chatID int) error
}

func NewHTTPTransport(service ChatsService) *ChatsHTTPTransport {
	return &ChatsHTTPTransport{
		service: service,
	}
}

func (t *ChatsHTTPTransport) CreateOrGetChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateOrGetChat handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		var request CreateOrGetChatRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		chat, err := t.service.CreateOrGetChat(r.Context(), userID, request.FriendID)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create or get chat")

			return
		}

		response := ChatDTOResponse{
			ID: chat.ID,
			FirstUser: UserDTOResponse{
				ID:       chat.FirstUser.ID,
				Username: chat.FirstUser.Username,
				IsOnline: chat.FirstUser.IsOnline,
			},
			SecondUser: UserDTOResponse{
				ID:       chat.SecondUser.ID,
				Username: chat.SecondUser.Username,
				IsOnline: chat.SecondUser.IsOnline,
			},
			LastMessageContent: chat.LastMessageContent,
			LastMessageAt:      chat.LastMessageAt,
		}

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *ChatsHTTPTransport) GetChatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke GetChats handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to autorize")

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

		chats, err := t.service.GetChats(r.Context(), userID, limit, offset)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get chats")

			return
		}

		response := ToDTOResponse(chats)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}

func (t *ChatsHTTPTransport) DeleteChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke DeleteChat handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to autorize")

			return
		}

		chatID, err := core_request.GetIntPathValue(r, "chat_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get chat id from path")

			return
		}

		if err := t.service.DeleteChat(context.Background(), userID, chatID); err != nil {
			responseHandler.ErrorResponse(err, "failed to delete chat")

			return
		}

		responseHandler.WriteHeader(http.StatusNoContent)
	}
}
