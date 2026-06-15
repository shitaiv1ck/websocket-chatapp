package messages_http_transport

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

type MessagesHTTPTransport struct {
	service MessagesService
}

type MessagesService interface {
	CreateMessage(ctx context.Context, message domains.Message) (domains.Message, error)
	GetMessages(ctx context.Context, userID int, chatID int) ([]domains.Message, error)
}

func NewHTTPTransport(service MessagesService) *MessagesHTTPTransport {
	return &MessagesHTTPTransport{
		service: service,
	}
}

func (t *MessagesHTTPTransport) CreateMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke CreateMessage handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		chatID, err := core_request.GetIntPathValue(r, "chat_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get chat id from path")

			return
		}

		var request CreateMessageRequest
		if err := core_request.DecodeAndValidate(r, &request); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate")

			return
		}

		message := domains.Message{
			ChatID:     chatID,
			SenderID:   userID,
			ReceiverID: request.ReceiverID,
			Content:    request.Content,
		}

		createdMessage, err := t.service.CreateMessage(r.Context(), message)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to create message")

			return
		}

		response := MessageDTOResponse{
			ID:         createdMessage.ID,
			ChatID:     createdMessage.ChatID,
			SenderID:   createdMessage.SenderID,
			ReceiverID: createdMessage.ReceiverID,
			Content:    createdMessage.Content,
			CreatedAt:  createdMessage.CreatedAt,
		}

		responseHandler.JsonResponse(response, http.StatusCreated)
	}
}

func (t *MessagesHTTPTransport) GetMessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())
		responseHandler := core_repsponse.NewResponseWriter(w)

		log.Debug("invoke GetMessages handler")

		userID, err := core_utils.GetIntFromContext(r.Context(), "user_id")
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorize, "failed to authorize")

			return
		}

		chatID, err := core_request.GetIntPathValue(r, "chat_id")
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get chat id from path")

			return
		}

		messages, err := t.service.GetMessages(r.Context(), userID, chatID)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get messages")

			return
		}

		response := ToDTOResponse(messages)

		responseHandler.JsonResponse(response, http.StatusOK)
	}
}
