package messages_http_transport

import (
	"time"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type CreateMessageRequest struct {
	ReceiverID int    `json:"receiver_id" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

type MessageDTOResponse struct {
	ID         int       `json:"id"`
	ChatID     int       `json:"chat_id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToDTOResponse(messages []domains.Message) []MessageDTOResponse {
	response := make([]MessageDTOResponse, len(messages))

	for i, message := range messages {
		response[i] = MessageDTOResponse{
			ID:         message.ID,
			ChatID:     message.ChatID,
			SenderID:   message.SenderID,
			ReceiverID: message.ReceiverID,
			Content:    message.Content,
			CreatedAt:  message.CreatedAt,
		}
	}

	return response
}
